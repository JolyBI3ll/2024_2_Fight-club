package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
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
				Email:    "newuser@example.com",
				Password: "password",
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"session_id": "mock-session-id",
				"user": map[string]interface{}{
					"id":       1,
					"username": "newuser",
					"email":    "newuser@example.com",
				},
			},
			existingUsers: map[string]Credentials{},
		},
		{
			name: "User Already Exists",
			input: Credentials{
				Username: "existinguser",
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

func TestLogoutUser_AfterLogin_Success(t *testing.T) {
	mockUser := Credentials{
		Username: "testuser",
		Password: "testpassword",
	}
	addUser(mockUser)

	loginBody := map[string]string{
		"username": mockUser.Username,
		"password": mockUser.Password,
	}
	loginBodyJSON, _ := json.Marshal(loginBody)
	req, err := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(loginBodyJSON))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	loginHandler := http.HandlerFunc(loginUser)
	loginHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Login failed: status code %v", rr.Code)
	}

	cookies := rr.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatalf("No cookies found after login")
	}

	req, err = http.NewRequest("DELETE", "/api/auth/logout", nil)
	if err != nil {
		t.Fatalf("Failed to create logout request: %v", err)
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	rr = httptest.NewRecorder()
	logoutHandler := http.HandlerFunc(logoutUser)
	logoutHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var result map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &result)
	if err != nil {
		t.Fatalf("Failed to parse response body as JSON: %v", err)
	}
	fmt.Printf("Parsed JSON response: %+v\n", result)
	response, ok := result["response"].(string)
	if !ok {
		t.Fatalf("Response key not found or not a string, actual response body: %+v", result)
	}

	response = strings.TrimSpace(response)
	expected := "Logout successfully"
	fmt.Printf("Expected string: '%s' (len: %d)\n", expected, len(expected))
	fmt.Printf("Actual string:   '%s' (len: %d)\n", response, len(response))
	// Дополнительная отладка: выводим байтовое представление строк
	fmt.Printf("Expected bytes: %v\n", []byte(expected))
	fmt.Printf("Actual bytes:   %v\n", []byte(response))
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", response, expected)
	}

	session, _ := store.Get(req, "session_id")
	if session.Options.MaxAge != -1 {
		t.Errorf("Session was not properly invalidated: expected MaxAge to be -1, got %v", session.Options.MaxAge)
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
