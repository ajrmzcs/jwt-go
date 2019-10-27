package driver

import (
	"database/sql"
	"github.com/ajrmzcs/books/driver"
	"log"
	"os"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ConnectDB() *sql.DB {
	db, err := sql.Open("mysql",
		os.Getenv("DB_USER")+":"+os.Getenv("DB_USER")+"@/"+os.Getenv("DB_NAME"))

	logFatal(err)

	return db
}