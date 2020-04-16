package data

import (
	"fmt"
	"time"

	"github.com/rs/xid"
)

/*InviteCodeInfo represents the invited code info.*/
type InviteCodeInfo struct {
	Code                string
	InviterUserID       int
	InvitedEmailAddress string
	Used                bool
	CreatedOn           time.Time
}

/*CreateInviteCode creates an invite code.*/
func CreateInviteCode(inviterUserID int, invitedEmail string) (string, error) {
	db, err := connectToDB()
	if err != nil {
		return "", &DBError{fmt.Sprintf("Cannot conenct to db to create new invite code. InviterUserID: %d, InvitedEmail: %s", inviterUserID, invitedEmail), err}
	}
	defer db.Close()
	inviteCode := xid.New().String()
	sql := "INSERT INTO invitecodes (code, inviteruserid, invitedemail, createdon) VALUES ($1, $2, $3, $4) RETURNING code"
	_, err = db.Exec(
		sql,
		inviteCode,
		inviterUserID,
		invitedEmail,
		time.Now())
	if err != nil {
		return "", &DBError{fmt.Sprintf("Cannot create a new intvite code. InviterUserID: %d, InvitedEmail: %s", inviterUserID, invitedEmail), err}
	}
	return inviteCode, nil
}

/*ExistsInviteCode checks whether invite code exists in user db or not*/
func ExistsInviteCode(inviteCode string) (exists bool, err error) {
	exists = false
	db, err := connectToDB()
	if err != nil {
		return exists, &DBError{fmt.Sprintf("Cannot conenct to db to if invite code is already existed. InviteCode: %s", inviteCode), err}
	}
	defer db.Close()
	sql := "SELECT COUNT(*) AS count FROM invitecodes WHERE code = $1"
	row := db.QueryRow(sql, inviteCode)
	recordCount := 0
	err = row.Scan(&recordCount)
	if err != nil {
		return exists, &DBError{fmt.Sprintf("Cannot get record count for inviteCode: %s", inviteCode), err}
	}
	if recordCount > 0 {
		exists = true
	}
	return exists, err
}

/*FindInviterEmailByInviteCode returns inviter email address by invite code.*/
func FindInviterEmailByInviteCode(inviteCode string) (string, error) {
	db, err := connectToDB()
	if err != nil {
		return "", err
	}
	defer db.Close()
	sql := "SELECT users.email FROM users INNER JOIN invitecodes ON users.id = invitecodes.userid where invitecodes.code =$1;"
	row := db.QueryRow(sql, inviteCode)
	var email string
	err = row.Scan(&email)
	if err != nil {
		return "", &DBError{fmt.Sprintf("Cannot read inviter semail by invite code from db. InviteCode: %s", inviteCode), err}
	}
	return email, nil
}

/*MarkInviteCodeAsUsed marks invite code as used*/
func MarkInviteCodeAsUsed(inviteCode string) error {
	db, err := connectToDB()
	if err != nil {
		return &DBError{fmt.Sprintf("Cannot connect to db to mark invite code as used. InviteCode: %s", inviteCode), err}
	}
	defer db.Close()
	sql := "UPDATE invitecodes SET used = true WHERE code = $1"
	_, err = db.Exec(
		sql,
		inviteCode)
	if err != nil {
		return &DBError{fmt.Sprintf("Cannot update invite code as used. InviteCode: %s", inviteCode), err}
	}
	return nil
}

/*IsInviteCodeUsed checks whether invite code is already used or not*/
func IsInviteCodeUsed(inviteCode string) (used bool, err error) {
	used = false
	db, err := connectToDB()
	if err != nil {
		return used, &DBError{fmt.Sprintf("Cannot connec to db to check if invite code is already used. InviteCode: %s", inviteCode), err}
	}
	defer db.Close()
	sql := "SELECT used FROM invitecodes WHERE code = $1"
	row := db.QueryRow(sql, inviteCode)
	err = row.Scan(&used)
	if err != nil {
		return used, &DBError{fmt.Sprintf("Cannot get record count for inviteCode: %s", inviteCode), err}
	}
	return used, err
}

/*GetInviteCodeInfoByCode gets invite code info form database.*/
func GetInviteCodeInfoByCode(inviteCode string) (*InviteCodeInfo, error) {
	db, err := connectToDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	sql := "SELECT code, inviteruserid, invitedemail, used, createdon FROM invitecodes WHERE code = $1"
	row := db.QueryRow(sql, inviteCode)
	inviteCodeInfo, err := MapSQLRowToInviteCodeInfo(row)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot read invite code info by code from db. InviteCode: %s", inviteCode), err}
	}
	return inviteCodeInfo, nil
}
