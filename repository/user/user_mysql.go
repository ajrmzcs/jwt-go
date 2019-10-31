package userRepository

import (
	"database/sql"
	"github.com/ajrmzcs/jwt-go/models"
	"log"
)

type UserRepository struct{}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (u UserRepository) SignUp(db * sql.DB, user models.User) (models.User, error) {
	res, err := db.Exec("INSERT INTO users (email, password) VALUES(?,?)",
		user.Email, user.Password)

	userId, err := res.LastInsertId()

	user.Id= userId
	user.Password=""

	return user, err
}

func (u UserRepository) Login(db * sql.DB, user models.User) (models.User, error) {
	row := db.QueryRow("SELECT * FROM users WHERE email=?", user.Email)

	err := row.Scan(&user.Id, &user.Email, &user.Password)

	if err != nil {
		return user, err
	}

	return user, nil
}
