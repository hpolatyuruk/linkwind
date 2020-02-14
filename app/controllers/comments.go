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

/*UpvoteCommentHandler upvotes a comment if not upvoted by same user*/
func UpvoteCommentHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return fmt.Errorf("Reqeuest should be POSTm request")
	}

	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("Error occured when parse from. Error : %v", err)
	}

	userIDStr := r.FormValue("userID")
	commentIDStr := r.FormValue("commentID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return fmt.Errorf("Error occured when convert string userID to int userID. Error : %v", err)
	}
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		return fmt.Errorf("Error occured when convert string commentID to int commentID. Error : %v", err)
	}

	isUpvoted, err := data.CheckIfCommentUpVotedByUser(userID, commentID)
	if err != nil {
		return fmt.Errorf("Error occured when check user upvote comment. Error : %v", err)
	}

	if isUpvoted {
		return nil
	}

	err = data.UpVoteComment(userID, commentID)
	if err != nil {
		return fmt.Errorf("Error occured when upvote comment. Error : %v", err)
	}
	return nil
}

/*UnvoteCommentHandler unvotes a comment if upvoted by same user*/
func UnvoteCommentHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return fmt.Errorf("Reqeuest should be POSTm request")
	}

	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("Error occured when parse from. Error : %v", err)
	}

	userIDStr := r.FormValue("userID")
	commentIDStr := r.FormValue("commentID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return fmt.Errorf("Error occured when convert string userID to int userID. Error : %v", err)
	}
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		return fmt.Errorf("Error occured when convert string commentID to int commentID. Error : %v", err)
	}

	isUpvoted, err := data.CheckIfCommentUpVotedByUser(userID, commentID)
	if err != nil {
		return fmt.Errorf("Error occured when check user upvote comment. Error : %v", err)
	}

	if isUpvoted == false {
		return nil
	}

	err = data.UnVoteComment(userID, commentID)
	if err != nil {
		return fmt.Errorf("Error occured when unvote comment. Error : %v", err)
	}
	return nil
}
