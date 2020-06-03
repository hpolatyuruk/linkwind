package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"linkwind/app/models"
	"linkwind/app/shared"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/getsentry/sentry-go"
	"github.com/oxtoacart/bpool"
)

const templatesDir = "templates/"

var bufpool *bpool.BufferPool
var templates map[string]*template.Template

// Load templates on program initialisation
func init() {
	bufpool = bpool.NewBufferPool(64)

	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	layouts, err := listAllHtmlsRecursively(templatesDir + "layouts")
	if err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
	}

	partials, err := listAllHtmlsRecursively(templatesDir + "partials")
	if err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
	}

	layoutPath := path.Join(templatesDir, "layouts", "layout.html")

	// Generate our templates map from our layouts/ and partials/ directories
	for _, layout := range layouts {
		fileName := filepath.Base(layout)
		files := append(partials, layout)
		files = append(files, layoutPath)
		templates[fileName] = template.Must(template.ParseFiles(files...))
	}
}

// RenderInLayout is a wrapper around template.ExecuteTemplate.
// It writes into a bytes.Buffer before writing to the http.ResponseWriter to catch
// any errors resulting from populating the template.
func RenderInLayout(
	w http.ResponseWriter,
	r *http.Request,
	name string,
	data models.BaseViewModelInterface) error {

	tmpl, ok := templates[name]
	if !ok {
		return fmt.Errorf("The template %s does not exist", name)
	}

	userCtx := shared.GetUserFromContext(r)
	customerCtx := shared.GetCustomerFromContext(r)

	capitailizedPlatform := strings.Title(customerCtx.Platform)

	data.SetLayout(capitailizedPlatform, customerCtx.Logo)
	data.SetSignedInUser(userCtx)

	// Create a buffer to temporarily write to and check if any errors were encounted.
	buf := bufpool.Get()
	defer bufpool.Put(buf)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tmpl.ExecuteTemplate(buf, name, data)
	if err != nil {
		return err
	}
	// Set the header and write the buffer to the http.ResponseWriter
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
	return nil
}

/*RenderFile renders a template file excluded from base template.*/
func RenderFile(
	w http.ResponseWriter,
	tmplPath string,
	data interface{}) error {

	tmpl, err := template.ParseFiles(path.Join(templatesDir, tmplPath))
	if err != nil {
		return err
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		return err
	}
	return nil
}

/*RenderAsString parses the template and return result as string.*/
func RenderAsString(tmplPath string, tmplName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(path.Join(templatesDir, tmplPath))
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = t.ExecuteTemplate(buf, tmplName, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func listAllHtmlsRecursively(dir string) ([]string, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}

		// check if it is a regular file (not dir)
		if info.Mode().IsRegular() {
			if filepath.Ext(path) == ".html" {
				paths = append(paths, path)
			}
		}
		return nil
	})
	return paths, err
}
