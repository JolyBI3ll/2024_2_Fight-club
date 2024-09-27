package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

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

var places = []Place{
	{ID: 1, Title: "Уютный диван в центре города", Description: "Привет! Я предлагаю место на своем диване для путешественников.", Location: "Moscow", Host: Host{ID: 1, Username: "johndoe", Email: "johndoe@example.com"}, AvailableDates: []string{"2024-05-01", "2024-05-15"}, Rating: 9.1},
	{ID: 1, Title: "Приглашаю иностранцев к себе", Description: "Хаюшки, приезжайте все ко мне!.", Location: "Sochi", Host: Host{ID: 2, Username: "oleg", Email: "oleg228@example.com"}, AvailableDates: []string{"2024-05-01", "2024-05-15"}, Rating: 10},
	{ID: 1, Title: "Нет места, где переночевать?", Description: "Приючу у себя людей на пару дней.", Location: "Chita", Host: Host{ID: 3, Username: "kerla", Email: "kerla1337@example.com"}, AvailableDates: []string{"2024-05-01", "2024-05-15"}, Rating: 8.5},
	{ID: 1, Title: "Хочу поболтать с японцами", Description: "Охае, приезжайте ко мне, анимешники", Location: "Khabarovsk", Host: Host{ID: 4, Username: "animeLover", Email: "nikitasuper@example.com"}, AvailableDates: []string{"2024-05-01", "2024-05-15"}, Rating: 8.8},
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

	router.HandleFunc("/ads", getAllPlaces).Methods("GET")

	http.Handle("/", router)
	fmt.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error on starting server: %s", err)
	}
}
