package controllers

import (
	"net/http"
	"turkdev/app/src/templates"
	"turkdev/data"
)

/*UserProfileHandler handles showing user profile detail*/
func UserProfileHandler(w http.ResponseWriter, r *http.Request) error {

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

	templates.RenderInLayout(
		w,
		"settings.html",
		nil,
	)
	return nil
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
