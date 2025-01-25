package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper function to create a mock JWT token (just for testing purposes)
func generateMockJWT(userID string, scopes []string) string {
	token, _ := generateJWT(userID, []string{"user"}, scopes) // Mock JWT for the given user and scopes
	return token
}

// Test the "create account" endpoint
func TestCreateAccount(t *testing.T) {
	r := setupRouter()

	tests := []struct {
		name          string
		userID        string
		scopes        []string
		payload       Account
		expectedCode  int
		expectedError string
	}{
		{
			name:          "Create own account",
			userID:        "user1",
			scopes:        []string{"user:write:self"},
			payload:       Account{ID: "3", UserID: "user1", Name: "Account 3"},
			expectedCode:  http.StatusCreated,
			expectedError: "",
		},
		{
			name:          "Forbidden when trying to create an account for another user",
			userID:        "user1",
			scopes:        []string{"user:write:self"},
			payload:       Account{ID: "4", UserID: "user2", Name: "Account 4"},
			expectedCode:  http.StatusForbidden,
			expectedError: "Permission denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new request
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(reqBody))

			// TODO :

			// Set Authorization header with a mock token
			req.Header.Set("Authorization", "Bearer "+generateMockJWT(tt.userID, tt.scopes))

			// Perform the request
			r.ServeHTTP(w, req)

			// Assert response code
			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode != http.StatusCreated {
				// Check if the error message is present
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			}
		})
	}
}

// Test the "get account" endpoint
func TestGetUserAccount(t *testing.T) {
	r := setupRouter()

	tests := []struct {
		name          string
		userID        string
		url           string
		scopes        []string
		accountID     string
		expectedCode  int
		expectedError string
	}{
		{
			name:          "Get own account",
			userID:        "user3",
			url:           "/users/user3/accounts/3",
			scopes:        []string{"user:read:self"},
			accountID:     "3",
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
		{
			name:          "Forbidden when trying to access another user's account",
			userID:        "user3",
			url:           "/users/user1/accounts/1",
			scopes:        []string{"user:read:self"},
			accountID:     "1",
			expectedCode:  http.StatusForbidden,
			expectedError: "Permission denied",
		},
		{
			name:          "Account not found",
			userID:        "user4",
			url:           "/users/user4/accounts/999",
			scopes:        []string{"user:read:self"},
			accountID:     "999",
			expectedCode:  http.StatusNotFound,
			expectedError: "Account not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tt.url, nil)

			// Set Authorization header with a mock token
			req.Header.Set("Authorization", "Bearer "+generateMockJWT(tt.userID, tt.scopes))

			// Perform the request
			r.ServeHTTP(w, req)

			// Assert response code
			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode != http.StatusOK {
				// Check if the error message is present
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			}
		})
	}
}

// Test the "get account" endpoint
func TestGetAccount(t *testing.T) {
	r := setupRouter()

	tests := []struct {
		name          string
		userID        string
		scopes        []string
		accountID     string
		expectedCode  int
		expectedError string
	}{
		{
			name:          "Get own account",
			userID:        "user1",
			scopes:        []string{"user:read:self"},
			accountID:     "1",
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
		{
			name:          "Forbidden when trying to access another user's account",
			userID:        "user1",
			scopes:        []string{"user:read:self"},
			accountID:     "2",
			expectedCode:  http.StatusForbidden,
			expectedError: "Permission denied",
		},
		{
			name:          "Account not found",
			userID:        "user1",
			scopes:        []string{"user:read:self"},
			accountID:     "999",
			expectedCode:  http.StatusNotFound,
			expectedError: "Account not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/accounts/"+tt.accountID, nil)

			// Set Authorization header with a mock token
			req.Header.Set("Authorization", "Bearer "+generateMockJWT(tt.userID, tt.scopes))

			// Perform the request
			r.ServeHTTP(w, req)

			// Assert response code
			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode != http.StatusOK {
				// Check if the error message is present
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			}
		})
	}
}

// Test the "update account" endpoint
func TestUpdateAccount(t *testing.T) {
	r := setupRouter()

	tests := []struct {
		name          string
		userID        string
		scopes        []string
		accountID     string
		payload       Account
		expectedCode  int
		expectedError string
	}{
		{
			name:          "Update own account",
			userID:        "user1",
			scopes:        []string{"user:write:self"},
			accountID:     "1",
			payload:       Account{ID: "1", UserID: "user1", Name: "Updated Account 1"},
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
		{
			name:          "Forbidden when trying to update another user's account",
			userID:        "user1",
			scopes:        []string{"user:write:self"},
			accountID:     "2",
			payload:       Account{ID: "2", UserID: "user2", Name: "Updated Account 2"},
			expectedCode:  http.StatusForbidden,
			expectedError: "Permission denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new request
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest(http.MethodPut, "/accounts/"+tt.accountID, bytes.NewReader(reqBody))

			// Set Authorization header with a mock token
			req.Header.Set("Authorization", "Bearer "+generateMockJWT(tt.userID, tt.scopes))

			// Perform the request
			r.ServeHTTP(w, req)

			// Assert response code
			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode != http.StatusOK {
				// Check if the error message is present
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			}
		})
	}
}

// Test the "delete account" endpoint
func TestDeleteAccount(t *testing.T) {
	r := setupRouter()

	tests := []struct {
		name          string
		userID        string
		scopes        []string
		accountID     string
		expectedCode  int
		expectedError string
	}{
		{
			name:          "Delete own account",
			userID:        "user1",
			scopes:        []string{"user:write:self"},
			accountID:     "1",
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
		{
			name:          "Forbidden when trying to delete another user's account",
			userID:        "user1",
			scopes:        []string{"user:write:self"},
			accountID:     "2",
			expectedCode:  http.StatusForbidden,
			expectedError: "Permission denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, "/accounts/"+tt.accountID, nil)

			// Set Authorization header with a mock token
			req.Header.Set("Authorization", "Bearer "+generateMockJWT(tt.userID, tt.scopes))

			// Perform the request
			r.ServeHTTP(w, req)

			// Assert response code
			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode != http.StatusOK {
				// Check if the error message is present
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			}
		})
	}
}
