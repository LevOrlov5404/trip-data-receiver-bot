package repository

import (
	"database/sql"
	"fmt"
)

var (
	host     string = "172.17.0.2"//"localhost"
	port     int    = 5432
	user     string = "bot_user"
	password string = "12345"
	dbname   string = "trip-data"
)

func ConnectToDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", initConnection())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func initConnection() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}
