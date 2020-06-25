package controllers

import (
	"fmt"
	"linkwind/app/data"
	"linkwind/app/shared"
	"linkwind/app/templates"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const authExpirationMinutes = 1440
const userNameMaxCharCount = 15

/*SignInViewModel represents the data which is needed on sigin UI.*/
type SignInViewModel struct {
	EmailOrUserName string
	Password        string
	Errors          map[string]string
}

/*SignUpViewModel represents the data which is needed on sigup UI.*/
type SignUpViewModel struct {
	UserName   string
	Email      string
	Password   string
	InviteCode string
	Errors     map[string]string
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
	} else {
		if shared.IsPasswordValid(model.NewPassword) == false {
			model.Errors["NewPassword"] = "The password is not valid. A password should contan at least 1 uppercase, 1 lowercase, 1 digit, one of #$+=!*@& special characters and have a length of at least of 8."
		}
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
	} else {
		if shared.IsPasswordValid(model.NewPassword) == false {
			model.Errors["NewPassword"] = "The password is not valid. A password should contan at least 1 uppercase, 1 lowercase, 1 digit, one of #$+=!*@& special characters and have a length of at least of 8."
		}
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
func SignInHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		isAuthenticated, _, err := shared.IsAuthenticated(r)
		if err != nil {
			panic(err)
		}
		if isAuthenticated {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		handleSignInGET(w, r)
	case "POST":
		handleSignInPOST(w, r)
	default:
		handleSignInGET(w, r)
	}
}

func handleSignInGET(w http.ResponseWriter, r *http.Request) {
	err := templates.RenderFile(
		w,
		"layouts/users/signin.html",
		SignInViewModel{},
	)
	if err != nil {
		panic(err)
	}
}

func handleSignInPOST(w http.ResponseWriter, r *http.Request) {
	model := &SignInViewModel{
		EmailOrUserName: r.FormValue("emailOrUserName"),
		Password:        r.FormValue("password"),
	}
	if model.Validate() == false {
		err := templates.RenderFile(w, "/layouts/users/signin.html", model)
		if err != nil {
			panic(err)
		}
		return
	}
	var err error
	var user *data.User
	var fnExistsUser existsUser = data.ExistsUserByUserName
	var fnFindUser findUser = data.FindUserByUserNameAndPassword

	if shared.IsEmailAdressValid(model.EmailOrUserName) {
		fnExistsUser = data.ExistsUserByEmail
		fnFindUser = data.FindUserByEmailAndPassword
	}
	user, err = checkUser(fnExistsUser,
		fnFindUser,
		model.EmailOrUserName,
		model.Password)
	if err != nil {
		panic(err)
	}

	customerCtx := shared.GetCustomerFromContext(r)

	//
	// In case user exists but customer id from context (coming from subdomain) and user's customer id are different, it means that user wants to login to someone else's platform. We don't allow this to happen
	//

	if user == nil || customerCtx.ID != user.CustomerID {
		model.Errors["General"] = "User does not exist!"
		err = templates.RenderFile(w, "/layouts/users/signin.html", model)
		if err != nil {
			panic(err)
		}
		return
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(authExpirationMinutes * time.Minute)

	token, err := shared.GenerateAuthToken(*user, expirationTime)
	if err != nil {
		panic(err)
	}
	shared.SetAuthCookie(w, token, expirationTime)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type findUser func(userNameOrEmail, password string) (*data.User, error)
type existsUser func(userNameOrEmail string) (bool, error)

func checkUser(fnExistsUser existsUser, fnFindUser findUser, userNameOrEmail, password string) (*data.User, error) {
	exists, err := fnExistsUser(userNameOrEmail)
	if err != nil {
		return nil, err
	}
	if exists == false {
		return nil, nil
	}
	user, err := fnFindUser(userNameOrEmail, password)
	if err != nil {
		return nil, err
	}
	return user, err
}

/*SignUpHandler handles user signup operations*/
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

func handleSignUpGET(w http.ResponseWriter, r *http.Request) {
	// Only invited users can create an account
	inviteCode := r.URL.Query().Get("invitecode")
	if strings.TrimSpace(inviteCode) == "" {
		templates.RenderFile(w, "layouts/users/forbidden-signup.html", &SignUpViewModel{})
		return
	}
	invideCodeInfo, err := data.GetInviteCodeInfoByCode(inviteCode)
	if err != nil {
		panic(err)
	}
	if invideCodeInfo == nil {
		http.Error(w, "Invite code could not be found!", http.StatusBadRequest)
		panic(nil)
	}
	if invideCodeInfo.Used {
		http.Error(w, "The invite code is already used!", http.StatusBadRequest)
		return
	}
	templates.RenderFile(
		w,
		"layouts/users/signup.html",
		SignUpViewModel{
			InviteCode: invideCodeInfo.Code,
			Email:      invideCodeInfo.InvitedEmailAddress,
		},
	)
}

func handleSignUpPOST(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	signUpHTMLPath := "layouts/users/signup.html"
	model := &SignUpViewModel{
		UserName:   r.FormValue("userName"),
		Email:      r.FormValue("email"),
		Password:   r.FormValue("password"),
		InviteCode: r.FormValue("inviteCode"),
	}
	if model.Validate() == false {
		templates.RenderFile(w, signUpHTMLPath, model)
		return
	}
	if strings.TrimSpace(model.InviteCode) == "" {
		model.Errors["General"] = "Missing invite code!"
		templates.RenderFile(w, signUpHTMLPath, model)
		return
	}
	invitedCodeInfo, err := data.GetInviteCodeInfoByCode(model.InviteCode)
	if err != nil {
		panic(err)
	}
	if invitedCodeInfo == nil {
		model.Errors["General"] = "Invite code could not be found. Please make sure that you have a valid invite code."
		templates.RenderFile(w, signUpHTMLPath, model)
		return
	}
	if invitedCodeInfo.Used {
		model.Errors["General"] = "The invite code is already used!"
		return
	}
	if invitedCodeInfo.InvitedEmailAddress != model.Email {
		model.Errors["General"] = "The email address you entered does not match the invited email address."
		templates.RenderFile(w, signUpHTMLPath, model)
		return
	}
	exists, err := data.ExistsUserByUserName(model.UserName)
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
	inviterUser, err := data.GetUserByID(invitedCodeInfo.InviterUserID)
	if err != nil {
		panic(err)
	}
	if inviterUser == nil {
		model.Errors["General"] = "The inviter user could not be found!"
		templates.RenderFile(w, signUpHTMLPath, model)
		return
	}
	var user data.User
	user.UserName = model.UserName
	user.Email = model.Email
	user.Password = model.Password
	user.Karma = 0
	user.RegisteredOn = time.Now()
	user.CustomerID = inviterUser.CustomerID
	user.InviteCode = model.InviteCode
	userID, err := data.CreateUser(&user)
	if err != nil {
		panic(err)
	}
	user.ID = *userID
	err = data.MarkInviteCodeAsUsed(model.InviteCode)
	if err != nil {
		panic(err)
	}
	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(authExpirationMinutes * time.Minute)

	token, err := shared.GenerateAuthToken(user, expirationTime)
	if err != nil {
		panic(err)
	}
	shared.SetAuthCookie(w, token, expirationTime)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

/*SignOutHandler handles user singout operations.*/
func SignOutHandler(w http.ResponseWriter, r *http.Request) {
	shared.SetAuthCookie(w, "", time.Now())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

/*ResetPasswordHandler handles user  reset password operations*/
func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleResetPasswordGET(w, r)
	case "POST":
		handleResetPasswordPOST(w, r)
	default:
		handleResetPasswordGET(w, r)
	}
}

func handleResetPasswordGET(w http.ResponseWriter, r *http.Request) {
	err := templates.RenderFile(
		w,
		"layouts/users/reset-password.html",
		ResetPasswordViewModel{},
	)
	if err != nil {
		panic(err)
	}
}

func handleResetPasswordPOST(w http.ResponseWriter, r *http.Request) {
	model := &ResetPasswordViewModel{
		EmailOrUserName: r.FormValue("emailOrUserName"),
	}

	if model.Validate() == false {
		err := templates.RenderFile(w, "layouts/users/reset-password.html", model)
		if err != nil {
			panic(err)
		}
		return
	}

	var err error
	var email string
	var userName string
	var user *data.User
	model.SuccessMessage = "Password recovery message sent. If you don't see it, you might want to check your spam folder."
	if shared.IsEmailAdressValid(model.EmailOrUserName) {
		email = model.EmailOrUserName
		userName, err = data.GetUserNameByEmail(model.EmailOrUserName)
	} else {
		userName = model.EmailOrUserName
	}

	user, err = data.GetUserByUserName(userName)
	if user != nil {
		email = user.Email
	}

	if err != nil {
		// If submitted mail does not exist in db, error return
		// "no rows in result set". If error return this message,
		// we response success message because we do not want to reveal db
		// records. Otherwise, return panic(500)
		if strings.Contains(err.Error(), "no rows in result set") {
			err = templates.RenderFile(w, "layouts/users/reset-password.html", model)
			if err != nil {
				panic(err)
			}
			return
		}
		panic(err)
	}

	domain, err := data.GetCustomerDomainByUserName(userName)
	if err != nil {
		panic(err)
	}

	token := shared.GenerateResetPasswordToken()

	customerCtx := shared.GetCustomerFromContext(r)

	query := shared.ResetPasswordMailInfo{
		Email:    email,
		UserName: userName,
		Domain:   domain,
		Platform: customerCtx.Platform,
		Token:    token,
	}
	err = shared.SendResetPasswordMail(query)
	if err != nil {
		panic(err)
	}

	data.SaveResetPasswordToken(token, user.ID)
	err = templates.RenderFile(w, "layouts/users/reset-password.html", model)
	if err != nil {
		panic(err)
	}
}

/*SetNewPasswordHandler handles set new password operations*/
func SetNewPasswordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleSetNewPasswordGET(w, r)
	case "POST":
		handleSetNewPasswordPOST(w, r)
	default:
		handleSetNewPasswordGET(w, r)
	}
}

func handleSetNewPasswordGET(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if strings.TrimSpace(token) == "" {
		http.Error(w, "Missing Token! ", http.StatusBadRequest)
		return
	}

	user, err := data.GetUserByResetPasswordToken(token)
	if err != nil {
		panic(err)
	}
	if user == nil {
		http.Error(w, "Token is not valid! ", http.StatusBadRequest)
		return
	}
	err = templates.RenderFile(
		w,
		"layouts/users/set-new-password.html",
		&SetNewPasswordViewModel{
			UserName: user.UserName,
		},
	)
	if err != nil {
		panic(err)
	}
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

/*ChangePasswordHandler handles change password operations*/
func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	_, claims, err := shared.IsAuthenticated(r)
	if err != nil {
		panic(err)
	}
	switch r.Method {
	case "GET":
		handleChangePasswordGET(w, r)
	case "POST":
		handleChangePasswordPOST(w, r, claims.ID)
	default:
		handleChangePasswordGET(w, r)
	}
}

func handleChangePasswordGET(w http.ResponseWriter, r *http.Request) {
	err := templates.RenderFile(
		w,
		"layouts/users/change-password.html",
		ChangePasswordViewModel{},
	)
	if err != nil {
		panic(err)
	}
}

func handleChangePasswordPOST(w http.ResponseWriter, r *http.Request, userID int) {
	model := &ChangePasswordViewModel{
		CurrentPassword: r.FormValue("currentPassword"),
		NewPassword:     r.FormValue("newPassword"),
		ConfirmPassword: r.FormValue("confirmPassword"),
	}

	if model.Validate() == false {
		err := templates.RenderFile(w, "layouts/users/change-password.html", model)
		if err != nil {
			panic(err)
		}
		return
	}

	matched, err := data.ConfirmPasswordMatch(userID, model.CurrentPassword)
	if !matched {
		model.Errors["General"] = "User does not exist!"
		err = templates.RenderFile(w, "layouts/users/change-password.html", model)
		if err != nil {
			panic(err)
		}
		return
	}
	err = data.ChangePassword(userID, model.NewPassword)
	if err != nil {
		panic(err)
	}

	model.SuccessMessage = "Password successfuly changed"
	err = templates.RenderFile(w, "layouts/users/change-password.html", model)
	if err != nil {
		panic(err)
	}
}

/*GenerateInviteCodeHandler generate the invite code to invite an user to join the system*/
func GenerateInviteCodeHandler(w http.ResponseWriter, r *http.Request) {
	inviterUserID, _ := strconv.Atoi(r.URL.Query().Get("userid"))
	invitedEmail := r.URL.Query().Get("invitedemail")
	inviteCode, err := data.CreateInviteCode(inviterUserID, invitedEmail)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte(inviteCode))
}

// SetAuthTokenHandler sets the auth cookie and redirect user to main page
func SetAuthTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only http get allowed", http.StatusMethodNotAllowed)
		return
	}
	platformName := r.URL.Query().Get("customer")
	if strings.TrimSpace(platformName) == "" {
		http.Error(w, "Missing customer param.", http.StatusBadRequest)
		return
	}
	customer, err := data.GetCustomerByName(platformName)
	if err != nil {
		panic(err)
	}
	if customer == nil {
		http.Error(w, "Customer does not exist.", http.StatusNotFound)
		return
	}
	authToken := r.URL.Query().Get("auth")
	if strings.TrimSpace(authToken) != "" {
		// Declare the expiration time of the token
		// here, we have kept it as 5 minutes
		expirationTime := time.Now().Add(authExpirationMinutes * time.Minute)
		shared.SetAuthCookie(w, authToken, expirationTime)
	}
	http.Redirect(
		w,
		r,
		fmt.Sprintf("https://%s.linkwind.co", platformName),
		http.StatusSeeOther)
}
