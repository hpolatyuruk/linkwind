package controllers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"
	"turkdev/app/models"
	"turkdev/app/src/templates"
	"turkdev/data"
	"turkdev/shared"
)

const (
	/*DefaultStoryCountPerPage represents story count to be listed per page*/
	DefaultStoryCountPerPage = 20
)

/*StoriesHandler handles showing the popular published stories*/
func StoriesHandler(w http.ResponseWriter, r *http.Request) error {
	title := "Turk Dev"
	user := models.User{"Anil Yuzener"}

	var customerID int = 1
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

	templates.Render(
		w,
		"stories/index.html",
		models.StoryPageData{
			User:    user,
			Title:   title,
			Stories: *mapStoriesToStoryViewModel(stories),
			//Stories: []data.Story{{URL: "1"}, {URL: "2"}},
		},
	)
	return nil
}

/*RecentStoriesHandler handles showing recently published stories*/
func RecentStoriesHandler(w http.ResponseWriter, r *http.Request) error {
	title := "Recent Stories | Turk Dev"
	user := models.User{"Anil Yuzener"}

	var customerID int = 1 // TODO: get actual customer id here
	var page int = 0
	strPage := r.URL.Query().Get("page")
	if len(strPage) > 0 {
		page, _ = strconv.Atoi(strPage)
	}
	stories, err := data.GetRecentStories(customerID, page, DefaultStoryCountPerPage)
	if err != nil {
		// TODO(Anil): Show error page here

	}

	templates.Render(
		w,
		"stories/index.html",
		models.StoryPageData{
			Title:   title,
			User:    user,
			Stories: *mapStoriesToStoryViewModel(stories),
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
	switch r.Method {
	case "GET":
		return handlesSubmitGET(w, r)
	case "POST":
		return handleSubmitPOST(w, r)
	default:
		return handlesSubmitGET(w, r)
	}
}

func handlesSubmitGET(w http.ResponseWriter, r *http.Request) error {

	title := "Submit Story | Turk Dev"
	user := models.User{"Anil Yuzener"}

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

func handleSubmitPOST(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	title := r.FormValue("title")
	url := r.FormValue("url")
	text := r.FormValue("text")

	if len(url) > 0 {
		fetchedTitle, err := shared.FetchURL(url)
		if err != nil {
			return err
		}
		title = fetchedTitle
	}

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
		return err
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func mapStoriesToStoryViewModel(stories *[]data.Story) *[]models.StoryViewModel {
	var viewModels []models.StoryViewModel

	for _, story := range *stories {
		var viewModel = models.StoryViewModel{
			ID:              story.ID,
			Title:           story.Title,
			URL:             story.URL,
			Points:          story.UpVotes, // TODO: call point calculation function here
			UserID:          story.UserID,
			UserName:        story.UserName,
			CommentCount:    story.CommentCount,
			SubmittedOnText: generateSubmittedOnText(story.SubmittedOn),
		}
		viewModels = append(viewModels, viewModel)
	}
	return &viewModels
}

func generateSubmittedOnText(submittedOn time.Time) string {
	var text string = ""
	diff := time.Now().Sub(submittedOn)

	if diff.Hours() < 1 {
		mins := int(math.Round(diff.Minutes()))
		text = fmt.Sprintf("%d minutes ago", mins)
		if mins == 1 {
			text = fmt.Sprintf("%d minute ago", mins)
		}
	} else if diff.Hours() < 24 {
		hours := int(math.Round(diff.Hours()))
		text = fmt.Sprintf("%d hours ago", hours)
		if hours == 1 {
			text = fmt.Sprintf("%d hour ago", hours)
		}
	} else {
		days := math.Round(diff.Hours() / 24)

		if days == 1 {
			text = fmt.Sprintf("%d day ago", int(days))
		} else if days > 1 && days < 30 {
			text = fmt.Sprintf("%d days ago", int(days))
		} else if days > 30 && days < 365 {
			months := int(math.Round(days / 30))
			text = fmt.Sprintf("%d months ago", months)
			if months == 1 {
				text = fmt.Sprintf("%d month ago", months)
			}

		} else {
			years := int(math.Round(days / 365))
			text = fmt.Sprintf("%d years ago", years)
			if years == 1 {
				text = fmt.Sprintf("%d year ago", years)
			}
		}
	}
	return text
}
