package shared

import (
	"fmt"
	"net/smtp"
)

type InviteMailQuery struct {
	Domain     string
	InviteCode string
	Email      string
	UserName   string
	Memo       string
	Platform   string
}

type ResetPasswordMailQuery struct {
	Email    string
	UserName string
	Domain   string
	Platform string
	Token    string
}

/*SetInviteMailBody combine parameters and return body for UserInviteMail*/
func SetInviteMailBody(m InviteMailQuery, platformName string) string {
	content := ""
	content += "<p>Hello: " + m.Email + "</p>"
	content += "<p>" + m.UserName + " invited you to " + platformName + ".</p>"
	if m.Memo != "" {
		content += "<p><i>Mesaj: " + m.Memo + "</i></p>"
	}

	domain := platformName + ".linkwind.co"
	if m.Domain != "" {
		domain = m.Domain
	}

	content += "<p>To join " + platformName + ", you can create an account by clicking the link below.</p>"
	content += "<p>" + domain + "/signup?invitecode=" + m.InviteCode + "</p>"

	return content
}

/*SendInvitemail send mail for invite to join*/
func SendInvitemail(m InviteMailQuery) error {
	pass := "Sedat.1242"
	from := "sedata38@gmail.com"
	to := m.Email
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	platformName := m.Platform
	if m.Domain != "" {
		platformName = m.Domain
	}

	subject := "Subject: " + "[" + platformName + "] You have been invited to join " + platformName + "\n"

	body := SetInviteMailBody(m, platformName)
	msg := []byte(subject + mime + "\n" + body)

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, msg)

	if err != nil {
		return err
	}
	return nil
}

/*SendResetPasswordMail send to mail for reset password with resetPassword token*/
func SendResetPasswordMail(r ResetPasswordMailQuery) error {
	pass := "Sedat.1242"
	from := "sedata38@gmail.com"
	to := r.Email
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	platformName := r.Platform
	if r.Domain != "" {
		platformName = r.Domain
	}

	subject := "Subject: " + "[" + platformName + "] Reset Your Password\n"

	body := setResetPasswordMailBody(r, platformName)
	msg := []byte(subject + mime + "\n" + body)

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, msg)

	if err != nil {
		return fmt.Errorf("An error occured when send forgot password mail : %s", err)
	}
	return nil
}

func setResetPasswordMailBody(r ResetPasswordMailQuery, platformName string) string {
	content := ""
	fontColour := "blue"
	content += "<p>Hello </><font color=" + fontColour + ">" + r.UserName + "</font>"
	content += "<p>You have requested a password renewal.</p>"
	content += "<p>You can reset your password by clicking the link below.</p>"
	content += "<p>If you did not make such a request, do not care about this message.</p>"

	domain := platformName + ".linkwind.co"
	if r.Domain != "" {
		domain = r.Domain
	}

	url := domain + "/set-new-password?token=" + r.Token

	content += "<a href=" + url + ">" + url + "</a>"

	return content
}
