package main

type Credentials struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Name     string `json:"name,omitempty"`
	Score    float64 `json:"score,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
	Sex      byte   `json:"sex,omitempty"`
	GuestCount int  `json:"guestCount,omitempty"`
	Age      int    `json:"age,omitempty"`
	Address  string `json:"address,omitempty"`
	IsHost   bool   `json:"isHost,omitempty"`
}

type Author struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Score       float64 `json:"score"`
	Avatar      string  `json:"avatar"`
	Sex         byte    `json:"sex"`
	GuestCount  int     `json:"guestCount"`
}

type Place struct {
	ID             int      `json:"id"`
	LocationMain          string   `json:"locationMain"`
	LocationStreet    string   `json:"locationStreet"`
	Position []float64 `json:"position"`
	Pictures       []string   `json:"pictures"`
	Author           Author     `json:"author"`
	PublicationDate string `json:"publicationDate"`
	AvailableDates []string `json:"avaibleDates"`
	Distance float64 `json:"distance"`
}

var Users = []Credentials{
	{ID: 1, Username: "johndoe", Password: "password123", Email: "johndoe@example.com", Name: "Leo D.", Score: 4.98, Avatar: "", Sex: 1, GuestCount: 50},
	{ID: 2, Username: "oleg", Password: "oleg123", Email: "oleg228@example.com", Name: "Oleg S.", Score: 4.50, Avatar: "", Sex: 1, GuestCount: 20},
}

var Places = []Place{
	{
		ID: 1, LocationMain: "Moscow", LocationStreet: "Tverskaya", Position: []float64{55.7558, 37.6173},
		Pictures: []string{"images/pic1.jpg", "images/pic2.jpg"}, Author: Author{ID: 1, Name: "Leo D.", Score: 4.98, Avatar: "", Sex: 1, GuestCount: 50},
		PublicationDate: "2024-09-01", AvailableDates: []string{"2024-09-15", "2024-09-20"}, Distance: 5.0,
	},
	{
		ID: 2, LocationMain: "Sochi", LocationStreet: "Kurortny Ave", Position: []float64{43.5855, 39.7231},
		Pictures: []string{"images/pic3.jpg", "images/pic4.jpg"}, Author: Author{ID: 2, Name: "Oleg S.", Score: 4.50, Avatar: "", Sex: 1, GuestCount: 20},
		PublicationDate: "2024-09-10", AvailableDates: []string{"2024-09-18", "2024-09-22"}, Distance: 3.2,
	},
}