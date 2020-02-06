package models

type StoryPageData struct {
	Title   string
	User    User
	Stories []StoryViewModel
}

type StoryViewModel struct {
	ID              int
	Title           string
	URL             string
	UserID          int
	UserName        string
	Points          int
	CommentCount    int
	SubmittedOnText string
}

type ViewModel struct {
	Title string
	User  User
	Data  map[string]interface{}
}
