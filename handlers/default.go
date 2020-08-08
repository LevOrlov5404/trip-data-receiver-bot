package handlers

import (
	"fmt"
	"strings"

	"github.com/LevOrlov5404/trip-data-receiver-bot/domain/models"
	"github.com/LevOrlov5404/trip-data-receiver-bot/domain/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

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
		dbUser, err := repository.GetUser(user.ID)
		if err != nil {
			return "", err
		}

		if dbUser != nil {
			user.FullName = *dbUser.FullName

			blockedStatus, err := repository.GetUserIsBlockedStatus(user.ID)
			if err != nil {
				return "", err
			}

			if blockedStatus {
				return fmt.Sprintf("%s, вы заблокированы.", user.FullName), nil
			}

			user.MessageHandlerNum = len(user.MessageHandlersArray)

			return fmt.Sprintf("%s, проверил, вы уже зарегистрированы. Отправить данные о поездке: /report", user.FullName), nil
		}

		user.MessageHandlersArray = GetRegistrationHandlers()

		return "Регистрация. Напишите ФИО (через пробел, по-русски).", nil
	}

	if requestMsg == "/report" {
		dbUser, err := repository.GetUser(user.ID)
		if err != nil {
			return "", err
		}

		if dbUser == nil {
			return "Похоже вы не зарегистрированы. Пройдите процедуру регистрации: /reg", nil
		}

		user.FullName = *dbUser.FullName

		blockedStatus, err := repository.GetUserIsBlockedStatus(user.ID)
		if err != nil {
			return "", err
		}

		if blockedStatus {
			return fmt.Sprintf("%s, вы заблокированы.", user.FullName), nil
		}

		tripInfoID, err := repository.GetNotFininishedTripInfoID(user.ID)
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
		dbUser, err := repository.GetUser(user.ID)
		if err != nil {
			return "", err
		}

		if dbUser == nil {
			return "Похоже вы не зарегистрированы. Пройдите процедуру регистрации: /reg", nil
		}

		user.FullName = *dbUser.FullName

		blockedStatus, err := repository.GetUserIsBlockedStatus(user.ID)
		if err != nil {
			return "", err
		}

		if blockedStatus {
			return fmt.Sprintf("%s, вы заблокированы.", user.FullName), nil
		}

		user.MessageHandlersArray = GetChangeNameHandlers()

		return "Смена ФИО. Напишите ФИО (через пробел, по-русски).", nil
	}

	return "", nil
}
