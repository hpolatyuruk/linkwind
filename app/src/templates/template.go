package templates

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

var templatesPath = "./app/src/templates/"
var categories = []string{"comments", "stories", "users"}
var pseudoTmpl string = `{{define "main"}}{{template "base" .}}{{end}}`
var templates map[string]*template.Template

/*Initialize initializes all base and child templates*/
func Initialize() error {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	baseTemplates, err := filepath.Glob(templatesPath + "base/*.html")
	if err != nil {
		return err
	}

	base := template.New("main")
	base, err = base.Parse(pseudoTmpl)
	if err != nil {
		return err
	}

	for _, category := range categories {
		pages, err := filepath.Glob(templatesPath + category + "/*.html")
		if err != nil {
			return err
		}
		for _, page := range pages {
			f := category + "/" + filepath.Base(page)
			files := append(baseTemplates, page)
			templates[f], err = base.Clone()
			if err != nil {
				return err
			}
			templates[f], err = templates[f].ParseFiles(files...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

/*Render renders templates by given parameters*/
func Render(w http.ResponseWriter, tmpl string, data interface{}) {
	t, found := templates[tmpl]

	if !found {
		http.Error(w, "Cannot find template", http.StatusInternalServerError)
		return
	}

	err := t.Execute(w, data)
	if err != nil {
		fmt.Print(err)
		http.Error(w, "Cannot parse template", http.StatusInternalServerError)
	}
}
