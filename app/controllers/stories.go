package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"turkdev/app/models"
	"turkdev/app/src/templates"
	"turkdev/data"
)

const (
	/*DefaultStoryCountPerPage represents story count to be listed per page*/
	DefaultStoryCountPerPage = 20
)

/*StoriesHandler handles showing the popular published stories*/
func StoriesHandler(w http.ResponseWriter, r *http.Request) {
	title := "Turk Dev"
	user := models.User{"Anil Yuzener"}

	var customerID int = 0
	var page int = 0
	strPage := r.URL.Query().Get("page")
	if len(strPage) > 0 {
		page, _ = strconv.Atoi(strPage)
	}
	stories, err := data.GetStories(customerID, page, DefaultStoryCountPerPage)
	if err != nil {
		// TODO(Anil): Show error page here
	}
	if stories == nil || len(*stories) == 0 {
		// TODO(Anil): There is no story yet. Show appropriate message here
	}

	pageData := map[string]interface{}{
		"Stories": stories,
	}

	templates.Render(
		w,
		"stories/index.html",
		models.ViewModel{
			title,
			user,
			pageData,
		},
	)
}

/*RecentStoriesHandler handles showing recently published stories*/
func RecentStoriesHandler(w http.ResponseWriter, r *http.Request) {
	title := "Recent Stories | Turk Dev"
	user := models.User{"Anil Yuzener"}

	var customerID int = 0
	var page int = 0
	strPage := r.URL.Query().Get("page")
	if len(strPage) > 0 {
		page, _ = strconv.Atoi(strPage)
	}
	stories, err := data.GetRecentStories(customerID, page, DefaultStoryCountPerPage)
	if err != nil {
		// TODO(Anil): Show error page here

	}

	data := map[string]interface{}{
		"Content": "Recent Stories",
		"Stories": stories,
	}

	templates.Render(
		w,
		"stories/index.html",
		models.ViewModel{
			title,
			user,
			data,
		},
	)
}

/*SavedStoriesHandler handles showing the saved stories of a user*/
func SavedStoriesHandler(w http.ResponseWriter, r *http.Request) {
	title := "Saved Stories | Turk Dev"
	user := models.User{"Anil Yuzener"}

	var userID int = 0
	strUserID := r.URL.Query().Get("userID")
	if len(strUserID) == 0 {
		// TODO (Anil): Show user not found message.
		return
	}
	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		// TODO(Anil): Cannot parse to int. Show user not found message.
		return
	}
	var page int = 0
	strPage := r.URL.Query().Get("page")
	if len(strPage) > 0 {
		page, _ = strconv.Atoi(strPage)
	}
	stories, err := data.GetUserSavedStories(userID, page, DefaultStoryCountPerPage)
	if err != nil {
		// TODO(Anil): Show error page here
	}

	data := map[string]interface{}{
		"Content": "Saved Stories",
		"Stories": stories,
	}

	templates.Render(
		w,
		"stories/index.html",
		models.ViewModel{
			title,
			user,
			data,
		},
	)
}

/*SubmitStoryHandler handles to submit a new story*/
func SubmitStoryHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		handleSubmitGET(w, r)
	case "POST":
		handleSubmitPOST(w, r)
	default:
		handleSubmitGET(w, r)
	}
}

func handleSubmitGET(w http.ResponseWriter, r *http.Request) {
	title := "Submit Story | Turk Dev"
	user := models.User{"Anil Yuzener"}

	// TODO(Anil): Get story data from html form and map them to below story struct

	var story data.Story = data.Story{}
	err := data.CreateStory(&story)
	if err != nil {
		// TODO(Anil): Show error page here
	}

	data := map[string]interface{}{
		"Content": "Submit Story",
	}

	templates.Render(
		w,
		"stories/submit.html",
		models.ViewModel{
			title,
			user,
			data,
		},
	)
}

func handleSubmitPOST(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		// TODO: log err here
		fmt.Fprintf(w, "Story submit form parsing error: %v", err)
	}

	title := r.FormValue("title")
	url := r.FormValue("url")
	text := r.FormValue("text")

	var story data.Story
	story.Title = title
	story.URL = url
	story.Text = text
	story.CommentCount = 0
	story.UpVotes = 0
	story.SubmittedOn = time.Now()
	story.UserID = 2 // TODO: use actual user id here

	err := data.CreateStory(&story)

	if err != nil {
		// TODO: log error here
		fmt.Fprintf(w, "Error creating story: %v", err)
		return
	}
	fmt.Fprintf(w, "Succeeded!")
}
