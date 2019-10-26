package data

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	/* To install postgresql driver. Check more here: https://www.calhoun.io/why-we-import-sql-drivers-with-the-blank-identifier/ */
	_ "github.com/lib/pq"
)

/*User represents the user in database*/
type User struct {
	UserName     string
	FullName     string
	Email        string
	RegisteredOn time.Time
	Password     string
	Website      string
	About        string
	Invitedby    string
	InviteCode   string
	Karma        int
}

/*UserError contains the error and user data which caused to error*/
type UserError struct {
	Msg         string
	User        *User
	OriginalErr error
}

func (err *UserError) Error() string {
	return fmt.Sprintf(
		"%s | OriginalError: %v | User: %+v",
		err.Msg,
		err.OriginalErr,
		err.User)
}

/*CreateUser creates a user*/
func CreateUser(user *User) (err error) {
	db, err := connectToDB()
	defer db.Close()
	if err != nil {
		return &UserError{"Cannot connect to db", user, err}
	}
	sql := "INSERT INTO users (username, fullname, email, registeredon," + "password, website, about, invitedby, invitecode, karma) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return &UserError{"Cannot encrypt user password", user, err}
	}
	_, err = db.Exec(
		sql,
		user.UserName,
		user.FullName,
		user.Email,
		user.RegisteredOn,
		string(encryptedPassword),
		user.Website,
		user.About,
		user.Invitedby,
		user.Invitedby,
		user.Karma)
	if err != nil {
		return &UserError{"Cannot insert user to the database!", user, err}
	}
	return nil
}

/*ExistsInviteCode checks whether invite code exists in user db or not*/
func ExistsInviteCode(inviteCode string) (exists bool, err error) {
	exists = false
	db, err := connectToDB()
	defer db.Close()
	if err != nil {
		return exists, err
	}
	sql := "SELECT COUNT(*) AS count FROM users WHERE invitecode = $1"
	row := db.QueryRow(sql, inviteCode)
	var recordCount int = 0
	err = row.Scan(&recordCount)
	if err != nil {
		return exists, &DBError{fmt.Sprintf("Cannot get record count for inviteCode: %s", inviteCode), err}
	}
	if recordCount > 0 {
		exists = true
	}
	return exists, err
}
