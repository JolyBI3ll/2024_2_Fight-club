package main

import (
	"encoding/json"
	"github.com/gorilla/sessions"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("super-secret-key"))

func getSessionData(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session_id")

	if session.IsNew {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		errorResponse := map[string]string{"error": "No active session"}
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return

	}

	ID := session.Values["id"]
	Avatar, okAvatar := session.Values["avatar"].(string)
	body := map[string]interface{}{}
	if okAvatar {
		body = map[string]interface{}{
			"id":     ID,
			"avatar": Avatar,
		}
	} else {
		body = map[string]interface{}{
			"id":     ID,
			"avatar": "",
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
