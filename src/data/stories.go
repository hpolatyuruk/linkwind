package data

import (
	"fmt"
	"math"
	"time"
	"turkdev/src/enums"

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
	DownVotes    int
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
	sql := "INSERT INTO stories (url, title, text, tags, upvotes, downvotes,  commentcount, userid, submittedon) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)"
	_, err = db.Exec(
		sql,
		story.URL,
		story.Title,
		story.Text,
		pq.Array(story.Tags),
		0,
		0,
		0,
		story.UserID,
		time.Now())
	if err != nil {
		return &StoryError{"Cannot create story!", story, err}
	}
	return nil
}

/*GetStories returns story list according to customer id by provided paging parameters*/
func GetStories(customerID, pageNumber, pageRowCount int) (*[]Story, error) {
	db, err := connectToDB()
	defer db.Close()
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	// TODO(Huseyin): Sort it by point algorithim when sedat finishes it
	sql := "SELECT stories.*, users.UserName FROM stories INNER JOIN users ON users.id = stories.userid WHERE users.customerid = $1 LIMIT $2 OFFSET $3"
	rows, err := db.Query(sql, customerID, pageRowCount, (pageNumber-1)*pageRowCount)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot get stories. PageNumber: %d, PageRowCount: %d", pageNumber, pageRowCount), err}
	}
	stories, err := MapSQLRowsToStories(rows)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot read rows. PageNumber: %d, PageRowCount: %d", pageNumber, pageRowCount), err}
	}
	return stories, nil
}

/*GetCustomerStoriesCount returns stories count number*/
func GetCustomerStoriesCount(customerID int) (int, error) {
	sql := "SELECT COUNT(*) FROM stories INNER JOIN users ON users.id = stories.userid WHERE users.customerid = $1"
	return count(sql, customerID)
}

/*GetStoryByID gets story by id from db*/
func GetStoryByID(storyID int) (*Story, error) {
	db, err := connectToDB()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	sql := "SELECT stories.*, users.UserName FROM stories INNER JOIN users ON users.id = stories.userid WHERE stories.id = $1"
	row := db.QueryRow(sql, storyID)
	story, err := MapSQLRowToStory(row)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot read row. Story Id: %d", storyID), err}
	}
	return story, nil
}

/*VoteStory votes (upvote, downvote) the story on database*/
func VoteStory(userID, storyID int, voteType enums.VoteType) error {
	voteUpdateSQL := "UPDATE stories SET upvotes = upvotes + 1 WHERE id = $1"
	if voteType == enums.DownVote {
		voteUpdateSQL = "UPDATE stories SET downvotes = downvotes + 1 WHERE id = $1"
	}
	db, err := connectToDB()
	defer db.Close()
	if err != nil {
		return &DBError{fmt.Sprintf("DB connection error. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	tran, err := db.Begin()
	if err != nil {
		return &DBError{fmt.Sprintf("Cannot begin transaction. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	sql := "INSERT INTO storyvotes(storyid, userid, votetype) VALUES($1, $2, $3)"
	_, err = tran.Exec(sql, storyID, userID, voteType)
	if err != nil {
		tran.Rollback()
		return &DBError{fmt.Sprintf("Error occurred while inserting storyvotes. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	_, err = tran.Exec(voteUpdateSQL, storyID)
	if err != nil {
		tran.Rollback()
		return &DBError{fmt.Sprintf("Error occurred while increasing story upvotes. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	if voteType == enums.UpVote {
		sql = "UPDATE users SET karma = karma + 1 WHERE id = (SELECT userid FROM stories WHERE id = $1)"
		_, err = tran.Exec(sql, storyID)
		if err != nil {
			tran.Rollback()
			return &DBError{fmt.Sprintf("Error occurred while increasing user's karma. UserID: %d, StoryID: %d", userID, storyID), err}
		}
	}
	err = tran.Commit()
	if err != nil {
		return &DBError{fmt.Sprintf("Cannot commit transaction. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	return nil
}

/*CheckIfStoryVotedByUser check if user already voted(upvote or downvote) to given story*/
func CheckIfStoryVotedByUser(userID, storyID int, voteType enums.VoteType) (bool, error) {
	db, err := connectToDB()
	defer db.Close()
	if err != nil {
		return false, &DBError{fmt.Sprintf("DB connection error. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	sql := "SELECT COUNT(*) as count FROM storyvotes WHERE userid = $1 AND storyid = $2 AND votetype = $3"
	row := db.QueryRow(sql, userID, storyID, voteType)
	count := 0
	err = row.Scan(&count)
	if err != nil {
		return false, &DBError{fmt.Sprintf("Cannot read db row. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

/*RemoveStoryVote removes the vote (upvote, downvote) of story on database*/
func RemoveStoryVote(userID, storyID int, voteType enums.VoteType) error {
	voteUpdateSQL := "UPDATE stories SET upvotes = upvotes - 1 WHERE id = $1"
	if voteType == enums.DownVote {
		voteUpdateSQL = "UPDATE stories SET downvotes = downvotes - 1 WHERE id = $1"
	}
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
	_, err = tran.Exec(voteUpdateSQL, storyID)
	if err != nil {
		tran.Rollback()
		return &DBError{fmt.Sprintf("Cannot update story's vote. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	if voteType == enums.UpVote {
		sql = "UPDATE users SET karma = karma - 1 WHERE id = (SELECT userid FROM stories WHERE id = $1)"
		_, err = tran.Exec(sql, storyID)
		if err != nil {
			tran.Rollback()
			return &DBError{fmt.Sprintf("Error occurred while decreasing user's karma. UserID: %d, StoryID: %d", userID, storyID), err}
		}
	}
	err = tran.Commit()
	if err != nil {
		return &DBError{fmt.Sprintf("Cannot commit transaction. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	return nil
}

/*SaveStory saves the given story to user's favorites*/
func SaveStory(userID int, storyID int) error {
	db, err := connectToDB()
	defer db.Close()
	if err != nil {
		return &DBError{fmt.Sprintf("DB connection error. UserID: %d, StoryID: %d", userID, storyID), err}
	}
	sql := "INSERT INTO saved (userid, storyid, savedon) VALUES ($1, $2, $3)"
	_, err = db.Exec(sql, userID, storyID, time.Now())
	if err != nil {
		return &DBError{fmt.Sprintf("Cannot save story to user's favorites. UserID: %d, StoryID: %d", userID, storyID), err}
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
	count := 0
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
func GetRecentStories(customerID, pageNumber, pageRowCount int) (*[]Story, error) {
	db, err := connectToDB()
	defer db.Close()
	if err != nil {
		return nil, &DBError{fmt.Sprintf("DB connection error. PageNumber: %d, PageRowCount: %d", pageNumber, pageRowCount), err}
	}
	sql := "SELECT stories.*, users.username FROM stories INNER JOIN users ON stories.userid = users.id WHERE users.customerid = $1 ORDER BY stories.submittedon DESC LIMIT $2 OFFSET $3"
	rows, err := db.Query(sql, customerID, pageRowCount, (pageNumber-1)*pageRowCount)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot get stories. PageNumber: %d, PageRowCount: %d", pageNumber, pageRowCount), err}
	}
	stories, err := MapSQLRowsToStories(rows)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot read rows. PageNumber: %d, PageRowCount: %d", pageNumber, pageRowCount), err}
	}
	return stories, nil
}

/*GetUserSavedStories returns the paging user's favorite stories*/
func GetUserSavedStories(userID int, pageNumber int, pageRowCount int) (*[]Story, error) {
	db, err := connectToDB()
	defer db.Close()
	if err != nil {
		return nil, &DBError{fmt.Sprintf("DB connection error. UserID: %d PageNo: %d, PageRowCount: %d", userID, pageNumber, pageRowCount), err}
	}
	sql := "SELECT stories.*, users.username FROM stories INNER JOIN saved ON stories.id = saved.storyid INNER JOIN users ON users.id = saved.userid WHERE saved.userid = $1 ORDER BY savedon DESC LIMIT $2 OFFSET $3"
	rows, err := db.Query(sql, userID, pageRowCount, (pageNumber-1)*pageRowCount)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot query user's saved stories. UserID: %d, PageNumber: %d, PageRowCount: %d", userID, pageNumber, pageRowCount), err}
	}
	stories, err := MapSQLRowsToStories(rows)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot map sql rows to story struct array. UserID: %d, PageNumber: %d, PageRowCount: %d", userID, pageNumber, pageRowCount), err}
	}
	return stories, nil
}

/*GetUserSavedStoriesCount gets the total number of user's saved stories.*/
func GetUserSavedStoriesCount(userID int) (int, error) {
	sql := "SELECT COUNT(stories.id) FROM stories INNER JOIN saved ON stories.id = saved.storyid INNER JOIN users ON users.id = saved.userid WHERE saved.userid = $1"
	return count(sql, userID)
}

/*GetUserUpvotedStories returns the paging user's upvoted stories*/
func GetUserUpvotedStories(userID int, pageNumber int, pageRowCount int) (*[]Story, error) {
	db, err := connectToDB()
	defer db.Close()
	if err != nil {
		return nil, &DBError{fmt.Sprintf("DB connection error. UserID: %d PageNo: %d, PageRowCount: %d", userID, pageNumber, pageRowCount), err}
	}
	sql := "SELECT stories.*, users.username FROM stories INNER JOIN storyvotes ON stories.id = storyvotes.storyid INNER JOIN users ON users.id = storyvotes.userid WHERE storyvotes.userid = $1 ORDER BY stories.submittedon DESC LIMIT $2 OFFSET $3"
	rows, err := db.Query(sql, userID, pageRowCount, (pageNumber-1)*pageRowCount)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot query user's saved stories. UserID: %d, PageNumber: %d, PageRowCount: %d", userID, pageNumber, pageRowCount), err}
	}
	stories, err := MapSQLRowsToStories(rows)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot map sql rows to story struct array. UserID: %d, PageNumber: %d, PageRowCount: %d", userID, pageNumber, pageRowCount), err}
	}
	return stories, nil
}

/*GetUserUpvotedStoriesCount gets the total number of user's upvoted stories.*/
func GetUserUpvotedStoriesCount(userID int) (int, error) {
	sql := "SELECT COUNT(stories.id) FROM stories INNER JOIN storyvotes ON stories.id = storyvotes.storyid INNER JOIN users ON users.id = storyvotes.userid WHERE storyvotes.userid = $1"
	return count(sql, userID)
}

/*GetUserSubmittedStories get user's stories from db according to userID*/
func GetUserSubmittedStories(userID int, pageNumber int, pageRowCount int) (*[]Story, error) {
	db, err := connectToDB()
	defer db.Close()
	if err != nil {
		return nil, err
	}

	sql := "SELECT stories.*, users.username FROM stories INNER JOIN users ON users.id = stories.userid WHERE stories.userid = $1 ORDER BY submittedon DESC LIMIT $2 OFFSET $3"
	rows, err := db.Query(sql, userID, pageRowCount, (pageNumber-1)*pageRowCount)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot query user's posted stories. UserID: %d, PageNumber: %d, PageRowCount: %d", userID, pageNumber, pageRowCount), err}
	}

	stories, err := MapSQLRowsToStories(rows)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot map sql rows to story struct array. UserID: %d", userID), err}
	}
	return stories, nil
}

/*GetUserSubmittedStoriesCount gets the total number of user's submissions.*/
func GetUserSubmittedStoriesCount(userID int) (int, error) {
	sql := "SELECT COUNT(stories.id) FROM stories INNER JOIN users ON users.id = stories.userid WHERE stories.userid = $1"
	return count(sql, userID)
}

/*CalculateStoryPenalty calculates story's penalty. If commentCount rises, penalty downs*/
func CalculateStoryPenalty(commentCount int) int {
	penalty := 40
	for i := 0; i < 40; i++ {
		if commentCount == i {
			return penalty - i
		}
	}
	return 0
}

/*CalculateStoryRank calcualte story's rank according to this formula :
	 	  ((upVotes-downVotes)-1)^0.8
Score =	———————————————————————————————— * penalty
		((submittedTime-time.Now)+2)^1.8
More details : http://www.righto.com/2013/11/how-hacker-news-ranking-really-works.html*/
func CalculateStoryRank(penalty, votes, timeDiff int) int {
	floatScore := (math.Pow(float64(votes-1), 0.8) / math.Pow(float64(timeDiff+2), 1.8)) * float64(penalty)

	return int(floatScore)
}

func count(sql string, args ...interface{}) (int, error) {
	var count int
	db, err := connectToDB()
	defer db.Close()
	if err != nil {
		return count, err
	}
	row := db.QueryRow(sql, args...)
	err = row.Scan(&count)
	if err != nil {
		return count, &DBError{fmt.Sprintf("Cannot read sql row to get count."), err}
	}
	return count, nil
}
