package main

import (
	"fmt"
	"linkwind/app/controllers"
	"linkwind/app/middlewares"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
)

func init() {
	envFileName := ".env.dev"
	env := os.Getenv("APP_ENV")
	if env == "production" {
		envFileName = ".env"
	}
	err := godotenv.Load(envFileName)
	if err != nil {
		log.Fatalf("Error loading .env file. Error: %v", err)
	}

	fmt.Println(fmt.Sprintf("App is initialized in %s mode", env))

	err = sentry.Init(sentry.ClientOptions{
		// Either set your DSN here or set the SENTRY_DSN environment variable.
		Dsn: os.Getenv("SENTRY_DSN"),
		// Enable printing of SDK debug messages.
		// Useful when getting started or trying to figure something out.
		Debug: false,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	defer sentry.Flush(2 * time.Second)
}

func main() {
	router := http.NewServeMux()
	configuredRouter := configureRouter(router)

	port, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	fmt.Println(fmt.Sprintf("Application is work on port %d", port))
	// Start our HTTP server
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), configuredRouter); err != nil {
		sentry.CaptureException(err)
		os.Exit(1)
	}
}

func configureRouter(router *http.ServeMux) http.Handler {

	profilePath := "/users/profile"
	changePasswordPath := "/change-password"
	profileEditPath := "/profile-edit"
	userInvitePath := "/users/invite"
	adminPath := "/admin"
	storyVotePath := "/stories/vote"
	storyVoteRemovePath := "/stories/remove/vote"
	storySavePath := "/stories/save"
	storyUnsavePath := "/stories/unsave"
	submitStoryPath := "/submit"
	addCommentPath := "/comments/add"
	voteCommentPath := "/comments/vote"
	removeCommentVotePath := "/comments/remove/vote"
	replyCommentPath := "/comments/reply"
	userSavedStoryPath := "/users/stories/saved"
	userSubmittedStoryPath := "/users/stories/submitted"
	userUpvotedStoryPath := "/users/stories/upvoted"

	authorizedPaths := []string{
		profilePath,
		changePasswordPath,
		profileEditPath,
		userInvitePath,
		adminPath,
		storyVotePath,
		storyVoteRemovePath,
		storySavePath,
		storyUnsavePath,
		submitStoryPath,
		addCommentPath,
		voteCommentPath,
		removeCommentVotePath,
		replyCommentPath,
		userSavedStoryPath,
		userSubmittedStoryPath,
		userUpvotedStoryPath,
	}

	staticFileServer := http.FileServer(http.Dir("public/"))
	router.Handle("/public/", http.StripPrefix("/public/", staticFileServer))

	router.HandleFunc("/", controllers.StoriesHandler)
	router.HandleFunc("/recent", controllers.RecentStoriesHandler)
	router.HandleFunc("/signup", controllers.SignUpHandler)
	router.HandleFunc("/signin", controllers.SignInHandler)
	router.HandleFunc("/signout", controllers.SignOutHandler)
	router.HandleFunc("/reset-password", controllers.ResetPasswordHandler)
	router.HandleFunc("/set-new-password", controllers.SetNewPasswordHandler)
	router.HandleFunc("/stories/detail", controllers.StoryDetailHandler)
	router.HandleFunc("/invitecodes/generate", controllers.GenerateInviteCodeHandler)
	router.HandleFunc(profilePath, controllers.UserProfileHandler)
	router.HandleFunc(changePasswordPath, controllers.ChangePasswordHandler)
	router.HandleFunc(profileEditPath, controllers.UserProfileHandler)
	router.HandleFunc(userInvitePath, controllers.InviteUserHandler)
	router.HandleFunc(adminPath, controllers.AdminHandler)
	router.HandleFunc(storyVotePath, controllers.VoteStoryHandler)
	router.HandleFunc(storyVoteRemovePath, controllers.RemoveStoryVoteHandler)
	router.HandleFunc(storySavePath, controllers.SaveStoryHandler)
	router.HandleFunc(storyUnsavePath, controllers.UnSaveStoryHandler)
	router.HandleFunc(submitStoryPath, controllers.SubmitStoryHandler)
	router.HandleFunc(addCommentPath, controllers.AddCommentHandler)
	router.HandleFunc(voteCommentPath, controllers.VoteCommentHandler)
	router.HandleFunc(removeCommentVotePath, controllers.RemoveCommentVoteHandler)
	router.HandleFunc(replyCommentPath, controllers.ReplyToCommentHandler)
	router.HandleFunc(userSavedStoryPath, controllers.UserSavedStoriesHandler)
	router.HandleFunc(userSubmittedStoryPath, controllers.UserSubmittedStoriesHandler)
	router.HandleFunc(userUpvotedStoryPath, controllers.UserUpvotedStoriesHandler)

	errorMiddleware := middlewares.ErrorMiddleWare()
	errorHandledRouter := errorMiddleware(router)

	authMiddleware := middlewares.AuthMiddleWare(authorizedPaths)
	authHandledRouter := authMiddleware(errorHandledRouter)

	return authHandledRouter
}
