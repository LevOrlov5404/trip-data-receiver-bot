package models

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

type (
	UserMessageHandler func(message *tgbotapi.Message, user *User) (string, error)
	User               struct {
		ID                   int
		Registrated          bool
		FullName             string
		CurrentFail          int
		MessageHandlersArray []UserMessageHandler
		MessageHandlerNum    int
	}
	Users map[int]*User
)

func CreateUser(id int) *User {
	return &User{
		ID: id,
	}
}
