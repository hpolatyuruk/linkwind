package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

/*DBError represents the database error*/
type DBError struct {
	Message       string
	OriginalError error
}

func (err *DBError) Error() string {
	return fmt.Sprintf("%s | Args: | OriginalError: %v", err.Message, err.OriginalError)
}

func connectionString() (conStr string, err error) {
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading staging.env file. Error: %v", err)
	}

	host := os.Getenv("dbHost")
	port, err := strconv.Atoi(os.Getenv("dbPort"))
	if err != nil {
		log.Fatalf("Cannot parse db port. Error: %v", err)
	}
	user := os.Getenv("dbUser")
	password := os.Getenv("dbPassword")
	name := os.Getenv("dbName")

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
