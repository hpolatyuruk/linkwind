package data

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
)

const (
	slash              = "`"
	regexForEmailValid = `^(?("")("".+?(?<!\\)""@)|(([0-9a-z]((\.(?!\.))|[-!#\$%&'\*\+/=\?\^` + slash + `\{\}\|~\w])*)(?<=[0-9a-z])@))(?(\[)(\[(\d{1,3}\.){3}\d{1,3}\])|(([0-9a-z][-\w]*[0-9a-z]*\.)+[a-z0-9][\-a-z0-9]{0,22}[a-z0-9]))$`
	charset            = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

/*DBError represents the database error*/
type DBError struct {
	Message       string
	OriginalError error
}

func (err *DBError) Error() string {
	return fmt.Sprintf("%s | OriginalError: %v", err.Message, err.OriginalError)
}

func connectionString() (conStr string, err error) {
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading staging.env file. Error: %v", err)
	}

	host := os.Getenv("DBHost")
	port, err := strconv.Atoi(os.Getenv("DBPort"))
	if err != nil {
		log.Fatalf("Cannot parse db port. Error: %v", err)
	}
	user := os.Getenv("DBUser")
	password := os.Getenv("DBPassword")
	name := os.Getenv("DBName")

	conStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, name)

	return conStr, err
}

func connectToDB() (db *sql.DB, err error) {
	connStr, err := connectionString()
	if err != nil {
		return nil, &DBError{"Cannot read connection string!", err}
	}
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, &DBError{"Cannot open db connection!", err}
	}
	err = db.Ping()
	if err != nil {
		return nil, &DBError{"Cannot ping db!", err}
	}
	return db, err
}

/*MapSQLRowToUser creates an user struct object by sql row*/
func MapSQLRowToUser(row *sql.Row) (user *User, err error) {
	var _user User
	err = row.Scan(
		&_user.ID,
		&_user.UserName,
		&_user.FullName,
		&_user.Email,
		&_user.RegisteredOn,
		&_user.Password,
		&_user.Website,
		&_user.About,
		&_user.InvitedBy,
		&_user.InviteCode,
		&_user.Karma)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot map sql row to user struct"), err}
	}
	user = &_user
	return user, nil
}

/*MapSQLRowsToStories creates a story struct array by sql rows*/
func MapSQLRowsToStories(rows *sql.Rows) (stories *[]Story, err error) {
	var _stories []Story = []Story{}
	var username string
	for rows.Next() {
		var story Story
		err = rows.Scan(
			&story.ID,
			&story.URL,
			&story.Title,
			&story.Text,
			&story.UpVotes,
			&story.DownVotes,
			&story.CommentCount,
			&story.UserID,
			&story.SubmittedOn,
			pq.Array(&story.Tags),
			&username)
		if err != nil {
			return nil, &DBError{fmt.Sprintf("Cannot read rows"), err}
		}
		story.UserName = username
		_stories = append(_stories, story)
	}
	return &_stories, nil
}

/*MapSQLRowsToComments creates a comment struct array by sql rows*/
func MapSQLRowsToComments(rows *sql.Rows) (comments *[]Comment, err error) {
	var _comments []Comment = []Comment{}
	for rows.Next() {
		var comment Comment = Comment{}
		var parentID sql.NullInt32
		var userName string
		err = rows.Scan(
			&comment.Comment,
			&comment.UpVotes,
			&comment.DownVotes,
			&comment.StoryID,
			&parentID,
			&comment.ReplyCount,
			&comment.UserID,
			&comment.CommentedOn,
			&comment.ID,
			&userName)
		if err != nil {
			return nil, &DBError{"Cannot read comment row.", err}
		}
		if parentID.Valid {
			comment.ParentID = int(parentID.Int32)
		} else {
			comment.ParentID = CommentRootID
		}
		comment.UserName = userName
		_comments = append(_comments, comment)
	}
	return &_comments, nil
}

/*MapSQLRowsToReplies creates a reply struct array by sql rows*/
func MapSQLRowsToReplies(rows *sql.Rows) (replies *[]Reply, err error) {
	var _replies []Reply = []Reply{}
	for rows.Next() {
		var reply Reply = Reply{}
		var comment Comment = Comment{}
		var storyID int
		var storyTitle string
		var userName string
		var parentID sql.NullInt32
		err = rows.Scan(
			&comment.Comment,
			&comment.UpVotes,
			&comment.DownVotes,
			&comment.StoryID,
			&parentID,
			&comment.ReplyCount,
			&comment.UserID,
			&comment.CommentedOn,
			&comment.ID,
			&storyTitle,
			&storyID,
			&userName)
		if err != nil {
			return nil, &DBError{"Cannot read comment row.", err}
		}
		if parentID.Valid {
			comment.ParentID = int(parentID.Int32)
		} else {
			comment.ParentID = CommentRootID
		}
		reply.Comment = &comment
		reply.StoryID = storyID
		reply.StoryTitle = storyTitle
		reply.UserName = userName
		_replies = append(_replies, reply)
	}
	return &_replies, nil
}

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

/*SetResetPasswordMailBody set mail's body with resetToken*/
func SetResetPasswordMailBody(token string) string {

	content := ""
	content += "<p>Şifre yenileme isteğinde bulunduz.</p>"
	content += "<p>Aşağıdaki linke tıklayarak şifrenizi sıfırlayabilirsiniz.</p>"
	content += "<p>Böyle bir istekte bulunmadıysanız, bu mesajı önemsemeyin.</p>"
	content += "<p>https://turkdev.com/login/set_new_password?token=" + token + "</p>"

	return content
}

/*IsEmailAdrressValid take mail adrres, if adrress is valid return true.*/
func IsEmailAdrressValid(email string) bool {
	Re := regexp.MustCompile(regexForEmailValid)
	return Re.MatchString(email)
}

/*InviteCodeGenerator generates the invite code*/
func InviteCodeGenerator() string {

	c := ""
	for i := 0; i < 4; i++ {
		i := seededRand.Intn(10)
		s := stringWithCharset(1)
		c = c + strconv.Itoa(i) + s
	}
	return c
}

/*GenerateResetPasswordToken generates the token for password reset*/
func GenerateResetPasswordToken() string {

	c := ""
	for i := 0; i < 4; i++ {
		s := stringWithCharset(1)
		i := seededRand.Intn(10)
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

/*Send request to url that given as parameter and fetch title from HTML code.*/
func FetchURL(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	dataInBytes, err := ioutil.ReadAll(response.Body)
	pageContent := string(dataInBytes)

	titleStartIndex := strings.Index(pageContent, "<title>")
	if titleStartIndex == -1 {
		return "", fmt.Errorf("No title element found")
	}
	// The start index of the title is the index of the first
	// character, the < symbol. We don't want to include
	// <title> as part of the final value, so let's offset
	// the index by the number of characers in <title>
	titleStartIndex += 7

	titleEndIndex := strings.Index(pageContent, "</title>")
	if titleEndIndex == -1 {
		fmt.Println("No closing tag for title found.")
		os.Exit(0)
	}

	pageTitle := []byte(pageContent[titleStartIndex:titleEndIndex])
	return string(pageTitle), nil
}

func setCookie(w http.ResponseWriter, userNameOrEmail, password string) {
	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{
		Name:    userNameOrEmail,
		Value:   password,
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
}

func LoginUser(emailOrUserName, password string) (int, error) {

	if IsEmailAdrressValid(emailOrUserName) {
		exists, err := ExistsUserByEmail(emailOrUserName)
		if err != nil {
			return LoginError, err
		}

		if exists {
			user, err := FindUserByEmailAndPassword(emailOrUserName, password)
			if err != nil {
				return LoginError, err
			}
			if user == nil {
				return WrongPassword, nil
			}
			return LoginSuccessful, nil
		}
		return NoUserWithEmail, nil
	}

	exists, err := ExistsUserByUserName(emailOrUserName)
	if err != nil {
		return LoginError, err
	}

	if exists {
		user, err := FindUserByUserNameAndPassword(emailOrUserName, password)
		if err != nil {
			return LoginError, err
		}

		if user == nil {
			return WrongPassword, nil
		}
		return LoginSuccessful, nil
	}
	return NoUserWithUserName, nil
}

/*SendInvitemail send mail for invite to join*/
/*TODO : These configurations are not perminant. These conf for gmail.We should add pass and etc*/
func SendInvitemail(mailAddress, memo, inviteCode, userName string) {

	pass := "...."
	from := "our smtp mail adrress"
	to := mailAddress
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + "TurkDev'e katılmaya davet edildiniz\n"

	body := SetInviteMailBody(to, userName, memo, InviteCodeGenerator())
	msg := []byte(subject + mime + "\n" + body)

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, msg)

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}

/*SendForgotPasswordMail send to mail for reset password with resetPassword token*/
//TODO: In lobsters they add coming ip for reset pass request. Should we do that? Do not forget to change "pass" and "to" variables.
func SendForgotPasswordMail(mailAddress string) {

	pass := "..."
	from := "our smtp mail adrress"
	to := mailAddress
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + "Şifre Sıfırlama\n"

	token := GenerateResetPasswordToken()
	body := SetResetPasswordMailBody(token)
	msg := []byte(subject + mime + "\n" + body)

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, msg)

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}

func CalculateKarma(userID int) int {

	stories, err := GetUserStoriesNotPaging(userID)
	sVotes := 0
	for _, s := range *stories {
		sVotes += (s.UpVotes - s.DownVotes)
	}

	comments, err := GetUserCommentsNotPaging(userID)
	cVotes := 0
	for _, c := range *comments {
		cVotes += (c.UpVotes - c.DownVotes)
	}

	return sVotes + cVotes
}
