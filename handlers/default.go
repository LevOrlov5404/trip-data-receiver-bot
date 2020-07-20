package handlers

import (
	"strings"

	"github.com/LevOrlov5404/trip-data-receiver-bot/models"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
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

	if strings.Contains(requestMsg, "Привет") || strings.Contains(requestMsg, "привет") || strings.Contains(requestMsg, "/start") {
		return "Привет. Я бот, принимающий данные о командировке. Напиши /help", nil
	}

	if requestMsg == "/help" {
		return "Вначале зарегистрируйся, если еще не сделал этого. Затем можно отправлять данные о командировке.\n" +
			"Зарегистрироваться: /registrate\n" +
			"Отправить отчет: /report", nil
	}

	if requestMsg == "/registrate" {
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

	return /*"Не знаю, как ответить"*/"", nil
}
