package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"turkdev/app/models"
	"turkdev/app/templates"
	"turkdev/data"
	"turkdev/shared"
)

const (
	/*DefaultStoryCountPerPage represents story count to be listed per page*/
	DefaultStoryCountPerPage = 20
)

/*StoriesHandler handles showing the popular published stories*/
func StoriesHandler(w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path == "/" || r.URL.Path == "/index.html" {
		title := "Turk Dev"
		user := models.User{"Anil Yuzener"}

		var page int = 0
		strPage := r.URL.Query().Get("page")
		if len(strPage) > 0 {
			page, _ = strconv.Atoi(strPage)
		}
		stories, err := data.GetStories(1, page, DefaultStoryCountPerPage)
		if err != nil {
			// TODO(Anil): Show error page here
		}
		if stories == nil || len(*stories) == 0 {
			// TODO(Anil): There is no story yet. Show appropriate message here
		}

		pageData := map[string]interface{}{
			"Content": "Stories",
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
		return nil
	}

	if r.URL.Path == "/robots.txt" {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("User-agent: *\nDisallow: /"))
		return nil
	}

	byteValue, err := shared.ReadFile("app/static/html/404.html")
	if err != nil {
		return fmt.Errorf("Error occured in readFile func. Original err : %v", err)
	}

	w.WriteHeader(404)
	w.Write(byteValue)
	return nil
}

/*RecentStoriesHandler handles showing recently published stories*/
func RecentStoriesHandler(w http.ResponseWriter, r *http.Request) error {
	title := "Recent Stories | Turk Dev"
	user := models.User{"Anil Yuzener"}

	var page int = 0
	strPage := r.URL.Query().Get("page")
	if len(strPage) > 0 {
		page, _ = strconv.Atoi(strPage)
	}
	stories, err := data.GetRecentStories(1, page, DefaultStoryCountPerPage)
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
	return nil
}

/*SavedStoriesHandler handles showing the saved stories of a user*/
func SavedStoriesHandler(w http.ResponseWriter, r *http.Request) error {
	title := "Saved Stories | Turk Dev"
	user := models.User{"Anil Yuzener"}

	var userID int = 0
	strUserID := r.URL.Query().Get("userID")
	if len(strUserID) == 0 {
		// TODO (Anil): Show user not found message.
		return nil
	}
	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		// TODO(Anil): Cannot parse to int. Show user not found message.
		return nil
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
	return nil
}

/*SubmitStoryHandler handles to submit a new story*/
func SubmitStoryHandler(w http.ResponseWriter, r *http.Request) error {
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
	return nil
}
