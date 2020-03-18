package main

import (
	"net/http"
	"turkdev/src/controllers"
	"turkdev/src/middlewares"
)

func main() {

	registerHandlers()

	staticFileServer := http.FileServer(http.Dir("dist/"))
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
		"/stories/vote",
		middlewares.Middleware(
			controllers.VoteStoryHandler,
			middlewares.Auth))

	http.HandleFunc(
		"/stories/remove/vote",
		middlewares.Middleware(
			controllers.RemoveStoryVoteHandler,
			middlewares.Auth))

	http.HandleFunc(
		"/stories/save",
		middlewares.Middleware(
			controllers.SaveStoryHandler,
			middlewares.Auth))

	http.HandleFunc(
		"/stories/unsave",
		middlewares.Middleware(
			controllers.UnSaveStoryHandler,
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
		"/comments/vote",
		middlewares.Middleware(
			controllers.VoteCommentHandler,
			middlewares.Auth))

	http.HandleFunc(
		"/comments/remove/vote",
		middlewares.Middleware(
			controllers.RemoveCommentVoteHandler,
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
