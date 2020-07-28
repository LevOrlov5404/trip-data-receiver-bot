package models

type Config struct {
	TelegramBotToken string
	PostgresDb       PostgresDb
}

type PostgresDb struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}
