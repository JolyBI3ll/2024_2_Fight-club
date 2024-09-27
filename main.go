package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("super-secret-key"))

type Credentials struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Host struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Place struct {
	ID             int      `json:"id"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Location       string   `json:"location"`
	Host           Host     `json:"host"`
	AvailableDates []string `json:"avaibleDates"`
	Rating         float64  `json:"rating"`
}

var users = []Credentials{
	{ID: 1, Username: "johndoe", Password: "password123", Email: "johndoe@example.com"},
	{ID: 2, Username: "oleg", Password: "oleg123", Email: "oleg228@example.com"},
	{ID: 3, Username: "kerla", Password: "kerla123", Email: "kerla1337@example.com"},
	{ID: 4, Username: "animeLover", Password: "neruto", Email: "nikitasuper@example.com"},
}

var places = []Place{
	{ID: 1, Title: "Уютный диван в центре города", Description: "Привет! Я предлагаю место на своем диване для путешественников.", Location: "Moscow", Host: Host{ID: 1, Username: "johndoe", Email: "johndoe@example.com"}, AvailableDates: []string{"2024-05-01", "2024-05-15"}, Rating: 9.1},
	{ID: 1, Title: "Приглашаю иностранцев к себе", Description: "Хаюшки, приезжайте все ко мне!.", Location: "Sochi", Host: Host{ID: 2, Username: "oleg", Email: "oleg228@example.com"}, AvailableDates: []string{"2024-05-01", "2024-05-15"}, Rating: 10},
	{ID: 1, Title: "Нет места, где переночевать?", Description: "Приючу у себя людей на пару дней.", Location: "Chita", Host: Host{ID: 3, Username: "kerla", Email: "kerla1337@example.com"}, AvailableDates: []string{"2024-05-01", "2024-05-15"}, Rating: 8.5},
	{ID: 1, Title: "Хочу поболтать с японцами", Description: "Охае, приезжайте ко мне, анимешники", Location: "Khabarovsk", Host: Host{ID: 4, Username: "animeLover", Email: "nikitasuper@example.com"}, AvailableDates: []string{"2024-05-01", "2024-05-15"}, Rating: 8.8},
}

func findUserByUsername(username string) (Credentials, bool) {
	for _, user := range users {
		if user.Username == username {
			return user, true
		}
	}
	return Credentials{}, false
}

func addUser(creds Credentials) {
	users = append(users, creds)
}

var userIDCounter = users[len(users)-1].ID + 1 //уникальные id

func generateSessionID() string {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Failed to generate random session ID: %v", err)
	}

	return base64.StdEncoding.EncodeToString(b)
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

	creds.ID = userIDCounter
	userIDCounter++
	addUser(creds)

	session, _ := store.Get(r, "session-id")
	sessionID := generateSessionID()
	session.Values["session_id"] = sessionID
	session.Values["username"] = creds.Username
	session.Values["email"] = creds.Email
	err := session.Save(r, w)
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

	session, _ := store.Get(r, "session-id")

	sessionID := generateSessionID()

	session.Values["session_id"] = sessionID
	session.Values["username"] = requestedUser.Username
	session.Values["email"] = requestedUser.Email

	err := session.Save(r, w)
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
	session, _ := store.Get(r, "session-id")
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
		"places": places,
	}

	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/ads", getAllPlaces).Methods("GET")

	router.HandleFunc("/api/auth/register", registerUser).Methods("POST")
	router.HandleFunc("/api/auth/login", loginUser).Methods("POST")
	router.HandleFunc("/api/auth/logout", logoutUser).Methods("POST")

	http.Handle("/", router)
	fmt.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error on starting server: %s", err)
	}
}
