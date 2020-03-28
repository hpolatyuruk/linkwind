package data

import (
	"database/sql"
	"fmt"
	"time"
	"turkdev/src/enums"
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
	UserName    string
	ParentID    int
	UpVotes     int
	DownVotes   int
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

/*Reply represents the comments which belongs to a story.*/
type Reply struct {
	Comment    *Comment
	StoryTitle string
	StoryID    int
	UserName   string
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

/*WriteComment insert a comment to database.*/
func WriteComment(comment *Comment) (*int, error) {
	var commentID int
	db, err := connectToDB()
	if err != nil {
		return &commentID, &CommentError{"DB connection error", comment, err}
	}
	defer db.Close()
	tran, err := db.Begin()
	if err != nil {
		return &commentID, &CommentError{"Can not start the transaction.", comment, err}
	}
	sql := "INSERT INTO comments (storyid, userid, parentid, upvotes, downvotes, replycount, comment, commentedon) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"
	err = tran.QueryRow(
		sql,
		comment.StoryID,
		comment.UserID,
		nullCommentParentID(comment.ParentID),
		comment.UpVotes,
		comment.DownVotes,
		comment.ReplyCount,
		comment.Comment,
		comment.CommentedOn).Scan(&commentID)
	if err != nil {
		tran.Rollback()
		return &commentID, &CommentError{"Cannot insert comment to the db.", comment, err}
	}
	sql = "UPDATE stories SET commentcount = commentcount + 1 WHERE id = $1"
	_, err = tran.Exec(sql, comment.StoryID)
	if err != nil {
		tran.Rollback()
		return &commentID, &CommentError{"Cannot increase story's comment count.", comment, err}
	}
	if comment.ParentID != CommentRootID {
		sql = "UPDATE comments SET replycount = replycount + 1 WHERE id = $1"
		_, err = tran.Exec(sql, comment.ParentID)
		if err != nil {
			tran.Rollback()
			return &commentID, &CommentError{"Cannot increase comment's reply count.", comment, err}
		}
	}
	err = tran.Commit()
	if err != nil {
		return &commentID, &CommentError{"Cannot commit transaction.", comment, err}
	}
	return &commentID, nil
}

/*GetComments retunrs comment list by provided story id*/
func GetComments(storyID int) (comments *[]Comment, err error) {
	db, err := connectToDB()

	if err != nil {
		return nil, &DBError{fmt.Sprintf("DB connection error. StoryID: %d.", storyID), err}
	}
	defer db.Close()
	// TODO(Huseyin): Order by special algorithm when Sedat finishes it
	sql := "SELECT comments.*, users.username FROM comments INNER JOIN users ON users.id = comments.userid WHERE storyid = $1"
	rows, err := db.Query(sql, storyID)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot query comments. StoryID: %d.", storyID), err}
	}
	comments, err = MapSQLRowsToComments(rows)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot read comment row. StoryID: %d.", storyID), err}
	}
	return comments, nil
}

/*GetRootCommentsByStoryID retunrs only root comments by provided story id*/
func GetRootCommentsByStoryID(storyID int) (comments *[]Comment, err error) {
	db, err := connectToDB()
	if err != nil {
		return nil, &DBError{fmt.Sprintf("DB connection error. StoryID: %d.", storyID), err}
	}
	defer db.Close()
	sql := "SELECT comments.*, users.username FROM comments INNER JOIN users ON users.id = comments.userid WHERE storyid = $1 AND parentid is null ORDER BY commentedon ASC"
	rows, err := db.Query(sql, storyID)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot query comments. StoryID: %d.", storyID), err}
	}
	comments, err = MapSQLRowsToComments(rows)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot read comment row. StoryID: %d.", storyID), err}
	}
	return comments, nil
}

/*GetCommentsByParentIDAndStoryID retunrs only root comments by provided story id*/
func GetCommentsByParentIDAndStoryID(parentID int, storyID int) (comments *[]Comment, err error) {
	db, err := connectToDB()
	if err != nil {
		return nil, &DBError{fmt.Sprintf("DB connection error. ParentID: %d, StoryID: %d.", parentID, storyID), err}
	}
	defer db.Close()
	sql := "SELECT comments.*, users.username FROM comments INNER JOIN users ON users.id = comments.userid WHERE storyid = $1 AND parentid = $2 ORDER BY commentedon ASC"
	rows, err := db.Query(sql, storyID, parentID)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot query comments by parent id and story id. ParentID: %d and StoryID: %d.", parentID, storyID), err}
	}
	comments, err = MapSQLRowsToComments(rows)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot read comment row. ParentID: %d, StoryID: %d.", parentID, storyID), err}
	}
	return comments, nil
}

/*VoteComment votes (upvote, downvote) for comment on database*/
func VoteComment(userID int, commentID int, voteType enums.VoteType) error {
	updateQuery := "UPDATE comments SET upvotes = upvotes + 1 WHERE id = $1"
	if voteType == enums.DownVote {
		updateQuery = "UPDATE comments SET downvotes = downvotes + 1 WHERE id = $1"
	}
	db, err := connectToDB()
	if err != nil {
		return &DBError{fmt.Sprintf("DB connection error. UserID: %d, CommentID: %d", userID, commentID), err}
	}
	defer db.Close()
	tran, err := db.Begin()
	if err != nil {
		return &DBError{fmt.Sprintf("Cannot begin transaction. UserID: %d, CommentID: %d", userID, commentID), err}
	}
	sql := "INSERT INTO commentvotes (userid, commentid, votetype) VALUES ($1, $2, $3)"
	_, err = tran.Exec(sql, userID, commentID, voteType)
	if err != nil {
		tran.Rollback()
		return &DBError{fmt.Sprintf("Cannot insert commentvotes. UserID: %d, CommentID: %d", userID, commentID), err}
	}
	_, err = tran.Exec(updateQuery, commentID)
	if err != nil {
		tran.Rollback()
		return &DBError{fmt.Sprintf("Cannot update comment's upvotes. UserID: %d, CommentID: %d", userID, commentID), err}
	}
	if voteType == enums.UpVote {
		sql = "UPDATE users SET karma = karma + 1 WHERE id = (SELECT userid FROM comments WHERE id = $1)"
		_, err = tran.Exec(sql, commentID)
		if err != nil {
			tran.Rollback()
			return &DBError{fmt.Sprintf("Error occurred while increasing user's karma. UserID: %d, CommentID: %d", userID, commentID), err}
		}
	}
	err = tran.Commit()
	if err != nil {
		return &DBError{fmt.Sprintf("Cannot commit transaction. UserID: %d, CommentID: %d", userID, commentID), err}
	}
	return nil
}

/*RemoveCommentVote unvotes (upvote, downvote) the comment on database*/
func RemoveCommentVote(userID int, commentID int, voteType enums.VoteType) error {
	updateQuery := "UPDATE comments SET upvotes = upvotes - 1 WHERE id = $1"
	if voteType == enums.DownVote {
		updateQuery = "UPDATE comments SET downvotes = downvotes - 1 WHERE id = $1"
	}
	db, err := connectToDB()
	if err != nil {
		return &DBError{fmt.Sprintf("DB connection error. UserID: %d, CommentID: %d", userID, commentID), err}
	}
	defer db.Close()
	tran, err := db.Begin()
	if err != nil {
		return &DBError{fmt.Sprintf("Cannot begin transaction. UserID: %d, CommentID: %d", userID, commentID), err}
	}
	sql := "DELETE FROM commentvotes WHERE userid = $1 AND commentid = $2"
	_, err = tran.Exec(sql, userID, commentID)
	if err != nil {
		tran.Rollback()
		return &DBError{fmt.Sprintf("Cannot delete commentvotes. UserID: %d, CommentID: %d", userID, commentID), err}
	}
	_, err = tran.Exec(updateQuery, commentID)
	if err != nil {
		tran.Rollback()
		return &DBError{fmt.Sprintf("Cannot update comment's upvotes. UserID: %d, CommentID: %d", userID, commentID), err}
	}
	if voteType == enums.UpVote {
		sql = "UPDATE users SET karma = karma - 1 WHERE id = (SELECT userid FROM comments WHERE id = $1)"
		_, err = tran.Exec(sql, commentID)
		if err != nil {
			tran.Rollback()
			return &DBError{fmt.Sprintf("Error occurred while decreasing user's karma. UserID: %d, CommentID: %d", userID, commentID), err}
		}
	}
	err = tran.Commit()
	if err != nil {
		return &DBError{fmt.Sprintf("Cannot commit transaction. UserID: %d, CommentID: %d", userID, commentID), err}
	}
	return nil
}

/*GetCommentVoteByUser gets type of vote to given comment by user*/
func GetCommentVoteByUser(userID int, commentID int) (*enums.VoteType, error) {
	db, err := connectToDB()
	if err != nil {
		return nil, &DBError{fmt.Sprintf("DB connection error. UserID: %d, CommentID: %d", userID, commentID), err}
	}
	defer db.Close()
	query := "SELECT votetype FROM commentvotes WHERE userid = $1 and commentid = $2"
	row := db.QueryRow(query, userID, commentID)
	var voteType *enums.VoteType = nil
	err = row.Scan(&voteType)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, &DBError{fmt.Sprintf("Cannot read db row. UserID: %d, CommentID: %d", userID, commentID), err}
	}
	return voteType, nil
}

/*GetUserReplies returns reply list by provided user id and paging parameters*/
func GetUserReplies(userID int) (replies *[]Reply, err error) {
	db, err := connectToDB()
	if err != nil {
		return nil, &DBError{fmt.Sprintf("DB connection error. UserID: %d.", userID), err}
	}
	defer db.Close()
	sql := "SELECT comments.*, stories.title, stories.id, users.username FROM comments INNER JOIN stories ON comments.storyid = stories.id INNER JOIN users ON users.id = stories.userid WHERE stories.userid = $1 ORDER BY comments.commentedon DESC"
	rows, err := db.Query(sql, userID)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot query replies. UserID: %d.", userID), err}
	}
	replies, err = MapSQLRowsToReplies(rows)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot read reply row. UserID: %d.", userID), err}
	}
	return replies, nil
}

/*GetUserCommentsNotPaging get user's comments from db according to userID and not paging*/
func GetUserCommentsNotPaging(userID int) (comments *[]Comment, err error) {
	db, err := connectToDB()
	if err != nil {
		return nil, &DBError{fmt.Sprintf("DB connection error. UserID: %d.", userID), err}
	}
	defer db.Close()
	// TODO(Huseyin): Order by special algorithm when Sedat finishes it
	sql := "SELECT * FROM public.comments WHERE userid = $1"
	rows, err := db.Query(sql, userID)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot query comments. UserID: %d.", userID), err}
	}
	comments, err = MapSQLRowsToComments(rows)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot read comment row. UserID: %d.", userID), err}
	}
	return comments, nil
}
