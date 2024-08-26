package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var (
	urlMapping = make(map[string]string)
	mu         sync.Mutex
)

func shortenURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		originalURL := r.FormValue("url")
		if originalURL == "" {
			http.Error(w, "URL is required", http.StatusBadRequest)
			return
		}

		shortURL := fmt.Sprintf("http://localhost:8080/%d", len(urlMapping)+1)

		mu.Lock()
		urlMapping[shortURL] = originalURL
		mu.Unlock()

		fmt.Fprintf(w, "Shortened URL: %s\n", shortURL)
		fmt.Printf("INFO Shortened URL: %s\n", shortURL)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path
	mu.Lock()
	originalURL, exists := urlMapping[shortURL]
	mu.Unlock()

	if !exists {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusPermanentRedirect)
}

func mapHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(urlMapping)
}

func main() {
	http.HandleFunc("/shorten", shortenURLHandler)
	http.HandleFunc("/", redirectHandler)
	http.HandleFunc("/map", mapHandler)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
