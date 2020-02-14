package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"turkdev/app/models"
	"turkdev/app/src/templates"
	"turkdev/data"
	"turkdev/shared"
)

/*CommentVoteModel represents the data in http request body to upvote comment.*/
type CommentVoteModel struct {
	CommentID int
	UserID    int
}

/*ReplyModel represents the data to reply to comment.*/
type ReplyModel struct {
	ParentCommentID int
	StoryID         int
	ReplyText       string
}

/*AddCommentHandler adds comment to the story. */
func AddCommentHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		http.Error(w, "Only POST method is supported.", http.StatusMethodNotAllowed)
		return nil
	}
	isAuthenticated, signedInUser, err := shared.IsAuthenticated(r)
	if err != nil {
		return err
	}
	if isAuthenticated == false {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return nil
	}
	commentText := r.FormValue("comment")
	strStoryID := r.FormValue("storyID")
	storyURL := fmt.Sprintf("/stories/detail?id=%s", strStoryID)
	if strings.TrimSpace(commentText) == "" ||
		strings.TrimSpace(strStoryID) == "" {
		http.Redirect(w, r, storyURL, http.StatusSeeOther)
		return nil
	}
	storyID, err := strconv.Atoi(strStoryID)
	if err != nil {
		return err
	}
	comment := &data.Comment{
		StoryID:     storyID,
		UserID:      signedInUser.ID,
		UserName:    signedInUser.UserName,
		ParentID:    data.CommentRootID,
		UpVotes:     0,
		DownVotes:   0,
		ReplyCount:  0,
		Comment:     commentText,
		CommentedOn: time.Now(),
	}
	_, err = data.WriteComment(comment)
	if err != nil {
		return err
	}
	http.Redirect(w, r, storyURL, http.StatusSeeOther)
	return nil
}

/*ReplyToCommentHandler write a reply to comment.*/
func ReplyToCommentHandler(w http.ResponseWriter, r *http.Request) {
	isAuth, user, err := shared.IsAuthenticated(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("An error occurred. Error:%v", err), http.StatusInternalServerError)
		return
	}
	if isAuth == false {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if r.Method == "GET" {
		http.Error(w, "Unsupported method. Only post method is supported.", http.StatusMethodNotAllowed)
		return
	}
	var model ReplyModel
	err = json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		http.Error(w, "Cannot parse json.", http.StatusBadRequest)
		return
	}
	comment := &data.Comment{
		UserID:      user.ID,
		UserName:    user.UserName,
		StoryID:     model.StoryID,
		ParentID:    model.ParentCommentID,
		Comment:     model.ReplyText,
		CommentedOn: time.Now(),
		UpVotes:     0,
		DownVotes:   0,
		ReplyCount:  0,
	}
	commentID, err := data.WriteComment(comment)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong :(", http.StatusInternalServerError)
		return
	}
	comment.ID = *commentID
	signedUserModel := &models.SignedInUserViewModel{
		UserID:     user.ID,
		UserName:   user.UserName,
		CustomerID: user.CustomerID,
		Email:      user.Email,
	}
	output, err := templates.RenderAsString("partials/comment.html", "comment",
		mapCommentToCommentViewModel(comment, signedUserModel))
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot render comment template. Error: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(output))
}

/*RepliesHandler handles showing user's replies*/
func RepliesHandler(w http.ResponseWriter, r *http.Request) error {
	title := "Repli,es | Turk Dev"
	user := models.User{"Anil Yuzener"}

	strUserID := r.URL.Query().Get("userid")
	if len(strUserID) == 0 {
		// TODO(Anil): User cannot be found. Show appropriate message here
		return nil
	}
	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		// TODO(Anil): Cannot parse to int. Show user not found message.
		return nil
	}
	replies, err := data.GetUserReplies(userID)
	if err != nil {
		// TODO(Anil): show error page here
	}

	data := map[string]interface{}{
		"Content": "Replies",
		"Replies": replies,
	}

	templates.RenderInLayout(
		w,
		"index.html",
		models.ViewModel{
			title,
			user,
			data,
		},
	)
	return nil
}

/*UpvoteCommentHandler runs when click to upvote comment button. If not upvoted before by user, upvotes that comment*/
func UpvoteCommentHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated, _, err := shared.IsAuthenticated(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if isAuthenticated == false {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if r.Method == "GET" {
		http.Error(w, "Unsupported method. Only post method is supported.", http.StatusMethodNotAllowed)
		return
	}

	var model CommentVoteModel
	err = json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isUpvoted, err := data.CheckIfCommentUpVotedByUser(model.UserID, model.CommentID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error occured while checking if user already upvoted. Error : %v", err), http.StatusInternalServerError)
		return
	}
	if isUpvoted {
		res, _ := json.Marshal(&JSONResponse{
			Result: "AlreadyUpvoted",
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	}

	err = data.UpVoteComment(model.UserID, model.CommentID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, fmt.Sprintf("Error occured while upvoting story. Error : %v", err), http.StatusInternalServerError)
		return
	}
	res, _ := json.Marshal(&JSONResponse{
		Result: "Upvoted",
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

/*UnvoteCommentHandler handles unvote button. If a comment voted by user before, this handler undo that operation*/
func UnvoteCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Error(w, "Unsupported request method. Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}

	var model CommentVoteModel
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isUpvoted, err := data.CheckIfCommentUpVotedByUser(model.UserID, model.CommentID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error occured while checking user story vote. Error : %v", err), http.StatusInternalServerError)
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

	err = data.UnVoteComment(model.UserID, model.CommentID)
	if err != nil {
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
