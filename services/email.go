package services

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

/*SetInviteMailBody combine parameters and return body for UserInviteMail*/
func SetInviteMailBody(to, userName, memo, inviteCode string) string {
	content := ""
	content += "<p>Merhaba: " + to + "</p>"
	content += "<p>" + userName + " adlı kullanıcı sizi TurkDev'e davet etti.</p>"
	if memo != "" {
		content += "<p><i>Mesaj: " + memo + "</i></p>"
	}

	content += "<p>TurkDev'e katılmak için aşağıdaki linke tıklayarak hesap oluşturabilirsiniz.</p>"
	content += "<p>https://turkdev.com/davet/" + inviteCode + "</p>"

	return content
}

/*SendInvitemail send mail for invite to join*/
/*TODO : These configurations are not perminant. These conf for gmail.We should add pass and etc*/
func SendInvitemail(mailAddress, memo, inviteCode, userName string) error {
	pass := "...."
	from := "our smtp mail adrress"
	to := mailAddress
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + "TurkDev'e katılmaya davet edildiniz\n"

	body := SetInviteMailBody(to, userName, memo, inviteCodeGenerator())
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
//TODO: In lobsters they add coming ip for reset pass request. Should we do that? Do not forget to change "pass" and "to" variables.
func SendResetPasswordMail(email, userName, domain string) error {
	pass := "..."
	from := "..."
	to := email
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + "[" + domain + "] Reset Your Password\n"

	token := generateResetPasswordToken()
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

func generateResetPasswordToken() string {
	c := ""
	for i := 0; i < 4; i++ {
		s := stringWithCharset(1)
		i := seededRand.Intn(10)
		c = c + strconv.Itoa(i) + s
	}
	return c
}

func inviteCodeGenerator() string {
	c := ""
	for i := 0; i < 4; i++ {
		i := seededRand.Intn(10)
		s := stringWithCharset(1)
		c = c + strconv.Itoa(i) + s
	}
	return c
}

func stringWithCharset(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
