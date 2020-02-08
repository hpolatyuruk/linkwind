package controllers

import (
	"net/http"
	"turkdev/app/models"
	"turkdev/app/src/templates"
	"turkdev/data"

	"github.com/gorilla/sessions"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

/*UserSettingsHandler handles showing user profile detail*/
func UserSettingsHandler(w http.ResponseWriter, r *http.Request) error {
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

	templates.RenderWithBase(
		w,
		"user/settings.html",
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

	templates.RenderWithBase(
		w,
		"user/sign-up.html",
		models.ViewModel{
			title,
			user,
			data,
		},
	)
	return nil
}
