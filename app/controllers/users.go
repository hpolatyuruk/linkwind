package controllers

import (
	"net/http"
	"turkdev/app/models"
	"turkdev/app/src/templates"
	"turkdev/data"
)

/*UserProfileHandler handles showing user profile detail*/
func UserProfileHandler(w http.ResponseWriter, r *http.Request) error {
	title := "User Settings | Turk Dev"
	userViewModel := models.User{"Anil Yuzener"}

	userName := r.URL.Query().Get("username")
	if len(userName) == 0 {
		// TODO(Anil): There is no user. Show appropriate message here
		return nil
	}
	user, err := data.GetUserByUserName(userName)
	if err != nil {
		// TODO(Anil): Show error page here
	}
	if user != nil {
		// TODO(Anil): User does not exist. Show appropriate message here
	}

	// TODO(Anil): Maybe map user struct to viewmodel here? up to you.

	data := map[string]interface{}{
		"Content": "Settings",
		"User":    user,
	}

	templates.RenderInLayout(
		w,
		"settings.html",
		models.ViewModel{
			title,
			userViewModel,
			data,
		},
	)
	return nil
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
