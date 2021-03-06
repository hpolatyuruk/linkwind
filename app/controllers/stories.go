package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"linkwind/app/data"
	"linkwind/app/enums"
	"linkwind/app/models"
	"linkwind/app/shared"
	"linkwind/app/templates"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
)

const (
	/*DefaultPageSize represents story count to be listed per page*/
	DefaultPageSize = 15
	/*MinKarmaToDownVote represents minimum number of karma a user needs to downvote a story or comment*/
	MinKarmaToDownVote = 100
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

/*JSONResponse respresents the json response.*/
type JSONResponse struct {
	Result string
}

type getStoriesPaged func(customerID, pageNo, storyCountPerPage int) (*[]data.Story, error)

/*StoriesHandler handles showing the popular published stories*/
func StoriesHandler(w http.ResponseWriter, r *http.Request) {
	renderStoriesPage("Stories", data.GetStories, w, r)
}

/*RecentStoriesHandler handles showing recently published stories*/
func RecentStoriesHandler(w http.ResponseWriter, r *http.Request) {
	renderStoriesPage("Recent Stories", data.GetRecentStories, w, r)
}

func renderStoriesPage(title string, fnGetStories getStoriesPaged, w http.ResponseWriter, r *http.Request) {
	var model = &models.StoryPageViewModel{Title: title}
	customerCtx := shared.GetCustomerFromContext(r)
	user := shared.GetUserFromContext(r)

	var page int = getPage(r)
	stories, err := fnGetStories(customerCtx.ID, page, DefaultPageSize)
	if err != nil {
		panic(err)
	}

	storiesCount, err := data.GetCustomerStoriesCount(customerCtx.ID)
	if err != nil {
		panic(err)
	}

	pagingModel, err := setPagingViewModel(customerCtx.ID, page, storiesCount)
	if err != nil {
		panic(err)
	}
	model.Page = pagingModel
	if stories != nil && len(*stories) > 0 {
		model.Stories = *mapStoriesToStoryViewModel(stories, user)
	}
	templates.RenderInLayout(w, r, "stories.html", model)
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
func UserSavedStoriesHandler(w http.ResponseWriter, r *http.Request) {
	var model = &models.StoryPageViewModel{
		Title: "Saved Stories"}

	user := shared.GetUserFromContext(r)

	var page int = getPage(r)
	stories, err := data.GetUserSavedStories(user.ID, page, DefaultPageSize)
	if err != nil {
		panic(err)
	}
	storiesCount, err := data.GetUserSavedStoriesCount(user.ID)
	if err != nil {
		panic(err)
	}
	pagingModel, err := setPagingViewModel(user.CustomerID, page, storiesCount)
	if err != nil {
		panic(err)
	}
	model.Page = pagingModel
	if stories == nil || len(*stories) > 0 {
		model.Stories = *mapStoriesToStoryViewModel(stories, user)
	}
	err = templates.RenderInLayout(w, r, "stories.html", model)
	if err != nil {
		panic(err)
	}
}

/*UserSubmittedStoriesHandler handles user's submitted stories*/
func UserSubmittedStoriesHandler(w http.ResponseWriter, r *http.Request) {
	var model = &models.StoryPageViewModel{
		Title: "Submitted Stories",
	}
	user := shared.GetUserFromContext(r)
	userID := user.ID
	strUserID := r.URL.Query().Get("userid")
	if strings.TrimSpace(strUserID) != "" {
		userID, _ = strconv.Atoi(strUserID)
	}
	var page int = getPage(r)
	stories, err := data.GetUserSubmittedStories(userID, page, DefaultPageSize)
	if err != nil {
		panic(err)
	}
	storiesCount, err := data.GetUserSubmittedStoriesCount(userID)
	if err != nil {
		panic(err)
	}
	pagingModel, err := setPagingViewModel(user.CustomerID, page, storiesCount)
	if err != nil {
		panic(err)
	}
	model.Page = pagingModel
	if stories == nil || len(*stories) > 0 {
		model.Stories = *mapStoriesToStoryViewModel(stories, user)
	}
	err = templates.RenderInLayout(w, r, "stories.html", model)
	if err != nil {
		panic(err)
	}
}

/*UserUpvotedStoriesHandler handles showing the upvoted stories by user*/
func UserUpvotedStoriesHandler(w http.ResponseWriter, r *http.Request) {
	var model = &models.StoryPageViewModel{
		Title: "Upvoted Stories",
	}

	user := shared.GetUserFromContext(r)

	var page int = getPage(r)
	stories, err := data.GetUserUpvotedStories(user.ID, page, DefaultPageSize)
	if err != nil {
		panic(err)
	}
	storiesCount, err := data.GetUserUpvotedStoriesCount(user.ID)
	if err != nil {
		panic(err)
	}
	pagingModel, err := setPagingViewModel(user.CustomerID, page, storiesCount)
	if err != nil {
		panic(err)
	}
	model.Page = pagingModel
	if stories == nil || len(*stories) > 0 {
		model.Stories = *mapStoriesToStoryViewModel(stories, user)
	}
	err = templates.RenderInLayout(w, r, "stories.html", model)
	if err != nil {
		panic(err)
	}
}

/*SubmitStoryHandler handles to submit a new story*/
func SubmitStoryHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handlesSubmitGET(w, r)
	case "POST":
		handleSubmitPOST(w, r)
	default:
		handlesSubmitGET(w, r)
	}
}

func handlesSubmitGET(w http.ResponseWriter, r *http.Request) {
	model := &models.StorySubmitModel{}
	templates.RenderInLayout(w, r, "submit.html", model)
}

func handleSubmitPOST(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	model := &models.StorySubmitModel{
		URL:   r.FormValue("url"),
		Title: r.FormValue("title"),
		Text:  r.FormValue("text"),
	}
	if model.Validate() == false {
		templates.RenderInLayout(w, r, "submit.html", model)
		return
	}
	if strings.TrimSpace(model.URL) != "" {
		fetchedTitle, err := shared.FetchURL(model.URL)
		if err != nil {
			fmt.Println(err)
			model.Errors["URL"] = "Something went wrong while fetching URL. Please make sure that you entered a valid URL."
			templates.RenderInLayout(w, r, "submit.html", model)
			return
		}

		if strings.TrimSpace(model.Title) == "" {
			model.Title = fetchedTitle
		}
	}

	user := shared.GetUserFromContext(r)
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
		panic(err)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

/*StoryDetailHandler handles showing comments by giving story id*/
func StoryDetailHandler(w http.ResponseWriter, r *http.Request) {
	strStoryID := r.URL.Query().Get("id")
	if len(strStoryID) == 0 {
		err := templates.RenderFile(w, "errors/404.html", nil)
		if err != nil {
			panic(err)
		}
		return
	}
	storyID, err := strconv.Atoi(strStoryID)
	if err != nil {
		panic(fmt.Errorf("Cannot convert string StoryID to int. Original err : %v", err))
	}
	story, err := data.GetStoryByID(storyID)
	if err != nil {
		panic(fmt.Errorf("Cannot get story from db (StoryID : %d). Original err : %v", storyID, err))
	}
	if story == nil {
		err = templates.RenderFile(w, "errors/404.html", nil)
		if err != nil {
			panic(err)
		}
		return
	}
	comments, err := data.GetRootCommentsByStoryID(storyID)
	if err != nil {
		panic(fmt.Errorf("Cannot get comments from db (StoryID : %d). Original err : %v", storyID, err))
	}
	model := &models.StoryDetailPageViewModel{
		Title: story.Title,
	}
	user := shared.GetUserFromContext(r)
	model.Story = mapStoryToStoryViewModel(story, user)
	model.Comments = mapCommentsToViewModelsWithChildren(comments, user, storyID)

	templates.RenderInLayout(w, r, "detail.html", model)
}

func mapStoriesToStoryViewModel(stories *[]data.Story, userClaims *shared.SignedInUserClaims) *[]models.StoryViewModel {
	var viewModels []models.StoryViewModel

	for _, story := range *stories {
		viewModel := mapStoryToStoryViewModel(&story, userClaims)
		viewModels = append(viewModels, *viewModel)
	}
	return &viewModels
}

func mapStoryToStoryViewModel(story *data.Story, userClaims *shared.SignedInUserClaims) *models.StoryViewModel {
	uri, _ := url.Parse(story.URL)
	var viewModel = models.StoryViewModel{
		ID:              story.ID,
		Title:           story.Title,
		URL:             story.URL,
		Text:            template.HTML(strings.ReplaceAll(story.Text, "\n", "<br />")),
		Host:            uri.Hostname(),
		Points:          story.UpVotes,
		UserID:          story.UserID,
		UserName:        story.UserName,
		CommentCount:    story.CommentCount,
		IsUpvoted:       false,
		IsDownvoted:     false,
		IsSaved:         false,
		ShowDownvoteBtn: false,
		SubmittedOnText: shared.DateToString(story.SubmittedOn),
	}

	if userClaims != nil {
		viewModel.SignedInUser = mapUserClaimsToSignedUserViewModel(userClaims)
		viewModel.ShowDownvoteBtn = userClaims.Karma > MinKarmaToDownVote
		voteType, err := data.GetStoryVoteByUser(userClaims.ID, story.ID)
		if err != nil {
			sentry.CaptureException(err)
		}
		if voteType != nil {
			if *voteType == enums.UpVote {
				viewModel.IsUpvoted = true
			} else if *voteType == enums.DownVote {
				viewModel.IsDownvoted = true
			}
		}

		isSaved, err := data.CheckIfUserSavedStory(userClaims.ID, story.ID)

		if err != nil {
			sentry.CaptureException(err)
			isSaved = false
		}
		viewModel.IsSaved = isSaved
	}
	return &viewModel
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
		sentry.CaptureException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	voteType, err := data.GetStoryVoteByUser(model.UserID, model.StoryID)
	if err != nil {
		sentry.CaptureMessage(fmt.Sprintf("Error occured while getting user's current vote. UserID: %d, StoryID: %d, VoteType: %d,  Error : %v", model.UserID, model.StoryID, model.VoteType, err))
		http.Error(w, fmt.Sprintf("Error occured while getting user's current vote. UserID: %d, StoryID: %d, VoteType: %d,  Error : %v", model.UserID, model.StoryID, model.VoteType, err), http.StatusInternalServerError)
		return
	}
	if voteType != nil {
		if model.VoteType == *voteType {
			res, _ := json.Marshal(&JSONResponse{
				Result: "AlreadyVoted",
			})
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(res)
			return
		}
		// Remove story's previous vote given by user to make sure that there can be only one type of vote at a time
		err = data.RemoveStoryVote(model.UserID, model.StoryID, *voteType)
		if err != nil {
			sentry.CaptureMessage(fmt.Sprintf("Error occured while removing user's previous vote. UserID: %d, StoryID: %d, VoteType: %d,  Error : %v", model.UserID, model.StoryID, *voteType, err))
			http.Error(w, fmt.Sprintf("Error occured while removing user's previous vote. UserID: %d, StoryID: %d, VoteType: %d,  Error : %v", model.UserID, model.StoryID, *voteType, err), http.StatusInternalServerError)
			return
		}
	}
	err = data.VoteStory(model.UserID, model.StoryID, model.VoteType)
	if err != nil {
		sentry.CaptureMessage(fmt.Sprintf("Error occured while voting story. UserID: %d, StoryID: %d, VoteType: %d, Error : %v", model.UserID, model.StoryID, model.VoteType, err))
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
		sentry.CaptureException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	voteType, err := data.GetStoryVoteByUser(model.UserID, model.StoryID)
	if err != nil {
		sentry.CaptureException(err)
		http.Error(w, fmt.Sprintf("Error occured while checking user story vote. Error : %v", err), http.StatusInternalServerError)
		return
	}
	if voteType == nil {
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
		sentry.CaptureException(err)
		http.Error(w, fmt.Sprintf("Error occured while unvoting story. Error : %v", err), http.StatusInternalServerError)
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
		sentry.CaptureException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isSaved, err := data.CheckIfUserSavedStory(model.UserID, model.StoryID)
	if err != nil {
		sentry.CaptureException(err)
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
		sentry.CaptureException(err)
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
		sentry.CaptureException(err)
		http.Error(w, "An error occured while parsing json.", http.StatusBadRequest)
		return
	}

	isSaved, err := data.CheckIfUserSavedStory(model.UserID, model.StoryID)
	if err != nil {
		sentry.CaptureException(err)
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
		sentry.CaptureException(err)
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

func mapCommentToCommentViewModel(comment *data.Comment, userClaims *shared.SignedInUserClaims) *models.CommentViewModel {
	model := &models.CommentViewModel{
		ID:              comment.ID,
		ParentID:        comment.ParentID,
		StoryID:         comment.StoryID,
		Comment:         comment.Comment,
		Points:          comment.UpVotes,
		UserID:          comment.UserID,
		UserName:        comment.UserName,
		CommentedOnText: shared.DateToString(comment.CommentedOn),
	}
	if comment.ParentID == data.CommentRootID {
		model.IsRoot = true
	}
	if userClaims != nil {
		model.SignedInUser = mapUserClaimsToSignedUserViewModel(userClaims)
		model.ShowDownvoteBtn = userClaims.Karma > MinKarmaToDownVote
		voteType, err := data.GetCommentVoteByUser(userClaims.ID, comment.ID)
		if err != nil {
			sentry.CaptureException(err)
		}
		if voteType != nil {
			if *voteType == enums.UpVote {
				model.IsUpvoted = true
			}
			if *voteType == enums.DownVote {
				model.IsDownvoted = true
			}
		}
	}
	return model
}

func mapCommentsToViewModelsWithChildren(comments *[]data.Comment, userClaims *shared.SignedInUserClaims, storyID int) *[]models.CommentViewModel {
	var viewModels []models.CommentViewModel
	for _, comment := range *comments {
		viewModel := *mapCommentToCommentViewModel(&comment, userClaims)
		childComments, err := data.GetCommentsByParentIDAndStoryID(comment.ID, storyID)
		if err != nil {
			sentry.CaptureException(err)
			continue
		}
		viewModel.ChildComments = *mapCommentsToViewModelsWithChildren(childComments, userClaims, storyID)
		viewModels = append(viewModels, viewModel)
	}
	return &viewModels
}

func mapUserClaimsToSignedUserViewModel(signedInUserClaims *shared.SignedInUserClaims) *models.SignedInUserViewModel {
	return &models.SignedInUserViewModel{
		UserID:     signedInUserClaims.ID,
		CustomerID: signedInUserClaims.CustomerID,
		Email:      signedInUserClaims.Email,
		UserName:   signedInUserClaims.UserName,
		Karma:      signedInUserClaims.Karma,
	}
}
