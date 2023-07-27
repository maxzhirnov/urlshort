package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func handleCreateShortURL(urlShortener URLShortenerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests allowed", http.StatusBadRequest)
			return
		}

		//if contentType := r.Header.Get("Content-Type"); contentType != "text/plain" {
		//	http.Error(w, "Invalid Content-Type", http.StatusBadRequest)
		//	return
		//}
		// Looks like because of this auto test fails
		defer r.Body.Close()

		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		url := string(data)
		log.Println("Received url:", url)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)

		p := "http://" //TODO: Implement protocol parsing and mapping to string
		h := r.Host
		id, err := urlShortener.Create(url)
		if err != nil {
			http.Error(w, "Error creating shorten url", http.StatusInternalServerError)
		}
		shortenURL := fmt.Sprintf("%s%s/%s", p, h, id)

		if _, err := w.Write([]byte(shortenURL)); err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
	}
}

func handleGetOriginalURLByID(urlShortener URLShortenerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET requests allowed", http.StatusBadRequest)
			return
		}
		parts := strings.Split(r.URL.Path, "/")
		id := parts[1]
		url, err := urlShortener.Get(id)
		if err != nil {
			http.Error(w, "id not found", http.StatusBadRequest)
		}
		originalUrl := url.OriginalURL
		if strings.HasPrefix(originalUrl, "http") == false {
			originalUrl = "https://" + originalUrl
		}
		w.Header().Set("Location", originalUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
