package models

import (
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type (
	UserMessageHandler func(bot *tgbotapi.BotAPI, message *tgbotapi.Message, user *User) (string, error)
	UserTripInfo       struct {
		NotFinishedTripInfoID int64
		TelegramFileID        string
		Km                    int
	}
	User struct {
		ID                   int
		Registrated          bool
		FullName             string
		CurrentFail          int
		MessageHandlersArray []UserMessageHandler
		MessageHandlerNum    int
		TripInfo             UserTripInfo
		FullNameToChange     string
		Mux                  sync.Mutex
	}
	Users map[int]*User
)

func CreateUser(id int) *User {
	return &User{
		ID: id,
	}
}
