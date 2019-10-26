package main

import (
	"net/http"
	"turkdev/app/controllers"
	"turkdev/app/templates"
)

func main() {

	/*user := data.User{
		UserName:     "hpy",
		FullName:     "Huseyin Polat Yuruk",
		Email:        "h.poaltyuruk@gmail.com",
		RegisteredOn: time.Now(),
		Password:     "111111",
		Website:      "http://huseyinpolatyuruk.com",
		About:        "Software Developer",
		Invitedby:    "anil",
		InviteCode:   "abcdef",
		Karma:        12,
	}
	err := data.CreateUser(&user)
	if err != nil {
		fmt.Println(err)
	}
	story := data.Story{
		URL:          "http://huseyinpolatyuruk.com",
		Title:        "Test Title",
		Text:         "Test Text",
		Tags:         []string{"programming", "coding", "web"},
		UpVotes:      0,
		DownVotes:    0,
		CommentCount: 0,
		UserID:       16,
		SubmittedOn:  time.Now(),
	}
	err := data.CreateStory(&story)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Succeded!")*/
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
