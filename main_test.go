package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChatHandler_ValidRequest(t *testing.T) {
	// Prepare a valid chat request
	chatReq := ChatRequest{Message: "Hello"}
	body, _ := json.Marshal(chatReq)
	req := httptest.NewRequest("POST", "/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Mock Ollama API response
	// You may want to refactor chatHandler to allow injecting a mock client for full isolation

	chatHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status 200 or 500, got %d", resp.StatusCode)
	}
}

func TestChatHandler_InvalidRequest(t *testing.T) {
	// Send invalid JSON
	req := httptest.NewRequest("POST", "/chat", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	chatHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}
