package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID uint) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // 1 day expiry
	}

	log.Println("JWT_SECRET during sign:", os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ParseJWT(tokenStr string) (uint, error) {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		log.Println("JWT parse error:", err)
		return 0, fmt.Errorf("token invalido ou expirado")
	}

	if !token.Valid {
		log.Println("JWT invalid token")
		return 0, fmt.Errorf("token invalido ou expirado")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("JWT claims cast failed")
		return 0, fmt.Errorf("invalido token claims")
	}

	uidFloat, ok := claims["user_id"].(float64)
	if !ok {
		log.Println("user_id missing in claims")
		return 0, fmt.Errorf("user_id not found in token")
	}
	log.Println("Parsed user_id:", uidFloat)

	// save the user_id in context for later use
	if uidFloat <= 0 {
		log.Println("Invalid user_id in claims:", uidFloat)
		return 0, fmt.Errorf("invalid user_id in token")
	}

	return uint(uidFloat), nil
}

// ValidateJWT checks if the JWT is valid and returns the user ID if it is.
func ValidateJWT(tokenStr string) (uint, error) {
	userID, err := ParseJWT(tokenStr)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
