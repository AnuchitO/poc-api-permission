package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Claims structure representing the payload of a JWT token
type Claims struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
	Scopes []string `json:"scopes"`
	jwt.StandardClaims
}

// Generate a sample JWT token for a user
func generateJWT(userID string, roles []string, scopes []string) (string, error) {
	claims := Claims{
		UserID: userID,
		Roles:  roles,
		Scopes: scopes,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "keycloak",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret")) // Use a secure secret key in production
}

// Extract claims from the token
func extractClaimsFromToken(authHeader string) (*Claims, error) {
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		return nil, fmt.Errorf("token missing")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil // Use a secure secret key in production
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("could not extract claims")
	}
	return claims, nil
}

// Check if a user has a specific scope
func hasScope(scopes []string, requiredScope string) bool {
	for _, scope := range scopes {
		if scope == requiredScope {
			return true
		}
	}
	return false
}

// Authorization middleware to verify permissions at the middleware level
func authorizationMiddleware(permissionRequired string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token missing"})
			c.Abort()
			return
		}

		// Extract claims from the JWT token
		claims, err := extractClaimsFromToken(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check if the user has the required permission (scope)
		if !hasScope(claims.Scopes, permissionRequired) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}

		// Attach user_id to context for later use in handlers
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

type Account struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

type Transaction struct {
	ID        string  `json:"id"`
	AccountID string  `json:"account_id"`
	Amount    float64 `json:"amount"`
	CreatedAt string  `json:"created_at"`
}

// Mock data
var accounts = []Account{
	{ID: "1", UserID: "user1", Name: "Account 1"},
	{ID: "2", UserID: "user2", Name: "Account 2"},
}

var transactions = []Transaction{
	{ID: "tx1", AccountID: "1", Amount: 100, CreatedAt: "2024-01-01"},
	{ID: "tx2", AccountID: "2", Amount: 50, CreatedAt: "2024-01-02"},
}

// Create an account (only admin or the owner)
func createAccount(c *gin.Context) {
	claims, _ := c.Get("user_id")
	userID := claims.(string)

	var newAccount Account
	if err := c.ShouldBindJSON(&newAccount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if admin or the user is the owner
	if newAccount.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	// Add to accounts list
	accounts = append(accounts, newAccount)
	c.JSON(http.StatusCreated, newAccount)
}

// Get an account (only admin or the owner)
func getAccount(c *gin.Context) {
	accountID := c.Param("id")
	claims, _ := c.Get("user_id")
	userID := claims.(string)

	var account *Account
	for _, a := range accounts {
		if a.ID == accountID {
			account = &a
			break
		}
	}

	if account == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	// Check if the user is admin or the owner of the account
	if account.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	c.JSON(http.StatusOK, account)
}

// Update an account (only admin or the owner)
func updateAccount(c *gin.Context) {
	accountID := c.Param("id")
	claims, _ := c.Get("user_id")
	userID := claims.(string)

	var updatedAccount Account
	if err := c.ShouldBindJSON(&updatedAccount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var account *Account
	for i, a := range accounts {
		if a.ID == accountID {
			account = &accounts[i]
			break
		}
	}

	if account == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	// Check if admin or the user is the owner
	if account.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	account.Name = updatedAccount.Name
	c.JSON(http.StatusOK, account)
}

// Delete an account (only admin or the owner)
func deleteAccount(c *gin.Context) {
	accountID := c.Param("id")
	claims, _ := c.Get("user_id")
	userID := claims.(string)

	var account *Account
	for i, a := range accounts {
		if a.ID == accountID {
			account = &accounts[i]
			accounts = append(accounts[:i], accounts[i+1:]...)
			break
		}
	}

	if account == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	// Check if admin or the user is the owner
	if account.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted"})
}

func main() {
	r := gin.Default()

	// Define routes with authorization checks

	// Account routes - employee can only manage their own accounts
	r.POST("/accounts", authorizationMiddleware("user:write:self"), createAccount)
	r.GET("/accounts/:id", authorizationMiddleware("user:read:self"), getAccount)
	r.PUT("/accounts/:id", authorizationMiddleware("user:write:self"), updateAccount)
	r.DELETE("/accounts/:id", authorizationMiddleware("user:write:self"), deleteAccount)

	// Define transactions routes - same authorization model applies
	// Similar CRUD for transactions

	r.Run(":8888")
}
