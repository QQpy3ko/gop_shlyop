package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func fibonacci(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

func fiboHandler(w http.ResponseWriter, r *http.Request) {
	nStr := r.URL.Query().Get("N")
	if nStr == "" {
		http.Error(w, "Параметр N не указан", http.StatusBadRequest)
		return
	}

	n, err := strconv.Atoi(nStr)
	if err != nil {
		http.Error(w, "Некорректное значение для N", http.StatusBadRequest)
		return
	}

	if n < 0 {
		http.Error(w, "N не может быть отрицательным", http.StatusBadRequest)
		return
	}

	if n > 40 {
		http.Error(w, "N слишком большое", http.StatusBadRequest)
		return
	}

	result := fibonacci(n)
	fmt.Fprintf(w, "Число Фибоначчи для N=%d: %d\n", n, result)
}

func main() {
	http.HandleFunc("/fibo", fiboHandler)

	fmt.Println("Сервер запускается на порту :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
