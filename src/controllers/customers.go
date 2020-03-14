package controllers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
	"turkdev/src/data"
	"turkdev/src/models"
	"turkdev/src/services"
	"turkdev/src/shared"
	"turkdev/src/templates"
)

const (
	maxPlatformNameLength = 25
)

/*CustomerSignUpViewModel represents the data which is needed on sigin UI.*/
type CustomerSignUpViewModel struct {
	Name     string
	Email    string
	UserName string
	Password string
	Errors   map[string]string
}

/*InviteUserViewModel represents the data which is needed on sigin UI.*/
type InviteUserViewModel struct {
	EmailAddress   string
	SuccessMessage string
	Memo           string
	SignedInUser   *models.SignedInUserViewModel
	Errors         map[string]string
}

/*CustomerAdminViewModel represents the data which is needed on sigin UI.*/
type CustomerAdminViewModel struct {
	Name              string
	Domain            string
	LogoImageAsBase64 string
	Errors            map[string]string
	SuccessMessage    string
	SignedInUser      *models.SignedInUserViewModel
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

/*Validate validates the InviteUserViewModel*/
func (model *InviteUserViewModel) Validate(email string) (bool, error) {
	model.Errors = make(map[string]string)

	if strings.TrimSpace(model.EmailAddress) == "" {
		model.Errors["Email"] = "Email is required!"
	} else {
		if shared.IsEmailAdrressValid(model.EmailAddress) == false {
			model.Errors["Email"] = "Please enter a valid email address!"
		}
	}

	exists, err := data.ExistsUserByEmail(email)
	if err != nil {
		return false, err
	}
	if exists {
		model.Errors["Email"] = "This email address is already in use!"
	}
	return len(model.Errors) == 0, nil
}

/*Validate validates the InviteUserViewModel*/
func (model *CustomerAdminViewModel) Validate() bool {
	model.Errors = make(map[string]string)

	if strings.TrimSpace(model.Name) == "" {
		model.Errors["Name"] = "Name is required!"
	} else {
		if len(model.Name) > maxPlatformNameLength {
			model.Errors["Name"] = "Name cannot be longer than 25 characters"
		}
		if strings.Contains(model.Name, " ") {
			model.Errors["Name"] = "Name cannot contain spaces"
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
	model, err := setCustomerSignUpViewModel(r)
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

/*InviteUserHandler handles user invite operations*/
func InviteUserHandler(w http.ResponseWriter, r *http.Request) error {
	user := shared.GetUser(r)
	switch r.Method {
	case "GET":
		return handleInviteUserGET(w, r, user)
	case "POST":
		return handleInviteUserPOST(w, r, user)
	default:
		return handleInviteUserGET(w, r, user)
	}
}

func handleInviteUserGET(w http.ResponseWriter, r *http.Request, user *shared.SignedInUserClaims) error {
	err := templates.RenderInLayout(
		w,
		"invite.html",
		InviteUserViewModel{
			SignedInUser: &models.SignedInUserViewModel{
				UserName:   user.UserName,
				UserID:     user.ID,
				CustomerID: user.CustomerID,
				Email:      user.Email,
			},
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func handleInviteUserPOST(w http.ResponseWriter, r *http.Request, user *shared.SignedInUserClaims) error {
	model := &InviteUserViewModel{
		EmailAddress: r.FormValue("email"),
		Memo:         r.FormValue("memo"),
		SignedInUser: &models.SignedInUserViewModel{
			UserName:   user.UserName,
			UserID:     user.ID,
			CustomerID: user.CustomerID,
			Email:      user.Email,
		},
	}

	inviteHTMLPath := "invite.html"
	isValid, err := model.Validate(model.EmailAddress)
	if err != nil {
		return err
	}
	if !isValid {
		err := templates.RenderInLayout(w, inviteHTMLPath, model)
		if err != nil {
			return err
		}
		return nil
	}

	inviteCode, err := data.CreateInviteCode(user.ID, model.EmailAddress)
	if err != nil {
		return err
	}

	err = services.SendInvitemail(model.EmailAddress, model.Memo, inviteCode, user.UserName)
	if err != nil {
		return err
	}
	model.SuccessMessage = "Inivitation mail successfully sent to " + model.EmailAddress
	err = templates.RenderInLayout(w, inviteHTMLPath, model)
	if err != nil {
		return err
	}
	return nil
}

/*AdminHandler handles admin operations*/
func AdminHandler(w http.ResponseWriter, r *http.Request) error {
	user := shared.GetUser(r)
	switch r.Method {
	case "GET":
		return handleAdminGET(w, r, user)
	case "POST":
		return handleAdminPOST(w, r, user)
	default:
		return handleAdminGET(w, r, user)
	}
}

func handleAdminGET(w http.ResponseWriter, r *http.Request, user *shared.SignedInUserClaims) error {
	isAdmin, err := data.IsUserAdmin(user.ID)
	if err != nil {
		return err
	}

	if !isAdmin {
		err = templates.RenderFile(w, "errors/500.html", nil)
		if err != nil {
			return err
		}
		return nil
	}

	customer, err := data.GetCustomerByID(user.CustomerID)
	if err != nil {
		return err
	}

	model, err := setCustomerAdminViewModel(customer, user)
	if err != nil {
		return err
	}

	if model.LogoImageAsBase64 != "" {
		imageasB64, err := decodeLogoImageToBase64(customer.LogoImage)
		if err != nil {
			return err
		}
		model.LogoImageAsBase64 = imageasB64
	}

	err = templates.RenderInLayout(w, "admin.html", model)
	if err != nil {
		return err
	}
	return nil
}

func handleAdminPOST(w http.ResponseWriter, r *http.Request, user *shared.SignedInUserClaims) error {
	file, err := getImageFile(r)
	if err != nil {
		return err
	}

	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	model := &CustomerAdminViewModel{
		Name:              r.FormValue("name"),
		Domain:            r.FormValue("domain"),
		LogoImageAsBase64: base64.StdEncoding.EncodeToString(fileBytes),
	}

	adminHTMLPath := "admin.html"
	if model.Validate() == false {
		err := templates.RenderInLayout(w, adminHTMLPath, model)
		if err != nil {
			return err
		}
		return nil
	}

	customer, err := data.GetCustomerByID(user.CustomerID)
	if err != nil {
		return err
	}

	if customer.Name != model.Name {
		exists, err := data.ExistsCustomerByName(model.Name)
		if err != nil {
			return err
		}
		if exists {
			model.Errors["Name"] = "This name is already taken"
			err := templates.RenderInLayout(w, adminHTMLPath, model)
			if err != nil {
				return err
			}
			return nil
		}

	}

	if model.Domain != "" {
		if customer.Domain != model.Domain {
			exists, err := data.ExistsCustomerByDomain(model.Domain)
			if err != nil {
				return err
			}
			if exists {
				model.Errors["Domain"] = "This domain is already taken"
				err := templates.RenderInLayout(w, adminHTMLPath, model)
				if err != nil {
					return err
				}
				return nil
			}
		}
	}

	setUpdatedCustomerByModel(model, customer)
	err = data.UpdateCustomer(customer)
	if err != nil {
		return err
	}

	model.SuccessMessage = "Account updated successfuly"
	err = templates.RenderInLayout(w, adminHTMLPath, model)
	if err != nil {
		return err
	}
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

func setCustomerSignUpViewModel(r *http.Request) (*CustomerSignUpViewModel, error) {
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

func setCustomerAdminViewModel(customer *data.Customer, user *shared.SignedInUserClaims) (*CustomerAdminViewModel, error) {
	var model CustomerAdminViewModel
	model.Domain = customer.Domain
	model.LogoImageAsBase64 = base64.StdEncoding.EncodeToString(customer.LogoImage)
	model.Name = customer.Name

	var s models.SignedInUserViewModel
	s.CustomerID = customer.ID
	s.Email = customer.Email
	s.UserID = user.ID
	s.UserName = user.UserName
	model.SignedInUser = &s
	return &model, nil
}

func setUpdatedCustomerByModel(model *CustomerAdminViewModel, customer *data.Customer) {
	customer.Name = model.Name
	customer.Domain = model.Domain
	if model.LogoImageAsBase64 == "" {
		return
	}
	customer.LogoImage = []byte(model.LogoImageAsBase64)
}

func getImageFile(r *http.Request) (multipart.File, error) {
	file, _, err := r.FormFile("logo")
	if err != nil {
		return file, err
	}
	return file, err
}

func decodeLogoImageToBase64(logoImage []byte) (string, error) {
	var img image.Image
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(string(logoImage)))
	img, _, err := image.Decode(reader)
	if err != nil {
		return "", err
	}
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, img, nil); err != nil {
		log.Fatalln("unable to encode image.")
	}
	return base64.StdEncoding.EncodeToString(buffer.Bytes()), err

}
