package model

import (
	"database/sql"
	"fmt"
)

type DbValue interface {
	ScanRow(row *sql.Rows) (bool)
}

type User struct {
	Username string
	Dob      string
	Location string
}

// func (u User) ScanRow(row *sql.Rows, index int, results *[]DbValue) bool {
func (u *User) ScanRow(row *sql.Rows) (success bool) {
	var username, location, dob string

	if scanError := row.Scan(&username, &location, &dob); scanError != nil {
		fmt.Println("Encountered an error when processing query results.")
		return false;
	}

	fmt.Printf("Query Result: {username: %v, date of birth: %v, location: %v}\n", username, location, dob);
	fmt.Printf("User item before assigning values: {%v}\n", u)
	fmt.Printf("Actual value of User item: {%v}\n", &u)

	u.Username = username;
	u.Location = location;
	u.Dob = dob;

	fmt.Printf("User item after assigning values: {%v}\n", u)
	fmt.Printf("Actual value of User item: {%v}\n", &u)

	return true;
}