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
    fmt.Fprintf(w, "HeloWorld\n");

	if s.connection.Database != nil {
		fmt.Fprintf(w, s.connection.QueryUsers());
	} else {
		fmt.Fprintf(w, "The database connection was unusable.\n");
	}
}