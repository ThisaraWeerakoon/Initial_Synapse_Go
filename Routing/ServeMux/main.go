package main

import (
	"fmt"
	"log"
	"net/http"
)



func main(){
	router := http.NewServeMux()
	router.HandleFunc("/user",func(w http.ResponseWriter,r *http.Request){
		fmt.Fprint(w,"/user")
	})

	router.HandleFunc("/user/",func(w http.ResponseWriter,r *http.Request){
		fmt.Fprint(w,"/user/")
	})

	server:= http.Server{
		Addr: ":8080",
		Handler: router,
	}

	server.ListenAndServe()
	err := server.ListenAndServe() // Capture the error

	if err != nil {
			log.Fatalf("Server failed to start: %v", err)
	}

}