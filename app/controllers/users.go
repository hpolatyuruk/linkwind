package controllers

import (
	"net/http"
	"strings"
	"time"
	"turkdev/app/models"
	"turkdev/app/src/templates"
	"turkdev/data"
	"turkdev/shared"
)

type UserProfileViewModel struct {
	About          string
	Email          string
	FullName       string
	Karma          int
	RegisteredOn   time.Time
	UserName       string
	Errors         map[string]string
	SuccessMessage string
}

/*Validate validates the UserProfileViewModel*/
func (model *UserProfileViewModel) Validate() bool {
	model.Errors = make(map[string]string)

	if strings.TrimSpace(model.Email) == "" {
		model.Errors["Email"] = "Email is required!"
	}
	return len(model.Errors) == 0
}

/*UserProfileHandler handles showing user profile detail*/
func UserProfileHandler(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return handleUserProfileGET(w, r)
	case "POST":
		return handleUserProfilePOST(w, r)
	default:
		return handleUserProfileGET(w, r)
	}
}

/*InviteUserHandler handles sending invitations to user*/
func InviteUserHandler(w http.ResponseWriter, r *http.Request) error {
	title := "Invite a new user | Turk Dev"
	user := models.User{"Anil Yuzener"}
	data := map[string]interface{}{
		"Content": "Invite a new user",
	}

	templates.RenderInLayout(
		w,
		"signup.html",
		models.ViewModel{
			title,
			user,
			data,
		},
	)
	return nil
}

func handleUserProfileGET(w http.ResponseWriter, r *http.Request) error {
	userName := r.URL.Query().Get("user")
	if len(userName) == 0 {
		err := templates.RenderFile(w, "errors/404.html", nil)
		if err != nil {
			return err
		}
		return nil
	}

	if userName == "" {
		isAuthenticated, claims, err := shared.IsAuthenticated(r)
		if err != nil {
			return err
		}

		if !isAuthenticated {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return nil
		}

		user, err := data.GetUserByID(claims.ID)
		if err != nil {
			return err
		}

		model := &UserProfileViewModel{
			About:        user.About,
			Email:        user.Email,
			FullName:     user.FullName,
			Karma:        user.Karma,
			RegisteredOn: user.RegisteredOn,
			UserName:     user.UserName,
		}

		err = templates.RenderInLayout(w, "layouts/users/profile-edit.html", model)
		if err != nil {
			return err
		}
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

	model := &UserProfileViewModel{
		About:        user.About,
		Email:        user.Email,
		FullName:     user.FullName,
		Karma:        user.Karma,
		RegisteredOn: user.RegisteredOn,
		UserName:     user.UserName,
	}

	err = templates.RenderInLayout(w, "profile-edit.html", model)
	if err != nil {
		return err
	}
	return nil
}

func handleUserProfilePOST(w http.ResponseWriter, r *http.Request) error {
	model := &UserProfileViewModel{
		Email: r.FormValue("email"),
		About: r.FormValue("about"),
	}

	if model.Validate() == false {
		err := templates.RenderInLayout(w, "profile-edit.html", model)
		if err != nil {
			return err
		}
		return nil
	}

	if !shared.IsEmailAdrressValid(model.Email) {
		model.Errors["General"] = "E-mail address is not valid!"
		err := templates.RenderInLayout(w, "profile-edit.html", model)
		if err != nil {
			return err
		}
		return nil
	}

	user, err := data.GetUserByUserName(r.FormValue("userName"))
	if err != nil {
		return err
	}

	if user.Email != model.Email {
		exists, err := data.ExistsUserByEmail(user.Email)
		if err != nil {
			return err
		}
		if exists {
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
	err = data.UpdateUser(user)
	if err != nil {
		return err
	}

	model.SuccessMessage = "User infos updated successfuly!"
	err = templates.RenderInLayout(w, "profile-edit.html", model)
	if err != nil {
		return err
	}
	return nil
}
