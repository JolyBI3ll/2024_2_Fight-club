package main

import (
	"2024_2_FIGHT-CLUB/ds"
	"encoding/json"
	"net/http"
)

func getAllPlaces(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ads []ds.Ad
	if err := db.Find(&ads).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body := map[string]interface{}{
		"places": ads,
	}

	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
