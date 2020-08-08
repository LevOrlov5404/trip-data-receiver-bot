package handlers

import (
	"fmt"
	"strings"

	"github.com/LevOrlov5404/trip-data-receiver-bot/domain/models"
	"github.com/LevOrlov5404/trip-data-receiver-bot/domain/repository"
	"github.com/LevOrlov5404/trip-data-receiver-bot/infrastructure"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var changeNameHandlers []models.UserMessageHandler

func GetChangeNameHandlers() []models.UserMessageHandler {
	if changeNameHandlers == nil {
		initChangeNameHandlers()
	}
	return changeNameHandlers
}
func GetFullNameToChange(bot *tgbotapi.BotAPI, message *tgbotapi.Message, user *models.User) (string, error) {
	if message.Text == "" {
		return handleFail(user)
	}

	msgParts := strings.Split(message.Text, " ")
	if len(msgParts) != 3 || !infrastructure.CheckStrHasOnlyRuSymbols(msgParts[0]) ||
		!infrastructure.CheckStrHasOnlyRuSymbols(msgParts[1]) || !infrastructure.CheckStrHasOnlyRuSymbols(msgParts[2]) {
		return handleFail(user)
	}

	user.FullNameToChange = message.Text

	return fmt.Sprintf("Сменить ФИО с %s на %s?  /yes  |  /no", user.FullName, user.FullNameToChange), nil
}
func GetAnswerToChangeName(bot *tgbotapi.BotAPI, message *tgbotapi.Message, user *models.User) (string, error) {
	if message.Text == "" {
		return handleFail(user)
	}

	answer := strings.ToLower(message.Text)
	if answer != "/yes" && answer != "/no" {
		return handleFail(user)
	}

	if answer == "/no" {
		user.FullNameToChange = ""
		return fmt.Sprintf("ФИО осталось прежнее: %s.", user.FullName), nil
	}

	err := repository.SetUserName(user.ID, user.FullNameToChange)
	if err != nil {
		return "", err
	}

	user.FullName = user.FullNameToChange
	user.FullNameToChange = ""

	return fmt.Sprintf("Сменил ФИО на: %s.", user.FullName), nil
}
func initChangeNameHandlers() {
	changeNameHandlers = []models.UserMessageHandler{}
	changeNameHandlers = append(changeNameHandlers, GetFullNameToChange)
	changeNameHandlers = append(changeNameHandlers, GetAnswerToChangeName)
}
