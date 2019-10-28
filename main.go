package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ajrmzcs/jwt-go/driver"
	"github.com/davecgh/go-spew/spew"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
)

type User struct {
	Id 			int64    `json:"id"`
	Email 		string `json:"email"`
	Password 	string `json:"password"`
}

type JWT struct {
	Token string `json:"token"`
}

type Error struct {
	Message string `json:"message"`
}

var db *sql.DB

func init() {
	_ = gotenv.Load()
}

func main() {
	db = driver.ConnectDB()
	r := mux.NewRouter()

	r.HandleFunc("/signup", signup).Methods("POST")
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/protected", TokenVerifyMiddleware(protectedEndpoint)).Methods("GET")

	log.Println("Listen on port: 8000...")
	log.Fatal(http.ListenAndServe(":8000",r))
}

func respondWithError(w http.ResponseWriter, status int, error Error) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(error)
}

func respondJson(w http.ResponseWriter, data interface{}) {
	_ = json.NewEncoder(w).Encode(data)
}

func signup(w http.ResponseWriter, r *http.Request) {
	var user User
	var error Error
	_ = json.NewDecoder(r.Body).Decode(&user)
	spew.Dump(user)

	if user.Email == "" {
		error.Message = "Email is missing"
		respondWithError(w, http.StatusBadRequest, error)
		return
	}

	if user.Password == "" {
		error.Message = "Password is missing"
		respondWithError(w, http.StatusBadRequest, error)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		log.Fatal(err)
	}

	user.Password = string(hash)

	res, err := db.Exec("INSERT INTO users (email, password) VALUES(?,?)",
		user.Email, user.Password)

	if err != nil {
		error.Message = "Server error"
		respondWithError(w, http.StatusInternalServerError, error)
		return
	}

	userId, err := res.LastInsertId()
	user.Id= userId
	user.Password=""

	w.Header().Set("Content-type", "application/json")
	respondJson(w, user)
}

func GenerateToken(user User) (string, error) {
	secret := "secret"

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

func login(w http.ResponseWriter, r *http.Request) {

	var user User
	var jwt JWT
	var error Error

	_ = json.NewDecoder(r.Body).Decode(&user)

	if user.Email == "" {
		error.Message = "Email is missing"
		respondWithError(w, http.StatusBadRequest, error)
		return
	}

	if user.Password == "" {
		error.Message = "Password is missing"
		respondWithError(w, http.StatusBadRequest, error)
		return
	}

	p := user.Password

	row := db.QueryRow("SELECT * FROM users WHERE email=?", user.Email)

	err := row.Scan(&user.Id, &user.Email, &user.Password)

	if err !=nil {
		if err == sql.ErrNoRows {
			error.Message = "User does not exist"
			respondWithError(w, http.StatusNotFound, error)
			return
		} else {
			log.Fatal(err)
		}
	}

	hashedP := user.Password

	err = bcrypt.CompareHashAndPassword([]byte(hashedP), []byte(p))

	if err != nil {
		error.Message = "Invalid password"
		respondWithError(w, http.StatusUnauthorized, error)
		return
	}

	token, err := GenerateToken(user)

	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	jwt.Token = token
	respondJson(w, jwt)
}

func protectedEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("protectedEndpoint invoked.")
}

func TokenVerifyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		var eObj Error
		authHeader := r.Header.Get("Authorization")

		bearerToken := strings.Split(authHeader, " ")

		if len(bearerToken) == 2 {
			authToken := bearerToken[1]

			token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}

				return []byte("secret"), nil
			})

			if err != nil {
				eObj.Message = err.Error()
				respondWithError(w, http.StatusInternalServerError, eObj)
				return
			}

			if token.Valid {
				next.ServeHTTP(w, r)
			} else {
				eObj.Message = err.Error()
				respondWithError(w, http.StatusUnauthorized, eObj)
				return
			}
		} else {
			eObj.Message = "Invalid token"
			respondWithError(w, http.StatusUnauthorized, eObj)
			return
		}
	})
}


