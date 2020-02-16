package controllers

import (
	"net/http"
	"strings"
	"time"
	"turkdev/app/src/templates"
	"turkdev/data"
	"turkdev/services"
	"turkdev/shared"
)

/*SignInViewModel represents the data which is needed on sigin UI.*/
type SignInViewModel struct {
	EmailOrUserName string
	Password        string
	Errors          map[string]string
}

/*SignUpViewModel represents the data which is needed on sigup UI.*/
type SignUpViewModel struct {
	UserName string
	Email    string
	Password string
	Errors   map[string]string
}

/*ResetPasswordViewModel represents the data which is needed on reset password UI*/
type ResetPasswordViewModel struct {
	EmailOrUserName string
	Errors          map[string]string
	SuccessMessage  string
}

/*SetNewPasswordViewModel represents the data which is needed on set new password UI*/
type SetNewPasswordViewModel struct {
	UserName        string
	NewPassword     string
	ConfirmPassword string
	Errors          map[string]string
	SuccessMessage  string
}

/*ChangePasswordViewModel represents the data which is needed on change password UI*/
type ChangePasswordViewModel struct {
	CurrentPassword string
	NewPassword     string
	ConfirmPassword string
	Errors          map[string]string
	SuccessMessage  string
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

/*Validate validates the SignUpViewModel*/
func (model *SignUpViewModel) Validate() bool {
	model.Errors = make(map[string]string)

	if strings.TrimSpace(model.UserName) == "" {
		model.Errors["UserName"] = "User name is required!"
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
	}
	return len(model.Errors) == 0
}

/*Validate validates the ResetPasswordViewModel*/
func (model *ResetPasswordViewModel) Validate() bool {
	model.Errors = make(map[string]string)

	if strings.TrimSpace(model.EmailOrUserName) == "" {
		model.Errors["EmailOrUserName"] = "Email or user name is required!"
	}

	return len(model.Errors) == 0
}

/*Validate validates the SetNewPasswordViewModel*/
func (model *SetNewPasswordViewModel) Validate() bool {
	model.Errors = make(map[string]string)

	if strings.TrimSpace(model.NewPassword) == "" {
		model.Errors["NewPassword"] = "New password is required!"
	}
	if strings.TrimSpace(model.ConfirmPassword) == "" {
		model.Errors["ConfirmPassword"] = "Confirm password is required!"
	}
	if model.ConfirmPassword != model.NewPassword {
		model.Errors["NotEqual"] = "New password and confirm password is not equal!"
	}
	return len(model.Errors) == 0
}

/*Validate validates the ChangePasswordViewModel*/
func (model *ChangePasswordViewModel) Validate() bool {
	model.Errors = make(map[string]string)

	if strings.TrimSpace(model.CurrentPassword) == "" {
		model.Errors["CurrentPassword"] = "Current password is required!"
	}
	if strings.TrimSpace(model.NewPassword) == "" {
		model.Errors["NewPassword"] = "New password is required!"
	}
	if strings.TrimSpace(model.ConfirmPassword) == "" {
		model.Errors["ConfirmPassword"] = "Confirm password is required!"
	}
	if model.ConfirmPassword != model.NewPassword {
		model.Errors["NotEqual"] = "New password and confirm password is not equal!"
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

/*ResetPasswordHandler handles user  reset password operations*/
func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		isAuthenticated, _, err := shared.IsAuthenticated(r)
		if err != nil {
			return err
		}
		if !isAuthenticated {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return nil
		}
		return handleResetPasswordGET(w, r)
	case "POST":
		return handleResetPasswordPOST(w, r)
	default:
		return handleResetPasswordGET(w, r)
	}
}

/*SetNewPasswordHandler handles set new password operations*/
func SetNewPasswordHandler(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		isAuthenticated, _, err := shared.IsAuthenticated(r)
		if err != nil {
			return err
		}
		if !isAuthenticated {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return nil
		}
		return handleSetNewPasswordGET(w, r)
	case "POST":
		return handleSetNewPasswordPOST(w, r)
	default:
		return handleSetNewPasswordGET(w, r)
	}
}

/*ChangePasswordHandler handles change password operations*/
func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) error {
	isAuthenticated, claims, err := shared.IsAuthenticated(r)
	if err != nil {
		return err
	}
	switch r.Method {
	case "GET":
		if !isAuthenticated {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return nil
		}
		return handleChangePasswordGET(w, r)
	case "POST":
		return handleChangePasswordPOST(w, r, claims.ID)
	default:
		return handleChangePasswordGET(w, r)
	}
}

func handleSignUpGET(w http.ResponseWriter, r *http.Request) error {
	templates.RenderFile(
		w,
		"layouts/users/signup.html",
		SignUpViewModel{},
	)
	return nil
}

func handleSignUpPOST(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	model := &SignUpViewModel{
		UserName: r.FormValue("userName"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	if model.Validate() == false {
		templates.RenderFile(w, "layouts/users/signup.html", model)
		return nil
	}
	var user data.User
	user.UserName = model.UserName
	user.Email = model.Email
	user.Password = model.Password
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
		err := templates.RenderFile(w, "/layouts/users/signin.html", model)
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
		err = templates.RenderFile(w, "/layouts/users/signin.html", model)
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

func handleResetPasswordGET(w http.ResponseWriter, r *http.Request) error {
	err := templates.RenderFile(
		w,
		"layouts/users/reset-password.html",
		ResetPasswordViewModel{},
	)
	if err != nil {
		return err
	}
	return nil
}

func handleResetPasswordPOST(w http.ResponseWriter, r *http.Request) error {
	model := &ResetPasswordViewModel{
		EmailOrUserName: r.FormValue("emailOrUserName"),
	}

	if model.Validate() == false {
		err := templates.RenderFile(w, "layouts/users/reset-password.html", model)
		if err != nil {
			return err
		}
		return nil
	}

	var exists bool
	var err error
	var email string
	var userName string
	var user *data.User
	if shared.IsEmailAdrressValid(model.EmailOrUserName) {
		exists, err = data.ExistsUserByEmail(model.EmailOrUserName)
		email = model.EmailOrUserName
		userName, err = data.GetUserNameByEmail(model.EmailOrUserName)
	} else {
		exists, err = data.ExistsUserByUserName(model.EmailOrUserName)
		if exists {
			user, err = data.GetUserByUserName(model.EmailOrUserName)
			userName = user.UserName
			email = user.Email
		}
	}
	if err != nil {
		return err
	}

	if !exists {
		model.Errors["General"] = "User does not exist!"
		err = templates.RenderFile(w, "layouts/users/reset-password.html", model)
		if err != nil {
			return err
		}
		return nil
	}
	domain, err := data.GetCustomerDomainByUserName(userName)
	if err != nil {
		return err
	}

	err = services.SendResetPasswordMail(email, userName, domain)
	if err != nil {
		return err
	}

	model.SuccessMessage = "Password recovery message sent. If you don't see it, you might want to check your spam folder."
	err = templates.RenderFile(w, "layouts/users/reset-password.html", model)
	if err != nil {
		return err
	}
	return nil
}

func handleSetNewPasswordGET(w http.ResponseWriter, r *http.Request) error {
	err := templates.RenderFile(
		w,
		"layouts/users/set-new-password.html",
		SetNewPasswordViewModel{},
	)
	if err != nil {
		return err
	}
	return nil
}

func handleSetNewPasswordPOST(w http.ResponseWriter, r *http.Request) error {
	model := &SetNewPasswordViewModel{
		UserName:        r.FormValue("userName"),
		NewPassword:     r.FormValue("newPassword"),
		ConfirmPassword: r.FormValue("confirmPassword"),
	}

	if model.Validate() == false {
		err := templates.RenderFile(w, "layouts/users/set-new-password.html", model)
		if err != nil {
			return err
		}
		return nil
	}

	exists, err := data.ExistsUserByUserName(model.UserName)
	if !exists {
		model.Errors["General"] = "User does not exist!"
		err = templates.RenderFile(w, "layouts/users/set-new-password.html", model)
		if err != nil {
			return err
		}
		return nil
	}
	user, err := data.GetUserByUserName(model.UserName)
	if err != nil {
		return err
	}

	err = data.ChangePassword(user.ID, model.NewPassword)
	if err != nil {
		return err
	}

	model.SuccessMessage = "Password successfuly changed"
	err = templates.RenderFile(w, "layouts/users/set-new-password.html", model)
	if err != nil {
		return err
	}
	return nil
}

func handleChangePasswordGET(w http.ResponseWriter, r *http.Request) error {
	err := templates.RenderFile(
		w,
		"layouts/users/change-password.html",
		ChangePasswordViewModel{},
	)
	if err != nil {
		return err
	}
	return nil
}

func handleChangePasswordPOST(w http.ResponseWriter, r *http.Request, userID int) error {
	model := &ChangePasswordViewModel{
		CurrentPassword: r.FormValue("currentPassword"),
		NewPassword:     r.FormValue("newPassword"),
		ConfirmPassword: r.FormValue("confirmPassword"),
	}

	if model.Validate() == false {
		err := templates.RenderFile(w, "layouts/users/change-password.html", model)
		if err != nil {
			return err
		}
		return nil
	}

	matched, err := data.ConfirmPasswordMatch(userID, model.CurrentPassword)
	if !matched {
		model.Errors["General"] = "User does not exist!"
		err = templates.RenderFile(w, "layouts/users/change-password.html", model)
		if err != nil {
			return err
		}
		return nil
	}
	err = data.ChangePassword(userID, model.NewPassword)
	if err != nil {
		return err
	}

	model.SuccessMessage = "Password successfuly changed"
	err = templates.RenderFile(w, "layouts/users/change-password.html", model)
	if err != nil {
		return err
	}
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
