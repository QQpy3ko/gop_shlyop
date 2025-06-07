package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

    // get q query param
	queryParam := r.URL.Query().Get("q")

	if queryParam == "" {
		http.Error(w, "Параметр 'q' обязателен", http.StatusBadRequest)
		return
	}

	log.Printf("q query param = %s\n", queryParam)

	// get stuff from Ollama
	ollamaResp, err := queryOllama(queryParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка запроса к Ollama: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Ответ от Ollama:\n%s", ollamaResp)
}

func queryOllama(prompt string) (string, error) {
	reqBody := OllamaRequest{
		Model:  "codellama",
		Prompt: prompt,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var fullResponse string
	decoder := json.NewDecoder(resp.Body)

	for {
		var ollamaResp OllamaResponse
		if err := decoder.Decode(&ollamaResp); err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}
		fullResponse += ollamaResp.Response
	}

	return fullResponse, nil
}

func main() {
	http.HandleFunc("/", rootHandler)

	port := ":8080"
	log.Printf("Сервер запускается на порту %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server start failed: %s\n", err)
	}
}
