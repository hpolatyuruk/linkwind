package main

import (
	"fmt"
	"linkwind/app/controllers"
	"linkwind/app/middlewares"
	"linkwind/app/shared"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
)

type RouteData struct {
	Path         string
	Handler      http.HandlerFunc
	IsAuthorized bool
}

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

	routes := []RouteData{
		{"/", controllers.StoriesHandler, false},
		{"/recent", controllers.RecentStoriesHandler, false},
		{"/signup", controllers.SignUpHandler, false},
		{"/signin", controllers.SignInHandler, false},
		{"/signout", controllers.SignOutHandler, false},
		{"/reset-password", controllers.ResetPasswordHandler, false},
		{"/set-new-password", controllers.SetNewPasswordHandler, false},
		{"/stories/detail", controllers.StoryDetailHandler, false},
		{"/exists-custom-domain", controllers.ExistsCustomDomain, false},
		{"/customer-signup", controllers.CustomerSignUpHandler, false},
		{"/about", controllers.AboutHandler, false},
		{"/faq", controllers.FAQHandler, false},
		{"/privacy", controllers.PrivacyHandler, false},
		{"/auth", controllers.SetAuthTokenHandler, false},
		{"/invitecodes/generate", controllers.GenerateInviteCodeHandler, false},
		{"/users/profile", controllers.UserProfileHandler, true},
		{"/change-password", controllers.ChangePasswordHandler, true},
		{"/profile-edit", controllers.UserProfileHandler, true},
		{"/users/invite", controllers.InviteUserHandler, true},
		{"/admin", controllers.AdminHandler, true},
		{"/stories/vote", controllers.VoteStoryHandler, true},
		{"/stories/remove/vote", controllers.RemoveStoryVoteHandler, true},
		{"/stories/save", controllers.SaveStoryHandler, true},
		{"/stories/unsave", controllers.UnSaveStoryHandler, true},
		{"/submit", controllers.SubmitStoryHandler, true},
		{"/comments/add", controllers.AddCommentHandler, true},
		{"/comments/vote", controllers.VoteCommentHandler, true},
		{"/comments/remove/vote", controllers.RemoveCommentVoteHandler, true},
		{"/comments/reply", controllers.ReplyToCommentHandler, true},
		{"/users/stories/saved", controllers.UserSavedStoriesHandler, true},
		{"/users/stories/submitted", controllers.UserSubmittedStoriesHandler, true},
		{"/users/stories/upvoted", controllers.UserUpvotedStoriesHandler, true},
	}

	staticFileServer := http.FileServer(http.Dir("public/"))
	router.Handle(shared.StaticFolderPath, http.StripPrefix(shared.StaticFolderPath, staticFileServer))

	var authorizedPaths = []string{}
	var allPaths = []string{}

	for _, route := range routes {

		allPaths = append(allPaths, route.Path)
		router.HandleFunc(route.Path, route.Handler)

		if route.IsAuthorized {
			authorizedPaths = append(authorizedPaths, route.Path)
		}
	}

	authMiddleware := middlewares.AuthMiddleWare(authorizedPaths)
	authHandledRouter := authMiddleware(router)

	notFoundMiddleware := middlewares.NotFoundMiddleware(allPaths)
	notFoundHandledRouter := notFoundMiddleware(authHandledRouter)

	customerMiddleware := middlewares.CustomerMiddleware()
	customerHandledRouter := customerMiddleware(notFoundHandledRouter)

	errorMiddleware := middlewares.ErrorMiddleware()
	errorHandledRouter := errorMiddleware(customerHandledRouter)

	return errorHandledRouter
}
