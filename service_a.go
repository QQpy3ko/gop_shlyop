package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		secondServResp, err := http.Get("http://localhost:8081")
		if err != nil {
			http.Error(w, "Failed to call Service B", http.StatusInternalServerError)
			return
		}
		defer secondServResp.Body.Close()

		body, err := io.ReadAll(secondServResp.Body)
		if err != nil {
			http.Error(w, "Failed to read response from Service B", http.StatusInternalServerError)
			return
		}

		io.WriteString(w, "Crnogorski rap is rolling on Service A!\n")
		io.WriteString(w, string(body))
	})

	log.Println("Service A is running on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
} 