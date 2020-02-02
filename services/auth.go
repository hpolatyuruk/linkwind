package services

import (
	"fmt"
	"regexp"
	"time"
	"turkdev/data"

	"github.com/dgrijalva/jwt-go"
)

const (
	LoginError = iota
	LoginSuccessful
	WrongPassword
	NoUserWithEmail
	NoUserWithUserName
	TokenNotCreated
)

const (
	slash              = "`"
	regexForEmailValid = `^(?("")("".+?(?<!\\)""@)|(([0-9a-z]((\.(?!\.))|[-!#\$%&'\*\+/=\?\^` + slash + `\{\}\|~\w])*)(?<=[0-9a-z])@))(?(\[)(\[(\d{1,3}\.){3}\d{1,3}\])|(([0-9a-z][-\w]*[0-9a-z]*\.)+[a-z0-9][\-a-z0-9]{0,22}[a-z0-9]))$`
)

func LoginUser(emailOrUserName, password string) (int, string, error) {
	if data.IsEmailAdrressValid(emailOrUserName) {
		exists, err := ExistsUserByEmail(emailOrUserName)
		if err != nil {
			return LoginError,"", err
		}

		if exists {
			user, err := data.FindUserByEmailAndPassword(emailOrUserName, password)
			if err != nil {
				return LoginError,"", err
			}

			if user == nil {
				return WrongPassword,"", nil
			}

			jwtToken, err := user.GenerateAuthToken()
			if err != nil {
				return TokenNotCreated,"", err
			}
			return LoginSuccessful, jwtToken, nil
		}
		return NoUserWithEmail,"", nil
	}

	exists, err := data.ExistsUserByUserName(emailOrUserName)
	if err != nil {
		return LoginError,"", err
	}

	if exists {
		user, err := data.FindUserByUserNameAndPassword(emailOrUserName, password)
		if err != nil {
			return LoginError,"", err
		}

		if user == nil {
			return WrongPassword,"", nil
		}

		jwtToken, err := user.GenerateAuthToken()
		if err != nil {
			return TokenNotCreated,"", err
		}
		return LoginSuccessful,jwtToken, nil
	}
	return NoUserWithUserName,"", nil
}

/*IsEmailAdrressValid take mail adrres, if adrress is valid return true.*/
func IsEmailAdrressValid(email string) bool {
	Re := regexp.MustCompile(regexForEmailValid)
	return Re.MatchString(email)
}

/*GenerateAuthToken generate jwt token for user login. */
func (user *data.User)GenerateAuthToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"email":user.Email,
		"username": user.UserName,
		"fullname": user.FullName
	})

	privateKey := os.Getenv("JWTPrivateKey")
	return token.SignedString([]byte(privateKey))
}
