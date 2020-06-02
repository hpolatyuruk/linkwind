package models

import (
	"html/template"
	"linkwind/app/shared"
	"strings"
)

// StoryPageViewModel contains stories that will show on stroies page
type StoryPageViewModel struct {
	Title           string
	IsAuthenticated bool
	Stories         []StoryViewModel
	Page            *Paging
	BaseViewModel
}

/*SetLayout sets story page view model layout members.*/
func (model *StoryPageViewModel) SetLayout(platformName string, logo string) {
	model.Layout = generateLayoutViewModel(platformName, logo)
}

/*SetSignedInUser sets story page view model signed in user members.*/
func (model *StoryPageViewModel) SetSignedInUser(userClaims *shared.SignedInUserClaims) {

	model.IsAuthenticated = false
	if userClaims == nil {
		return
	}
	model.IsAuthenticated = true
	model.SignedInUser = generateSignedInUserViewModel(userClaims)
}

/*StorySubmitModel represents the data to submit a story.*/
type StorySubmitModel struct {
	URL    string
	Title  string
	Text   string
	Errors map[string]string
	BaseViewModel
}

/*Validate validates the StorySubmitModel*/
func (model *StorySubmitModel) Validate() bool {
	model.Errors = make(map[string]string)
	if strings.TrimSpace(model.URL) == "" &&
		strings.TrimSpace(model.Text) == "" &&
		strings.TrimSpace(model.Text) == "" {
		model.Errors["General"] = "Please enter a url or title/text."
		return false
	}
	if strings.TrimSpace(model.URL) == "" &&
		strings.TrimSpace(model.Text) != "" &&
		strings.TrimSpace(model.Title) == "" {
		model.Errors["Title"] = "Please enter a title."
		return false
	}
	return true
}

/*SetLayout sets story submit model layout members.*/
func (model *StorySubmitModel) SetLayout(platformName string, logo string) {
	model.Layout = generateLayoutViewModel(platformName, logo)
}

/*SetSignedInUser sets story submit model signed in user members.*/
func (model *StorySubmitModel) SetSignedInUser(userClaims *shared.SignedInUserClaims) {
	if userClaims == nil {
		return
	}
	model.SignedInUser = generateSignedInUserViewModel(userClaims)
}

/*StoryDetailPageViewModel represents the story page view model that contains indivual story informations*/
type StoryDetailPageViewModel struct {
	Title           string
	Story           *StoryViewModel
	Comments        *[]CommentViewModel
	IsAuthenticated bool
	BaseViewModel
}

/*SetLayout sets story detail page view model layout members.*/
func (model *StoryDetailPageViewModel) SetLayout(platformName string, logo string) {
	model.Layout = generateLayoutViewModel(platformName, logo)
}

/*SetSignedInUser sets story detail page view model signed in user members.*/
func (model *StoryDetailPageViewModel) SetSignedInUser(userClaims *shared.SignedInUserClaims) {

	model.IsAuthenticated = false
	if userClaims == nil {
		return
	}
	model.IsAuthenticated = true
	model.SignedInUser = generateSignedInUserViewModel(userClaims)
}

// StoryViewModel represents the an indiviudual story information
type StoryViewModel struct {
	ID              int
	Title           string
	URL             string
	Text            template.HTML
	Host            string
	UserID          int
	UserName        string
	Points          int
	CommentCount    int
	SubmittedOnText string
	IsSaved         bool
	IsUpvoted       bool
	IsDownvoted     bool
	ShowDownvoteBtn bool
	SignedInUser    *SignedInUserViewModel
}
