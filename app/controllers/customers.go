package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"turkdev/app/src/templates"
	"turkdev/data"
	"turkdev/shared"
)

/*CustomerSignUpViewModel represents the data which is needed on sigin UI.*/
type CustomerSignUpViewModel struct {
	Name     string
	Email    string
	UserName string
	Password string
	Errors   map[string]string
}

/*Validate validates the CustomerSignUpViewModel*/
func (model *CustomerSignUpViewModel) Validate() bool {
	model.Errors = make(map[string]string)

	if strings.TrimSpace(model.Name) == "" {
		model.Errors["Name"] = "Name is required!"
	}
	if strings.TrimSpace(model.UserName) == "" {
		model.Errors["UserName"] = "User name is required!"
	} else {
		if len(model.UserName) > userNameMaxCharCount {
			model.Errors["UserName"] = fmt.Sprintf("User name can not be longer than %d!", userNameMaxCharCount)
		}
	}
	if strings.TrimSpace(model.Email) == "" {
		model.Errors["Email"] = "Email is required!"
	} else {
		if shared.IsEmailAdrressValid(model.Email) == false {
			model.Errors["Email"] = "Please enter a valid email address!"
		}
	}
	if strings.TrimSpace(model.Password) == "" {
		model.Errors["Password"] = "Password is required!"
	} else {
		if shared.IsPasswordValid(model.Password) == false {
			model.Errors["Password"] = "The password is not valid. A password should contan at least 1 uppercase, 1 lowercase, 1 digit, one of #$+=!*@& special characters and have a length of at least of 8."
		}
	}
	return len(model.Errors) == 0
}

/*CustomerSignUpHandler handles customer signup operations*/
func CustomerSignUpHandler(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return handleCustomerSignUpGET(w, r)
	case "POST":
		return handleCustomerSignUpPOST(w, r)
	default:
		return handleCustomerSignUpGET(w, r)
	}
}

func handleCustomerSignUpGET(w http.ResponseWriter, r *http.Request) error {
	err := templates.RenderFile(
		w,
		"layouts/customers/signup.html",
		CustomerSignUpViewModel{},
	)
	if err != nil {
		return err
	}
	return nil
}

func handleCustomerSignUpPOST(w http.ResponseWriter, r *http.Request) error {
	model, err := setCustomerSingUpViewModel(r)
	if err != nil {
		return err
	}

	signUpHTMLPath := "layouts/customers/signup.html"
	if model.Validate() == false {
		err := templates.RenderFile(w, signUpHTMLPath, model)
		if err != nil {
			return err
		}
		return nil
	}
	// TODO: Db sorgulamaları daha sonra refoctor edilecek.
	exists, err := data.ExistsCustomerByName(model.Name)
	if err != nil {
		return err
	}
	if exists {
		model.Errors["Name"] = "Name is already taken!"
		templates.RenderFile(w, signUpHTMLPath, model)
		return nil
	}

	exists, err = data.ExistsCustomerByEmail(model.Email)
	if err != nil {
		return err
	}
	if exists {
		model.Errors["Email"] = "The user associated with this email already exists!"
		templates.RenderFile(w, signUpHTMLPath, model)
		return nil
	}

	exists, err = data.ExistsUserByUserName(model.UserName)
	if err != nil {
		return err
	}
	if exists {
		model.Errors["UserName"] = "User name is already taken!"
		templates.RenderFile(w, signUpHTMLPath, model)
		return nil
	}

	exists, err = data.ExistsUserByEmail(model.Email)
	if err != nil {
		return err
	}
	if exists {
		model.Errors["Email"] = "The user associated with this email already exists!"
		templates.RenderFile(w, signUpHTMLPath, model)
		return nil
	}

	customer := setCustomerByModel(model)
	err = data.CreateCustomer(&customer)
	if err != nil {
		return err
	}

	addedCustomer, err := data.GetCustomerByName(model.Name)
	if err != nil {
		return err
	}

	user := setUserByModel(model, addedCustomer)
	err = data.CreateUser(&user)
	if err != nil {
		return err
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(authExpirationMinutes * time.Minute)

	token, err := shared.GenerateAuthToken(user, expirationTime)
	if err != nil {
		return err
	}
	shared.SetAuthCookie(w, token, expirationTime)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func setCustomerByModel(model *CustomerSignUpViewModel) data.Customer {
	var customer data.Customer
	customer.Email = model.Email
	customer.Name = model.Name
	customer.RegisteredOn = time.Now()
	customer.LogoImage = nil
	return customer
}

func setUserByModel(model *CustomerSignUpViewModel, addedCustomer *data.Customer) data.User {
	var user data.User
	user.CustomerID = addedCustomer.ID
	user.Email = model.Email
	user.Karma = 0
	user.Password = model.Password
	user.RegisteredOn = time.Now()
	user.UserName = model.UserName
	return user
}

func setCustomerSingUpViewModel(r *http.Request) (*CustomerSignUpViewModel, error) {
	var model CustomerSignUpViewModel
	if err := r.ParseForm(); err != nil {
		return &model, err
	}
	model.Name = r.FormValue("name")
	model.Email = r.FormValue("email")
	model.UserName = r.FormValue("userName")
	model.Password = r.FormValue("password")
	return &model, nil
}
