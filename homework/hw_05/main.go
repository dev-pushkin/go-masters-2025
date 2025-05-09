package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Для запуска llm API используйте команду:
// docker run -d -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama
// docker exec -it ollama ollama run llama3

const (
	// LLM API URL
	LLMAPIURL = "http://localhost:11434/api/generate"
	// LLM Model Name
	LLMModel = "llama3"
)

func main() {
	http.HandleFunc("/generate", handler)
	fmt.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	promt := r.URL.Query().Get("q")
	if promt == "" {
		http.Error(w, "missing query parameter 'q'", http.StatusBadRequest)
		return
	}

	reqData := struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
	}{
		Model:  LLMModel,
		Prompt: promt,
	}

	body, err := json.Marshal(reqData)
	if err != nil {
		http.Error(w, "failed to marshal request body", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequestWithContext(r.Context(), "POST", LLMAPIURL, bytes.NewReader(body))
	if err != nil {
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "failed to send request", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "LLM API returned non-200 status", resp.StatusCode)
		return
	}

	// Parse the response from the LLM API
	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	err = parseAndWriteLlmStream(resp, w)
	if err != nil {
		http.Error(w, "failed to parse LLM response", http.StatusInternalServerError)
		return
	}
}

type ResponseChunk struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

func parseAndWriteLlmStream(in *http.Response, out http.ResponseWriter) error {
	defer in.Body.Close()
	resBody := in.Body

	decoder := json.NewDecoder(resBody)
	flusher, ok := out.(http.Flusher)
	if !ok {
		return fmt.Errorf("response writer does not support flushing")
	}

	for {
		var chunk ResponseChunk
		if err := decoder.Decode(&chunk); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("decode error: %w", err)
		}

		if chunk.Response != "" {
			fmt.Fprint(out, chunk.Response)
			flusher.Flush()
		}

		if chunk.Done {
			break
		}
	}

	return nil
}
