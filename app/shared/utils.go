package shared

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
	"unicode"

	"github.com/getsentry/sentry-go"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

const (
	Charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

/*SeededRand is help to create random values by time*/
var SeededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

/*IsEmailAdressValid takes mail address, if address is valid return true.*/
func IsEmailAdressValid(email string) bool {
	Re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return Re.MatchString(email)
}

/*IsPasswordValid takes password, if password is valid return true.*/
func IsPasswordValid(password string) bool {

	/*
		Note: Since golang uses different regex validation library, regex for password validation causes exception for this pattern: "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[#$^+=!*()@%&]).{8,}$"
		Since we don't have time for that, we use the below method to verify password
		 *
		 * Password rules:
		 * at least 8 letters
		 * at least 1 number
		 * at least 1 upper case
		 * at least 1 lower case
		 * at least one of #$+=!*@&_ special characters
	*/
	var number, upperCase, lowerCase, special, eigthOrMore bool
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upperCase = true
		case unicode.IsLower(c):
			lowerCase = true
		case c == '#' ||
			c == '$' ||
			c == '&' ||
			c == '+' ||
			c == '=' ||
			c == '!' ||
			c == '@' ||
			c == '*' ||
			c == '_':
			special = true
		default:
			//return false, false, false, false
		}
	}
	eigthOrMore = len(password) >= 8
	return number && upperCase && lowerCase && special && eigthOrMore
}

/*FetchURL send request to url that given as parameter and fetch title from HTML code.*/
func FetchURL(url string) (title string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Error occured when get url response. Error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("response status code: %d", resp.StatusCode)
		return
	}

	r, _, _, err := convertHTMLToUTF8(resp.Body)

	doc, err := html.Parse(r)
	if err != nil {
		return "", fmt.Errorf("Fail to parse html. Error: %v", err)
	}

	title, ok := traverse(doc)
	if !ok {
		return "", fmt.Errorf("Cannot parse title")
	}
	fmt.Println(title)
	return title, nil
}

/*
For more info: https://siongui.github.io/2018/10/27/auto-detect-and-convert-html-encoding-to-utf8-in-go/
*/
func convertHTMLToUTF8(body io.Reader) (r io.Reader, name string, certain bool, err error) {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	e, name, certain, err := determineEncodingFromReader(bytes.NewReader(b))
	if err != nil {
		return
	}

	r = transform.NewReader(bytes.NewReader(b), e.NewDecoder())
	return
}

func determineEncodingFromReader(r io.Reader) (e encoding.Encoding, name string, certain bool, err error) {
	b, err := bufio.NewReader(r).Peek(1024)
	if err != nil {
		return
	}

	e, name, certain = charset.DetermineEncoding(b, "")
	return
}

func traverse(n *html.Node) (string, bool) {
	if isTitleElement(n) {
		return n.FirstChild.Data, true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := traverse(c)
		if ok {
			return result, ok
		}
	}

	return "", false
}

func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
}

/*ReadFile reads file content from given path.*/
func ReadFile(filePath string) ([]byte, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		return nil, fmt.Errorf("Error occured when open %s file. Original err: %v", filePath, err)
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("Error occured when read %s file. Original err: %v", filePath, err)
	}

	return byteValue, err
}

func setCookie(w http.ResponseWriter, userNameOrEmail, password string) {
	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{
		Name:    userNameOrEmail,
		Value:   password,
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
}

/*DateToString converts date to user friendly string. eg: 15 days ago*/
func DateToString(submittedOn time.Time) string {
	var text string = ""
	diff := time.Now().Sub(submittedOn)

	if diff.Hours() < 1 {
		mins := int(math.Round(diff.Minutes()))
		text = fmt.Sprintf("%d minutes ago", mins)
		if mins == 1 {
			text = fmt.Sprintf("%d minute ago", mins)
		}
	} else if diff.Hours() < 24 {
		hours := int(math.Round(diff.Hours()))
		text = fmt.Sprintf("%d hours ago", hours)
		if hours == 1 {
			text = fmt.Sprintf("%d hour ago", hours)
		}
	} else {
		days := math.Round(diff.Hours() / 24)

		if days == 1 {
			text = fmt.Sprintf("%d day ago", int(days))
		} else if days > 1 && days < 30 {
			text = fmt.Sprintf("%d days ago", int(days))
		} else if days > 30 && days < 365 {
			months := int(math.Round(days / 30))
			text = fmt.Sprintf("%d months ago", months)
			if months == 1 {
				text = fmt.Sprintf("%d month ago", months)
			}

		} else {
			years := int(math.Round(days / 365))
			text = fmt.Sprintf("%d years ago", years)
			if years == 1 {
				text = fmt.Sprintf("%d year ago", years)
			}
		}
	}
	return text
}

/*GenerateResetPasswordToken generates password token for reset*/
func GenerateResetPasswordToken() string {
	c := ""
	for i := 0; i < 4; i++ {
		s := StringWithCharset(1)
		i := SeededRand.Intn(10)
		c = c + strconv.Itoa(i) + s
	}
	return c
}

/*StringWithCharset generate random string by length*/
func StringWithCharset(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = Charset[SeededRand.Intn(len(Charset))]
	}
	return string(b)
}

/*ReturnNotFoundTemplate writes 404 not found html template to given response*/
func ReturnNotFoundTemplate(w http.ResponseWriter) {
	byteValue, err := ReadFile("templates/errors/404.html")
	if err != nil {
		sentry.CaptureException(err)
		http.Error(w, "Unexpected error!", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write(byteValue)
}
