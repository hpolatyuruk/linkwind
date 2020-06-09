package data

import (
	"database/sql"
	"fmt"
	"time"

	/* To install postgresql driver. Check more here: https://www.calhoun.io/why-we-import-sql-drivers-with-the-blank-identifier/ */
	_ "github.com/lib/pq"
)

const (
	/*CustomerDefaultDomain represents the default value for null csutomer domain. Because we cannot make string as nil*/
	CustomerDefaultDomain string = ""
)

/*Customer represents the customer in database*/
type Customer struct {
	ID           int
	Email        string
	Name         string
	RegisteredOn time.Time
	Domain       string
	LogoImage    []byte
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
	if err != nil {
		return &CustomerError{"Cannot connect to db", customer, err}
	}
	defer db.Close()
	query := "INSERT INTO customers (email, name, domain, registeredon) VALUES ($1, $2, $3, $4)"
	_, err = db.Exec(
		query,
		customer.Email,
		customer.Name,
		nullCustomerDomain(customer.Domain),
		customer.RegisteredOn)
	if err != nil {
		return &CustomerError{"Cannot insert customer to the database!", customer, err}
	}
	return nil
}

/*UpdateCustomer updates the provided customer on database*/
func UpdateCustomer(customer *Customer) error {
	db, err := connectToDB()
	if err != nil {
		return &CustomerError{"Db connection error", customer, err}
	}
	defer db.Close()
	sql := "UPDATE customers SET name = $1, email = $2, domain = $3, registeredon = $4, imglogo = $5 WHERE id = $6"
	_, err = db.Exec(
		sql,
		customer.Name,
		customer.Email,
		nullCustomerDomain(customer.Domain),
		customer.RegisteredOn,
		customer.LogoImage,
		customer.ID)
	if err != nil {
		return &CustomerError{"Cannot update customer!", customer, err}
	}
	return nil
}

/*ExistsCustomerByName check if customer associated with name exists on database*/
func ExistsCustomerByName(name string) (exists bool, err error) {
	exists = false
	db, err := connectToDB()
	if err != nil {
		return exists, err
	}
	defer db.Close()
	sql := "SELECT COUNT(*) AS count FROM customers WHERE name = $1"
	row := db.QueryRow(sql, name)
	recordCount := 0
	err = row.Scan(&recordCount)
	if err != nil {
		return exists, &DBError{fmt.Sprintf("Cannot read record count. Name: %s", name), err}
	}
	if recordCount > 0 {
		exists = true
	}
	return exists, nil
}

/*ExistsCustomerByEmail check if customer associated with email exists on database*/
func ExistsCustomerByEmail(email string) (exists bool, err error) {
	exists = false
	db, err := connectToDB()
	if err != nil {
		return exists, err
	}
	defer db.Close()
	sql := "SELECT COUNT(*) AS count FROM customers WHERE email = $1"
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

/*ExistsCustomerByDomain check if customer associated with domain exists on database*/
func ExistsCustomerByDomain(domain string) (exists bool, err error) {
	exists = false
	db, err := connectToDB()
	if err != nil {
		return exists, err
	}
	defer db.Close()
	sql := "SELECT COUNT(*) AS count FROM customers WHERE domain = $1"
	row := db.QueryRow(sql, domain)
	recordCount := 0
	err = row.Scan(&recordCount)
	if err != nil {
		return exists, &DBError{fmt.Sprintf("Cannot read record count. Domain: %s", domain), err}
	}
	if recordCount > 0 {
		exists = true
	}
	return exists, nil
}

/*GetCustomerByName gets customer associated with name from database*/
func GetCustomerByName(name string) (customer *Customer, err error) {
	db, err := connectToDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	query := "SELECT id, name, email, registeredon, domain, imglogo FROM customers WHERE name = $1"
	row := db.QueryRow(query, name)
	customer, err = MapSQLRowToCustomer(row)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot read customer by name from db. Name: %s", name), err}
	}
	return customer, nil
}

/*GetCustomerByID gets customer associated with id from database*/
func GetCustomerByID(id int) (customer *Customer, err error) {
	db, err := connectToDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	sql := "SELECT id, name, email, registeredon, domain, imglogo FROM customers WHERE id = $1"
	row := db.QueryRow(sql, id)
	customer, err = MapSQLRowToCustomer(row)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot read customer by id from db. ID: %d", id), err}
	}
	return customer, nil
}

/*GetCustomerByDomain gets customer associated with domain from database*/
func GetCustomerByDomain(domain string) (customer *Customer, err error) {
	db, err := connectToDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	query := "SELECT id, name, email, registeredon, domain, imglogo FROM customers WHERE domain = $1"
	row := db.QueryRow(query, domain)
	customer, err = MapSQLRowToCustomer(row)
	if err != nil {
		return nil, &DBError{fmt.Sprintf("Cannot read customer by domain from db. Name: %s", domain), err}
	}
	return customer, nil
}

func nullCustomerDomain(domain string) sql.NullString {
	if domain == "" {
		return sql.NullString{}
	}
	return sql.NullString{
		String: domain,
		Valid:  true,
	}
}
