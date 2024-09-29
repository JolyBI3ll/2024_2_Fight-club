package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("super-secret-key"))

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
	var creds Credentials

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if _, foundUser := findUserByUsername(creds.Username); foundUser {
		http.Error(w, "User already exist", http.StatusConflict)
		return
	}

	if creds.Username == "" || creds.Password == "" || creds.Email == "" {
		http.Error(w, "Username, password, and email are required", http.StatusBadRequest)
		return
	}

	creds.ID = userIDCounter
	userIDCounter++
	addUser(creds)
	session, _ := store.Get(r, "session_id")

	session.Values["id"] = creds.ID
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
		http.Error(w, "Failed to generate session ID", http.StatusInternalServerError)
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
			"id":       creds.ID,
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
	var creds Credentials

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	requestedUser, foundUser := findUserByUsername(creds.Username)
	if !foundUser || requestedUser.Password != creds.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	session, _ := store.Get(r, "session_id")

	session.Values["id"] = requestedUser.ID
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
		http.Error(w, "Failed to generate session ID", http.StatusInternalServerError)
		return
	}
	session.Values["session_id"] = sessionID

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"session_id": sessionID,
		"user": map[string]interface{}{
			"id":       requestedUser.ID,
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
		http.Error(w, "No such session", http.StatusBadRequest)
		return
	}
	session.Options.MaxAge = -1

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to leave session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Logout successfully")); err != nil {
		http.Error(w, "Failed to leave session", http.StatusInternalServerError)
	}
}

func getAllPlaces(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body := map[string]interface{}{
		"places": Places,
	}

	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getSessionData(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session_id")

	if session.IsNew {
		http.Error(w, "No active session", http.StatusUnauthorized)
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

func main() {
	store.Options.HttpOnly = true
	router := mux.NewRouter()
	api := "/api"

	router.HandleFunc(api+"/ads", getAllPlaces).Methods("GET")

	router.HandleFunc(api+"/auth/register", registerUser).Methods("POST")
	router.HandleFunc(api+"/auth/login", loginUser).Methods("POST")
	router.HandleFunc(api+"/auth/logout", logoutUser).Methods("DELETE")

	router.HandleFunc(api+"/getSessionData", getSessionData).Methods("GET")

	http.Handle("/", router)
	fmt.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error on starting server: %s", err)
	}
}
