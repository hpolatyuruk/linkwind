package controllers

import (
	"encoding/json"
	"fmt"
	"linkwind/app/data"
	"linkwind/app/enums"
	"linkwind/app/models"
	"linkwind/app/shared"
	"linkwind/app/templates"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
)

/*CommentVoteModel represents the data in http request body to upvote comment.*/
type CommentVoteModel struct {
	CommentID int
	UserID    int
	VoteType  enums.VoteType
}

/*ReplyModel represents the data to reply to comment.*/
type ReplyModel struct {
	ParentCommentID int
	StoryID         int
	ReplyText       string
}

/*AddCommentHandler adds comment to the story. */
func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Error(w, "Only POST method is supported.", http.StatusMethodNotAllowed)
		return
	}
	signedInUser := shared.GetUser(r)
	commentText := r.FormValue("comment")
	strStoryID := r.FormValue("storyID")
	storyURL := fmt.Sprintf("/stories/detail?id=%s", strStoryID)
	if strings.TrimSpace(commentText) == "" ||
		strings.TrimSpace(strStoryID) == "" {
		http.Redirect(w, r, storyURL, http.StatusSeeOther)
		return
	}
	storyID, err := strconv.Atoi(strStoryID)
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
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
		sentry.CaptureException(err)
		panic(err)
	}
	http.Redirect(w, r, storyURL, http.StatusSeeOther)
}

/*ReplyToCommentHandler write a reply to comment.*/
func ReplyToCommentHandler(w http.ResponseWriter, r *http.Request) {
	user := shared.GetUser(r)
	if r.Method == "GET" {
		http.Error(w, "Unsupported method. Only post method is supported.", http.StatusMethodNotAllowed)
		return
	}
	var model ReplyModel
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		sentry.CaptureException(err)
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
		sentry.CaptureException(err)
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
		sentry.CaptureException(err)
		http.Error(w, fmt.Sprintf("Cannot render comment template. Error: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(output))
}

/*VoteCommentHandler runs when click to upvote and downvote comment button. If it's not voted before by user, votes that comment*/
func VoteCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Error(w, "Unsupported method. Only post method is supported.", http.StatusMethodNotAllowed)
		return
	}
	var model CommentVoteModel
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		sentry.CaptureException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	voteType, err := data.GetCommentVoteByUser(model.UserID, model.CommentID)
	if err != nil {
		sentry.CaptureException(err)
		http.Error(w, fmt.Sprintf("Error occured while checking if user already upvoted. Error : %v", err), http.StatusInternalServerError)
		return
	}
	if voteType != nil {
		if model.VoteType == *voteType {
			res, _ := json.Marshal(&JSONResponse{
				Result: "AlreadyUpvoted",
			})
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(res)
			return
		}
		//Remove comment's previous vote given by user to make sure that there can be only one type of vote at a time
		err = data.RemoveCommentVote(model.UserID, model.CommentID, *voteType)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error occured while removing user's previous vote. UserID: %d, CommentID: %d, VoteType: %d,  Error : %v", model.UserID, model.CommentID, *voteType, err), http.StatusInternalServerError)
			return
		}
	}
	err = data.VoteComment(model.UserID, model.CommentID, model.VoteType)
	if err != nil {
		sentry.CaptureException(err)
		http.Error(w, fmt.Sprintf("Error occured while upvoting story. Error : %v", err), http.StatusInternalServerError)
		return
	}
	res, _ := json.Marshal(&JSONResponse{
		Result: "Voted",
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

/*RemoveCommentVoteHandler handles unvote button. If a comment voted by user before, this handler undo that operation*/
func RemoveCommentVoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Error(w, "Unsupported request method. Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}
	var model CommentVoteModel
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		sentry.CaptureException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	voteType, err := data.GetCommentVoteByUser(model.UserID, model.CommentID)
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

	err = data.RemoveCommentVote(model.UserID, model.CommentID, model.VoteType)
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
