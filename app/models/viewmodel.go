package models

import "html/template"

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
	ID                      int
	Title                   string
	URL                     string
	Text                    template.HTML
	Host                    string
	UserID                  int
	UserName                string
	Points                  int
	CommentCount            int
	SubmittedOnText         string
	IsSavedBySignedInUser   bool
	IsUpvotedBySignedInUser bool
	SignedInUser            *SignedInUserViewModel
}

type ViewModel struct {
	Title string
	User  User
	Data  map[string]interface{}
}

type CommentViewModel struct {
	ID                      int
	UserID                  int
	UserName                string
	StoryID                 int
	Points                  int
	Comment                 string
	CommentedOnText         string
	IsUpvotedBySignedInUser bool
	IsRoot                  bool
	ParentID                int
	ChildComments           []CommentViewModel
	SignedInUser            *SignedInUserViewModel
}

type StoryDetailPageViewModel struct {
	Title           string
	Story           *StoryViewModel
	Comments        *[]CommentViewModel
	SignedInUser    *SignedInUserViewModel
	IsAuthenticated bool
}
