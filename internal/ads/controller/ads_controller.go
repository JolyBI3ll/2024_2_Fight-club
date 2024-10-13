package controller

import (
	"2024_2_FIGHT-CLUB/internal/ads/usecase"
	"encoding/json"
	"net/http"
)

type AdHandler struct {
	adUseCase usecase.AdUseCase
}

func NewAdHandler(adUseCase usecase.AdUseCase) *AdHandler {
	return &AdHandler{
		adUseCase: adUseCase,
	}
}
func (h *AdHandler) GetAllPlaces(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	places, err := h.adUseCase.GetAllPlaces()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := map[string]string{"error": err.Error()}
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	body := map[string]interface{}{
		"places": places,
	}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
