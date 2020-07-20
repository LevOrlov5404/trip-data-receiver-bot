package handlers

import (
	"strings"

	"github.com/LevOrlov5404/trip-data-receiver-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// type smth struct {
// 	lol UserMessageHandler
// }

// var (
// 	smth_lol smth = smth{
// 		lol: GetFullNameToRegistrate,
// 	}
// )

func HandleTextMessage(message *tgbotapi.Message, user *models.User) (string, error) {
	requestMsg := message.Text

	reqMsgLowerCase := strings.ToLower(requestMsg)
	if strings.Contains(reqMsgLowerCase, "привет") || strings.Contains(reqMsgLowerCase, "доров") ||
		strings.Contains(reqMsgLowerCase, "добрый день") || strings.Contains(requestMsg, "/start") {
		return "Добрый день. Я бот, принимающий данные о командировке. Напишите /help", nil
	}

	if requestMsg == "/help" {
		return "Вначале зарегистрируйтесь, если еще не сделали этого. Затем можно отправлять данные о командировке.\n" +
			"Зарегистрироваться: /reg\n" +
			"Отправить отчет: /report", nil
	}

	if requestMsg == "/reg" {
		if user.Registrated {
			return "Вы уже зарегистрированы.", nil
		}
		user.MessageHandlersArray = GetRegistrationHandlers()
		// user.MessageHandlerNum = 0 ?

		return "Напишите ФИО (через пробел, по-русски).", nil
	}

	if requestMsg == "/report" {
		return "smth with report", nil
	}

	return /*"Не знаю, как ответить"*/ "", nil
}
