package controllers

import (
	"net/http"
	"turkdev/app/models"
	"turkdev/app/templates"
)

func CommentsHandler(w http.ResponseWriter, r *http.Request) {
	title := "Comments | Turk Dev"
	user := models.User{"Anil Yuzener"}
	data := map[string]interface{}{
		"Content": "Comments",
	}

	templates.Render(
		w,
		"comments/index.html",
		models.ViewModel{
			title,
			user,
			data,
		},
	)
}

func RepliesHandler(w http.ResponseWriter, r *http.Request) {
	title := "Replies | Turk Dev"
	user := models.User{"Anil Yuzener"}
	data := map[string]interface{}{
		"Content": "Replies",
	}

	templates.Render(
		w,
		"comments/index.html",
		models.ViewModel{
			title,
			user,
			data,
		},
	)
}
