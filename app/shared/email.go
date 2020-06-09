package shared

import (
	"fmt"
	"net/smtp"
	"os"
	"strconv"
)

/*InviteMailInfo represents InviteMail parameters*/
type InviteMailInfo struct {
	Domain     *string
	InviteCode string
	Email      string
	UserName   string
	Memo       string
	Platform   string
}

// SMTPConfig contains SMTP server and credentials
type SMTPConfig struct {
	Server   string
	Port     int
	UserName string
	Password string
	Address  string
}

/*ResetPasswordMailInfo represents ResetPasswordMail parameters*/
type ResetPasswordMailInfo struct {
	Email    string
	UserName string
	Domain   string
	Platform string
	Token    string
}

/*SetInviteMailBody combine parameters and return body for UserInviteMail*/
func SetInviteMailBody(m InviteMailInfo, platformName string) string {
	content := ""
	content += "<p>Hello: " + m.Email + "</p>"
	content += "<p>" + m.UserName + " invited you to " + platformName + ".</p>"
	if m.Memo != "" {
		content += "<p><i>Mesaj: " + m.Memo + "</i></p>"
	}

	domain := platformName + ".linkwind.co"
	if m.Domain != nil {
		domain = *m.Domain
	}

	content += "<p>To join " + platformName + ", you can create an account by clicking the link below.</p>"
	content += "<p>" + domain + "/signup?invitecode=" + m.InviteCode + "</p>"

	return content
}

/*SendEmailInvitation send mail for invite to join*/
func SendEmailInvitation(m InviteMailInfo) error {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	platformName := m.Platform
	if m.Domain != nil {
		platformName = *m.Domain
	}

	subject := "Subject: " + "[" + platformName + "] You have been invited to join " + platformName + "\n"

	body := SetInviteMailBody(m, platformName)
	msg := []byte(subject + mime + "\n" + body)

	from := "www.linkwind.co@gmail.com"
	to := m.Email
	smtpConfig := getSMTPConfig()

	err := smtp.SendMail(
		smtpConfig.Address,
		smtp.PlainAuth(
			"",
			smtpConfig.UserName,
			smtpConfig.Password,
			smtpConfig.Server),
		from,
		[]string{to},
		msg)

	if err != nil {
		return err
	}
	return nil
}

/*SendResetPasswordMail send to mail for reset password with resetPassword token*/
func SendResetPasswordMail(r ResetPasswordMailInfo) error {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	platformName := r.Platform
	if r.Domain != "" {
		platformName = r.Domain
	}

	subject := "Subject: " + "[" + platformName + "] Reset Your Password\n"

	body := generateResetPasswordMailBody(r, platformName)
	msg := []byte(subject + mime + "\n" + body)

	from := "www.linkwind.co@gmail.com"
	to := r.Email
	smtpConfig := getSMTPConfig()

	err := smtp.SendMail(
		smtpConfig.Address,
		smtp.PlainAuth(
			"",
			smtpConfig.UserName,
			smtpConfig.Password,
			smtpConfig.Server),
		from,
		[]string{to},
		msg)

	if err != nil {
		return fmt.Errorf("An error occured when send forgot password mail : %s", err)
	}
	return nil
}

func generateResetPasswordMailBody(r ResetPasswordMailInfo, platformName string) string {
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

func getSMTPConfig() *SMTPConfig {
	server := os.Getenv("SMTP_SERVER")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	userName := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	return &SMTPConfig{
		Server:   server,
		Port:     port,
		UserName: userName,
		Password: password,
		Address:  fmt.Sprintf("%s:%d", server, port),
	}
}
