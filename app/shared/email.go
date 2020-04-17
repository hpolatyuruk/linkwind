package shared

import (
	"fmt"
	"net/smtp"
	"strconv"
)

/*SetInviteMailBody combine parameters and return body for UserInviteMail*/
func SetInviteMailBody(to, userName, memo, inviteCode string) string {
	content := ""
	content += "<p>Merhaba: " + to + "</p>"
	content += "<p>" + userName + " adlı kullanıcı sizi LinkWind'e davet etti.</p>"
	if memo != "" {
		content += "<p><i>Mesaj: " + memo + "</i></p>"
	}

	content += "<p>LinkWind'e katılmak için aşağıdaki linke tıklayarak hesap oluşturabilirsiniz.</p>"
	content += "<p>https://linkwind.co/davet/" + inviteCode + "</p>"

	return content
}

/*SendInvitemail send mail for invite to join*/
/*TODO : These configurations are not perminant. These conf for gmail.We should add pass and etc*/
func SendInvitemail(mailAddress, memo, inviteCode, userName string) error {
	pass := "...."
	from := "our smtp mail adrress"
	to := mailAddress
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + "LinkWind'e katılmaya davet edildiniz\n"

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
func SendResetPasswordMail(email, userName, domain, token string) error {
	pass := "Sedat.1242"
	from := "sedata38@gmail.com"
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

func inviteCodeGenerator() string {
	c := ""
	for i := 0; i < 4; i++ {
		i := SeededRand.Intn(10)
		s := StringWithCharset(1)
		c = c + strconv.Itoa(i) + s
	}
	return c
}
