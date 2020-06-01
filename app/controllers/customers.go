package controllers

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	cache "linkwind/app/caching"
	"linkwind/app/data"
	"linkwind/app/models"
	"linkwind/app/shared"
	"linkwind/app/templates"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
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
		if shared.IsEmailAdressValid(model.Email) == false {
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

/*ExistsCustomDomain checks wheter provided custom domain from url query exists or not. If it exists returns 200, othwesise 404*/
func ExistsCustomDomain(w http.ResponseWriter, r *http.Request) {

	domain := r.URL.Query().Get("domain")
	if domain == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Domain cannot be empty."))
		return
	}
	exists, err := data.ExistsCustomerByDomain(domain)
	if err != nil {
		panic(err)
	}
	if exists == false {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

/*CustomerSignUpHandler handles customer signup operations*/
func CustomerSignUpHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleCustomerSignUpGET(w, r)
	case "POST":
		handleCustomerSignUpPOST(w, r)
	default:
		handleCustomerSignUpGET(w, r)
	}
}

func handleCustomerSignUpGET(w http.ResponseWriter, r *http.Request) {
	err := templates.RenderFile(
		w,
		"layouts/customers/signup.html",
		CustomerSignUpViewModel{},
	)
	if err != nil {
		panic(err)
	}
}

func handleCustomerSignUpPOST(w http.ResponseWriter, r *http.Request) {
	model, err := setCustomerSignUpViewModel(r)
	if err != nil {
		panic(err)
	}

	signUpHTMLPath := "layouts/customers/signup.html"
	if model.Validate() == false {
		err := templates.RenderFile(w, signUpHTMLPath, model)
		if err != nil {
			panic(err)
		}
		return
	}
	exists, err := data.ExistsCustomerByName(model.Name)
	if err != nil {
		panic(err)
	}
	if exists {
		model.Errors["Name"] = "Name is already taken!"
		templates.RenderFile(w, signUpHTMLPath, model)
		return
	}

	exists, err = data.ExistsCustomerByEmail(model.Email)
	if err != nil {
		panic(err)
	}
	if exists {
		model.Errors["Email"] = "The user associated with this email already exists!"
		templates.RenderFile(w, signUpHTMLPath, model)
		return
	}

	exists, err = data.ExistsUserByUserName(model.UserName)
	if err != nil {
		panic(err)
	}
	if exists {
		model.Errors["UserName"] = "User name is already taken!"
		templates.RenderFile(w, signUpHTMLPath, model)
		return
	}

	exists, err = data.ExistsUserByEmail(model.Email)
	if err != nil {
		panic(err)
	}
	if exists {
		model.Errors["Email"] = "The user associated with this email already exists!"
		templates.RenderFile(w, signUpHTMLPath, model)
		return
	}

	customer := setCustomerByModel(model)
	err = data.CreateCustomer(&customer)
	if err != nil {
		panic(err)
	}

	addedCustomer, err := data.GetCustomerByName(model.Name)
	if err != nil {
		panic(err)
	}

	user := setUserByModel(model, addedCustomer)
	userID, err := data.CreateUser(&user)
	if err != nil {
		panic(err)
	}
	user.ID = *userID

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(authExpirationMinutes * time.Minute)

	token, err := shared.GenerateAuthToken(user, expirationTime)
	if err != nil {
		panic(err)
	}
	shared.SetAuthCookie(w, token, expirationTime)
	http.Redirect(w, r, fmt.Sprintf("https://%s.linkwind.co", model.Name), http.StatusSeeOther)
}

/*InviteUserHandler handles user invite operations*/
func InviteUserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleInviteUserGET(w, r)
	case "POST":
		handleInviteUserPOST(w, r)
	default:
		handleInviteUserGET(w, r)
	}
}

func handleInviteUserGET(w http.ResponseWriter, r *http.Request) {
	err := templates.RenderInLayout(
		w,
		r,
		"invite.html",
		&models.InviteUserViewModel{},
	)
	if err != nil {
		panic(err)
	}
}

func handleInviteUserPOST(w http.ResponseWriter, r *http.Request) {
	model := &models.InviteUserViewModel{
		EmailAddress: r.FormValue("email"),
		Memo:         r.FormValue("memo"),
	}

	inviteHTMLPath := "invite.html"
	isValid, err := model.Validate(model.EmailAddress)
	if err != nil {
		panic(err)
	}
	if !isValid {
		err := templates.RenderInLayout(w, r, inviteHTMLPath, model)
		if err != nil {
			panic(err)
		}
		return
	}

	user := shared.GetUserFromContext(r)

	inviteCode, err := data.CreateInviteCode(user.ID, model.EmailAddress)
	if err != nil {
		panic(err)
	}

	domain, err := data.GetCustomerDomainByUserName(user.UserName)
	if err != nil {
		panic(err)
	}

	m := shared.InviteMailInfo{
		Domain:     domain,
		InviteCode: inviteCode,
		Email:      model.EmailAddress,
		UserName:   user.UserName,
		Memo:       model.Memo,
		Platform:   model.Layout.Platform,
	}
	err = shared.SendEmailInvitation(m)
	if err != nil {
		panic(err)
	}
	model.SuccessMessage = "Inivitation mail successfully sent to " + model.EmailAddress
	err = templates.RenderInLayout(w, r, inviteHTMLPath, model)
	if err != nil {
		panic(err)
	}
}

/*AdminHandler handles admin operations*/
func AdminHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleAdminGET(w, r)
	case "POST":
		handleAdminPOST(w, r)
	default:
		handleAdminGET(w, r)
	}
}

func handleAdminGET(w http.ResponseWriter, r *http.Request) {

	user := shared.GetUserFromContext(r)
	isAdmin, err := data.IsUserAdmin(user.ID)
	if err != nil {
		panic(err)
	}

	if !isAdmin {
		err = templates.RenderFile(w, "errors/500.html", nil)
		if err != nil {
			panic(err)
		}
		return
	}

	customer, err := data.GetCustomerByID(user.CustomerID)
	if err != nil {
		panic(err)
	}

	model, err := setCustomerAdminViewModel(customer, user)
	if err != nil {
		panic(err)
	}

	err = templates.RenderInLayout(w, r, "admin.html", model)
	if err != nil {
		panic(err)
	}
}

func handleAdminPOST(w http.ResponseWriter, r *http.Request) {

	user := shared.GetUserFromContext(r)

	customer, err := data.GetCustomerByID(user.CustomerID)
	if err != nil {
		panic(err)
	}

	adminHTMLPath := "admin.html"
	model := &models.CustomerAdminViewModel{
		Name:              r.FormValue("name"),
		Domain:            r.FormValue("domain"),
		LogoImageAsBase64: getImage(customer, r),
	}

	if model.Validate() == false {
		if customer.LogoImage != nil {
			imageasB64, err := shared.EncodeLogoImageToBase64(customer.LogoImage)
			if err != nil {
				panic(err)
			}
			model.LogoImageAsBase64 = imageasB64
			layout := &models.LayoutViewModel{
				Platform: model.Name,
				Logo:     model.LogoImageAsBase64,
			}
			model.Layout = layout
		}
		err := templates.RenderInLayout(w, r, adminHTMLPath, model)
		if err != nil {
			panic(err)
		}
		return
	}

	if customer.Name != model.Name {
		exists, err := data.ExistsCustomerByName(model.Name)
		if err != nil {
			panic(err)
		}
		if exists {
			model.Errors["Name"] = "This name is already taken"
			err := templates.RenderInLayout(w, r, adminHTMLPath, model)
			if err != nil {
				panic(err)
			}
			return
		}

	}

	if model.Domain != "" {
		if customer.Domain != model.Domain {
			exists, err := data.ExistsCustomerByDomain(model.Domain)
			if err != nil {
				panic(err)
			}
			if exists {
				model.Errors["Domain"] = "This domain is already taken"
				err := templates.RenderInLayout(w, r, adminHTMLPath, model)
				if err != nil {
					panic(err)
				}
				return
			}
		}
	}

	setUpdatedCustomerByModel(model, customer)
	err = data.UpdateCustomer(customer)
	if err != nil {
		panic(err)
	}

	//
	// We keep customer objects in cache and request context with name and domain keys. Since customer got updated we clear cache in order to get updated data
	//

	customerCtx := &cache.CustomerCtx{
		ID:       customer.ID,
		Platform: customer.Name,
		Logo:     model.LogoImageAsBase64,
	}

	ctx := context.WithValue(r.Context(), shared.CustomerContextKey, customerCtx)

	cache.SetCustomer(customer.Name, customerCtx)
	cache.SetCustomer(customer.Domain, customerCtx)

	model.SuccessMessage = "Account updated successfuly"
	err = templates.RenderInLayout(w, r.WithContext(ctx), adminHTMLPath, model)
	if err != nil {
		panic(err)
	}
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

func setCustomerAdminViewModel(customer *data.Customer, user *shared.SignedInUserClaims) (*models.CustomerAdminViewModel, error) {
	var model models.CustomerAdminViewModel
	if customer.Domain != "" {
		model.Domain = customer.Domain
	}

	if len(customer.LogoImage) != 0 {
		imageasB64, err := shared.EncodeLogoImageToBase64(customer.LogoImage)
		if err != nil {
			return nil, err
		}
		model.LogoImageAsBase64 = imageasB64
	}

	model.Name = customer.Name
	return &model, nil
}

func setUpdatedCustomerByModel(model *models.CustomerAdminViewModel, customer *data.Customer) {
	customer.Name = model.Name
	customer.Domain = model.Domain
	if model.LogoImageAsBase64 == "" {
		return
	}

	customer.LogoImage = []byte(model.LogoImageAsBase64)
}

func getImageFile(r *http.Request) (multipart.File, error) {
	file, _, err := r.FormFile("logo")
	if err != nil && err != http.ErrMissingFile {
		return file, err
	}
	return file, nil
}

func getImage(customer *data.Customer, r *http.Request) string {
	var logoImageAsBase64 string
	file, err := getImageFile(r)
	if err != nil {
		panic(err)
	}

	if file != nil {
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		logoImageAsBase64 = base64.StdEncoding.EncodeToString(fileBytes)
	} else {
		if len(customer.LogoImage) != 0 {
			imageasB64, err := shared.EncodeLogoImageToBase64(customer.LogoImage)
			if err != nil {
				panic(err)
			}
			logoImageAsBase64 = imageasB64
		}
	}
	return logoImageAsBase64
}
