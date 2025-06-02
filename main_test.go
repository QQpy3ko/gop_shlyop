package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFibonacci(t *testing.T) {
	testCases := []struct {
		name     string
		n        int
		expected int
	}{
		{"n=0", 0, 0},
		{"n=1", 1, 1},
		{"n=2", 2, 1},
		{"n=3", 3, 2},
		{"n=5", 5, 5},
		{"n=10", 10, 55},
		{"n=20", 20, 6765},
		{"n=-1", -1, 0},
		{"n=-5", -5, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := fibonacci(tc.n)
			
			if result != tc.expected {
				t.Errorf("fibonacci(%d) = %d; ожидалось %d", tc.n, result, tc.expected)
			}
		})
	}
}

func TestFiboHandler(t *testing.T) {
	testCases := []struct {
		name                 string
		queryParam           string
		expectedStatusCode   int
		expectedBodyContains string
		expectError          bool
	}{
		{"valid N=10", "N=10", http.StatusOK, "Число Фибоначчи для N=10: 55", false},
		{"valid N=0", "N=0", http.StatusOK, "Число Фибоначчи для N=0: 0", false},
		{"valid N=1", "N=1", http.StatusOK, "Число Фибоначчи для N=1: 1", false},
		{"missing N", "", http.StatusBadRequest, "Параметр N не указан", true},
		{"invalid N=abc", "N=abc", http.StatusBadRequest, "Некорректное значение для N", true},
		{"negative N=-5", "N=-5", http.StatusBadRequest, "N не может быть отрицательным", true},
		{"N too large", "N=41", http.StatusBadRequest, "N слишком большое", true},
		{"N at limit", "N=40", http.StatusOK, "Число Фибоначчи для N=40: 102334155", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := "/fibo"
			if tc.queryParam != "" {
				url = "/fibo?" + tc.queryParam
			}
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Fatalf("Не удалось создать запрос: %v", err)
			}

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(fiboHandler)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatusCode {
				t.Errorf("handler вернул неверный статус код: получили %v, ожидали %v. Тело ответа: %s",
					status, tc.expectedStatusCode, rr.Body.String())
			}

			responseBody := strings.TrimSpace(rr.Body.String())
			expectedBody := strings.TrimSpace(tc.expectedBodyContains)

			if !strings.Contains(responseBody, expectedBody) {
				t.Errorf("handler вернул неожиданное тело: получили '%s', ожидали содержание '%s'",
					responseBody, expectedBody)
			}
		})
	}
} 