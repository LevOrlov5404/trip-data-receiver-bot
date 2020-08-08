package handlers

import (
	"strings"

	"github.com/LevOrlov5404/trip-data-receiver-bot/domain/models"
	"github.com/LevOrlov5404/trip-data-receiver-bot/domain/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var reportChoiceHandlers []models.UserMessageHandler

func GetReportChoiceHandlers() []models.UserMessageHandler {
	if reportChoiceHandlers == nil {
		initReportChoiceHandlers()
	}
	return reportChoiceHandlers
}
func GetAnswerToContinueStartOrFinish(bot *tgbotapi.BotAPI, message *tgbotapi.Message, user *models.User) (string, error) {
	if message.Text == "" {
		return handleFail(user)
	}

	answer := strings.ToLower(message.Text)
	if answer != "/finish_info" && answer != "/new_start_info" {
		return handleFail(user)
	}

	user.MessageHandlerNum = -1

	if answer == "/finish_info" {
		user.MessageHandlersArray = GetReportFinishHandlers()
		return "Прием конечных данных о поездке. Пришлите фото или фото в виде файла с показаниями спидометра.", nil
	}

	err := repository.SetFinishedToTripInfo(user.TripInfo.NotFinishedTripInfoID)
	if err != nil {
		return "", err
	}

	user.MessageHandlersArray = GetReportStartHandlers()

	return "Прием начальных данных о поездке. Пришлите фото или фото в виде файла с показаниями спидометра.", nil
}
func initReportChoiceHandlers() {
	reportChoiceHandlers = []models.UserMessageHandler{}
	reportChoiceHandlers = append(reportChoiceHandlers, GetAnswerToContinueStartOrFinish)
}
