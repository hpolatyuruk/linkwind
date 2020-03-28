package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/lib/pq"
)

/*DBError represents the database error*/
type DBError struct {
	Message       string
	OriginalError error
}

func (err *DBError) Error() string {
	return fmt.Sprintf("%s | OriginalError: %v", err.Message, err.OriginalError)
}

func connectionString() (conStr string) {
	host := os.Getenv("DBHost")
	user := os.Getenv("DBUser")
	password := os.Getenv("DBPassword")
	dbName := os.Getenv("DBName")
	port := os.Getenv("DBPort")

	conStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	return conStr
}

func connectToDB() (*sql.DB, error) {
	connStr := connectionString()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, &DBError{"Cannot open db connection!", err}
	}
	err = db.Ping()
	if err != nil {
		return nil, &DBError{"Cannot ping db!", err}
	}
	return db, err
}

/*MapSQLRowToUser creates an user struct object by sql row*/
func MapSQLRowToUser(row *sql.Row) (user *User, err error) {
	var _user User
	err = row.Scan(
		&_user.ID,
		&_user.UserName,
		&_user.FullName,
		&_user.Email,
		&_user.RegisteredOn,
		&_user.Password,
		&_user.Website,
		&_user.About,
		&_user.InviteCode,
		&_user.Karma,
		&_user.CustomerID)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot map sql row to user struct"), err}
	}
	user = &_user
	return user, nil
}

/*MapSQLRowToCustomer creates an user struct object by sql row*/
func MapSQLRowToCustomer(row *sql.Row) (customer *Customer, err error) {
	var _customer Customer
	var domain sql.NullString
	err = row.Scan(
		&_customer.ID,
		&_customer.Name,
		&_customer.Email,
		&_customer.RegisteredOn,
		&domain,
		&_customer.LogoImage)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot map sql row to customer struct"), err}
	}
	if domain.Valid {
		_customer.Domain = string(domain.String)
	} else {
		_customer.Domain = CustomerDefaultDomain
	}
	customer = &_customer
	return customer, nil
}

/*MapSQLRowToInviteCodeInfo creates an invite code info struct object by sql row*/
func MapSQLRowToInviteCodeInfo(row *sql.Row) (inviteCodeInfo *InviteCodeInfo, err error) {
	var _inviteCodeInfo InviteCodeInfo
	err = row.Scan(
		&_inviteCodeInfo.Code,
		&_inviteCodeInfo.InviterUserID,
		&_inviteCodeInfo.InvitedEmailAddress,
		&_inviteCodeInfo.Used,
		&_inviteCodeInfo.CreatedOn)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot map sql row to invite code info struct"), err}
	}
	inviteCodeInfo = &_inviteCodeInfo
	return inviteCodeInfo, nil
}

/*MapSQLRowToStory creates a story struct by sql rows*/
func MapSQLRowToStory(rows *sql.Row) (story *Story, err error) {
	var _story Story
	var username string
	err = rows.Scan(
		&_story.ID,
		&_story.URL,
		&_story.Title,
		&_story.Text,
		&_story.UpVotes,
		&_story.CommentCount,
		&_story.UserID,
		&_story.SubmittedOn,
		pq.Array(&_story.Tags),
		&_story.DownVotes,
		&username)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot read rows"), err}
	}
	_story.UserName = username
	story = &_story
	return story, nil
}

/*MapSQLRowsToStories creates a story struct array by sql rows*/
func MapSQLRowsToStories(rows *sql.Rows) (stories *[]Story, err error) {
	_stories := []Story{}
	var username string
	for rows.Next() {
		var story Story
		err = rows.Scan(
			&story.ID,
			&story.URL,
			&story.Title,
			&story.Text,
			&story.UpVotes,
			&story.CommentCount,
			&story.UserID,
			&story.SubmittedOn,
			pq.Array(&story.Tags),
			&story.DownVotes,
			&username)
		if err != nil {
			return nil, &DBError{fmt.Sprintf("Cannot read rows"), err}
		}
		story.UserName = username
		_stories = append(_stories, story)
	}
	return &_stories, nil

}

/*MapSQLRowsToComments creates a comment struct array by sql rows*/
func MapSQLRowsToComments(rows *sql.Rows) (comments *[]Comment, err error) {
	_comments := []Comment{}
	for rows.Next() {
		comment := Comment{}
		var parentID sql.NullInt32
		err = rows.Scan(
			&comment.Comment,
			&comment.UpVotes,
			&comment.StoryID,
			&parentID,
			&comment.ReplyCount,
			&comment.UserID,
			&comment.CommentedOn,
			&comment.ID,
			&comment.DownVotes,
			&comment.UserName)
		if err != nil {
			return nil, &DBError{"Cannot read comment row.", err}
		}
		if parentID.Valid {
			comment.ParentID = int(parentID.Int32)
		} else {
			comment.ParentID = CommentRootID
		}
		_comments = append(_comments, comment)
	}
	return &_comments, nil
}

/*MapSQLRowsToUsers creates a user struct array by sql rows*/
func MapSQLRowsToUsers(rows *sql.Rows) (users *[]User, err error) {
	_users := []User{}
	for rows.Next() {
		user := User{}
		err = rows.Scan(
			&user.ID,
			&user.UserName,
			&user.FullName,
			&user.Email,
			&user.RegisteredOn,
			&user.Password,
			&user.Website,
			&user.About,
			&user.InviteCode,
			&user.Karma)
		if err != nil {
			return nil, &DBError{"Cannot read comment row.", err}
		}

		_users = append(_users, user)
	}
	return &_users, nil
}

/*MapSQLRowsToReplies creates a reply struct array by sql rows*/
func MapSQLRowsToReplies(rows *sql.Rows) (replies *[]Reply, err error) {
	_replies := []Reply{}
	for rows.Next() {
		reply := Reply{}
		comment := Comment{}
		var storyID int
		var storyTitle string
		var userName string
		var parentID sql.NullInt32
		err = rows.Scan(
			&comment.Comment,
			&comment.UpVotes,
			&comment.DownVotes,
			&comment.StoryID,
			&parentID,
			&comment.ReplyCount,
			&comment.UserID,
			&comment.CommentedOn,
			&comment.ID,
			&storyTitle,
			&storyID,
			&userName)
		if err != nil {
			return nil, &DBError{"Cannot read comment row.", err}
		}
		if parentID.Valid {
			comment.ParentID = int(parentID.Int32)
		} else {
			comment.ParentID = CommentRootID
		}
		reply.Comment = &comment
		reply.StoryID = storyID
		reply.StoryTitle = storyTitle
		reply.UserName = userName
		_replies = append(_replies, reply)
	}
	return &_replies, nil
}

/*CalculateKarma calculates user's karma by its upvotes and downvotes*/
// TODO : pageNumber ve pageRowCount 0 oalrak verdim. CalculateKarma bu parameterleri dışarıdan mı almalı, yoksa içeride başka bir değer mi verilmeli?
func CalculateKarma(userID int) (int, error) {
	pageNumber := 0
	pageRowCount := 0
	stories, err := GetUserSubmittedStories(userID, pageNumber, pageRowCount)
	if err != nil {
		log.Printf("An error occurred while calculating user's karma error: %s", err)
		return 0, err
	}
	sVotes := 0
	for _, s := range *stories {
		sVotes += (s.UpVotes - s.DownVotes)
	}

	comments, err := GetUserCommentsNotPaging(userID)
	cVotes := 0
	for _, c := range *comments {
		cVotes += (c.UpVotes - c.DownVotes)
	}

	return sVotes + cVotes, nil
}
