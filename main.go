package main

import (
	"log"
	"net/http"
	"turkdev/app/controllers"
	"turkdev/app/src/templates"
	"turkdev/shared"
)

func main() {

	templates.Initialize()

	http.HandleFunc("/users/settings", errorHandler(controllers.UserSettingsHandler))
	http.HandleFunc("/signup", controllers.SignUpHandler)
	http.HandleFunc("/", notFoundHandler(controllers.StoriesHandler))
	http.HandleFunc("/recent", errorHandler(controllers.RecentStoriesHandler))
	http.HandleFunc("/comments", errorHandler(controllers.CommentsHandler))
	http.HandleFunc("/submit", errorHandler(controllers.SubmitStoryHandler))
	http.HandleFunc("/saved", errorHandler(controllers.SavedStoriesHandler))
	http.HandleFunc("/invite", errorHandler(controllers.InviteUserHandler))
	http.HandleFunc("/replies", errorHandler(controllers.RepliesHandler))
	http.HandleFunc("/settings", errorHandler(controllers.UserSettingsHandler))
	http.HandleFunc("/login", errorHandler(controllers.SignInHandler))

	staticFileServer := http.FileServer(http.Dir("app/dist/"))

	http.Handle("/dist/", http.StripPrefix("/dist/", staticFileServer))

	http.ListenAndServe(":80", nil)
}

func errorHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
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
		}

		if r.URL.Path == "/robots.txt" {
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte("User-agent: *\nDisallow: /"))
			return
		}

		byteValue, err := shared.ReadFile("app/src/templates/errors/404.html")
		if err != nil {
			log.Printf("Error occured in readFile func. Original err : %v", err)
		}
		w.WriteHeader(404)
		w.Write(byteValue)
	}
}
