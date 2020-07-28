package repository

import (
	"database/sql"
	"fmt"

	"github.com/LevOrlov5404/trip-data-receiver-bot/domain/models"
)

var connectionString string

func ConnectToDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InitConnectionString(dbInfo models.PostgresDb) {
	connectionString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbInfo.Host, dbInfo.Port, dbInfo.User, dbInfo.Password, dbInfo.Name)
}
