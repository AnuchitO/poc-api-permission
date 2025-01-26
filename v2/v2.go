package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Role
type Role string

const (
	Admin Role = "admin"
	User  Role = "user"
)

// Scope
type Scope string

const (
	UserReadSelf  Scope = "user:read:self"
	UserWriteSelf Scope = "user:write:self"
	AdminReadAll  Scope = "admin:read:all"
	AdminWriteAll Scope = "admin:write:all"
)

// Claims struct to define the JWT claims
type Claims struct {
	Role   Role   `json:"role"`
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

// Mock data for accounts and profiles
var accounts = map[string]string{
	"1": "Account 1",
	"2": "Account 2",
}

var profiles = map[string]string{
	"1": "Profile 1",
	"2": "Profile 2",
}

// Middleware to extract JWT token and set it in the context
func jwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
			c.Abort()
			return
		}

		// Set the claims into context for later use
		c.Set("claims", claims)
		c.Next()
	}
}

// Middleware to check role
func defineRole(allowedRoles []Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		if claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No claims found"})
			c.Abort()
			return
		}

		userClaims := claims.(*Claims)

		// Check if the user's role matches any of the allowed roles
		roleAllowed := false
		for _, role := range allowedRoles {
			if userClaims.Role == role {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Middleware to check scope
func defineScope(allowedScopes []Scope) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		if claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No claims found"})
			c.Abort()
			return
		}

		userClaims := claims.(*Claims)

		// Check if the user's scope is allowed
		scopeAllowed := false
		for _, scope := range allowedScopes {
			if scope == "user:read:self" && userClaims.UserID == c.Param("id") {
				scopeAllowed = true
				break
			} else if scope == "admin:read:all" && userClaims.Role == Admin {
				scopeAllowed = true
				break
			}
		}

		if !scopeAllowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Scope does not allow this action"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Handlers
func getAccountsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, accounts)
}

func getAccountByIDHandler(c *gin.Context) {
	id := c.Param("id")
	account, exists := accounts[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"account": account})
}

func getProfilesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, profiles)
}

func getProfileByIDHandler(c *gin.Context) {
	id := c.Param("id")
	profile, exists := profiles[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"profile": profile})
}

func createProfileHandler(c *gin.Context) {
	var profileData map[string]string
	if err := c.ShouldBindJSON(&profileData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	profileID := fmt.Sprintf("%d", len(profiles)+1)
	profiles[profileID] = profileData["profile"]
	c.JSON(http.StatusCreated, gin.H{"profileID": profileID})
}

func createAccountHandler(c *gin.Context) {
	var accountData map[string]string
	if err := c.ShouldBindJSON(&accountData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	accountID := fmt.Sprintf("%d", len(accounts)+1)
	accounts[accountID] = accountData["account"]
	c.JSON(http.StatusCreated, gin.H{"accountID": accountID})
}

// Main Router Setup
func setupRouter() *gin.Engine {
	r := gin.Default()

	// Apply JWT middleware globally
	r.Use(jwtMiddleware())

	// Define routes
	r.GET("/api/v1/accounts", defineRole([]Role{Admin}), getAccountsHandler)
	r.GET("/api/v1/accounts/:id", defineRole([]Role{User, Admin}), defineScope([]Scope{UserReadSelf, AdminReadAll}), getAccountByIDHandler)

	r.GET("/api/v1/profiles", defineRole([]Role{User, Admin}), getProfilesHandler)
	r.GET("/api/v1/profiles/:id", defineRole([]Role{User, Admin}), defineScope([]Scope{UserReadSelf, AdminReadAll}), getProfileByIDHandler)
	r.POST("/api/v1/profiles", defineRole([]Role{Admin}), createProfileHandler)
	r.POST("/api/v1/accounts", defineRole([]Role{Admin}), createAccountHandler)

	return r
}

func main() {
	fmt.Println("Starting server...")
	r := setupRouter()

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	fmt.Println("Server running on port:", port)
	r.Run(":8080")
}
