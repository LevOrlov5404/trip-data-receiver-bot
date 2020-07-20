package handlers

import (
	"fmt"
	"strings"

	"github.com/LevOrlov5404/trip-data-receiver-bot/infrastructure"
	"github.com/LevOrlov5404/trip-data-receiver-bot/models"
	"github.com/LevOrlov5404/trip-data-receiver-bot/repository"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

var registrationHandlers []models.UserMessageHandler

func GetRegistrationHandlers() []models.UserMessageHandler {
	if registrationHandlers == nil {
		initRegistrationHandlers()
	}
	return registrationHandlers
}
func GetFullNameToRegistrate(message *tgbotapi.Message, user *models.User) (string, error) {
	if message.Text == "" {
		return handleFail(user)
	}

	msgParts := strings.Split(message.Text, " ")
	if len(msgParts) != 3 || !infrastructure.CheckStrHasOnlyRuSymbols(msgParts[0]) ||
		!infrastructure.CheckStrHasOnlyRuSymbols(msgParts[1]) || !infrastructure.CheckStrHasOnlyRuSymbols(msgParts[2]) {
		return handleFail(user)
	}

	db, err := repository.ConnectToDB()
	if err != nil {
		return "", err
	}
	defer db.Close()

	dbUser, err := repository.GetUser(db, user.ID)
	if err != nil {
		return "", err
	}

	if dbUser != nil {
		user.FullName = *dbUser.FullName
		user.MessageHandlerNum = len(user.MessageHandlersArray)
		return "Проверил, вы есть в базе. Значит вы уже зарегистрированы.", nil
	}

	user.FullName = message.Text

	return "Введите пароль для регистрации", nil
}
func GetPasswordToRegistrate(message *tgbotapi.Message, user *models.User) (string, error) {
	if message.Text == "" || message.Text != "12345" {
		return handleFail(user)
	}

	db, err := repository.ConnectToDB()
	if err != nil {
		return "", err
	}
	defer db.Close()

	err = repository.AddUser(db, user.ID, user.FullName)
	if err != nil {
		return "", err
	}

	user.Registrated = true

	return "Вы успешно зарегистрированы", nil
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
