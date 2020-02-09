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
	var model = &models.StoryPageViewModel{Title: "Stories"}

	var userID int = -1
	var customerID int = 1 // TODO: get actual customer id from registered customer website

	isAuthenticated, user, err := shared.IsAuthenticated(r)

	if err != nil {
		return err
	}

	if isAuthenticated {
		userID = user.ID
		customerID = user.CustomerID
		model.IsAuthenticated = isAuthenticated
		model.SignedInUser = models.SignedInUserViewModel{
			UserID:     user.ID,
			UserName:   user.UserName,
			CustomerID: user.CustomerID,
			Email:      user.Email,
		}
	}

	var page int = 0
	strPage := r.URL.Query().Get("page")
	if len(strPage) > 0 {
		page, _ = strconv.Atoi(strPage)
	}
	stories, err := data.GetStories(customerID, page, DefaultStoryCountPerPage)
	if err != nil {
		return err
	}
	if stories == nil || len(*stories) == 0 {
		// TODO: There is no story yet. Show appropriate message here
		return nil
	}
	model.Stories = *mapStoriesToStoryViewModel(stories, userID)
	templates.RenderWithBase(w, "stories/index.html", model)
	return nil
}

/*RecentStoriesHandler handles showing recently published stories*/
func RecentStoriesHandler(w http.ResponseWriter, r *http.Request) error {
	var model = &models.StoryPageViewModel{Title: "Recent Stories | Turk Dev"}

	var userID int = -1    // TODO: get actual user id here
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
	model.Stories = *mapStoriesToStoryViewModel(stories, userID)
	templates.RenderWithBase(w, "stories/index.html", model)
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

	templates.RenderWithBase(
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

/*StoryDetailHandler handles showing comments by giving story id*/
func StoryDetailHandler(w http.ResponseWriter, r *http.Request) error {
	title := "Comments | Turk Dev"
	user := models.User{"Anil Yuzener"}

	strStoryID := r.URL.Query().Get("storyid")
	if len(strStoryID) == 0 {
		// TODO: Story cannot be found. Show appropriate message here.
		return nil
	}
	storyID, err := strconv.Atoi(strStoryID)
	if err != nil {
		// TODO(Anil): Cannot parse to int. Show story not found message.
		return nil
	}
	comments, err := data.GetComments(storyID)
	if err != nil {
		// TODO(Anil): show error page here
	}
	if comments == nil || len(*comments) == 0 {
		// TODO(Anil): There is no comment yet. Show appropriate message here
	}

	data := map[string]interface{}{
		"Content":  "Comments",
		"Comments": comments,
	}

	templates.RenderWithBase(
		w,
		"stories/detail.html",
		models.ViewModel{
			title,
			user,
			data,
		},
	)
	return nil
}

/*UpvoteStoryHandler runs when click to upvote story button. If not upvoted before by user, upvotes that story*/
func UpvoteStoryHandler(w http.ResponseWriter, r *http.Request) error {

	isAuthenticated, _, err := shared.IsAuthenticated(r)
	if err != nil {
		return err
	}
	if isAuthenticated == false {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return nil
	}
	if r.Method == "GET" {
		return fmt.Errorf("Reqeuest should be POST request")
	}

	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("Error occured when parse from. Error : %v", err)
	}

	userIDStr := r.FormValue("userID")
	storyIDStr := r.FormValue("storyID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return fmt.Errorf("Error occured when convert string userID to int userID. Error : %v", err)
	}
	storyID, err := strconv.Atoi(storyIDStr)
	if err != nil {
		return fmt.Errorf("Error occured when convert string storyID to int storyID. Error : %v", err)
	}

	isUpvoted, err := data.CheckIfStoryUpVotedByUser(userID, storyID)
	if err != nil {
		return fmt.Errorf("Error occured when check user save story. Error : %v", err)
	}
	if isUpvoted {
		return nil
	}

	err = data.UpVoteStory(userID, storyID)
	if err != nil {
		return fmt.Errorf("Error occured when upvote story. Error : %v", err)
	}
	return nil
}

/*UnvoteStoryHandler handles unvote button. If a story voted by user before, this handler undo that operation*/
func UnvoteStoryHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return fmt.Errorf("Reqeuest should be POSTm request")
	}

	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("Error occured when parse from. Error : %v", err)
	}

	userIDStr := r.FormValue("userID")
	storyIDStr := r.FormValue("storyID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return fmt.Errorf("Error occured when convert string userID to int userID. Error : %v", err)
	}
	storyID, err := strconv.Atoi(storyIDStr)
	if err != nil {
		return fmt.Errorf("Error occured when convert string storyID to int storyID. Error : %v", err)
	}

	isUpvoted, err := data.CheckIfStoryUpVotedByUser(userID, storyID)
	if err != nil {
		return fmt.Errorf("Error occured when check user save story. Error : %v", err)
	}
	if isUpvoted == false {
		return nil
	}

	err = data.UnVoteStory(userID, storyID)
	if err != nil {
		return fmt.Errorf("Error occured when unvote story. Error : %v", err)
	}
	return nil
}

/*SaveStoryHandler saves a story for user*/
func SaveStoryHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return fmt.Errorf("Reqeuest should be POSTm request")
	}

	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("Error occured when parse from. Error : %v", err)
	}

	userIDStr := r.FormValue("userID")
	storyIDStr := r.FormValue("storyID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return fmt.Errorf("Error occured when convert string userID to int userID. Error : %v", err)
	}
	storyID, err := strconv.Atoi(storyIDStr)
	if err != nil {
		return fmt.Errorf("Error occured when convert string storyID to int storyID. Error : %v", err)
	}

	isSaved, err := data.CheckIfUserSavedStory(userID, storyID)
	if err != nil {
		return fmt.Errorf("Error occured when check user save story. Error : %v", err)
	}

	if isSaved {
		return nil
	}

	err = data.SaveStory(userID, storyID)
	if err != nil {
		return fmt.Errorf("Error occured when save story. Error : %v", err)
	}
	return nil
}

/*UnSaveStoryHandler unsaves a story if user save that story*/
func UnSaveStoryHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return fmt.Errorf("Reqeuest should be POSTm request")
	}

	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("Error occured when parse from. Error : %v", err)
	}

	userIDStr := r.FormValue("userID")
	storyIDStr := r.FormValue("storyID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return fmt.Errorf("Error occured when convert string userID to int userID. Error : %v", err)
	}
	storyID, err := strconv.Atoi(storyIDStr)
	if err != nil {
		return fmt.Errorf("Error occured when convert string storyID to int storyID. Error : %v", err)
	}

	isSaved, err := data.CheckIfUserSavedStory(userID, storyID)
	if err != nil {
		return fmt.Errorf("Error occured when check user save story. Error : %v", err)
	}

	if isSaved == false {
		return nil
	}

	err = data.UnSaveStory(userID, storyID)
	if err != nil {
		return fmt.Errorf("Error occured when unsave story. Error : %v", err)
	}
	return nil
}

func handlesSubmitGET(w http.ResponseWriter, r *http.Request) error {

	title := "Submit Story | Turk Dev"
	user := models.User{"Anil Yuzener"}

	data := map[string]interface{}{
		"Content": "Submit Story",
	}

	templates.RenderWithBase(
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

func mapStoriesToStoryViewModel(stories *[]data.Story, signedInUserID int) *[]models.StoryViewModel {
	var viewModels []models.StoryViewModel

	for _, story := range *stories {
		var viewModel = models.StoryViewModel{
			ID:                    story.ID,
			Title:                 story.Title,
			URL:                   story.URL,
			Points:                story.UpVotes, // TODO: call point calculation function here
			UserID:                story.UserID,
			UserName:              story.UserName,
			CommentCount:          story.CommentCount,
			IsUpvotedSignedInUser: false,
			IsSavedBySignedInUser: false,
			SubmittedOnText:       generateSubmittedOnText(story.SubmittedOn),
		}

		if signedInUserID > -1 {
			isUpvoted, err := data.CheckIfStoryUpVotedByUser(signedInUserID, story.ID)

			if err != nil {
				// TODO: log error here
				isUpvoted = false
			}

			isSaved, err := data.CheckIfUserSavedStory(signedInUserID, story.ID)

			if err != nil {
				// TODO: log error here
				isSaved = false
			}
			viewModel.IsUpvotedSignedInUser = isUpvoted
			viewModel.IsSavedBySignedInUser = isSaved
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
