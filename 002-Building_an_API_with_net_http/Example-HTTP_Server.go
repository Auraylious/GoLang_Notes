package main

import (
	"fmt"
	"log"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request){
	// log request
	fmt.Println("Endpoint Hit: homePage")
	// send data
	fmt.Fprintf(w, "Welcome to the HomePage!")
}

func otherPage(w http.ResponseWriter, r *http.Request){
	// log request
	fmt.Println("Endpoint Hit: Other Page")
	// send data
	fmt.Fprintf(w, "Welcome to the Other Page!")
}

func main() {
	// register paths to functions
	http.HandleFunc("GET /", homePage)
	http.HandleFunc("GET /other", otherPage)

	// if you would like to write the handler here instead
	http.HandleFunc("GET /third", func (w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "This is the Third Page!")
	})

	//Listen on port 3000, logs errors to stderr and exits the program
	fmt.Println("Starting Server on Port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
