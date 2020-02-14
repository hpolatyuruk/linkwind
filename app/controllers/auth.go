package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"turkdev/app/models"
	"turkdev/app/src/templates"
	"turkdev/data"
	"turkdev/shared"
)

/*SignInViewModel represents the data which is needed on sigin UI.*/
type SignInViewModel struct {
	EmailOrUserName string
	Password        string
	Errors          map[string]string
}

/*Validate validates the SignInViewModel*/
func (model *SignInViewModel) Validate() bool {
	model.Errors = make(map[string]string)

	if strings.TrimSpace(model.EmailOrUserName) == "" {
		model.Errors["EmailOrUserName"] = "Email or user name is required!"
	}
	if strings.TrimSpace(model.Password) == "" {
		model.Errors["Password"] = "Password is required!"
	}

	return len(model.Errors) == 0
}

/*SignInHandler handles user signin operations.*/
func SignInHandler(w http.ResponseWriter, r *http.Request) error {

	switch r.Method {
	case "GET":
		isAuthenticated, _, err := shared.IsAuthenticated(r)
		if err != nil {
			return err
		}
		if isAuthenticated {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return nil
		}
		return handleSignInGET(w, r)
	case "POST":
		return handleSignInPOST(w, r)
	default:
		return handleSignInGET(w, r)
	}
}

/*SignUpHandler handles user signup operations*/
func SignUpHandler(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return handleSignUpGET(w, r)
	case "POST":
		return handleSignUpPOST(w, r)
	default:
		return handleSignUpGET(w, r)
	}
}

/*SignOutHandler handles user singout operations.*/
func SignOutHandler(w http.ResponseWriter, r *http.Request) error {
	shared.SetAuthCookie(w, "", time.Now())
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func handleSignUpGET(w http.ResponseWriter, r *http.Request) error {
	templates.RenderInLayout(
		w,
		"signup.html",
		models.ViewModel{
			Title: "Sign Up",
		},
	)
	return nil
}

func handleSignUpPOST(w http.ResponseWriter, r *http.Request) error {
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
		return err
	}
	return nil
}

func handleSignInGET(w http.ResponseWriter, r *http.Request) error {
	err := templates.RenderFile(
		w,
		"layouts/users/signin.html",
		SignInViewModel{},
	)
	if err != nil {
		return err
	}
	return nil
}

func handleSignInPOST(w http.ResponseWriter, r *http.Request) error {
	model := &SignInViewModel{
		EmailOrUserName: r.FormValue("emailOrUserName"),
		Password:        r.FormValue("password"),
	}

	if model.Validate() == false {
		err := templates.RenderFile(w, "users/signin.html", model)
		if err != nil {
			return err
		}
		return nil
	}

	var err error
	var user *data.User

	if shared.IsEmailAdrressValid(model.EmailOrUserName) {
		user, err = checkUserByEmail(model.EmailOrUserName, model.Password)
	} else {
		user, err = checkUserByUserName(model.EmailOrUserName, model.Password)
	}
	if err != nil {
		return err
	}
	if user == nil {
		model.Errors["General"] = "User does not exist!"
		err = templates.RenderFile(w, "users/signin.html", model)
		if err != nil {
			return err
		}
		return nil
	}
	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(1440 * time.Minute)

	token, err := shared.GenerateAuthToken(*user, expirationTime)
	if err != nil {
		return err
	}
	shared.SetAuthCookie(w, token, expirationTime)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func checkUserByEmail(email string, password string) (*data.User, error) {
	exists, err := data.ExistsUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if exists == false {
		return nil, nil
	}
	user, err := data.FindUserByEmailAndPassword(email, password)
	if err != nil {
		return nil, err
	}
	return user, err
}

func checkUserByUserName(userName string, password string) (*data.User, error) {
	exists, err := data.ExistsUserByUserName(userName)
	if err != nil {
		return nil, err
	}
	if exists == false {
		return nil, nil
	}
	user, err := data.FindUserByUserNameAndPassword(userName, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
