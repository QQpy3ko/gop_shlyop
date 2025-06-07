package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/valyala/fasthttp"
)

func main() {
    // net/http server
    go func() {
        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintf(w, "Farm your gold with net/http!")
        })
        log.Println("Standard HTTP server running on :8080")
        log.Fatal(http.ListenAndServe(":8080", nil))
    }()

    // fasthttp server
    go func() {
        handler := func(ctx *fasthttp.RequestCtx) {
            fmt.Fprintf(ctx, "Farm your gold with fasthttp!")
        }
        log.Println("fasthttp server running on :8081")
        log.Fatal(fasthttp.ListenAndServe(":8081", handler))
    }()

    // block main goroutine forever
    select {}
}
