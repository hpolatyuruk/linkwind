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
	DownVotes    int
	CommentCount int
	UserID       int
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
	sql := "INSERT INTO stories (url, title, text, tags, upvotes, downvotes, commentcount, userid, submittedon) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)"
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

/*GetStories retunrs story list by provided paging parameters*/
func GetStories(pageNumber int, pageRowCount int) ([]*Story, error) {
	db, err := connectToDB()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	// TODO(Huseyin): Sort it by point algorithim when sedat finishes it
	sql := "SELECT id, url, title, text, tags, upvotes, downvotes, commentcount, userid, submittedon FROM stories LIMIT $1 OFFSET $2"
	rows, err := db.Query(sql, pageRowCount, pageNumber*pageRowCount)
	if err != nil {
		return nil, &StoryError{fmt.Sprintf("Cannot get stories. PageNumber: %d, PageRowCount: %d", pageNumber, pageRowCount), nil, err}
	}
	stories := []*Story{}
	for rows.Next() {
		var story Story
		err = rows.Scan(
			&story.ID,
			&story.URL,
			&story.Title,
			&story.Text,
			pq.Array(&story.Tags),
			&story.UpVotes,
			&story.DownVotes,
			&story.CommentCount,
			&story.UserID,
			&story.SubmittedOn)

		if err != nil {
			return nil, &StoryError{fmt.Sprintf("Cannot read rows. PageNumber: %d, PageRowCount: %d", pageNumber, pageRowCount), nil, err}
		}
		stories = append(stories, &story)
	}
	return stories, nil
}
