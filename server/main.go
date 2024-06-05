package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/cat-fact", catFact)

	port := ":9001"
	log.Printf("Starting server on %s", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Failed to run HTTP server: ", err)
	}
}

func catFact(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request received: %s", r.URL.Path)

	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "name parameter is missing", http.StatusBadRequest)
		return
	}

	resp := fmt.Sprintf("Hello %s! Happy birthday!", name)
	greetResp := Response{
		Message: resp,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(greetResp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
