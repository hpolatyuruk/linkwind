package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"turkdev/app/models"
	"turkdev/app/src/templates"
	"turkdev/data"
)

/*RepliesHandler handles showing user's replies*/
func RepliesHandler(w http.ResponseWriter, r *http.Request) error {
	title := "Replies | Turk Dev"
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

	templates.Render(
		w,
		"comments/index.html",
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
