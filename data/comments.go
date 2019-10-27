package data

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	/*CommentRootID represents the default value for null parent id. Because we cannot make int as nil*/
	CommentRootID int = 0
)

/*Comment represents the comment object which can be written for a story by a user*/
type Comment struct {
	ID          int
	StoryID     int
	UserID      int
	ParentID    int
	UpVotes     int
	ReplyCount  int
	Comment     string
	CommentedOn time.Time
}

/*CommentError contains the error and comment data which caused to error*/
type CommentError struct {
	Message       string
	Comment       *Comment
	OriginalError error
}

func (err *CommentError) Error() string {
	return fmt.Sprintf(
		"%s | Comment: %v | OriginalError: %v",
		err.Message,
		err.Comment,
		err.OriginalError)
}

func nullCommentParentID(i int) sql.NullInt32 {
	if i == CommentRootID {
		return sql.NullInt32{}
	}
	return sql.NullInt32{
		Int32: int32(i),
		Valid: true,
	}
}

/*WriteComment insert a comment to database*/
func WriteComment(comment *Comment) error {
	db, err := connectToDB()
	defer db.Close()
	if err != nil {
		return &CommentError{"DB connection error", comment, err}
	}
	sql := "INSERT INTO comments (storyid, userid, parentid, upvotes, replycount, comment, commentedon) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	_, err = db.Exec(
		sql,
		comment.StoryID,
		comment.UserID,
		nullCommentParentID(comment.ParentID),
		comment.UpVotes,
		comment.ReplyCount,
		comment.Comment,
		comment.CommentedOn)
	if err != nil {
		return &CommentError{"Cannot inser comment to the db.", comment, err}
	}
	return nil
}

/*GetComments retunrs comment list by provided story id and paging parameters*/
func GetComments(storyID int, pageNumber int, pageRowCount int) (comments *[]Comment, err error) {
	db, err := connectToDB()
	defer db.Close()
	if err != nil {
		return nil, &DBError{fmt.Sprintf("DB connection error. StoryID: %d, PageNumber: %d, PageRowCount: %d", storyID, pageNumber, pageRowCount), err}
	}
	// TODO(Huseyin): Order by special algorithm when Sedat finishes it
	sql := "SELECT id, comment, upvotes, storyid, parentid, replycount, userid, commentedon FROM comments WHERE storyid = $1 LIMIT $2 OFFSET $3"
	rows, err := db.Query(sql, storyID, pageRowCount, pageNumber*pageRowCount)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot query comments. StoryID: %d, PageNumber: %d, PageRowCount: %d", storyID, pageNumber, pageRowCount), err}
	}
	comments, err = MapSQLRowsToComments(rows)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot read comment row. StoryID: %d, PageNumber: %d, PageRowCount: %d", storyID, pageNumber, pageRowCount), err}
	}
	return comments, nil
}
