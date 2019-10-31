package main

import (
	"database/sql"
	"github.com/ajrmzcs/jwt-go/controllers"
	"github.com/ajrmzcs/jwt-go/driver"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
	"log"
	"net/http"
)

var db *sql.DB

func init() {
	_ = gotenv.Load()
}

func main() {
	db = driver.ConnectDB()
	r := mux.NewRouter()
	controller := controllers.Controller{}

	r.HandleFunc("/signup", controller.SignUp(db)).Methods("POST")
	r.HandleFunc("/login", controller.Login(db)).Methods("POST")
	r.HandleFunc("/protected", controller.TokenVerifyMiddleware(controller.ProtectedEndpoint())).Methods("GET")

	log.Println("Listen on port: 8000...")
	log.Fatal(http.ListenAndServe(":8000",r))
}








