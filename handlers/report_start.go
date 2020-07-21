package handlers

// import (
// 	"fmt"
// 	"strings"

// 	"github.com/LevOrlov5404/trip-data-receiver-bot/infrastructure"
// 	"github.com/LevOrlov5404/trip-data-receiver-bot/models"
// 	"github.com/LevOrlov5404/trip-data-receiver-bot/repository"
// 	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
// )

// var reportStartHandlers []models.UserMessageHandler

// func GetReportStartHandlers() []models.UserMessageHandler {
// 	if reportStartHandlers == nil {
// 		initReportStartHandlers()
// 	}
// 	return reportStartHandlers
// }
// func GetFullNameToRegistrate(message *tgbotapi.Message, user *models.User) (string, error) {
// 	if message.Text == "" {
// 		return handleFail(user)
// 	}

// 	msgParts := strings.Split(message.Text, " ")
// 	if len(msgParts) != 3 || !infrastructure.CheckStrHasOnlyRuSymbols(msgParts[0]) ||
// 		!infrastructure.CheckStrHasOnlyRuSymbols(msgParts[1]) || !infrastructure.CheckStrHasOnlyRuSymbols(msgParts[2]) {
// 		return handleFail(user)
// 	}

// 	db, err := repository.ConnectToDB()
// 	if err != nil {
// 		return "", err
// 	}
// 	defer db.Close()

// 	dbUser, err := repository.GetUser(db, user.ID)
// 	if err != nil {
// 		return "", err
// 	}

// 	if dbUser != nil {
// 		user.FullName = *dbUser.FullName
// 		user.MessageHandlerNum = len(user.MessageHandlersArray)
// 		user.Registrated = true
// 		return "Проверил, вы есть в базе. Значит вы уже зарегистрированы.", nil
// 	}

// 	user.FullName = message.Text

// 	return "Введите пароль для регистрации", nil
// }
// func GetPasswordToRegistrate(message *tgbotapi.Message, user *models.User) (string, error) {
// 	if message.Text == "" || message.Text != "12345" {
// 		return handleFail(user)
// 	}

// 	db, err := repository.ConnectToDB()
// 	if err != nil {
// 		return "", err
// 	}
// 	defer db.Close()

// 	err = repository.AddUser(db, user.ID, user.FullName)
// 	if err != nil {
// 		return "", err
// 	}

// 	user.Registrated = true

// 	return "Вы успешно зарегистрированы", nil
// }
// func initReportStartHandlers() {
// 	reportStartHandlers = []models.UserMessageHandler{}
// 	reportStartHandlers = append(registrationHandlers, GetFullNameToRegistrate)
// 	reportStartHandlers = append(registrationHandlers, GetPasswordToRegistrate)
// }
