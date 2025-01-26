package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

// Helper function to generate a test JWT token
func generateTestJWT(role Role, userID string) string {
	claims := Claims{
		Role:   role,
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("secret"))
	return tokenString
}

// Helper function to make a request with JWT token
func makeRequestWithToken(r http.Handler, method, url, token string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func makePostRequestWithToken(r http.Handler, url, token string, payload io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodPost, url, payload)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

// Test Cases
func TestGetAccounts(t *testing.T) {
	r := setupRouter()
	token := generateTestJWT(Admin, "admin1")

	// Test for admin role
	resp := makeRequestWithToken(r, http.MethodGet, "/api/v1/accounts", token)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetAccountByID(t *testing.T) {
	r := setupRouter()
	token := generateTestJWT(User, "1")

	// Test for valid user
	resp := makeRequestWithToken(r, http.MethodGet, "/api/v1/accounts/1", token)
	assert.Equal(t, http.StatusOK, resp.Code)

	// Test for user accessing another user's account
	resp = makeRequestWithToken(r, http.MethodGet, "/api/v1/accounts/2", token)
	assert.Equal(t, http.StatusForbidden, resp.Code)
}

func TestGetProfileByID(t *testing.T) {
	r := setupRouter()
	token := generateTestJWT(User, "1")

	// Test for valid user profile access
	resp := makeRequestWithToken(r, http.MethodGet, "/api/v1/profiles/1", token)
	assert.Equal(t, http.StatusOK, resp.Code)

	// Test for user accessing another user's profile
	resp = makeRequestWithToken(r, http.MethodGet, "/api/v1/profiles/2", token)
	assert.Equal(t, http.StatusForbidden, resp.Code)
}

func TestCreateAccount(t *testing.T) {
	r := setupRouter()
	token := generateTestJWT(Admin, "admin1")

	// Test for admin creating an account
	resp := makePostRequestWithToken(r, "/api/v1/accounts", token, strings.NewReader(`{"account": "test"}`))
	assert.Equal(t, http.StatusCreated, resp.Code)
}

func TestCreateProfile(t *testing.T) {
	r := setupRouter()
	token := generateTestJWT(Admin, "admin1")

	// Test for admin creating a profile
	resp := makePostRequestWithToken(r, "/api/v1/profiles", token, strings.NewReader(`{"profile": "test"}`))
	assert.Equal(t, http.StatusCreated, resp.Code)
}
