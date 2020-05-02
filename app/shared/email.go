package shared

import (
	"fmt"
	"net/smtp"
)

/*SetInviteMailBody combine parameters and return body for UserInviteMail*/
func SetInviteMailBody(to, userName, memo, inviteCode, domain string) string {
	content := ""
	content += "<p>Merhaba: " + to + "</p>"
	content += "<p>" + userName + " invited you to LinkWind.</p>"
	if memo != "" {
		content += "<p><i>Mesaj: " + memo + "</i></p>"
	}

	content += "<p>To join LinkWind, you can create an account by clicking the link below.</p>"
	content += "<p>https://" + domain + "/signup?invitecode=" + inviteCode + "</p>"

	return content
}

/*SendInvitemail send mail for invite to join*/
func SendInvitemail(mailAddress, memo, inviteCode, userName, domain string) error {
	pass := "...."
	from := "our smtp mail address"
	to := mailAddress
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + "[" + domain + "] You have been invited to join LinkWind\n"

	body := SetInviteMailBody(to, userName, memo, inviteCode, domain)
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
func SendResetPasswordMail(email, userName, domain, token string) error {
	pass := "...."
	from := "our smtp mail addres"
	to := email
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + "[" + domain + "] Reset Your Password\n"

	body := setResetPasswordMailBody(token, userName, domain)
	msg := []byte(subject + mime + "\n" + body)

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, msg)

	if err != nil {
		return fmt.Errorf("An error occured when send forgot password mail : %s", err)
	}
	return nil
}

func setResetPasswordMailBody(token, userName, domain string) string {
	content := ""
	fontColour := "blue"
	content += "<p>Hello </><font color=" + fontColour + ">" + userName + "</font>"
	content += "<p>You have requested a password renewal.</p>"
	content += "<p>You can reset your password by clicking the link below.</p>"
	content += "<p>If you did not make such a request, do not care about this message.</p>"

	url := "http://" + domain + "/set-new-password?token=" + token

	content += "<a href=" + url + ">" + url + "</a>"

	return content
}
