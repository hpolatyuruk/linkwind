package models

type StoryPageData struct {
	Title      string
	UserID     int
	UserName   string
	CustomerID int
	IsSignedIn bool
	Stories    []StoryViewModel
}

type StoryViewModel struct {
	ID                    int
	Title                 string
	URL                   string
	UserID                int
	UserName              string
	Points                int
	CommentCount          int
	SubmittedOnText       string
	IsSavedBySignedInUser bool
	IsUpvotedSignedInUser bool
}

type ViewModel struct {
	Title string
	User  User
	Data  map[string]interface{}
}
