package main

import (
	"net/http"
	"turkdev/app/controllers"
	"turkdev/app/templates"
)

func main() {

	templates.Initialize()

	http.HandleFunc("/", controllers.StoriesHandler)
	http.HandleFunc("/recent", controllers.RecentStoriesHandler)
	http.HandleFunc("/comments", controllers.CommentsHandler)
	http.HandleFunc("/stories/new", controllers.SubmitStoryHandler)
	http.HandleFunc("/saved", controllers.SavedStoriesHandler)
	http.HandleFunc("/invite", controllers.InviteUserHandler)
	http.HandleFunc("/replies", controllers.RepliesHandler)
	http.HandleFunc("/settings", controllers.UserSettingsHandler)
	http.HandleFunc("/login", controllers.SignInHandler)

	staticFileServer := http.FileServer(http.Dir("app/static/"))

	http.Handle("/static/", http.StripPrefix("/static/", staticFileServer))

	http.ListenAndServe(":80", nil)
}
