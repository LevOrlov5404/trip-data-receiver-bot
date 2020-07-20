package repository

import (
	"database/sql"

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

func AddUser(db *sql.DB, userTelegramId int, userFullName string) (error) {
	_, err := db.Exec("select * from user_add($1, $2)", userTelegramId, userFullName)
	if err != nil {
		return err
	}

	return nil
}
