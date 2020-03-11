package main

import (
	"net/http"
	"turkdev/app/controllers"
	"turkdev/middlewares"
)

func main() {
	// TODO: UserProfileHandler hem /users/porfile'ı hem de /profile-edit 'i karşılıyor.
	http.HandleFunc("/users/profile", middlewares.Error(controllers.UserProfileHandler))
	http.HandleFunc("/signup", middlewares.Error(controllers.SignUpHandler))
	http.HandleFunc("/signin", middlewares.Error(controllers.SignInHandler))
	http.HandleFunc("/signout", middlewares.Error(controllers.SignOutHandler))
	http.HandleFunc("/", middlewares.NotFound(controllers.StoriesHandler))
	http.HandleFunc("/recent", middlewares.Error(controllers.RecentStoriesHandler))
	http.HandleFunc("/stories/detail", middlewares.Error(controllers.StoryDetailHandler))
	http.HandleFunc("/stories/upvote", controllers.UpvoteStoryHandler)
	http.HandleFunc("/stories/unvote", controllers.UnvoteStoryHandler)
	http.HandleFunc("/submit", middlewares.Error(controllers.SubmitStoryHandler))
	http.HandleFunc("/comments/add", middlewares.Error(controllers.AddCommentHandler))
	http.HandleFunc("/comments/upvote", controllers.UpvoteCommentHandler)
	http.HandleFunc("/comments/unvote", controllers.UnvoteCommentHandler)
	http.HandleFunc("/comments/reply", controllers.ReplyToCommentHandler)
	http.HandleFunc("/users/stories/saved", middlewares.Error(controllers.UserSavedStoriesHandler))
	http.HandleFunc("/users/stories/submitted", middlewares.Error(controllers.UserSubmittedStoriesHandler))
	http.HandleFunc("/users/stories/upvoted", middlewares.Error(controllers.UserUpvotedStoriesHandler))
	http.HandleFunc("/reset-password", middlewares.Error(controllers.ResetPasswordHandler))
	http.HandleFunc("/set-new-password", middlewares.Error(controllers.SetNewPasswordHandler))
	http.HandleFunc("/change-password", middlewares.Error(controllers.ChangePasswordHandler))
	http.HandleFunc("/profile-edit", middlewares.Error(controllers.UserProfileHandler))
	http.HandleFunc("/invitecodes/generate", middlewares.Error(controllers.GenerateInviteCodeHandler))
	http.HandleFunc("/customer-signup", middlewares.Error(controllers.CustomerSignUpHandler))
	http.HandleFunc("/users/invite", middlewares.Error(controllers.InviteUserHandler))
	http.HandleFunc("/admin", middlewares.Error(controllers.AdminHandler))

	staticFileServer := http.FileServer(http.Dir("app/dist/"))

	http.Handle("/dist/", http.StripPrefix("/dist/", staticFileServer))

	http.ListenAndServe(":80", nil)
}
