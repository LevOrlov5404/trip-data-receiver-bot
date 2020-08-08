package repository

import (
	"time"

	"github.com/LevOrlov5404/trip-data-receiver-bot/domain/models"
)

func GetPassword() (string, error) {
	rows, err := db.Query("select * from registration_password_get_latest()")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	password := ""
	if rows.Next() {
		err := rows.Scan(&password)
		if err != nil {
			return "", err
		}
	} else {
		return "", nil
	}

	if err = rows.Err(); err != nil {
		return "", err
	}

	return password, err
}

func GetUser(userTelegramID int) (*models.DbUser, error) {
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

func GetUserIsBlockedStatus(userTelegramID int) (bool, error) {
	rows, err := db.Query("select * from user_get_is_blocked_status($1)", userTelegramID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	isBlockedStatus := false
	if rows.Next() {
		err := rows.Scan(&isBlockedStatus)
		if err != nil {
			return false, err
		}
	} else {
		return false, nil
	}

	if err = rows.Err(); err != nil {
		return false, err
	}

	return isBlockedStatus, nil
}

func AddUser(userTelegramID int, userFullName string) error {
	_, err := db.Exec("select * from user_add($1, $2)", userTelegramID, userFullName)
	if err != nil {
		return err
	}

	return nil
}

func SetUserName(userTelegramID int, userFullName string) error {
	_, err := db.Exec("select * from user_set_name($1, $2)", userTelegramID, userFullName)
	if err != nil {
		return err
	}

	return nil
}

func AddTripInfoStart(userTelegramID int, date time.Time, km int, filePath string) error {
	_, err := db.Exec("select * from trip_info_start_add($1, $2, $3, $4)", userTelegramID, date, km, filePath)
	if err != nil {
		return err
	}

	return nil
}

func AddTripInfoFinish(tripInfoID int64, date time.Time, km int, filePath string) error {
	_, err := db.Exec("select * from trip_info_finish_add($1, $2, $3, $4)", tripInfoID, date, km, filePath)
	if err != nil {
		return err
	}

	return nil
}

func GetNotFininishedTripInfoID(userTelegramID int) (int64, error) {
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

func SetFinishedToTripInfo(tripInfoID int64) error {
	_, err := db.Exec("select * from set_finished_to_trip_info($1)", tripInfoID)
	if err != nil {
		return err
	}

	return nil
}

func GetKmDifferenceByTripInfoID(tripInfoID int64) (int, error) {
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
