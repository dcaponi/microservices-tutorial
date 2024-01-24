package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/ping", pongHandler)
	http.HandleFunc("/hello-world", helloWorldHandler)

	fmt.Println("Server is running on http://localhost:8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func pongHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"ping": "plong"}
	json.NewEncoder(w).Encode(response)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the Candidate Service.")
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World! This is the Candidate Finder Service.")
}
