package data

import (
	"fmt"
	"time"

	"github.com/lib/pq"
)

/*Story represents the story which contains shared article or link info*/
type Story struct {
	ID           int
	URL          string
	Title        string
	Text         string
	Tags         []string
	UpVotes      int
	CommentCount int
	UserID       int
	UserName     string
	SubmittedOn  time.Time
}

/*StoryError represents any error related to story*/
type StoryError struct {
	Message       string
	Story         *Story
	OriginalError error
}

func (err *StoryError) Error() string {
	return fmt.Sprintf(
		"%s | Story: %v | OriginalError: %v",
		err.Message,
		err.Story,
		err.OriginalError)
}

/*CreateStory creates a story on database*/
func CreateStory(story *Story) error {
	db, err := connectToDB()
	defer db.Close()
	if err != nil {
		return err
	}
	sql := "INSERT INTO stories (url, title, text, tags, upvotes, commentcount, userid, submittedon) VALUES($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err = db.Exec(
		sql,
		story.URL,
		story.Title,
		story.Text,
		pq.Array(story.Tags),
		0,
		0,
		story.UserID,
		time.Now())
	if err != nil {
		return &StoryError{"Cannot create story!", story, err}
	}
	return nil
}

/*GetStories retunrs story list by provided paging parameters*/
func GetStories(pageNumber int, pageRowCount int) (*[]Story, error) {
	db, err := connectToDB()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	// TODO(Huseyin): Sort it by point algorithim when sedat finishes it
	sql := "SELECT stories.*, users.UserName FROM stories INNER JOIN users ON users.id = stories.userid LIMIT $1 OFFSET $2"
	rows, err := db.Query(sql, pageRowCount, pageNumber*pageRowCount)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot get stories. PageNumber: %d, PageRowCount: %d", pageNumber, pageRowCount), err}
	}
	stories, err := MapSQLRowsToStories(rows)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot read rows. PageNumber: %d, PageRowCount: %d", pageNumber, pageRowCount), err}
	}
	return stories, nil
}

/*UpVoteStory increases votes for story on database*/
func UpVoteStory(userID int, storyID int) error {
	db, err := connectToDB()
	defer db.Close()
	if err != nil {
		return &DBError{fmt.Sprintf("DB connection error. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	tran, err := db.Begin()
	if err != nil {
		return &DBError{fmt.Sprintf("Cannot begin transaction. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	sql := "INSERT INTO storyvotes(storyid, userid) VALUES($1, $2)"
	_, err = tran.Exec(sql, storyID, userID)
	if err != nil {
		tran.Rollback()
		return &DBError{fmt.Sprintf("Error occurred while inserting storyvotes. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	sql = "UPDATE stories SET upvotes = upvotes + 1 WHERE id = $1"
	_, err = tran.Exec(sql, storyID)
	if err != nil {
		tran.Rollback()
		return &DBError{fmt.Sprintf("Error occurred while increasing story upvotes. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	err = tran.Commit()
	if err != nil {
		return &DBError{fmt.Sprintf("Cannot commit transaction. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	return nil
}

/*UnVoteStory unvotes the story on database*/
func UnVoteStory(userID int, storyID int) error {
	db, err := connectToDB()
	defer db.Close()
	if err != nil {
		return &DBError{fmt.Sprintf("DB connection error. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	tran, err := db.Begin()
	if err != nil {
		return &DBError{fmt.Sprintf("Cannot begin transaction. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	sql := "DELETE FROM storyvotes WHERE userid = $1 AND storyid = $2"
	_, err = tran.Exec(sql, userID, storyID)
	if err != nil {
		tran.Rollback()
		return &DBError{fmt.Sprintf("Cannot delete story vote. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	sql = "UPDATE stories SET upvotes = upvotes - 1 WHERE id = $1"
	_, err = tran.Exec(sql, storyID)
	if err != nil {
		tran.Rollback()
		return &DBError{fmt.Sprintf("Cannot update story's vote. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	err = tran.Commit()
	if err != nil {
		return &DBError{fmt.Sprintf("Cannot commit transaction. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	return nil
}

/*CheckIfStoryUpVotedByUser check if user already upvoted to given story*/
func CheckIfStoryUpVotedByUser(userID int, storyID int) (bool, error) {
	db, err := connectToDB()
	defer db.Close()
	if err != nil {
		return false, &DBError{fmt.Sprintf("DB connection error. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	sql := "SELECT COUNT(*) as count FROM storyvotes WHERE userid = $1 and storyid = $2"
	row := db.QueryRow(sql, userID, storyID)
	var count int = 0
	err = row.Scan(&count)
	if err != nil {
		return false, &DBError{fmt.Sprintf("Cannot read db row. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

/*SaveStory saves the given story to user's favorites*/
func SaveStory(userID int, storyID int) error {
	db, err := connectToDB()
	defer db.Close()
	if err != nil {
		return &DBError{fmt.Sprintf("DB connection error. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	sql := "INSERT INTO saved (userid, storyid) VALUES ($1, $2)"
	_, err = db.Exec(sql, userID, storyID)
	if err != nil {
		return &DBError{fmt.Sprintf("Cannot save story to user's favotires. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	return nil
}

/*UnSaveStory removes the given story from user's favorites*/
func UnSaveStory(userID int, storyID int) error {
	db, err := connectToDB()
	defer db.Close()
	if err != nil {
		return &DBError{fmt.Sprintf("DB connection error. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	sql := "DELETE FROM saved WHERE userid = $1 AND storyid = $2"
	_, err = db.Exec(sql, userID, storyID)
	if err != nil {
		return &DBError{fmt.Sprintf("Cannot remove the story from user's favorites. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	return nil
}

/*CheckIfUserSavedStory check if user already saved the story*/
func CheckIfUserSavedStory(userID int, storyID int) (bool, error) {
	db, err := connectToDB()
	defer db.Close()
	if err != nil {
		return false, &DBError{fmt.Sprintf("DB connection error. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	sql := "SELECT COUNT(*) as count FROM saved WHERE userid = $1 AND storyid = $2"
	row := db.QueryRow(sql, userID, storyID)
	var count int = 0
	err = row.Scan(&count)
	if err != nil {
		return false, &DBError{fmt.Sprintf("Cannot read count from db. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

/*GetRecentStories returns the paging recently published stories*/
func GetRecentStories(pageNumber int, pageRowCount int) (*[]Story, error) {
	db, err := connectToDB()
	defer db.Close()
	if err != nil {
		return nil, &DBError{fmt.Sprintf("DB connection error. PageNumber: %d, PageRowCount: %d", pageNumber, pageRowCount), err}
	}
	sql := "SELECT stories.*, users.username FROM stories INNER JOIN users ON stories.userid = users.id ORDER BY submittedon ASC LIMIT $1 OFFSET $2"
	rows, err := db.Query(sql, pageRowCount, pageNumber*pageRowCount)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot get stories. PageNumber: %d, PageRowCount: %d", pageNumber, pageRowCount), err}
	}
	stories, err := MapSQLRowsToStories(rows)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot read rows. PageNumber: %d, PageRowCount: %d", pageNumber, pageRowCount), err}
	}
	return stories, nil
}
