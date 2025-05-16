package main

import (
	"fmt"
	"gop_shlyop/internal/api"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting application on port 8080...")

	http.HandleFunc("/ping", api.PingHandler)
	http.HandleFunc("/test-error-not-found", api.ExampleErrorHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
