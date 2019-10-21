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
	password = "*****"
	dbname   = "test"
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
	result, err := db.Exec(sql, "test", "test", "test@test.com", time.Now(), "111111", "test.com", "hakkimda", "test", "test", 12)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.RowsAffected())
}
