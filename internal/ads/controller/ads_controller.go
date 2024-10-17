package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/ads/usecase"
	"2024_2_FIGHT-CLUB/internal/service/session"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

type AdHandler struct {
	adUseCase      usecase.AdUseCase
	sessionService *session.ServiceSession
}

func NewAdHandler(adUseCase usecase.AdUseCase, sessionService *session.ServiceSession) *AdHandler {
	return &AdHandler{
		adUseCase:      adUseCase,
		sessionService: sessionService,
	}
}
func (h *AdHandler) GetAllPlaces(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	places, err := h.adUseCase.GetAllPlaces()
	if err != nil {
		h.handleError(w, err)
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

func (h *AdHandler) GetOnePlace(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	adId := vars["adId"]
	place, err := h.adUseCase.GetOnePlace(adId)
	if err != nil {
		h.handleError(w, err)
		return
	}
	body := map[string]interface{}{
		"place": place,
	}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *AdHandler) CreatePlace(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var place domain.Ad
	if err := json.NewDecoder(r.Body).Decode(&place); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := h.sessionService.GetUserID(r, w)
	if err != nil {
		h.handleError(w, errors.New("no active session"))
		return
	}
	place.AuthorUUID = userID
	err = h.adUseCase.CreatePlace(&place)
	if err != nil {
		h.handleError(w, err)
		return
	}
	body := map[string]interface{}{
		"place": place,
	}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *AdHandler) UpdatePlace(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var place domain.Ad
	if err := json.NewDecoder(r.Body).Decode(&place); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	adId := vars["adId"]

	userID, err := h.sessionService.GetUserID(r, w)
	if err != nil {
		h.handleError(w, errors.New("no active session"))
		return
	}

	err = h.adUseCase.UpdatePlace(&place, adId, userID)
	if err != nil {
		h.handleError(w, err)
		return
	}
	body := map[string]interface{}{
		"place": place,
	}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *AdHandler) DeletePlace(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	adId := vars["adId"]

	userID, err := h.sessionService.GetUserID(r, w)
	if err != nil {
		h.handleError(w, errors.New("no active session"))
	}

	err = h.adUseCase.DeletePlace(adId, userID)
	if err != nil {
		h.handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func (h *AdHandler) GetPlacesPerCity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	city := vars["city"]
	var ads []domain.Ad
	ads, err := h.adUseCase.GetPlacesPerCity(city)
	if err != nil {
		h.handleError(w, err)
		return
	}
	body := map[string]interface{}{
		"places": ads,
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *AdHandler) handleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	errorResponse := map[string]string{"error": err.Error()}

	switch err.Error() {
	case "ad not found":
		w.WriteHeader(http.StatusNotFound)
	case "ad already exists":
		w.WriteHeader(http.StatusConflict)
	case "not owner of ad", "no active session":
		w.WriteHeader(http.StatusUnauthorized)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	if jsonErr := json.NewEncoder(w).Encode(errorResponse); jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
	}
}
