package data

import (
	"fmt"
	"time"

	/* To install postgresql driver. Check more here: https://www.calhoun.io/why-we-import-sql-drivers-with-the-blank-identifier/ */
	_ "github.com/lib/pq"
)

/*Customer represents the customer in database*/
type Customer struct {
	ID           int
	Email        string
	Name         string
	Description  string
	RegisteredOn time.Time
	Domain       string
}

/*CustomerError contains the error and customer data which caused to error*/
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

/*CreateCustomer creates a customer*/
func CreateCustomer(customer *Customer) (err error) {
	db, err := connectToDB()
	defer db.Close()
	if err != nil {
		return &CustomerError{"Cannot connect to db", customer, err}
	}
	sql := "INSERT INTO customers (email, name, description, domain,registeredon)" +
		"VALUES ($1, $2, $3, $4, $5)"

	_, err = db.Exec(
		sql,
		customer.Email,
		customer.Name,
		customer.Description,
		customer.Domain,
		customer.RegisteredOn)
	if err != nil {
		return &CustomerError{"Cannot insert customer to the database!", customer, err}
	}
	return nil
}

/*UpdateCustomer updates the provided customer on database*/
func UpdateCustomer(customer *Customer) error {
	db, err := connectToDB()
	defer db.Close()
	if err != nil {
		return &CustomerError{"Db connection error", customer, err}
	}

	sql := "UPDATE customers SET name = $1, email = $2, description = $3, domain = $4, registeredon = $5 WHERE id = $6"
	_, err = db.Exec(
		sql,
		customer.Name,
		customer.Email,
		customer.Description,
		customer.Domain,
		customer.RegisteredOn,
		customer.ID)
	if err != nil {
		return &CustomerError{"Cannot update customer!", customer, err}
	}
	return nil
}
