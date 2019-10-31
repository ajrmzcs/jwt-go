package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ajrmzcs/jwt-go/models"
	"github.com/ajrmzcs/jwt-go/utils"
	"github.com/ajrmzcs/jwt-go/repository/user"
	"github.com/davecgh/go-spew/spew"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"strings"
)

type Controller struct{}

func (c Controller) SignUp (db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var user models.User
		var e models.Error
		var uRepository userRepository.UserRepository

		_ = json.NewDecoder(r.Body).Decode(&user)
		spew.Dump(user)

		if user.Email == "" {
			e.Message = "Email is missing"
			utils.RespondWithError(w, http.StatusBadRequest, e)
			return
		}

		if user.Password == "" {
			e.Message = "Password is missing"
			utils.RespondWithError(w, http.StatusBadRequest, e)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

		if err != nil {
			log.Fatal(err)
		}

		user.Password = string(hash)

		newUser, err := uRepository.SignUp(db, user)

		if err != nil {
			e.Message = "Server error"
			utils.RespondWithError(w, http.StatusInternalServerError, e)
			return
		}

		w.Header().Set("Content-type", "application/json")
		utils.RespondJson(w, newUser)
	}
}

func (c Controller) Login(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {

		var user models.User
		var jwt models.JWT
		var e models.Error
		var uRepository userRepository.UserRepository

		_ = json.NewDecoder(r.Body).Decode(&user)

		if user.Email == "" {
			e.Message = "Email is missing"
			utils.RespondWithError(w, http.StatusBadRequest, e)
			return
		}

		if user.Password == "" {
			e.Message = "Password is missing"
			utils.RespondWithError(w, http.StatusBadRequest, e)
			return
		}

		p := user.Password

		user, err := uRepository.Login(db, user)

		if err !=nil {
			if err == sql.ErrNoRows {
				e.Message = "User does not exist"
				utils.RespondWithError(w, http.StatusNotFound, e)
				return
			} else {
				log.Fatal(err)
			}
		}

		hashedP := user.Password

		err = bcrypt.CompareHashAndPassword([]byte(hashedP), []byte(p))

		if err != nil {
			e.Message = "Invalid password"
			utils.RespondWithError(w, http.StatusUnauthorized, e)
			return
		}

		token, err := utils.GenerateToken(user)

		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		jwt.Token = token
		utils.RespondJson(w, jwt)
	}
}

func (c Controller) TokenVerifyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		var eObj models.Error
		authHeader := r.Header.Get("Authorization")

		bearerToken := strings.Split(authHeader, " ")

		if len(bearerToken) == 2 {
			authToken := bearerToken[1]

			token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("there was an error")
				}
				return []byte(os.Getenv("SECRET")), nil
			})

			if err != nil {
				eObj.Message = err.Error()
				utils.RespondWithError(w, http.StatusInternalServerError, eObj)
				return
			}

			if token.Valid {
				next.ServeHTTP(w, r)
			} else {
				eObj.Message = err.Error()
				utils.RespondWithError(w, http.StatusUnauthorized, eObj)
				return
			}
		} else {
			eObj.Message = "Invalid token"
			utils.RespondWithError(w, http.StatusUnauthorized, eObj)
			return
		}
	})
}


