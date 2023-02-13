package service

import (
	"database/sql"
	"fmt"
	"strings"

	configuration "microservice/src/config"

	// Bringing this in for the drivers.
	_ "github.com/lib/pq"
)

type Connection struct {
	Database *sql.DB
	Err error
}

// TODO: Make it clear that it is the caller's responsiblity to call defer database.Close() when using this method.
func Connect(databaseConfig configuration.DatabaseConfig) Connection {
	var connection Connection;

	host     := databaseConfig.Host
	port     := databaseConfig.Port
	user     := databaseConfig.User
	password := databaseConfig.Password
	dbname   := databaseConfig.DBname
	
	// TODO: Refactor this to enable sslmode.
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	fmt.Println(psqlInfo)

	connection.Database, connection.Err = sql.Open("postgres", psqlInfo)

	if connection.Err != nil {
		// TODO: Properly log this error.
		fmt.Println("Failed to connect to the database. Cannot continue.")
		panic(connection.Err)
	}

	/*
	 * It is vitally important that you call the Ping() method becuase the sql.Open() function
	 * call does not ever create a connection to the database. Instead, it simply validates the arguments provided.
	 */
	 connection.Err = connection.Database.Ping()

	if connection.Err != nil {
	  panic(connection.Err)
	}
  
	fmt.Println("Successfully connected!")
	return connection
}

func (c Connection) QueryUsers() string {
	var builder strings.Builder;
	// TODO: Determine if we should check for errors first or just warn the users they need to do that between calls to Connect() and QueryUsers().
	rows, queryError := c.Database.Query("SELECT * from users;")

	if queryError != nil {
		builder.WriteString("Failed to query database for users. Cannot continue\n");
		panic(queryError);
    }

	defer rows.Close()

	builder.WriteString("Values from the users table:\n");

	// Loop through rows, using Scan to assign column data to struct fields.
    for rows.Next() {
		var username, location, dateOfBirth string

		if scanError := rows.Scan(&username, &location, &dateOfBirth); scanError != nil {
			builder.WriteString("Encountered an error when processing query results.\n")
			continue;
        }

		builder.WriteString(fmt.Sprintf("username: %v, date of birth: %v, location: %v\n", username, location, dateOfBirth))
    }

	if rowError := rows.Err(); rowError != nil {
		builder.WriteString("Error processing query results.\n")
    }

	return builder.String();
}