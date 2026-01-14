package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Calls the analyzeHandler for the /analyze endpoint
	http.HandleFunc("/analyze", corsMiddleware(analyzeHandler))

	port := ":8080"
	fmt.Printf("Servidor iniciado en http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// analyzeHandler handles HTTP requests to analyze the TLS/SSL security of a domain
// using the SSL Labs API. It extracts the host and optional cache parameter from the
// query string, performs the analysis, and returns the results as JSON.
//
// Query Parameters:
//   - host: The domain name to analyze (e.g., "example.com")
//   - fromCache: Whether to use cached results from SSL Labs server ("true" or "false")
//
// Parameters:
//   - writer: http.ResponseWriter to write the HTTP response
//   - request: *http.Request containing the incoming HTTP request
//
// Returns:
//   - 200 OK with JSON body containing the analysis results
//   - 400 Bad Request if the host parameter is missing
//   - 500 Internal Server Error if the analysis fails
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

// corsMiddleware wraps an HTTP handler to enable Cross-Origin Resource Sharing (CORS).
// It allows requests from a specified origin, which can be configured via the ALLOWED_ORIGIN
// environment variable. If not set, it defaults to http://localhost:5173.
//
// Parameters:
//   - next: The http.HandlerFunc to be wrapped with CORS functionality
//
// Returns:
//   - http.HandlerFunc: A new handler function with CORS headers configured
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
