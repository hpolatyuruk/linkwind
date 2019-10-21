package data

import (
	"database/sql"
	"fmt"
	"time"

	/* To install postgresql driver. Check more here: https://www.calhoun.io/why-we-import-sql-drivers-with-the-blank-identifier/ */
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "hgjh55_FFF"
	dbname   = "turkdev"
)

// CreateUser: Creates a user
func CreateUser() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	defer db.Close()
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	sql := "INSERT INTO users (username, fullname, email, registeredon, password, website, about, invitedby, invitecode, karma) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
	result, err := db.Exec(sql, "hpy", "huseyin polat yuruk", "h.polatyuruk@gmail.com", time.Now(), "111111", "huseyinpolatyuruk.com", "hakkimda", "anil", "apt123", 12)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.RowsAffected())
}
