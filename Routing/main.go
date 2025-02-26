package main

import (
	"fmt"
	"net/http"
)

//Handler Interface
type welcome int

func (wc welcome) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the server!")
}

//Handler Func
func login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Logging in...")
}

func getJson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	switch r.Method {
		case "GET":
			w.Write([]byte(`{"message": "GET method requested"}`))
		case "POST":
			w.Write([]byte(`{"message": "POST method requested"}`))
	}
}


func main() {

	//Router 
	router := http.NewServeMux()

	var wc welcome
	router.Handle("/", wc)

	//Handler Funcs
	router.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Logging out...")
	})

	router.HandleFunc("/login", login)
	router.HandleFunc("/json", getJson)

	//Server
	server := http.Server{
		Addr: ":8080",
		Handler: router,
	}

	//Run Server
	server.ListenAndServe()
}