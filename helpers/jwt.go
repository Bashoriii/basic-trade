package helpers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = "your-256-bit-secret"

func VerifyToken(r *http.Request) (map[string]interface{}, error) {
	headerToken := r.Header.Get("Authorization")
	if !strings.HasPrefix(headerToken, "Bearer ") {
		return nil, errors.New("token required or invalid token format")
	}

	stringToken := strings.Split(headerToken, " ")[1]

	// Parse the token
	token, err := jwt.Parse(stringToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Validate the token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// Check token expiration
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().After(time.Unix(int64(exp), 0)) {
			return nil, errors.New("token is expired")
		}
	} else {
		return nil, errors.New("expire claim is missing")
	}

	return claims, nil
}

func GenerateToken(id uint, email string) string {
	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
		"exp":   time.Now().Add(time.Minute * 10).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return err.Error()
	}

	return signedToken
}
