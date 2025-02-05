package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTMiddleware is a middleware that verifies the JWT access token in the request header
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"status":"failed", "message":"No authorization header", "result": gin.H{"error": "No authorization header"}})
			c.Abort()
			return
		}

		// Extract token (Bearer <token>)
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // No "Bearer " prefix
			c.JSON(http.StatusUnauthorized, gin.H{"status":"failed", "message":"Invalid token format", "result": gin.H{"error": "Invalid token format"}})
			c.Abort()
			return
		}

		// Verify token
		_, err := VerifyJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status":"failed", "message":"Unable to verify token", "result": gin.H{"error": err.Error()}})
			c.Abort()
			return
		}

		// Proceed to the next handler
		c.Next()
	}
}

// VerifyJWT verifies the JWT token and checks expiration
func VerifyJWT(tokenString string) (*jwt.MapClaims, error) {
	secretKey := []byte(os.Getenv("JSON_SECRET_KEY")) // Secret key used to sign the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token has expired")
		}		
		return nil, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, errors.New("invalid token")
}
