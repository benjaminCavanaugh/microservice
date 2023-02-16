package service

import (
	"database/sql"
	"encoding/json"
	"fmt"

	configuration "microservice/src/config"
	dbModel "microservice/src/model"

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

type QueryResult struct {
	Success bool
	Message string
	Values []string
}

func (q QueryResult) String() string {
	b, err := json.Marshal(q)

	if err != nil {
		return "Error processing query results."
	}

	return string(b)
}

func userQuery[d dbModel.DbValue](c Connection, queryString string, queryParams ... string) QueryResult {
	var results QueryResult
	results.Success = true;
	var rows *sql.Rows
	var queryError error

	// TODO: Determine if we should check for errors first or just warn the users they need to do that between calls to Connect() and QueryUsers().

	if queryParams != nil || len(queryParams) <= 0 {
		// Before a slice of parameters can be supplied to *sql.DB.Query they must be converted to the "any" type.
		params := make([]any, len(queryParams));

		for index, current := range queryParams {
			params[index] = any(current);
		}

		rows, queryError = c.Database.Query(queryString, params...);
	} else {
		rows, queryError = c.Database.Query(queryString);
	}

	if queryError != nil {
		results.Success = false;
		results.Message = "Failed to query database for users. Cannot continue";
		fmt.Println(queryError)
		return results;
    }

	defer rows.Close()

	results.Message = "users";
	usersFound := parseResults[d](rows)

	if usersFound == nil || len(usersFound) <= 0 {
		results.Success = false;
		results.Message = "Error processing query results."
	}

	for index, currentPointer := range usersFound {
		fmt.Printf("\nusersFoundIndex: %v\n", index);

		// Here we are using a type assertion to convert currentResult back into a User.
		currentResult, ok := currentPointer.(*dbModel.User);

		fmt.Printf("Current value before conversion: %v\nCurrent value after conversion: %v\nConversion was successful:%v\n",
			currentPointer, currentResult, ok);

		if ok {
			results.Values = append(results.Values, fmt.Sprintf("username: %v, date of birth: %v, location: %v",
				currentResult.Username, currentResult.Location, currentResult.Dob))
		} else {
			// TODO: Convert this to use the generic type's name.
			fmt.Printf("Error. Failed to convert result %v to the type %v\n", currentPointer, "d")
		}
	}
	
	return results;
}

//func parseResults[d dbModel.DbValue](rows *sql.Rows) []d {
func parseResults[d dbModel.DbValue](rows *sql.Rows) []dbModel.DbValue {
	usersFound := make([]dbModel.DbValue, 0);
	index := 0;

	// Loop through rows, using Scan to assign column data to struct fields.
    for rows.Next() {
		var currentUser d;
		success := currentUser.ScanRow(rows, index, &usersFound);

		fmt.Printf("The query was successful: %v\n", success);
		fmt.Printf("The length of the slice is: %v\n", len(usersFound));
		fmt.Printf("The current user returned by the query is: %v\n", usersFound[index]);
		index++;

		if success == false {
			fmt.Println("Encountered an error when processing query results.")
			continue;
		}
    }

	if rowError := rows.Err(); rowError != nil {
		fmt.Println("Error processing query results.");
    }

	return usersFound;
}

func (c Connection) QueryUsers() QueryResult {
	return userQuery[dbModel.User](c, "SELECT * from users;");
}

func (c Connection) QueryUsersByName(userName string) QueryResult {
	return userQuery[dbModel.User](c, "SELECT * from users WHERE username = $1;", userName);
}