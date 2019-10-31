package utils

import (
	"encoding/json"
	"github.com/ajrmzcs/jwt-go/models"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"os"
)

func RespondWithError(w http.ResponseWriter, status int, error models.Error) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(error)
}

func RespondJson(w http.ResponseWriter, data interface{}) {
	_ = json.NewEncoder(w).Encode(data)
}

func GenerateToken(user models.User) (string, error) {
	secret := os.Getenv("SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"iss": "jwt-app",
	})

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
}
