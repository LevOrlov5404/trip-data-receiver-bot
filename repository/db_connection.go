package repository

import (
    "database/sql"
    "fmt"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "bot_user"
	password = "12345"
	dbname   = "trip-data"
)  

func ConnectToDB() *sql.DB {
	db, err := sql.Open("postgres", initConnection())
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func initConnection() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}
