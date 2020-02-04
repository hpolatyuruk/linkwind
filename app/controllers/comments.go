package controllers

import (
	"net/http"
	"strconv"
	"turkdev/app/models"
	"turkdev/app/templates"
	"turkdev/data"
)

/*CommentsHandler handles showing comments by giving story id*/
func CommentsHandler(w http.ResponseWriter, r *http.Request) error {
	title := "Comments | Turk Dev"
	user := models.User{"Anil Yuzener"}

	strStoryID := r.URL.Query().Get("storyid")
	if len(strStoryID) == 0 {
		// TODO(Anil): Story cannot be found. Show appropriate message here.
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
