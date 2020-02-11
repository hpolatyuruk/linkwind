package models

import "turkdev/data"

type SignedInUserViewModel struct {
	UserID     int
	CustomerID int
	UserName   string
	Email      string
}

type StoryPageViewModel struct {
	Title           string
	IsAuthenticated bool
	SignedInUser    *SignedInUserViewModel
	Stories         []StoryViewModel
}

type StoryViewModel struct {
	ID                    int
	Title                 string
	URL                   string
	Host                  string
	UserID                int
	UserName              string
	Points                int
	CommentCount          int
	SubmittedOnText       string
	IsSavedBySignedInUser bool
	IsUpvotedSignedInUser bool
	SignedInUser          *SignedInUserViewModel
}

type ViewModel struct {
	Title string
	User  User
	Data  map[string]interface{}
}

type StoryDetailPageViewModel struct {
	Title           string
	Story           *data.Story
	Comments        *[]data.Comment
	SignedInUser    SignedInUserViewModel
	IsAuthenticated bool
}
