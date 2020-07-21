package repository

import (
	"database/sql"
	"time"

	"github.com/LevOrlov5404/trip-data-receiver-bot/models"
)

func GetUser(db *sql.DB, userTelegramId int) (*models.DbUser, error) {
	rows, err := db.Query("select * from user_get($1)", userTelegramId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbUser := models.DbUser{}
	if rows.Next() {
		err := rows.Scan(&dbUser.TelegramId, &dbUser.FullName)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &dbUser, nil
}

func AddUser(db *sql.DB, userTelegramID int, userFullName string) (error) {
	_, err := db.Exec("select * from user_add($1, $2)", userTelegramID, userFullName)
	if err != nil {
		return err
	}

	return nil
}

func TripInfoStartAdd(db *sql.DB, userTelegramID int, date time.Time, km int, filePath string) (error) {
	_, err := db.Exec("select * from trip_info_start_add($1, $2, $3, $4)", userTelegramID, date, km, filePath)
	if err != nil {
		return err
	}

	return nil
}

func TripInfoFinishAdd(db *sql.DB, tripInfoID int, date time.Time, km int, filePath string) (error) {
	_, err := db.Exec("select * from trip_info_finish_add($1, $2, $3, $4)", tripInfoID, date, km, filePath)
	if err != nil {
		return err
	}

	return nil
}
