package model

import (
	"database/sql"
	"fmt"
)

type DbValue interface {
	ScanRow(row *sql.Rows, index int, results *[]DbValue) (bool)
}

type User struct {
	Username string
	Dob      string
	Location string
}

func (u User) ScanRow(row *sql.Rows, index int, results *[]DbValue) bool {
	userValue := new(User);
	var username, location, dob string
	if scanError := row.Scan(&username, &location, &dob); scanError != nil {
		fmt.Println("Encountered an error when processing query results.")
		return false;
	}

	fmt.Printf("Query Result: {username: %v, date of birth: %v, location: %v}\n", username, location, dob);

	userValue.Username = username;
	userValue.Location = location;
	userValue.Dob = dob;

	if index >= len(*results) {
		*results = append(*results, userValue);
		fmt.Printf("Appending here! The new length of the slice is: %v\n", len(*results))
	} else {
		(*results)[index] = userValue;
	}

	return true;
}