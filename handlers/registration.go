package handlers

import (
	"fmt"
	"strings"

	"github.com/LevOrlov5404/trip-data-receiver-bot/domain/models"
	"github.com/LevOrlov5404/trip-data-receiver-bot/domain/repository"
	"github.com/LevOrlov5404/trip-data-receiver-bot/infrastructure"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var registrationHandlers []models.UserMessageHandler

func GetRegistrationHandlers() []models.UserMessageHandler {
	if registrationHandlers == nil {
		initRegistrationHandlers()
	}
	return registrationHandlers
}
func GetFullNameToRegistrate(bot *tgbotapi.BotAPI, message *tgbotapi.Message, user *models.User) (string, error) {
	if message.Text == "" {
		return handleFail(user)
	}

	msgParts := strings.Split(message.Text, " ")
	if len(msgParts) != 3 || !infrastructure.CheckStrHasOnlyRuSymbols(msgParts[0]) ||
		!infrastructure.CheckStrHasOnlyRuSymbols(msgParts[1]) || !infrastructure.CheckStrHasOnlyRuSymbols(msgParts[2]) {
		return handleFail(user)
	}

	user.FullName = message.Text

	return "Введите пароль для регистрации", nil
}
func GetPasswordToRegistrate(bot *tgbotapi.BotAPI, message *tgbotapi.Message, user *models.User) (string, error) {
	if message.Text == "" {
		return handleFail(user)
	}

	db, err := repository.ConnectToDB()
	if err != nil {
		return "", err
	}
	defer db.Close()

	password, err := repository.GetPassword(db)
	if err != nil {
		return "", err
	}

	if message.Text != password {
		return handleFail(user)
	}

	err = repository.AddUser(db, user.ID, user.FullName)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s, вы успешно зарегистрированы. Отправить данные о поездке: /report", user.FullName), nil
}
func handleFail(user *models.User) (string, error) {
	user.CurrentFail++

	if user.CurrentFail == 3 {
		user.CurrentFail = 0
		user.MessageHandlersArray = nil
		user.MessageHandlerNum = 0

		return "", models.ClientError{
			Description: "Попытки на текущем шаге закончились. Операция прервана."}
	}

	return "", models.ClientError{
		Description: fmt.Sprintf("Попробуйте еще раз. Осталось попыток: %d.", 3-user.CurrentFail)}
}
func initRegistrationHandlers() {
	registrationHandlers = []models.UserMessageHandler{}
	registrationHandlers = append(registrationHandlers, GetFullNameToRegistrate)
	registrationHandlers = append(registrationHandlers, GetPasswordToRegistrate)
}
