package main

import (
	"fmt"
	"log"
	"net/http"
	"turkdev/app/controllers"
	"turkdev/shared"
)

func main() {

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
	http.HandleFunc("/invite", errorHandler(controllers.InviteUserHandler))
	http.HandleFunc("/replies", errorHandler(controllers.RepliesHandler))
	http.HandleFunc("/comments/add", errorHandler(controllers.AddCommentHandler))
	http.HandleFunc("/comments/upvote", controllers.UpvoteCommentHandler)
	http.HandleFunc("/comments/unvote", controllers.UnvoteCommentHandler)
	http.HandleFunc("/comments/reply", controllers.ReplyToCommentHandler)
	http.HandleFunc("/users/stories/saved", errorHandler(controllers.UserSavedStoriesHandler))
	http.HandleFunc("/users/stories/submitted", errorHandler(controllers.UserSubmittedStoriesHandler))
	http.HandleFunc("/users/stories/upvoted", errorHandler(controllers.UserUpvotedStoriesHandler))
	http.HandleFunc("/reset-password", errorHandler(controllers.ResetPasswordHandler))
	http.HandleFunc("/set-new-password", errorHandler(controllers.SetNewPasswordHandler))
	http.HandleFunc("/change-password", errorHandler(controllers.ChangePasswordHandler))
	http.HandleFunc("/profile-edit", errorHandler(controllers.UserProfileHandler))
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
