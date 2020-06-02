package models

import (
	"linkwind/app/data"
	"linkwind/app/shared"
	"strings"
)

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
	ID             int
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
