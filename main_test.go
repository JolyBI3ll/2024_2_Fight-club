package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindUserByUsername(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		expectedFoundUser Credentials
		expectedFound     bool
	}{
		{
			name:              "Existing User",
			input:             "johndoe",
			expectedFoundUser: Users[0],
			expectedFound:     true,
		},
		{
			name:              "Non-existing User",
			input:             "nonExistentUser",
			expectedFoundUser: Credentials{},
			expectedFound:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			foundUser, found := findUserByUsername(tt.input)

			assert.Equal(t, tt.expectedFound, found)
			assert.Equal(t, tt.expectedFoundUser, foundUser)
		})
	}
}

func TestRegisterUser(t *testing.T) {
	tests := []struct {
		name           string
		input          Credentials
		expectedStatus int
		expectedBody   map[string]interface{}
		existingUsers  map[string]Credentials
	}{
		{
			name: "Successful Registration",
			input: Credentials{
				Username: "newuser",
				Name:     "aboba123",
				Email:    "newuser@example.com",
				Password: "password",
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"session_id": "mock-session-id",
				"user": map[string]interface{}{
					"id":       1,
					"username": "newuser",
					"name":     "aboba123",
					"email":    "newuser@example.com",
				},
			},
			existingUsers: map[string]Credentials{},
		},
		{
			name: "User Already Exists",
			input: Credentials{
				Username: "existinguser",
				Name:     "existinguser",
				Email:    "existinguser@example.com",
				Password: "password",
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   nil,
			existingUsers: map[string]Credentials{
				"existinguser": {
					ID:       1,
					Username: "existinguser",
					Email:    "existinguser@example.com",
					Password: "password",
				},
			},
		},
		{
			name: "Invalid JSON",
			input: Credentials{
				Username: "",
				Email:    "",
				Password: "",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
			existingUsers:  map[string]Credentials{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Users = make([]Credentials, 0)
			for _, user := range tt.existingUsers {
				addUser(user)
			}
			userIDCounter = len(tt.existingUsers) + 1

			var body []byte
			if tt.name == "Invalid JSON" {
				body = []byte("{invalid json}")
			} else {
				body, _ = json.Marshal(tt.input)
			}

			req, err := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(registerUser)

			handler.ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatus, rr.Code)

			if rr.Code == http.StatusCreated {
				var responseBody map[string]interface{}
				err = json.NewDecoder(rr.Body).Decode(&responseBody)
				assert.NoError(t, err)
				assert.NotEmpty(t, responseBody["session_id"])
				assert.Equal(t, float64(userIDCounter-1), responseBody["user"].(map[string]interface{})["id"])
				assert.Equal(t, tt.expectedBody["user"].(map[string]interface{})["username"], responseBody["user"].(map[string]interface{})["username"])
				assert.Equal(t, tt.expectedBody["user"].(map[string]interface{})["email"], responseBody["user"].(map[string]interface{})["email"])
			}
		})
	}
}

func TestLoginUser(t *testing.T) {
	tests := []struct {
		name           string
		input          Credentials
		expectedStatus int
		expectedBody   map[string]interface{}
		existingUsers  map[string]Credentials
	}{
		{
			name: "Successful Login",
			input: Credentials{
				Username: "existinguser",
				Password: "password",
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"session_id": "mock-session-id",
				"user": map[string]interface{}{
					"id":       1,
					"username": "existinguser",
					"email":    "existinguser@example.com",
				},
			},
			existingUsers: map[string]Credentials{
				"existinguser": {
					ID:       1,
					Username: "existinguser",
					Email:    "existinguser@example.com",
					Password: "password",
				},
			},
		},
		{
			name: "Invalid Login Credentials",
			input: Credentials{
				Username: "invaliduser",
				Password: "invalidpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   nil,
			existingUsers: map[string]Credentials{
				"existinguser": {
					ID:       1,
					Username: "existinguser",
					Email:    "existinguser@example.com",
					Password: "password",
				},
			},
		},
		{
			name: "Invalid JSON",
			input: Credentials{
				Username: "",
				Email:    "",
				Password: "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
			existingUsers:  map[string]Credentials{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Users = make([]Credentials, 0)
			for _, user := range tt.existingUsers {
				addUser(user)
			}
			userIDCounter = len(tt.existingUsers) + 1

			var body []byte
			if tt.name == "Invalid JSON" {
				body = []byte("{invalid json}")
			} else {
				body, _ = json.Marshal(tt.input)
			}

			req, err := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(loginUser)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if rr.Code == http.StatusCreated {
				var responseBody map[string]interface{}
				err = json.NewDecoder(rr.Body).Decode(&responseBody)
				assert.NoError(t, err)
				assert.NotEmpty(t, responseBody["session_id"])
				assert.Equal(t, float64(userIDCounter-1), responseBody["user"].(map[string]interface{})["id"])
				assert.Equal(t, tt.expectedBody["user"].(map[string]interface{})["username"], responseBody["user"].(map[string]interface{})["username"])
				assert.Equal(t, tt.expectedBody["user"].(map[string]interface{})["email"], responseBody["user"].(map[string]interface{})["email"])
			}
		})
	}
}

func TestGetAllPlaces(t *testing.T) {
	type ResponsePlaces struct {
		Places []Place `json:"places"`
	}
	req, err := http.NewRequest("GET", "/api/ads", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(getAllPlaces)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Handler returned wrong content type: got %v want %v", contentType, "application/json")
	}

	var responseBody ResponsePlaces
	if err := json.NewDecoder(rr.Body).Decode(&responseBody); err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	expectedPlaces := Places

	if !reflect.DeepEqual(responseBody.Places, expectedPlaces) {
		t.Errorf("Handler returned unexpected body: got %v want %v", responseBody.Places, expectedPlaces)
	}
}

func TestLogoutUser_NoSession(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/api/auth/logout", nil)
	if err != nil {
		t.Fatalf("Failed to create logout request: %v", err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(logoutUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	var result map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &result)
	if err != nil {
		t.Fatalf("Failed to parse response body as JSON: %v", err)
	}

	response, ok := result["error"].(string)
	if !ok {
		t.Fatalf("Response key not found or not a string, actual response body: %+v", result)
	}

	expected := "No such session"
	if response != expected {
		t.Errorf("Handler returned unexpected response: got %v want %v", response, expected)
	}
}
