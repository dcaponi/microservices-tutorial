package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/hello-world", helloWorldHandler)
	http.HandleFunc("/candidate", candidateHandler)
	fmt.Println("Server is running on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func candidateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("requesting the candidate service")
	resp, err := http.Get("http://candidate-service/ping")
	if err != nil {
		log.Printf("Error calling candidate finder service: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(body)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the Resume Service.")
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hola, Mundo! This is the Resume Service.")
}
