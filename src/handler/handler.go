package handler

import (
	"fmt"
	connector "microservice/src/service"
	"net/http"
)


type HelloWorldhandler struct {
	message string
	connection connector.Connection;
}

func NewHandler(connection connector.Connection) HelloWorldhandler {
	var handler HelloWorldhandler;
	handler.connection = connection;
	return handler;
}

func (s HelloWorldhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.connection.Database != nil {
		// TODO: Expand this to accept query params and use them to populate the SQL query.
		fmt.Fprintf(w, s.connection.QueryUsers().String());
	} else {
		fmt.Fprintf(w, "The database connection was unusable.\n");
	}
}