package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Crnogorski rap is rolling on Service B!\n")
	})

	log.Println("Service B is running on :8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
} 