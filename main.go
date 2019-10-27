package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
	"log"
	"net/http"
)

type User struct {
	Id 			int    `json:"id"`
	Email 		string `json:"email"`
	password 	string `json:"password"`
}

type JWT struct {
	Token string `json:"token"`
}

type Error struct {
	Message string `json:"message"`
}

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

func signup(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("succefully called signup"))
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("succefully called login"))
}

func protectedEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("protectedEndpoint invoked.")
}

func TokenVerifyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	fmt.Println("TokenVerifyMiddleware invoked.")
	return nil
}


