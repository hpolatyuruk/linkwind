package main

import (
	"net/http"
	"turkdev/app/controllers"
	"turkdev/app/src/templates"
)

func main() {

	templates.Initialize()

	http.HandleFunc("/", controllers.StoriesHandler)
	http.HandleFunc("/recent", controllers.RecentStoriesHandler)
	http.HandleFunc("/comments", controllers.CommentsHandler)
	http.HandleFunc("/submit", controllers.SubmitStoryHandler)
	http.HandleFunc("/saved", controllers.SavedStoriesHandler)
	http.HandleFunc("/invite", controllers.InviteUserHandler)
	http.HandleFunc("/replies", controllers.RepliesHandler)
	http.HandleFunc("/users/settings", controllers.UserSettingsHandler)
	http.HandleFunc("/signup", controllers.SignUpHandler)
	http.HandleFunc("/login", controllers.SignInHandler)

	staticFileServer := http.FileServer(http.Dir("app/dist/"))

	http.Handle("/dist/", http.StripPrefix("/dist/", staticFileServer))

	http.ListenAndServe(":80", nil)
}
