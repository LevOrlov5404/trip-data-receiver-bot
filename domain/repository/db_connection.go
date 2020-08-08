package repository

import (
	"database/sql"
	"fmt"

	"github.com/LevOrlov5404/trip-data-receiver-bot/domain/models"
)

var (
	db *sql.DB
)

func InitDB(dbInfo models.PostgresDb) error {
	var err error
	db, err = ConnectToDB(dbInfo)
	return err
}
func ConnectToDB(dbInfo models.PostgresDb) (*sql.DB, error) {
	db, err := sql.Open("postgres", InitConnectionString(dbInfo))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(dbInfo.MaxOpenConns)

	return db, nil
}
func InitConnectionString(dbInfo models.PostgresDb) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbInfo.Host, dbInfo.Port, dbInfo.User, dbInfo.Password, dbInfo.Name)
}
