package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/analyze", corsMiddleware(analyzeHandler))

	port := ":8080"
	fmt.Printf("Servidor iniciado en http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func analyzeHandler(writer http.ResponseWriter, request *http.Request) {

	host := request.URL.Query().Get("host")
	fromCacheParam := request.URL.Query().Get("fromCache")
	if host == "" {
		http.Error(writer, "Parameter 'host' is required. Example: /analyze?host=example.com", http.StatusBadRequest)
		return
	}

	var fromCache bool
	if fromCacheParam == "true" {
		fromCache = true
	} else {
		fromCache = false
	}

	result, err := analyze(host, fromCache)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Error analyzing the host: %v", err), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(writer).Encode(result); err != nil {
		log.Printf("Error encoding response: %v", err)
		return
	}
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
		if allowedOrigin == "" {
			allowedOrigin = "http://localhost:5173"
		}
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}
