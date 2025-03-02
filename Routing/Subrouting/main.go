package main

import (
	"fmt"
	"net/http"
)

func main() {

	router := http.NewServeMux()

	version1Router := http.NewServeMux()
	version1Router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Loggging in to version 1")
	})

	version2Router := http.NewServeMux()
	version2Router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprint(w, "Logging in to version 2")
	})

	//ContextBasedVersionStrategy ("api/v1" and "api/v2")
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", version1Router))
	router.Handle("/api/v2/", http.StripPrefix("/api/v2", version2Router))

	//URLBasedVersionStrategy ("api1" and "api2")
	router.Handle("/api1/", http.StripPrefix("/api1", version1Router))
	router.Handle("/api2/", http.StripPrefix("/api2", version2Router))

	http.ListenAndServe(":8080", router)


}


// regex 
// performance


//gin


/api/v1/department/1/employee/2