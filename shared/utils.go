package shared

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"regexp"
	"time"
	"unicode"

	"golang.org/x/net/html"
)

/*IsEmailAdrressValid takes mail address, if address is valid return true.*/
func IsEmailAdrressValid(email string) bool {
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
func FetchURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("Error occured when get url response. Error: %v", err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Fail to parse html. Error: %v", err)
	}

	title, ok := traverse(doc)
	if !ok {
		return "", fmt.Errorf("Cannot parse title")
	}

	return title, nil
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

func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
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
