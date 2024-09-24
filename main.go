package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to main page")
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", testHandler).Methods("GET")

	http.Handle("/",router)
	fmt.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error on starting server: %s", err)
	}
}