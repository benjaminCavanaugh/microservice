package service

import (
	"database/sql"
	"fmt"

	// Bringing this in for the drivers.
	_ "github.com/lib/pq"
)

// TODO: Make it clear that it is the caller's responsiblity to call defer db.Close() when using this method.
func Connect(host string, port int, user string, password string, dbname string) (*sql.DB, error) {
	// TODO: Refactor this to enable sslmode.
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	fmt.Println(psqlInfo)

	db, connectionError := sql.Open("postgres", psqlInfo)

	if connectionError != nil {
		// TODO: Properly log this error.
		fmt.Println("Failed to connect to the database. Cannot continue.")
		panic(connectionError)
	}
	/*
	 * It is vitally important that you call the Ping() method becuase the sql.Open() function
	 * call does not ever create a connection to the database. Instead, it simply validates the arguments provided.
	 */
	connectionError = db.Ping()

	if connectionError != nil {
	  panic(connectionError)
	}
  
	fmt.Println("Successfully connected!")
	return db, nil;
}

func QueryUsers(db *sql.DB) {
	rows, queryError := db.Query("SELECT * from users;")

	if queryError != nil {
		fmt.Println("Failed to query database for users. Cannot continue");
		panic(queryError);
    }

	defer rows.Close()

	fmt.Println("Values from the users table:");

	// Loop through rows, using Scan to assign column data to struct fields.
    for rows.Next() {
		var username, location, dateOfBirth string

		if scanError := rows.Scan(&username, &location, &dateOfBirth); scanError != nil {
			fmt.Println("Encountered an error when processing query results.")
			continue;
        }

		fmt.Printf("username: %v, date of birth: %v, location: %v\n", username, location, dateOfBirth)
    }

	if rowError := rows.Err(); rowError != nil {
		fmt.Println("Error processing query results.")
    }
}