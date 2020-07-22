package repository

import (
	"database/sql"
	"time"

	"github.com/LevOrlov5404/trip-data-receiver-bot/models"
)

func GetUser(db *sql.DB, userTelegramID int) (*models.DbUser, error) {
	rows, err := db.Query("select * from user_get($1)", userTelegramID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbUser := models.DbUser{}
	if rows.Next() {
		err := rows.Scan(&dbUser.TelegramID, &dbUser.FullName)
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

func AddUser(db *sql.DB, userTelegramID int, userFullName string) error {
	_, err := db.Exec("select * from user_add($1, $2)", userTelegramID, userFullName)
	if err != nil {
		return err
	}

	return nil
}

func SetUserName(db *sql.DB, userTelegramID int, userFullName string) error {
	_, err := db.Exec("select * from user_set_name($1, $2)", userTelegramID, userFullName)
	if err != nil {
		return err
	}

	return nil
}

func AddTripInfoStart(db *sql.DB, userTelegramID int, date time.Time, km int, filePath string) error {
	_, err := db.Exec("select * from trip_info_start_add($1, $2, $3, $4)", userTelegramID, date, km, filePath)
	if err != nil {
		return err
	}

	return nil
}

func AddTripInfoFinish(db *sql.DB, tripInfoID int64, date time.Time, km int, filePath string) error {
	_, err := db.Exec("select * from trip_info_finish_add($1, $2, $3, $4)", tripInfoID, date, km, filePath)
	if err != nil {
		return err
	}

	return nil
}

func GetNotFininishedTripInfoID(db *sql.DB, userTelegramID int) (int64, error) {
	rows, err := db.Query("select * from get_not_finished_trip_info_id($1)", userTelegramID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var tripInfoID int64
	if rows.Next() {
		err := rows.Scan(&tripInfoID)
		if err != nil {
			return 0, err
		}
	} else {
		return 0, nil
	}

	if err = rows.Err(); err != nil {
		return 0, err
	}
	return tripInfoID, nil
}

func SetFinishedToTripInfo(db *sql.DB, tripInfoID int64) error {
	_, err := db.Exec("select * from set_finished_to_trip_info($1)", tripInfoID)
	if err != nil {
		return err
	}

	return nil
}

func GetKmDifferenceByTripInfoID(db *sql.DB, tripInfoID int64) (int, error) {
	rows, err := db.Query("select * from get_km_difference_by_trip_info_id($1)", tripInfoID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var kmDifference int
	if rows.Next() {
		err := rows.Scan(&kmDifference)
		if err != nil {
			return 0, err
		}
	} else {
		return 0, nil
	}

	if err = rows.Err(); err != nil {
		return 0, err
	}
	return kmDifference, nil
}
