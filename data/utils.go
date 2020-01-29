package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
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

func connectionString() (conStr string, err error) {
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading staging.env file. Error: %v", err)
	}

	host := os.Getenv("DBHost")
	port, err := strconv.Atoi(os.Getenv("DBPort"))
	if err != nil {
		log.Fatalf("Cannot parse db port. Error: %v", err)
	}
	user := os.Getenv("DBUser")
	password := os.Getenv("DBPassword")
	name := os.Getenv("DBName")

	conStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, name)

	return conStr, err
}

func connectToDB() (db *sql.DB, err error) {
	connStr, err := connectionString()
	if err != nil {
		return nil, &DBError{"Cannot read connection string!", err}
	}
	db, err = sql.Open("postgres", connStr)
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
		&_user.InvitedBy,
		&_user.InviteCode,
		&_user.Karma)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot map sql row to user struct"), err}
	}
	user = &_user
	return user, nil
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
			&story.DownVotes,
			&story.CommentCount,
			&story.UserID,
			&story.SubmittedOn,
			pq.Array(&story.Tags),
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
		var userName string
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
			&userName)
		if err != nil {
			return nil, &DBError{"Cannot read comment row.", err}
		}
		if parentID.Valid {
			comment.ParentID = int(parentID.Int32)
		} else {
			comment.ParentID = CommentRootID
		}
		comment.UserName = userName
		_comments = append(_comments, comment)
	}
	return &_comments, nil
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
func CalculateKarma(userID int) (int, error) {
	stories, err := GetUserStoriesNotPaging(userID)
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
