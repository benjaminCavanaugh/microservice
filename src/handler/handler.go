package handler

import (
	"database/sql"
	"fmt"
	"net/http"
)

var database *sql.DB;

type HelloWorldhandler struct {
	message string
}

func NewHandler(db *sql.DB) {
	database = db;
}

func (s HelloWorldhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
     fmt.Fprintf(w, "HeloWorld")
}