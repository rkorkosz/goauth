package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

const SECRET_KEY = "$j%vi-tpv$4v0s^)-ala4idan@io#oo&sn_pn^e%hylt-c5p8b"

func main() {
	r := mux.NewRouter().StrictSlash(true)
	r.Handle("/", loginRequired(Index)).Methods("GET")
	AuthRouter(r.PathPrefix("/auth").Subrouter())
	NewController(r, "project", Project{})
	//CSRF := csrf.Protect([]byte(SECRET_KEY))
	lh := handlers.LoggingHandler(os.Stdout, r)
	log.Fatal(http.ListenAndServeTLS(":8000", "cert.pem", "key.pem", lh))
}

func Index(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, "index", "Index")
}
