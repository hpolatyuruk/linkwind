package models

import "html/template"

type SignedInUserViewModel struct {
	UserID     int
	CustomerID int
	UserName   string
	Email      string
}

type Paging struct {
	CurrentPage    int
	PreviousPage   int
	NextPage       int
	IsFinalPage    bool
	TotalPageCount int
}

type StoryPageViewModel struct {
	Title           string
	IsAuthenticated bool
	SignedInUser    *SignedInUserViewModel
	Stories         []StoryViewModel
	Page            *Paging
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
