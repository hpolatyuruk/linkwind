package main

import (
	"fmt"
	"log"
	"net/http"
	"turkdev/app/controllers"
	"turkdev/app/src/templates"
	"turkdev/shared"
)

func main() {

	err := templates.Initialize()
	if err != nil {
		fmt.Printf("An error occurred while initializing templates. Error: %v", err)
		panic(err)
	}

	http.HandleFunc("/users/profile", errorHandler(controllers.UserProfileHandler))
	http.HandleFunc("/signup", errorHandler(controllers.SignUpHandler))
	http.HandleFunc("/signin", errorHandler(controllers.SignInHandler))
	http.HandleFunc("/signout", errorHandler(controllers.SignOutHandler))
	http.HandleFunc("/", notFoundHandler(controllers.StoriesHandler))
	http.HandleFunc("/recent", errorHandler(controllers.RecentStoriesHandler))
	http.HandleFunc("/stories/detail", errorHandler(controllers.StoryDetailHandler))
	http.HandleFunc("/stories/upvote", controllers.UpvoteStoryHandler)
	http.HandleFunc("/stories/unvote", controllers.UnvoteStoryHandler)
	http.HandleFunc("/submit", errorHandler(controllers.SubmitStoryHandler))
	http.HandleFunc("/saved", errorHandler(controllers.SavedStoriesHandler))
	http.HandleFunc("/invite", errorHandler(controllers.InviteUserHandler))
	http.HandleFunc("/replies", errorHandler(controllers.RepliesHandler))
	http.HandleFunc("/login", errorHandler(controllers.SignInHandler))

	staticFileServer := http.FileServer(http.Dir("app/dist/"))

	http.Handle("/dist/", http.StripPrefix("/dist/", staticFileServer))

	http.ListenAndServe(":80", nil)
}

func errorHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			fmt.Printf("Error: %v", err)
			//logger.Error("Returned 500 internal server error! - " + r.Host + r.RequestURI + " - " + err.Error())
			byteValue, err := shared.ReadFile("app/src/templates/errors/500.html")
			if err != nil {
				log.Printf("Error occured in readFile func. Original err : %v", err)
			}
			w.WriteHeader(500)
			w.Write(byteValue)
		}
	}
}

func notFoundHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			err := f(w, r)
			if err != nil {
				byteValue, err := shared.ReadFile("app/src/templates/errors/500.html")
				if err != nil {
					log.Printf("Error occured in readFile func. Original err : %v", err)
				}
				w.WriteHeader(500)
				w.Write(byteValue)
			}
		} else if r.URL.Path == "/robots.txt" {
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte("User-agent: *\nDisallow: /"))
		} else {
			byteValue, err := shared.ReadFile("app/src/templates/errors/404.html")
			if err != nil {
				log.Printf("Error occured in readFile func. Original err : %v", err)
			}
			w.WriteHeader(404)
			w.Write(byteValue)
		}
	}
}
