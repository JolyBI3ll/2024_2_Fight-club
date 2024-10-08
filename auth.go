package main

import (
	"2024_2_FIGHT-CLUB/ds"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
)

func findUserByUsername(username string) (Credentials, bool) {
	for _, user := range Users {
		if user.Username == username {
			return user, true
		}
	}
	return Credentials{}, false
}

func addUser(creds Credentials) {
	Users = append(Users, creds)
}

var userIDCounter = Users[len(Users)-1].ID + 1 //уникальные id

func generateSessionID() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	var creds ds.User

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if creds.Username == "" || creds.Password == "" || creds.Email == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := map[string]string{"error": "Username, password, and email are required"}
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	errorResponse := map[string]interface{}{
		"error":       "Incorrect data forms",
		"wrongFields": []string{},
	}
	var wrongFields []string
	if !ValidateLogin(creds.Username) {
		wrongFields = append(wrongFields, "username")
	}
	if !ValidateEmail(creds.Email) {
		wrongFields = append(wrongFields, "email")
	}
	if !ValidatePassword(creds.Password) {
		wrongFields = append(wrongFields, "password")
	}
	if !ValidateName(creds.Name) {
		wrongFields = append(wrongFields, "name")
	}
	if len(wrongFields) != 0 {
		errorResponse["wrongFields"] = wrongFields
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var existingUser ds.User
	if err := db.Where("username = ?", creds.Username).First(&existingUser).Error; err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		errorResponse := map[string]string{"error": "User already exists"}
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := db.Create(&creds).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, _ := store.Get(r, "session_id")

	session.Values["id"] = creds.UUID
	session.Values["username"] = creds.Username
	session.Values["email"] = creds.Email
	if creds.Name != "" {
		session.Values["name"] = creds.Name
	}
	if creds.Avatar != "" {
		session.Values["avatar"] = creds.Avatar
	}

	sessionID, err := generateSessionID()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["session_id"] = sessionID

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body := map[string]interface{}{
		"session_id": sessionID,
		"user": map[string]interface{}{
			"id":       creds.UUID,
			"username": creds.Username,
			"email":    creds.Email,
		},
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func loginUser(w http.ResponseWriter, r *http.Request) {
	var creds ds.User

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	errorResponse := map[string]interface{}{
		"error":       "Incorrect data forms",
		"wrongFields": []string{},
	}
	var wrongFields []string
	if !ValidateLogin(creds.Username) {
		wrongFields = append(wrongFields, "username")
	}
	if !ValidatePassword(creds.Password) {
		wrongFields = append(wrongFields, "password")
	}
	if len(wrongFields) != 0 {
		errorResponse["wrongFields"] = wrongFields
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var requestedUser ds.User
	if err := db.Where("username = ?", creds.Username).First(&requestedUser).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		errorResponse := map[string]string{"error": "Invalid credentials"}
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if requestedUser.Password != creds.Password {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		errorResponse := map[string]string{"error": "Invalid credentials"}
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	session, _ := store.Get(r, "session_id")

	session.Values["id"] = requestedUser.UUID
	session.Values["username"] = requestedUser.Username
	session.Values["email"] = requestedUser.Email
	if requestedUser.Name != "" {
		session.Values["name"] = requestedUser.Name
	}
	if requestedUser.Avatar != "" {
		session.Values["avatar"] = requestedUser.Avatar
	}

	sessionID, err := generateSessionID()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["session_id"] = sessionID

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"session_id": sessionID,
		"user": map[string]interface{}{
			"id":       requestedUser.UUID,
			"username": requestedUser.Username,
			"email":    requestedUser.Email,
		},
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func logoutUser(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session_id")
	if session.IsNew {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := map[string]string{"error": "No such session"}
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	session.Options.MaxAge = -1

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	logoutResponse := map[string]string{"response": "Logout successfully"}
	if err := json.NewEncoder(w).Encode(logoutResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
