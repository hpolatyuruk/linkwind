package controllers

import (
	"linkwind/app/data"
	"linkwind/app/models"
	"linkwind/app/shared"
	"linkwind/app/templates"
	"net/http"
	"strings"
)

/*UserProfileHandler handles showing user profile detail*/
func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleUserProfileGET(w, r)
	case "POST":
		handleUserProfilePOST(w, r)
	default:
		handleUserProfileGET(w, r)
	}
}

func handleUserProfileGET(w http.ResponseWriter, r *http.Request) {
	userCtx := shared.GetUserFromContext(r)
	model := &models.UserProfileViewModel{}
	renderFilePath := "readonly-profile.html"
	userName := r.URL.Query().Get("user")
	if strings.TrimSpace(userName) == "" {
		userName = userCtx.UserName
	}
	if userName == userCtx.UserName {
		renderFilePath = "profile-edit.html"
	}

	user, err := data.GetUserByUserName(userName)
	if err != nil {
		panic(err)
	}
	if user == nil {
		err := templates.RenderFile(w, "errors/404.html", nil)
		if err != nil {
			panic(err)
		}
	}

	isAdmin, err := data.IsUserAdmin(user.ID)
	if err != nil {
		panic(err)
	}
	setUserToModel(user, model, isAdmin)
	err = templates.RenderInLayout(w, r, renderFilePath, model)
	if err != nil {
		panic(err)
	}
}

func handleUserProfilePOST(w http.ResponseWriter, r *http.Request) error {
	model := &models.UserProfileViewModel{
		FullName: r.FormValue("fullName"),
		Email:    r.FormValue("email"),
		About:    r.FormValue("about"),
	}

	if model.Validate() == false {
		err := templates.RenderInLayout(w, r, "profile-edit.html", model)
		if err != nil {
			return err
		}
		return nil
	}

	userCtx := shared.GetUserFromContext(r)

	user, err := data.GetUserByUserName(userCtx.UserName)
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
			err := templates.RenderInLayout(w, r, "profile-edit.html", model)
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
	err = templates.RenderInLayout(w, r, "profile-edit.html", model)
	if err != nil {
		return err
	}
	return nil
}

func setUserToModel(user *data.User, model *models.UserProfileViewModel, isAdmin bool) {
	model.UserName = user.UserName
	model.FullName = user.FullName
	model.Karma = user.Karma
	model.RegisteredOn = shared.DateToString(user.RegisteredOn)
	model.About = user.About
	model.Email = user.Email
	model.IsAdmin = isAdmin
}
