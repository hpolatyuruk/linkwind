package shared

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	slash              = "`"
	regexForEmailValid = `^(?("")("".+?(?<!\\)""@)|(([0-9a-z]((\.(?!\.))|[-!#\$%&'\*\+/=\?\^` + slash + `\{\}\|~\w])*)(?<=[0-9a-z])@))(?(\[)(\[(\d{1,3}\.){3}\d{1,3}\])|(([0-9a-z][-\w]*[0-9a-z]*\.)+[a-z0-9][\-a-z0-9]{0,22}[a-z0-9]))$`
)

/*IsEmailAdrressValid take mail adrres, if adrress is valid return true.*/
func IsEmailAdrressValid(email string) bool {
	Re := regexp.MustCompile(regexForEmailValid)
	return Re.MatchString(email)
}

/*FetchURL send request to url that given as parameter and fetch title from HTML code.*/
func FetchURL(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	dataInBytes, err := ioutil.ReadAll(response.Body)
	pageContent := string(dataInBytes)

	titleStartIndex := strings.Index(pageContent, "<title>")
	if titleStartIndex == -1 {
		return "", fmt.Errorf("No title element found")
	}
	// The start index of the title is the index of the first
	// character, the < symbol. We don't want to include
	// <title> as part of the final value, so let's offset
	// the index by the number of characers in <title>
	titleStartIndex += 7

	titleEndIndex := strings.Index(pageContent, "</title>")
	if titleEndIndex == -1 {
		fmt.Println("No closing tag for title found.")
		os.Exit(0)
	}

	pageTitle := []byte(pageContent[titleStartIndex:titleEndIndex])
	return string(pageTitle), nil
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
