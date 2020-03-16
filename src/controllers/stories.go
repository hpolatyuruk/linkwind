package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"turkdev/src/data"
	"turkdev/src/enums"
	"turkdev/src/models"
	"turkdev/src/shared"
	"turkdev/src/templates"
)

const (
	/*DefaultPageSize represents story count to be listed per page*/
	DefaultPageSize = 15
)

/*StoryVoteModel represents the data in http request body to upvote story.*/
type StoryVoteModel struct {
	StoryID  int
	UserID   int
	VoteType enums.VoteType
}

/*StorySaveModel represents the data in http request body to save story.*/
type StorySaveModel struct {
	StoryID int
	UserID  int
}

/*StorySubmitModel represents the data to submit a story.*/
type StorySubmitModel struct {
	URL          string
	Title        string
	Text         string
	Errors       map[string]string
	SignedInUser *models.SignedInUserViewModel
}

/*JSONResponse respresents the json response.*/
type JSONResponse struct {
	Result string
}

type getStoriesPaged func(customerID, pageNo, storyCountPerPage int) (*[]data.Story, error)

/*Validate validates the StorySubmitModel*/
func (model *StorySubmitModel) Validate() bool {
	model.Errors = make(map[string]string)
	if strings.TrimSpace(model.URL) == "" &&
		strings.TrimSpace(model.Text) == "" &&
		strings.TrimSpace(model.Text) == "" {
		model.Errors["General"] = "Please enter a url or title/text."
		return false
	}
	if strings.TrimSpace(model.URL) == "" &&
		strings.TrimSpace(model.Text) != "" &&
		strings.TrimSpace(model.Title) == "" {
		model.Errors["Title"] = "Please enter a title."
		return false
	}
	return true
}

/*StoriesHandler handles showing the popular published stories*/
func StoriesHandler(w http.ResponseWriter, r *http.Request) error {
	return renderStoriesPage("Stories", data.GetStories, w, r)
}

/*RecentStoriesHandler handles showing recently published stories*/
func RecentStoriesHandler(w http.ResponseWriter, r *http.Request) error {
	return renderStoriesPage("Recent Stories", data.GetRecentStories, w, r)
}

func renderStoriesPage(title string, fnGetStories getStoriesPaged, w http.ResponseWriter, r *http.Request) error {
	var model = &models.StoryPageViewModel{Title: title}
	var customerID int = 1 // TODO: get actual customer id from registered customer website
	isAuthenticated, user, err := shared.IsAuthenticated(r)
	if err != nil {
		return err
	}
	if isAuthenticated {
		customerID = user.CustomerID
		model.IsAuthenticated = isAuthenticated
		model.SignedInUser = &models.SignedInUserViewModel{
			UserID:     user.ID,
			UserName:   user.UserName,
			CustomerID: user.CustomerID,
			Email:      user.Email,
		}
	}
	var page int = getPage(r)
	stories, err := fnGetStories(customerID, page, DefaultPageSize)
	if err != nil {
		return err
	}
	storiesCount, err := data.GetCustomerStoriesCount(customerID)
	if err != nil {
		return err
	}
	pagingModel, err := setPagingViewModel(customerID, page, storiesCount)
	if err != nil {
		return err
	}
	model.Page = pagingModel
	if stories != nil && len(*stories) > 0 {
		model.Stories = *mapStoriesToStoryViewModel(stories, model.SignedInUser)
	}
	templates.RenderInLayout(w, "stories.html", model)
	return nil
}

func setPagingViewModel(customerID, currentPage, storiesCount int) (*models.Paging, error) {
	totalPageCount := calcualteTotalPageCount(storiesCount)
	isFinalPage := currentPage == totalPageCount
	model := &models.Paging{
		CurrentPage:    currentPage,
		NextPage:       currentPage + 1,
		PreviousPage:   currentPage - 1,
		IsFinalPage:    isFinalPage,
		TotalPageCount: totalPageCount,
	}
	return model, nil
}

func calcualteTotalPageCount(storiesCount int) int {
	pageCount := math.Ceil(float64(storiesCount) / float64(DefaultPageSize))
	return int(pageCount)
}

/*UserSavedStoriesHandler handles showing the saved stories of a user*/
func UserSavedStoriesHandler(w http.ResponseWriter, r *http.Request) error {
	var model = &models.StoryPageViewModel{Title: "Saved Stories"}

	isAuthenticated, user, err := shared.IsAuthenticated(r)
	if err != nil {
		return err
	}

	if isAuthenticated {
		model.IsAuthenticated = isAuthenticated
		model.SignedInUser = &models.SignedInUserViewModel{
			UserID:     user.ID,
			UserName:   user.UserName,
			CustomerID: user.CustomerID,
			Email:      user.Email,
		}
	}

	var page int = getPage(r)
	stories, err := data.GetUserSavedStories(user.ID, page, DefaultPageSize)
	if err != nil {
		return err
	}
	storiesCount, err := data.GetUserSavedStoriesCount(user.ID)
	if err != nil {
		return err
	}
	pagingModel, err := setPagingViewModel(user.CustomerID, page, storiesCount)
	if err != nil {
		return err
	}
	model.Page = pagingModel
	if stories == nil || len(*stories) > 0 {
		model.Stories = *mapStoriesToStoryViewModel(stories, model.SignedInUser)
	}
	err = templates.RenderInLayout(w, "stories.html", model)
	if err != nil {
		return err
	}
	return nil
}

/*UserSubmittedStoriesHandler handles user's submitted stories*/
func UserSubmittedStoriesHandler(w http.ResponseWriter, r *http.Request) error {
	var model = &models.StoryPageViewModel{Title: "Submitted Stories"}
	isAuthenticated, user, err := shared.IsAuthenticated(r)

	if err != nil {
		return err
	}
	if isAuthenticated {
		model.IsAuthenticated = isAuthenticated
		model.SignedInUser = &models.SignedInUserViewModel{
			UserID:     user.ID,
			UserName:   user.UserName,
			CustomerID: user.CustomerID,
			Email:      user.Email,
		}
	}
	var page int = getPage(r)
	stories, err := data.GetUserSubmittedStories(user.ID, page, DefaultPageSize)
	fmt.Print(len(*stories))
	if err != nil {
		fmt.Println(err)
		return err
	}
	storiesCount, err := data.GetUserSubmittedStoriesCount(user.ID)
	if err != nil {
		return err
	}
	pagingModel, err := setPagingViewModel(user.CustomerID, page, storiesCount)
	if err != nil {
		return err
	}
	model.Page = pagingModel
	if stories == nil || len(*stories) > 0 {
		model.Stories = *mapStoriesToStoryViewModel(stories, model.SignedInUser)
	}
	err = templates.RenderInLayout(w, "stories.html", model)
	if err != nil {
		return err
	}
	return nil
}

/*UserUpvotedStoriesHandler handles showing the upvoted stories by user*/
func UserUpvotedStoriesHandler(w http.ResponseWriter, r *http.Request) error {
	var model = &models.StoryPageViewModel{Title: "Upvoted Stories"}

	var userID int = -1
	isAuthenticated, user, err := shared.IsAuthenticated(r)

	if err != nil {
		return err
	}
	if isAuthenticated {
		userID = user.ID
		model.IsAuthenticated = isAuthenticated
		model.SignedInUser = &models.SignedInUserViewModel{
			UserID:     user.ID,
			UserName:   user.UserName,
			CustomerID: user.CustomerID,
			Email:      user.Email,
		}
	}
	var page int = getPage(r)
	stories, err := data.GetUserUpvotedStories(userID, page, DefaultPageSize)
	if err != nil {
		return err
	}
	storiesCount, err := data.GetUserUpvotedStoriesCount(user.ID)
	if err != nil {
		return err
	}
	pagingModel, err := setPagingViewModel(user.CustomerID, page, storiesCount)
	if err != nil {
		return err
	}
	model.Page = pagingModel
	if stories == nil || len(*stories) > 0 {
		model.Stories = *mapStoriesToStoryViewModel(stories, model.SignedInUser)
	}
	err = templates.RenderInLayout(w, "stories.html", model)
	if err != nil {
		return err
	}
	return nil
}

/*SubmitStoryHandler handles to submit a new story*/
func SubmitStoryHandler(w http.ResponseWriter, r *http.Request) error {
	user := shared.GetUser(r)
	switch r.Method {
	case "GET":
		return handlesSubmitGET(w, r, user)
	case "POST":
		return handleSubmitPOST(w, r, user)
	default:
		return handlesSubmitGET(w, r, user)
	}
}

func handlesSubmitGET(w http.ResponseWriter, r *http.Request, user *shared.SignedInUserClaims) error {
	model := &StorySubmitModel{
		SignedInUser: &models.SignedInUserViewModel{
			UserName:   user.UserName,
			UserID:     user.ID,
			CustomerID: user.CustomerID,
			Email:      user.Email,
		},
	}

	templates.RenderInLayout(w, "submit.html", model)
	return nil
}

func handleSubmitPOST(w http.ResponseWriter, r *http.Request, user *shared.SignedInUserClaims) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	model := &StorySubmitModel{
		URL:   r.FormValue("url"),
		Title: r.FormValue("title"),
		Text:  r.FormValue("text"),
		SignedInUser: &models.SignedInUserViewModel{
			UserName:   user.UserName,
			UserID:     user.ID,
			CustomerID: user.CustomerID,
			Email:      user.Email,
		},
	}
	if model.Validate() == false {
		templates.RenderInLayout(w, "submit.html", model)
		return nil
	}
	if strings.TrimSpace(model.URL) != "" &&
		strings.TrimSpace(model.Title) == "" {
		fetchedTitle, err := shared.FetchURL(model.URL)
		if err != nil {
			model.Errors["URL"] = "Something went wrong while fetching URL. Please make sure that you entered a valid URL."
			templates.RenderInLayout(w, "submit.html", model)
			return nil
		}
		model.Title = fetchedTitle
	}
	var story data.Story
	story.Title = model.Title
	story.URL = model.URL
	story.Text = model.Text
	story.CommentCount = 0
	story.UpVotes = 0
	story.SubmittedOn = time.Now()
	story.UserID = user.ID

	err := data.CreateStory(&story)
	if err != nil {
		// TODO: log error here
		fmt.Fprintf(w, "Error creating story: %v", err)
		return err
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

/*StoryDetailHandler handles showing comments by giving story id*/
func StoryDetailHandler(w http.ResponseWriter, r *http.Request) error {
	strStoryID := r.URL.Query().Get("id")
	if len(strStoryID) == 0 {
		err := templates.RenderFile(w, "errors/404.html", nil)
		if err != nil {
			return err
		}
		return nil
	}
	storyID, err := strconv.Atoi(strStoryID)
	if err != nil {
		return fmt.Errorf("Cannot convert string StoryID to int. Original err : %v", err)
	}
	story, err := data.GetStoryByID(storyID)
	if err != nil {
		return fmt.Errorf("Cannot get story from db (StoryID : %d). Original err : %v", storyID, err)
	}
	if story == nil {
		err = templates.RenderFile(w, "errors/404.html", nil)
		if err != nil {
			return err
		}
		return nil
	}
	comments, err := data.GetRootCommentsByStoryID(storyID)
	if err != nil {
		return fmt.Errorf("Cannot get comments from db (StoryID : %d). Original err : %v", storyID, err)
	}
	isAuth, signedInUserClaims, err := shared.IsAuthenticated(r)
	if err != nil {
		return fmt.Errorf("An error occured when run IsAuthenticated func in StoryDetailHandler")
	}
	model := &models.StoryDetailPageViewModel{
		Title: story.Title,
	}
	if isAuth {
		model.IsAuthenticated = true
		model.SignedInUser = &models.SignedInUserViewModel{
			UserID:     signedInUserClaims.ID,
			CustomerID: signedInUserClaims.CustomerID,
			Email:      signedInUserClaims.Email,
			UserName:   signedInUserClaims.UserName,
		}
	}
	model.Story = mapStoryToStoryViewModel(story, model.SignedInUser)
	model.Comments = mapCommentsToViewModelsWithChildren(comments, model.SignedInUser, storyID)
	templates.RenderInLayout(w, "detail.html", model)
	return nil
}

/*VoteStoryHandler runs when click to upvote and downvote story button. If not voted before by user, votes that story*/
func VoteStoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Error(w, "Unsupported method. Only post method is supported.", http.StatusMethodNotAllowed)
		return
	}
	var model StoryVoteModel
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(model)
	isUpvoted, err := data.CheckIfStoryVotedByUser(model.UserID, model.StoryID, model.VoteType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error occured while checking if user already voted. UserID: %d, StoryID: %d, VoteType: %d,  Error : %v", model.UserID, model.StoryID, model.VoteType, err), http.StatusInternalServerError)
		return
	}
	if isUpvoted {
		res, _ := json.Marshal(&JSONResponse{
			Result: "AlreadyVoted",
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	}
	err = data.VoteStory(model.UserID, model.StoryID, model.VoteType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error occured while voting story. UserID: %d, StoryID: %d, VoteType: %d, Error : %v", model.UserID, model.StoryID, model.VoteType, err), http.StatusInternalServerError)
	}
	res, _ := json.Marshal(&JSONResponse{
		Result: "Voted",
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

/*RemoveStoryVoteHandler handles removing upvote and downvote button. If a story voted by user before, this handler undo that operation*/
func RemoveStoryVoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Error(w, "Unsupported request method. Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}

	var model StoryVoteModel
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isUpvoted, err := data.CheckIfStoryVotedByUser(model.UserID, model.StoryID, model.VoteType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error occured while checking user story vote. UserID: %d, StoryID: %d, VoteType: %d, Error : %v", model.UserID, model.StoryID, model.VoteType, err), http.StatusInternalServerError)
		return
	}
	if isUpvoted == false {
		res, _ := json.Marshal(&JSONResponse{
			Result: "AlreadyUnvoted",
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	}

	err = data.RemoveStoryVote(model.UserID, model.StoryID, model.VoteType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error occured while unvoting story. UserID: %d, StoryID: %d, VoteType: %d, Error : %v", model.UserID, model.StoryID, model.VoteType, err), http.StatusInternalServerError)
		return
	}
	res, _ := json.Marshal(&JSONResponse{
		Result: "Unvoted",
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

/*SaveStoryHandler saves a story for user*/
func SaveStoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Error(w, "Unsupported method. Only post method is supported.", http.StatusMethodNotAllowed)
		return
	}
	var model StorySaveModel
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isSaved, err := data.CheckIfUserSavedStory(model.UserID, model.StoryID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "An error occured while parsing json.", http.StatusInternalServerError)
		return
	}

	if isSaved {
		res, _ := json.Marshal(&JSONResponse{
			Result: "AlreadySaved",
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(res))
		return
	}

	err = data.SaveStory(model.UserID, model.StoryID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error occured while saving story", http.StatusInternalServerError)
		return
	}

	res, _ := json.Marshal(&JSONResponse{
		Result: "Saved",
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

/*UnSaveStoryHandler unsaves a story if user save that story*/
func UnSaveStoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Error(w, "Unsupported method. Only post method is supported.", http.StatusMethodNotAllowed)
		return
	}

	var model StorySaveModel
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "An error occured while parsing json.", http.StatusBadRequest)
		return
	}

	isSaved, err := data.CheckIfUserSavedStory(model.UserID, model.StoryID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "An error occured while unsaving json.", http.StatusInternalServerError)
		return
	}

	if isSaved == false {
		res, _ := json.Marshal(&JSONResponse{
			Result: "AlreadyUnsaved",
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(res))
		return
	}

	err = data.UnSaveStory(model.UserID, model.StoryID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error occured when unsave story", http.StatusInternalServerError)
		return
	}
	res, _ := json.Marshal(&JSONResponse{
		Result: "Unsaved",
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

func getPage(r *http.Request) int {
	var page int = 1
	strPage := r.URL.Query().Get("page")
	if len(strPage) > 0 {
		page, _ = strconv.Atoi(strPage)
	}
	return page
}

func mapStoriesToStoryViewModel(stories *[]data.Story, signedInUser *models.SignedInUserViewModel) *[]models.StoryViewModel {
	var viewModels []models.StoryViewModel

	for _, story := range *stories {
		viewModel := mapStoryToStoryViewModel(&story, signedInUser)
		viewModels = append(viewModels, *viewModel)
	}
	return &viewModels
}

func mapStoryToStoryViewModel(story *data.Story, signedInUser *models.SignedInUserViewModel) *models.StoryViewModel {
	uri, _ := url.Parse(story.URL)
	var viewModel = models.StoryViewModel{
		ID:                      story.ID,
		Title:                   story.Title,
		URL:                     story.URL,
		Text:                    template.HTML(strings.ReplaceAll(story.Text, "\n", "<br />")),
		Host:                    uri.Hostname(),
		Points:                  story.UpVotes, // TODO: call point calculation function here
		UserID:                  story.UserID,
		UserName:                story.UserName,
		CommentCount:            story.CommentCount,
		IsUpvotedBySignedInUser: false,
		IsSavedBySignedInUser:   false,
		SubmittedOnText:         shared.DateToString(story.SubmittedOn),
		SignedInUser:            signedInUser,
	}

	if signedInUser != nil {
		isUpvoted, err := data.CheckIfStoryVotedByUser(signedInUser.UserID, story.ID, enums.UpVote)
		if err != nil {
			// TODO: log error here
			isUpvoted = false
		}

		isSaved, err := data.CheckIfUserSavedStory(signedInUser.UserID, story.ID)

		if err != nil {
			// TODO: log error here
			isSaved = false
		}
		viewModel.IsUpvotedBySignedInUser = isUpvoted
		viewModel.IsSavedBySignedInUser = isSaved
	}
	return &viewModel
}

func mapCommentToCommentViewModel(comment *data.Comment, signedInUser *models.SignedInUserViewModel) *models.CommentViewModel {
	model := &models.CommentViewModel{
		ID:              comment.ID,
		ParentID:        comment.ParentID,
		StoryID:         comment.StoryID,
		Comment:         comment.Comment,
		Points:          comment.UpVotes,
		UserID:          comment.UserID,
		UserName:        comment.UserName,
		CommentedOnText: shared.DateToString(comment.CommentedOn),
		SignedInUser:    signedInUser,
	}
	if comment.ParentID == data.CommentRootID {
		model.IsRoot = true
	}
	if signedInUser != nil {
		isUpvoted, err := data.CheckIfCommentUpVotedByUser(signedInUser.UserID, comment.ID)
		if err != nil {
			isUpvoted = false
		}
		model.IsUpvotedBySignedInUser = isUpvoted
	}
	return model
}

func mapCommentsToViewModelsWithChildren(comments *[]data.Comment, signedInUser *models.SignedInUserViewModel, storyID int) *[]models.CommentViewModel {
	var viewModels []models.CommentViewModel
	for _, comment := range *comments {
		viewModel := *mapCommentToCommentViewModel(&comment, signedInUser)
		childComments, err := data.GetCommentsByParentIDAndStoryID(comment.ID, storyID)
		if err != nil {
			// TODO: Log error here
			continue
		}
		viewModel.ChildComments = *mapCommentsToViewModelsWithChildren(childComments, signedInUser, storyID)
		viewModels = append(viewModels, viewModel)
	}
	return &viewModels
}
