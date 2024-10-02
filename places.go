package main

import (
	"encoding/json"
	"net/http"
)

func getAllPlaces(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body := map[string]interface{}{
		"places": Places,
	}

	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
