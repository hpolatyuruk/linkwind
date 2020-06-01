package models

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"image"
	"linkwind/app/data"
	"linkwind/app/shared"
	"strings"
)

type LayoutViewModel struct {
	Platform string
	Logo     string
}

type SignedInUserViewModel struct {
	UserID     int
	CustomerID int
	UserName   string
	Email      string
	Karma      int
}

/*BaseViewModelInterface represents the base view model interface that contains methods to set BaseViewModel*/
type BaseViewModelInterface interface {
	SetLayout(platformName string, logo string)
	SetSignedInUser(userClaims *shared.SignedInUserClaims)
}

/*BaseViewModel represents the base view model that container layout and signedin user informations*/
type BaseViewModel struct {
	Layout       *LayoutViewModel
	SignedInUser *SignedInUserViewModel
}

func generateLayoutViewModel(platformName string, logo string) *LayoutViewModel {
	return &LayoutViewModel{
		Platform: platformName,
		Logo:     logo,
	}
}

func generateSignedInUserViewModel(userClaims *shared.SignedInUserClaims) *SignedInUserViewModel {
	return &SignedInUserViewModel{
		UserName:   userClaims.UserName,
		UserID:     userClaims.ID,
		CustomerID: userClaims.CustomerID,
		Email:      userClaims.Email,
		Karma:      userClaims.Karma,
	}
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

/*InviteUserViewModel represents the data which is needed on sigin UI.*/
type InviteUserViewModel struct {
	EmailAddress   string
	SuccessMessage string
	Memo           string
	Errors         map[string]string
	BaseViewModel
}

/*SetLayout sets story detail page view model layout members.*/
func (model *InviteUserViewModel) SetLayout(platformName string, logo string) {
	model.Layout = generateLayoutViewModel(platformName, logo)
}

/*SetSignedInUser sets story detail page view model signed in user members.*/
func (model *InviteUserViewModel) SetSignedInUser(userClaims *shared.SignedInUserClaims) {
	if userClaims == nil {
		return
	}
	model.SignedInUser = generateSignedInUserViewModel(userClaims)
}

/*Validate validates the InviteUserViewModel*/
func (model *InviteUserViewModel) Validate(email string) (bool, error) {
	model.Errors = make(map[string]string)

	if strings.TrimSpace(model.EmailAddress) == "" {
		model.Errors["Email"] = "Email is required!"
	} else {
		if shared.IsEmailAdressValid(model.EmailAddress) == false {
			model.Errors["Email"] = "Please enter a valid email address!"
		}
	}

	exists, err := data.ExistsUserByEmail(email)
	if err != nil {
		return false, err
	}
	if exists {
		model.Errors["Email"] = "This email address is already in use!"
	}
	return len(model.Errors) == 0, nil
}

/*CustomerAdminViewModel represents the data which is needed on sigin UI.*/
type CustomerAdminViewModel struct {
	Name              string
	Domain            string
	LogoImageAsBase64 string
	Errors            map[string]string
	SuccessMessage    string
	BaseViewModel
}

/*SetLayout sets story detail page view model layout members.*/
func (model *CustomerAdminViewModel) SetLayout(platformName string, logo string) {
	model.Layout = generateLayoutViewModel(platformName, logo)
}

/*SetSignedInUser sets story detail page view model signed in user members.*/
func (model *CustomerAdminViewModel) SetSignedInUser(userClaims *shared.SignedInUserClaims) {
	if userClaims == nil {
		return
	}
	model.SignedInUser = generateSignedInUserViewModel(userClaims)
}

func getImageInfos(file []byte) (int, int, string, error) {
	r := bytes.NewReader(file)
	im, format, err := image.DecodeConfig(r)
	if err != nil {
		return 0, 0, "", err
	}

	return im.Width, im.Height, format, nil
}

func checkImageValid(err error, format string) error {
	if err != nil {
		if err.Error() == "image: unknown format" {
			return errors.New("Image format should be jpg")
		}
		panic(err)
	}
	if format != "jpeg" {
		return errors.New("Image format should be jpg")
	}
	return nil
}

/*Validate validates the InviteUserViewModel*/
func (model *CustomerAdminViewModel) Validate() bool {

	const (
		maxPlatformNameLength = 25
		maxImageWidth         = 30
		maxImageLength        = 30
	)

	model.Errors = make(map[string]string)

	if strings.TrimSpace(model.Name) == "" {
		model.Errors["Name"] = "Name is required!"
	} else {
		if len(model.Name) > maxPlatformNameLength {
			model.Errors["Name"] = "Name cannot be longer than 25 characters"
		}
		if strings.Contains(model.Name, " ") {
			model.Errors["Name"] = "Name cannot contain spaces"
		}
	}

	if model.LogoImageAsBase64 != "" {
		decodingLogo, err := base64.StdEncoding.DecodeString(model.LogoImageAsBase64)
		if err != nil {
			panic(err)
		}

		width, height, format, errImg := getImageInfos(decodingLogo)
		err = checkImageValid(errImg, format)
		if err != nil {
			model.Errors["LogoImageAsBase64"] = err.Error()
		}

		if width > maxImageWidth || height > maxImageLength {
			model.Errors["LogoImageAsBase64"] =
				fmt.Sprintf("Image file size should be %d*%d",
					maxImageWidth, maxImageLength)
		}
	}
	return len(model.Errors) == 0
}

/*UserProfileViewModel represents the data which is needed on sigin UI.*/
type UserProfileViewModel struct {
	About          string
	Email          string
	FullName       string
	Karma          int
	RegisteredOn   string
	UserName       string
	Errors         map[string]string
	SuccessMessage string
	IsAdmin        bool
	BaseViewModel
}

/*Validate validates the UserProfileViewModel*/
func (model *UserProfileViewModel) Validate() bool {
	model.Errors = make(map[string]string)
	if strings.TrimSpace(model.Email) == "" {
		model.Errors["Email"] = "Email is required!"
	}
	if !shared.IsEmailAdressValid(model.Email) {
		model.Errors["General"] = "E-mail address is not valid!"
	}
	return len(model.Errors) == 0
}

/*SetLayout sets story detail page view model layout members.*/
func (model *UserProfileViewModel) SetLayout(platformName string, logo string) {
	model.Layout = generateLayoutViewModel(platformName, logo)
}

/*SetSignedInUser sets story detail page view model signed in user members.*/
func (model *UserProfileViewModel) SetSignedInUser(userClaims *shared.SignedInUserClaims) {
	if userClaims == nil {
		return
	}
	model.SignedInUser = generateSignedInUserViewModel(userClaims)
}
