package controllers

import (
	"net/http"
	"strings"
	"turkdev/src/data"
	"turkdev/src/models"
	"turkdev/src/shared"
	"turkdev/src/templates"
)

/*UserProfileViewModel represents the data which is needed on sigin UI.*/
type UserProfileViewModel struct {
	About          string
	Email          string
	FullName       string
	Karma          int
	RegisteredOn   string
	UserName       string
	Errors         map[string]string
	SuccessMessage string
	SignedInUser   *models.SignedInUserViewModel
	IsAdmin        bool
}

/*Validate validates the UserProfileViewModel*/
func (model *UserProfileViewModel) Validate() bool {
	model.Errors = make(map[string]string)
	if strings.TrimSpace(model.Email) == "" {
		model.Errors["Email"] = "Email is required!"
	}
	if !shared.IsEmailAdrressValid(model.Email) {
		model.Errors["General"] = "E-mail address is not valid!"
	}
	return len(model.Errors) == 0
}

/*UserProfileHandler handles showing user profile detail*/
func UserProfileHandler(w http.ResponseWriter, r *http.Request) error {
	user := shared.GetUser(r)
	switch r.Method {
	case "GET":
		return handleUserProfileGET(w, r, user)
	case "POST":
		return handleUserProfilePOST(w, r, user)
	default:
		return handleUserProfileGET(w, r, user)
	}
}

func handleUserProfileGET(w http.ResponseWriter, r *http.Request, userClaims *shared.SignedInUserClaims) error {
	model := &UserProfileViewModel{
		SignedInUser: &models.SignedInUserViewModel{
			UserName:   userClaims.UserName,
			UserID:     userClaims.ID,
			CustomerID: userClaims.CustomerID,
			Email:      userClaims.Email,
		},
	}
	renderFilePath := "readonly-profile.html"
	userName := r.URL.Query().Get("user")
	if strings.TrimSpace(userName) == "" {
		userName = userClaims.UserName
	}
	if userName == userClaims.UserName {
		renderFilePath = "profile-edit.html"
	}

	user, err := data.GetUserByUserName(userName)
	if err != nil {
		return err
	}
	if user == nil {
		err := templates.RenderFile(w, "errors/404.html", nil)
		if err != nil {
			return err
		}
		return nil
	}

	isAdmin, err := data.IsUserAdmin(user.ID)
	if err != nil {
		return err
	}
	setUserToModel(user, model, isAdmin)
	err = templates.RenderInLayout(w, renderFilePath, model)
	if err != nil {
		return err
	}
	return nil
}

func handleUserProfilePOST(w http.ResponseWriter, r *http.Request, userClaims *shared.SignedInUserClaims) error {
	model := &UserProfileViewModel{
		FullName: r.FormValue("fullName"),
		Email:    r.FormValue("email"),
		About:    r.FormValue("about"),
		SignedInUser: &models.SignedInUserViewModel{
			UserName:   userClaims.UserName,
			UserID:     userClaims.ID,
			CustomerID: userClaims.CustomerID,
			Email:      userClaims.Email,
		},
	}

	if model.Validate() == false {
		err := templates.RenderInLayout(w, "profile-edit.html", model)
		if err != nil {
			return err
		}
		return nil
	}

	user, err := data.GetUserByUserName(userClaims.UserName)
	if err != nil {
		return err
	}
	isAdmin, err := data.IsUserAdmin(user.ID)
	if err != nil {
		return err
	}
	if user.Email != model.Email {
		exists, err := data.ExistsUserByEmail(user.Email)
		if err != nil {
			return err
		}
		if exists {
			model.Email = user.Email
			user.About = model.About
			setUserToModel(user, model, isAdmin)
			model.Errors["Email"] = "Entered e-mail address exists in db!"
			err := templates.RenderInLayout(w, "profile-edit.html", model)
			if err != nil {
				return err
			}
			return nil
		}
	}

	user.Email = model.Email
	user.About = model.About
	user.FullName = model.FullName
	err = data.UpdateUser(user)
	if err != nil {
		return err
	}

	setUserToModel(user, model, isAdmin)
	model.SuccessMessage = "User infos updated successfuly!"
	err = templates.RenderInLayout(w, "profile-edit.html", model)
	if err != nil {
		return err
	}
	return nil
}

func setUserToModel(user *data.User, model *UserProfileViewModel, isAdmin bool) {
	model.UserName = user.UserName
	model.FullName = user.FullName
	model.Karma = user.Karma
	model.RegisteredOn = shared.DateToString(user.RegisteredOn)
	model.About = user.About
	model.Email = user.Email
	model.IsAdmin = isAdmin
}
