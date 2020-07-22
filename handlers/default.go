package handlers

import (
	"fmt"
	"strings"

	"github.com/LevOrlov5404/trip-data-receiver-bot/models"
	"github.com/LevOrlov5404/trip-data-receiver-bot/repository"
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
		return "Добрый день. Я бот, принимающий данные о пробеге в начале и в конце поездки. Напишите /help", nil
	}

	if requestMsg == "/help" {
		return "Вначале зарегистрируйтесь, если еще не сделали этого. Затем можно отправлять данные о командировке.\n" +
			"Зарегистрироваться: /reg\n" +
			"Сменить ФИО: /change_name\n" +
			"Отправить отчет: /report", nil
	}

	if requestMsg == "/reg" {
		if user.Registrated {
			return "Вы уже зарегистрированы.", nil
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
			user.Registrated = true
			return fmt.Sprintf("%s, проверил, вы уже зарегистрированы. Отправить данные о поездке: /report", user.FullName), nil
		}

		user.MessageHandlersArray = GetRegistrationHandlers()
		// user.MessageHandlerNum = 0 ?

		return "Регистрация. Напишите ФИО (через пробел, по-русски).", nil
	}

	if requestMsg == "/report" {
		db, err := repository.ConnectToDB()
		if err != nil {
			return "", err
		}
		defer db.Close()

		if !user.Registrated {
			dbUser, err := repository.GetUser(db, user.ID)
			if err != nil {
				return "", err
			}

			if dbUser == nil {
				return "Похоже вы не зарегистрированы. Пройдите процедуру регистрации: /reg", nil
			}
			user.FullName = *dbUser.FullName
			user.Registrated = true
		}

		tripInfoID, err := repository.GetNotFininishedTripInfoID(db, user.ID)
		if err != nil {
			return "", err
		}

		if tripInfoID != 0 {
			user.TripInfo.NotFinishedTripInfoID = tripInfoID
			user.MessageHandlersArray = GetReportChoiceHandlers()
			return "Найдена незавершенная запись о поездке.\n" +
				"Хотите ввести для нее конечные данные? /finish_info\n" +
				"Хотите ввести начальные данные о новой поездке? /new_start_info", nil
		}

		user.MessageHandlersArray = GetReportStartHandlers()

		return "Прием начальных данных о поездке. Пришлите фото или фото в виде файла с показаниями спидометра.", nil
	}

	if requestMsg == "/change_name" {
		db, err := repository.ConnectToDB()
		if err != nil {
			return "", err
		}
		defer db.Close()

		if !user.Registrated {
			dbUser, err := repository.GetUser(db, user.ID)
			if err != nil {
				return "", err
			}

			if dbUser == nil {
				return "Похоже вы не зарегистрированы. Пройдите процедуру регистрации: /reg", nil
			}
			user.FullName = *dbUser.FullName
			user.Registrated = true
		}

		user.MessageHandlersArray = GetChangeNameHandlers()

		return "Смена ФИО. Напишите ФИО (через пробел, по-русски).", nil
	}

	return /*"Не знаю, как ответить"*/ "", nil
}
