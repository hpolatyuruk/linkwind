package controllers

import (
	"net/http"
	"turkdev/app/models"
	"turkdev/app/templates"
)

func StoriesHandler(w http.ResponseWriter, r *http.Request) {
	title := "Turk Dev"
	user := models.User{"Anil Yuzener"}
	data := map[string]interface{}{
		"Content": "Stories",
	}

	templates.Render(
		w,
		"stories/index.html",
		models.ViewModel{
			title,
			user,
			data,
		},
	)
}

func RecentStoriesHandler(w http.ResponseWriter, r *http.Request) {
	title := "Recent Stories | Turk Dev"
	user := models.User{"Anil Yuzener"}
	data := map[string]interface{}{
		"Content": "Recent Stories",
	}

	templates.Render(
		w,
		"stories/index.html",
		models.ViewModel{
			title,
			user,
			data,
		},
	)
}

func SavedStoriesHandler(w http.ResponseWriter, r *http.Request) {
	title := "Saved Stories | Turk Dev"
	user := models.User{"Anil Yuzener"}
	data := map[string]interface{}{
		"Content": "Saved Stories",
	}

	templates.Render(
		w,
		"stories/index.html",
		models.ViewModel{
			title,
			user,
			data,
		},
	)
}

func SubmitStoryHandler(w http.ResponseWriter, r *http.Request) {
	title := "Submit Story | Turk Dev"
	user := models.User{"Anil Yuzener"}
	data := map[string]interface{}{
		"Content": "Submit Story",
	}

	templates.Render(
		w,
		"stories/submit.html",
		models.ViewModel{
			title,
			user,
			data,
		},
	)
}
