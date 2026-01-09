package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/analyze", analyzeHandler)

	port := ":8080"
	fmt.Printf("Servidor iniciado en http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func analyzeHandler(writer http.ResponseWriter, request *http.Request) {

	host := request.URL.Query().Get("host")
	if host == "" {
		http.Error(writer, "Parameter 'host' is required. Example: /analyze?host=example.com", http.StatusBadRequest)
		return
	}

	result, err := analyze(host)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Error analyzing the host: %v", err), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(writer).Encode(result); err != nil {
		http.Error(writer, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}
