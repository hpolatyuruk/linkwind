package data

import (
	"fmt"
	"time"

	/* To install postgresql driver. Check more here: https://www.calhoun.io/why-we-import-sql-drivers-with-the-blank-identifier/ */
	_ "github.com/lib/pq"
)

/*Customer represents the user in database*/
type Customer struct {
	ID           int
	Email        string
	Name         string
	Description  string
	RegisteredOn time.Time
	Domain       string
}

/*UserError contains the error and user data which caused to error*/
type CustomerError struct {
	Message       string
	Customer      *Customer
	OriginalError error
}

func (err *CustomerError) Error() string {
	return fmt.Sprintf(
		"%s | OriginalError: %v | Customer: %+v",
		err.Message,
		err.OriginalError,
		err.Customer)
}
