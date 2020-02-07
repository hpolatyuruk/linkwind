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

func SignInHandler(w http.ResponseWriter, r *http.Request) error {
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

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleSignUpGET(w, r)
	case "POST":
		handleSignUpPOST(w, r)
	default:
		handleSignUpGET(w, r)
	}
}

/*LogoutHandler set cookie authenticated is false.*/
func LogoutHandler(w http.ResponseWriter, r *http.Request) error {
	session, _ := store.Get(r, "cookie-name")
	session.Values["authenticated"] = false
	session.Save(r, w)

	// TODO : Redirect bla bla.
	return nil
}

func handleSignUpGET(w http.ResponseWriter, r *http.Request) {
	title := "Sign Up | Turk Dev"
	user := models.User{"Anil Yuzener"}

	data := map[string]interface{}{
		"Content": "Sign Up",
	}

	templates.Render(
		w,
		"users/signup.html",
		models.ViewModel{
			title,
			user,
			data,
		},
	)
}

func handleSignUpPOST(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Handling sign up post error: %v", err)
	}

	userName := r.FormValue("userName")
	email := r.FormValue("email")
	password := r.FormValue("password")

	var user data.User
	user.UserName = userName
	user.Email = email
	user.Password = password
	user.Karma = 0
	user.RegisteredOn = time.Now()
	user.CustomerID = 1 // TODO: get it real customer id
	err := data.CreateUser(&user)

	if err != nil {
		// TODO: log it here
		fmt.Fprintf(w, "Error creating user: %v", err)
	}
	fmt.Fprintf(w, "Succeded!")
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
