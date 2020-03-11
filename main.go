package main

import (
	"net/http"
	"turkdev/app/controllers"
	"turkdev/middlewares"
)

func main() {

	registerHandlers()

	staticFileServer := http.FileServer(http.Dir("app/dist/"))
	http.Handle("/dist/", http.StripPrefix("/dist/", staticFileServer))
	http.ListenAndServe(":80", nil)
}

func registerHandlers() {
	http.HandleFunc(
		"/users/profile",
		middlewares.Middleware(
			middlewares.Error(controllers.UserProfileHandler),
			middlewares.Auth))

	http.HandleFunc(
		"/signup",
		middlewares.Error(controllers.SignUpHandler))

	http.HandleFunc(
		"/signin",
		middlewares.Error(controllers.SignInHandler))

	http.HandleFunc(
		"/signout",
		middlewares.Error(controllers.SignOutHandler))

	http.HandleFunc(
		"/",
		middlewares.NotFound(controllers.StoriesHandler))

	http.HandleFunc(
		"/recent",
		middlewares.Error(controllers.RecentStoriesHandler))

	http.HandleFunc(
		"/stories/detail",
		middlewares.Error(controllers.StoryDetailHandler))

	http.HandleFunc(
		"/stories/upvote",
		middlewares.Middleware(
			controllers.UpvoteStoryHandler,
			middlewares.Auth))

	http.HandleFunc(
		"/stories/unvote",
		middlewares.Middleware(
			controllers.UnvoteStoryHandler,
			middlewares.Auth))

	http.HandleFunc(
		"/submit",
		middlewares.Middleware(
			middlewares.Error(controllers.SubmitStoryHandler),
			middlewares.Auth))

	http.HandleFunc(
		"/comments/add",
		middlewares.Middleware(
			middlewares.Error(controllers.AddCommentHandler),
			middlewares.Auth))

	http.HandleFunc(
		"/comments/upvote",
		middlewares.Middleware(
			controllers.UpvoteCommentHandler,
			middlewares.Auth))

	http.HandleFunc(
		"/comments/unvote",
		middlewares.Middleware(
			controllers.UnvoteCommentHandler,
			middlewares.Auth))

	http.HandleFunc(
		"/comments/reply",
		middlewares.Middleware(
			controllers.ReplyToCommentHandler,
			middlewares.Auth))

	http.HandleFunc(
		"/users/stories/saved",
		middlewares.Middleware(
			middlewares.Error(controllers.UserSavedStoriesHandler),
			middlewares.Auth))

	http.HandleFunc(
		"/users/stories/submitted",
		middlewares.Middleware(
			middlewares.Error(controllers.UserSubmittedStoriesHandler), middlewares.Auth))

	http.HandleFunc(
		"/users/stories/upvoted",
		middlewares.Middleware(
			middlewares.Error(controllers.UserUpvotedStoriesHandler),
			middlewares.Auth))

	http.Handle(
		"/reset-password",
		middlewares.Error(controllers.ResetPasswordHandler))

	http.HandleFunc(
		"/set-new-password",
		middlewares.Error(controllers.SetNewPasswordHandler))

	http.HandleFunc(
		"/change-password",
		middlewares.Middleware(
			middlewares.Error(controllers.ChangePasswordHandler),
			middlewares.Auth))

	http.HandleFunc(
		"/profile-edit",
		middlewares.Middleware(
			middlewares.Error(controllers.UserProfileHandler),
			middlewares.Auth))

	http.HandleFunc(
		"/invitecodes/generate",
		middlewares.Error(controllers.GenerateInviteCodeHandler))

	http.HandleFunc(
		"/customer-signup",
		middlewares.Error(controllers.CustomerSignUpHandler))

	http.HandleFunc(
		"/users/invite",
		middlewares.Middleware(
			middlewares.Error(controllers.InviteUserHandler),
			middlewares.Auth))

	http.HandleFunc(
		"/admin",
		middlewares.Middleware(
			middlewares.Error(controllers.AdminHandler),
			middlewares.Auth))
}
