package utilities

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/paulantezana/review/provider"
	"github.com/paulantezana/review/models"
	"log"
	"time"
)

// Claim model use un JWT Authentication
type Claim struct {
	User models.User `json:"user"`
	jwt.StandardClaims
}

// GenerateJWT generate token custom claims
func GenerateJWT(user models.User) string {
	// Set custom claims
	claims := &Claim{
		user,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
			Issuer:    "paulantezana",
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	result, err := token.SignedString([]byte(provider.GetConfig().Server.Key))
	if err != nil {
		log.Fatal("No se pudo firmar el token")
	}
	return result
}
