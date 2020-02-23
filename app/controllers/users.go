package controllers

import (
	"net/http"
	"strings"
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
	RegisteredOn   string
	UserName       string
	Errors         map[string]string
	SuccessMessage string
	SignedInUser   *models.SignedInUserViewModel
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
	isAuthenticated, user, err := shared.IsAuthenticated(r)
	if err != nil {
		return nil
	}

	if !isAuthenticated {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return nil
	}

	switch r.Method {
	case "GET":
		return handleUserProfileGET(w, r, user)
	case "POST":
		return handleUserProfilePOST(w, r, user)
	default:
		return handleUserProfileGET(w, r, user)
	}
}

/*InviteUserHandler handles sending invitations to user*/
func InviteUserHandler(w http.ResponseWriter, r *http.Request) error {

	templates.RenderInLayout(
		w,
		"signup.html",
		nil,
	)
	return nil
}

func handleUserProfileGET(w http.ResponseWriter, r *http.Request, userClaims *shared.SignedInUserClaims) error {
	//TODO : Burada model oluşturulup aşağıda set edilmeler değiştirlmeli. Çok fazla logic tekrarı yapılmış
	model := &UserProfileViewModel{
		SignedInUser: &models.SignedInUserViewModel{
			IsSigned: true,
			UserName: userClaims.UserName,
		},
	}

	userName := r.URL.Query().Get("user")
	if len(userName) == 0 {
		err := templates.RenderFile(w, "errors/404.html", nil)
		if err != nil {
			return err
		}
		return nil
	}

	if userName != userClaims.UserName {
		user, err := data.GetUserByUserName(userName)
		if err != nil {
			return err
		}
		setUserToModel(user, model)

		err = templates.RenderInLayout(w, "readonly-profile.html", model)
		if err != nil {
			return err
		}
		return nil
	}
	if userName == "" {
		user, err := data.GetUserByID(userClaims.ID)
		if err != nil {
			return err
		}

		setUserToModel(user, model)

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

	setUserToModel(user, model)

	err = templates.RenderInLayout(w, "profile-edit.html", model)
	if err != nil {
		return err
	}
	return nil
}

func handleUserProfilePOST(w http.ResponseWriter, r *http.Request, userClaims *shared.SignedInUserClaims) error {
	model := &UserProfileViewModel{
		Email: r.FormValue("email"),
		About: r.FormValue("about"),
		SignedInUser: &models.SignedInUserViewModel{
			IsSigned: true,
			UserName: userClaims.UserName,
		},
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

	user, err := data.GetUserByUserName(userClaims.UserName)
	if err != nil {
		return err
	}

	if user.Email == model.Email {
		exists, err := data.ExistsUserByEmail(user.Email)
		if err != nil {
			return err
		}
		if exists {
			model.Email = user.Email
			user.About = model.About
			setUserToModel(user, model)
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

	setUserToModel(user, model)
	model.SuccessMessage = "User infos updated successfuly!"
	err = templates.RenderInLayout(w, "profile-edit.html", model)
	if err != nil {
		return err
	}
	return nil
}

func setUserToModel(user *data.User, model *UserProfileViewModel) {
	model.UserName = user.UserName
	model.FullName = user.FullName
	model.Karma = user.Karma
	model.RegisteredOn = shared.DateToString(user.RegisteredOn)
	model.About = user.About
	model.Email = user.Email
}
