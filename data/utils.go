package data

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/lib/pq"
)

var (
	DBHost        = "localhost"
	DBPort        = "5432"
	DBUser        = "postgres"
	DBPassword    = "3842"
	DBName        = "postgres"
	JWTPrivateKey = "jwtprivatekeyfordebug"
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
	port, err := strconv.Atoi(DBPort)
	if err != nil {
		log.Fatalf("Cannot parse db port. Error: %v", err)
	}

	conStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DBHost, port, DBUser, DBPassword, DBName)

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

/*MapSQLRowToCustomer creates an user struct object by sql row*/
func MapSQLRowToCustomer(row *sql.Row) (customer *Customer, err error) {
	var _customer Customer
	err = row.Scan(
		&_customer.ID,
		&_customer.Name,
		&_customer.Email,
		&_customer.Description,
		&_customer.RegisteredOn,
		&_customer.Domain)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot map sql row to customer struct"), err}
	}
	customer = &_customer
	return customer, nil
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
			&user.InvitedBy,
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
