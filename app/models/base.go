package models

import (
	"linkwind/app/shared"
)

// LayoutViewModel contains base informations related to layout
type LayoutViewModel struct {
	Platform string
	Logo     string
}

// SignedInUserViewModel contains informations related to signed in user
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

// Paging represents the paging model
type Paging struct {
	CurrentPage    int
	PreviousPage   int
	NextPage       int
	IsFinalPage    bool
	TotalPageCount int
}
