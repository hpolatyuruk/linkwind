package models

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"linkwind/app/shared"
	"strings"
)

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

/*AboutViewModel contains informations related to about page */
type AboutViewModel struct {
	BaseViewModel
}

/*SetLayout sets about page view model layout members.*/
func (model *AboutViewModel) SetLayout(platformName string, logo string) {
	model.Layout = generateLayoutViewModel(platformName, logo)
}

/*SetSignedInUser sets about page view model signed in user members.*/
func (model *AboutViewModel) SetSignedInUser(userClaims *shared.SignedInUserClaims) {
	if userClaims == nil {
		return
	}
	model.SignedInUser = generateSignedInUserViewModel(userClaims)
}

/*FAQViewModel contains informations related to faq page */
type FAQViewModel struct {
	BaseViewModel
}

/*SetLayout sets about page view model layout members.*/
func (model *FAQViewModel) SetLayout(platformName string, logo string) {
	model.Layout = generateLayoutViewModel(platformName, logo)
}

/*SetSignedInUser sets about page view model signed in user members.*/
func (model *FAQViewModel) SetSignedInUser(userClaims *shared.SignedInUserClaims) {
	if userClaims == nil {
		return
	}
	model.SignedInUser = generateSignedInUserViewModel(userClaims)
}

/*PrivacyViewModel contains informations related to privacy page */
type PrivacyViewModel struct {
	BaseViewModel
}

/*SetLayout sets about page view model layout members.*/
func (model *PrivacyViewModel) SetLayout(platformName string, logo string) {
	model.Layout = generateLayoutViewModel(platformName, logo)
}

/*SetSignedInUser sets about page view model signed in user members.*/
func (model *PrivacyViewModel) SetSignedInUser(userClaims *shared.SignedInUserClaims) {
	if userClaims == nil {
		return
	}
	model.SignedInUser = generateSignedInUserViewModel(userClaims)
}
