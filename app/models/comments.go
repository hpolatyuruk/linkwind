package models

// CommentViewModel represents the individual comment information
type CommentViewModel struct {
	ID              int
	UserID          int
	UserName        string
	StoryID         int
	Points          int
	Comment         string
	CommentedOnText string
	IsUpvoted       bool
	IsDownvoted     bool
	ShowDownvoteBtn bool
	IsRoot          bool
	ParentID        int
	ChildComments   []CommentViewModel
	SignedInUser    *SignedInUserViewModel
}
