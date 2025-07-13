package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
)

// ChatRequest represents a message sent from the frontend
type ChatRequest struct {
	Message string `json:"message"`
}

// ChatResponse represents a message from LLaMA3
type ChatResponse struct {
	Response string `json:"response"`
}

var promptQueue = struct {
	m map[string]string
	sync.Mutex
}{m: make(map[string]string)}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Invalid request:", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	log.Println("Received message:", req.Message)

	// Proxy to Ollama API (assuming Ollama is running locally)
	ollamaReq := map[string]interface{}{"prompt": req.Message, "model": "llama3"}
	ollamaBody, _ := json.Marshal(ollamaReq)
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", ioutil.NopCloser(bytes.NewReader(ollamaBody)))
	if err != nil {
		log.Println("Ollama API error:", err)
		http.Error(w, "Ollama API error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	responseText := ""
	decoder := json.NewDecoder(resp.Body)
	for {
		var chunk map[string]interface{}
		if err := decoder.Decode(&chunk); err != nil {
			break // End of stream
		}
		if val, ok := chunk["response"].(string); ok {
			responseText += val
		}
		if done, ok := chunk["done"].(bool); ok && done {
			break
		}
	}
	log.Println("Sending response:", responseText)

	json.NewEncoder(w).Encode(ChatResponse{Response: responseText})
}

func chatPromptHandler(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Invalid request:", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	id := uuid.New().String()
	promptQueue.Lock()
	promptQueue.m[id] = req.Message
	promptQueue.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func chatStreamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}
	promptQueue.Lock()
	prompt, ok := promptQueue.m[id]
	if ok {
		delete(promptQueue.m, id)
	}
	promptQueue.Unlock()
	if !ok {
		http.Error(w, "Prompt not found", http.StatusNotFound)
		return
	}
	ollamaReq := map[string]interface{}{"prompt": prompt, "model": "llama3"}
	ollamaBody, _ := json.Marshal(ollamaReq)
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", ioutil.NopCloser(bytes.NewReader(ollamaBody)))
	if err != nil {
		log.Println("Ollama API error:", err)
		http.Error(w, "Ollama API error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	for {
		var chunk map[string]interface{}
		if err := decoder.Decode(&chunk); err != nil {
			break // End of stream
		}
		if val, ok := chunk["response"].(string); ok {
			fmt.Fprintf(w, "data: %s\n\n", val)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		}
		if done, ok := chunk["done"].(bool); ok && done {
			break
		}
	}
}

func main() {
	staticDir, err := filepath.Abs("./static")
	if err != nil {
		log.Fatalf("Failed to get static dir: %v", err)
	}
	http.HandleFunc("/api/chat", chatHandler)
	http.HandleFunc("/api/chat-prompt", chatPromptHandler)
	http.HandleFunc("/api/chat-stream", chatStreamHandler)
	http.Handle("/", http.FileServer(http.Dir(staticDir)))
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
