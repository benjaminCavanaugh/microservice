package handler

import (
	"fmt"
	"net/http"
)

type HelloWorldhandler struct {
	message string
}

func (s HelloWorldhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
     fmt.Fprintf(w, "HeloWorld")
}