package controllers

import (
	"net/http"
	"turkdev/app/models"
	"turkdev/app/templates"
)

func UserSettingsHandler(w http.ResponseWriter, r *http.Request) {
	title := "User Settings | Turk Dev"
	user := models.User{"Anil Yuzener"}
	data := map[string]interface{}{
		"Content": "Settings",
	}

	templates.Render(
		w,
		"user/settings.html",
		models.ViewModel{
			title,
			user,
			data,
		},
	)
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	title := "Sign-In | Turk Dev"
	user := models.User{"Anil Yuzener"}
	data := map[string]interface{}{
		"Content": "User sign-in",
	}

	templates.Render(
		w,
		"user/sign-in.html",
		models.ViewModel{
			title,
			user,
			data,
		},
	)
}

func InviteUserHandler(w http.ResponseWriter, r *http.Request) {
	title := "Invite a new user | Turk Dev"
	user := models.User{"Anil Yuzener"}
	data := map[string]interface{}{
		"Content": "Invite a new user",
	}

	templates.Render(
		w,
		"user/sign-up.html",
		models.ViewModel{
			title,
			user,
			data,
		},
	)
}
