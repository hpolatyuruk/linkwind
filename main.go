package main

import (
	"fmt"
	"turkdev/data"
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
		URL:          "http://huseyinpolatyuruk6.com",
		Title:        "Test Title6",
		Text:         "Test Text6",
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
	stories, err := data.GetStories(1, 3)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(len(stories))
	/*templates.Initialize()

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

	http.ListenAndServe(":80", nil)*/
}
