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
	ID           int
	UserName     string
	FullName     string
	Email        string
	RegisteredOn time.Time
	Password     string
	Website      string
	About        string
	InviteCode   string
	Karma        int
	CustomerID   int
}

/*UserError contains the error and user data which caused to error*/
type UserError struct {
	Message       string
	User          *User
	OriginalError error
}

func (err *UserError) Error() string {
	return fmt.Sprintf(
		"%s | OriginalError: %v | User: %+v",
		err.Message,
		err.OriginalError,
		err.User)
}

/*CreateUser creates a user*/
func CreateUser(user *User) (*int, error) {
	db, err := connectToDB()
	if err != nil {
		return nil, &UserError{"Cannot connect to db", user, err}
	}
	defer db.Close()
	sql := "INSERT INTO users (username, fullname, email, registeredon," + "password, website, about, invitecode, karma, customerid) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id"

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, &UserError{"Cannot encrypt user password", user, err}
	}
	var lastInsertedID int
	err = db.QueryRow(
		sql,
		user.UserName,
		user.FullName,
		user.Email,
		user.RegisteredOn,
		string(encryptedPassword),
		user.Website,
		user.About,
		user.InviteCode,
		user.Karma,
		user.CustomerID).Scan(&lastInsertedID)
	if err != nil {
		return nil, &UserError{"Cannot insert user to the database!", user, err}
	}
	userID := int(lastInsertedID)
	return &userID, nil
}

/*ChangePassword changes user password associated with provided user id*/
func ChangePassword(userID int, newPassword string) error {
	db, err := connectToDB()
	if err != nil {
		return err
	}
	defer db.Close()
	sql := "UPDATE users SET password = $1 WHERE Id = $2"
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return &DBError{fmt.Sprintf("Cannot encrypt new password. UserId: %d", userID), err}
	}
	_, err = db.Exec(sql, string(encryptedPassword), userID)
	if err != nil {
		return &DBError{fmt.Sprintf("Cannot update user's new password. UserId: %d", userID), err}
	}
	return nil
}

/*ConfirmPasswordMatch checks whether provided password are equal to user's password*/
func ConfirmPasswordMatch(userID int, password string) (matched bool, err error) {
	matched = false
	db, err := connectToDB()
	if err != nil {
		return matched, err
	}
	defer db.Close()
	sql := "SELECT password FROM users WHERE Id = $1"
	row := db.QueryRow(sql, userID)
	var userPassword string
	err = row.Scan(&userPassword)
	if err != nil {
		return matched, &DBError{fmt.Sprintf("Cannot read password for password match. UserID: %v", userID), err}
	}
	ok := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password))
	matched = ok == nil
	return matched, nil
}

/*UpdateUser updates the provided user on database*/
func UpdateUser(user *User) error {
	db, err := connectToDB()
	if err != nil {
		return &UserError{"Db connection error", user, err}
	}
	defer db.Close()
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return &UserError{"Cannot encrypt passoword", user, err}
	}
	sql := "UPDATE users SET username = $1, fullName = $2, email = $3, password = $4, website = $5, about = $6 WHERE Id = $7"
	_, err = db.Exec(
		sql,
		user.UserName,
		user.FullName,
		user.Email,
		string(encryptedPassword),
		user.Website,
		user.About,
		user.ID)
	if err != nil {
		return &UserError{"Cannot update user!", user, err}
	}
	return nil
}

/*ExistsUserByEmail checks if user associated with email exists on database*/
func ExistsUserByEmail(email string) (exists bool, err error) {
	exists = false
	db, err := connectToDB()
	if err != nil {
		return exists, err
	}
	defer db.Close()
	sql := "SELECT COUNT(*) AS count FROM users WHERE email = $1"
	row := db.QueryRow(sql, email)
	recordCount := 0
	err = row.Scan(&recordCount)
	if err != nil {
		return exists, &DBError{fmt.Sprintf("Cannot read record count. Email: %s", email), err}
	}
	if recordCount > 0 {
		exists = true
	}
	return exists, nil
}

/*ExistsUserByUserName checks if user associated with user name exists on database*/
func ExistsUserByUserName(userName string) (exists bool, err error) {
	exists = false
	db, err := connectToDB()
	if err != nil {
		return exists, err
	}
	defer db.Close()
	sql := "SELECT COUNT(*) AS count FROM users WHERE username = $1"
	row := db.QueryRow(sql, userName)
	recordCount := 0
	err = row.Scan(&recordCount)
	if err != nil {
		return exists, &DBError{fmt.Sprintf("Cannot read record count. UserName: %s", userName), err}
	}
	if recordCount > 0 {
		exists = true
	}
	return exists, nil
}

/*GetUserByUserName gets user associated with user name from database*/
func GetUserByUserName(userName string) (user *User, err error) {
	db, err := connectToDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	sql := "SELECT id, username, fullname, email, registeredon, password, website, about, invitecode, karma, customerid FROM users WHERE username = $1"
	row := db.QueryRow(sql, userName)
	user, err = MapSQLRowToUser(row)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot read user by user name from db. UserName: %s", userName), err}
	}
	return user, nil
}

/*GetUserByID gets user associated with user id from database*/
func GetUserByID(userID int) (user *User, err error) {
	db, err := connectToDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	sql := "SELECT id, username, fullname, email, registeredon, password, website, about, invitecode, karma, customerid FROM users WHERE id = $1"
	row := db.QueryRow(sql, userID)
	user, err = MapSQLRowToUser(row)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot read user by user name from db. UserID: %d", userID), err}
	}
	return user, nil
}

/*GetUsersByCustomerID retunrs users list by provided customerID parameter*/
func GetUsersByCustomerID(customerID int) (*[]User, error) {
	db, err := connectToDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	sql := "SELECT id, username, fullname, email, registeredon, password, website, about, invitecode, karma, customerid FROM users WHERE customerID = $1"
	rows, err := db.Query(sql, customerID)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot get users. CustomerID: %d", customerID), err}
	}
	users, err := MapSQLRowsToUsers(rows)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot read rows. CustomerID: %d", customerID), err}
	}
	return users, nil
}

/*SaveResetPasswordToken inserts given token and user id to database*/
func SaveResetPasswordToken(token string, userID int) error {
	db, err := connectToDB()
	if err != nil {
		return err
	}
	defer db.Close()
	sql := "INSERT INTO resetpasswordtokens (token, userid) VALUES ($1, $2)"
	_, err = db.Exec(sql, token, userID)
	if err != nil {
		return &DBError{fmt.Sprintf("Cannot save reset password token. Token: %s, UserID: %d", token, userID), err}
	}
	return nil
}

/*GetUserByResetPasswordToken gets user associated with token from database*/
func GetUserByResetPasswordToken(token string) (user *User, err error) {
	db, err := connectToDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	sql := "SELECT id, username, fullname, email, registeredon, password, website, about,  invitecode, karma,customerid FROM users INNER JOIN resetpasswordtokens on users.Id = resetpasswordtokens.userid WHERE resetpasswordtokens.token = $1"
	row := db.QueryRow(sql, token)
	user, err = MapSQLRowToUser(row)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot read user by reset password token from db. Token: %s", token), err}
	}
	return user, nil
}

/*FindUserByEmailAndPassword returns user associated with email and password from database*/
func FindUserByEmailAndPassword(email string, password string) (user *User, err error) {
	db, err := connectToDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	sql := "SELECT id, username, fullname, email, registeredon, password, website, about, invitecode, karma, customerid FROM users WHERE email = $1"
	row := db.QueryRow(sql, email)
	user, err = MapSQLRowToUser(row)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot read user by email and password from db. Email: %s", email), err}
	}
	result := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if result != nil && result == bcrypt.ErrMismatchedHashAndPassword {
		return nil, nil
	}
	return user, nil
}

/*FindUserByUserNameAndPassword returns user associated with user name and password from database*/
func FindUserByUserNameAndPassword(userName string, password string) (user *User, err error) {
	db, err := connectToDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	sql := "SELECT id, username, fullname, email, registeredon, password, website, about, invitecode, karma, customerid FROM users WHERE username = $1"
	row := db.QueryRow(sql, userName)
	user, err = MapSQLRowToUser(row)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot read user by email and password from db. UserName: %s", userName), err}
	}
	result := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if result != nil && result == bcrypt.ErrMismatchedHashAndPassword {
		return nil, nil
	}
	return user, nil
}

/*GetUserNameByEmail returns username by email*/
func GetUserNameByEmail(email string) (string, error) {
	db, err := connectToDB()
	if err != nil {
		return "", err
	}
	defer db.Close()
	sql := "SELECT username FROM users where email =$1;"
	row := db.QueryRow(sql, email)
	var username string
	err = row.Scan(
		&username)
	if err != nil {
		return "", &DBError{fmt.Sprintf("Cannot read email by username from db. Email: %s", email), err}
	}
	return username, nil
}

/*GetCustomerDomainByUserName returns customer's domain by user name*/
func GetCustomerDomainByUserName(userName string) (string, error) {
	db, err := connectToDB()
	if err != nil {
		return "", err
	}
	defer db.Close()
	sql := "SELECT domain FROM customers INNER JOIN users ON users.customerid = customers.id WHERE users.username = $1 ;"
	row := db.QueryRow(sql, userName)
	var domain string
	err = row.Scan(
		&domain)
	if err != nil {
		return "", &DBError{fmt.Sprintf("Cannot read rows"), err}
	}
	if err != nil {
		return "", &DBError{fmt.Sprintf("Cannot read domain by username from db. Domain: %s", domain), err}
	}
	return domain, nil
}

/*IsUserAdmin return true, if user is admin*/
func IsUserAdmin(userID int) (isAdmin bool, err error) {
	isAdmin = false
	db, err := connectToDB()
	if err != nil {
		return isAdmin, err
	}
	defer db.Close()
	sql := "SELECT COUNT(*) AS count FROM customers INNER JOIN users ON users.email = customers.email WHERE users.id = $1"
	row := db.QueryRow(sql, userID)
	recordCount := 0
	err = row.Scan(&recordCount)
	if err != nil {
		return isAdmin, &DBError{fmt.Sprintf("Cannot read record count. UserID: %d", userID), err}
	}
	if recordCount > 0 {
		isAdmin = true
	}
	return isAdmin, nil
}
