package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/LevOrlov5404/trip-data-receiver-bot/domain/models"
	"github.com/LevOrlov5404/trip-data-receiver-bot/domain/repository"
	"github.com/LevOrlov5404/trip-data-receiver-bot/infrastructure"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var reportFinishHandlers []models.UserMessageHandler

func GetReportFinishHandlers() []models.UserMessageHandler {
	if reportFinishHandlers == nil {
		initReportFinishHandlers()
	}
	return reportFinishHandlers
}
func GetAnswerToWriteTripInfoFinish(bot *tgbotapi.BotAPI, message *tgbotapi.Message, user *models.User) (string, error) {
	if message.Text == "" {
		return handleFail(user)
	}

	answer := strings.ToLower(message.Text)
	if answer != "/yes" && answer != "/no" {
		return handleFail(user)
	}

	if answer == "/no" {
		user.TripInfo.NotFinishedTripInfoID = 0
		user.TripInfo.TelegramFileID = ""
		user.TripInfo.Km = 0
		return "Сбросил данные. Снова отправить данные о поездке: /report", nil
	}

	fileBytes, err := infrastructure.GetFileFromTelegramByFileID(bot, user.TripInfo.TelegramFileID)
	if err != nil {
		return "", nil
	}

	filePath, err := infrastructure.NewUserFile(fileBytes, user.FullName+"/Конец поездки", time.Now().Format("02_01_2006_15_04_05"))
	if err != nil {
		return "", err
	}

	err = repository.AddTripInfoFinish(user.TripInfo.NotFinishedTripInfoID, time.Now(), user.TripInfo.Km, filePath)
	if err != nil {
		return "", err
	}

	kmDifference, err := repository.GetKmDifferenceByTripInfoID(user.TripInfo.NotFinishedTripInfoID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Записал конечные данные о поездке. Пробег составил: %d км.", kmDifference), nil
}
func initReportFinishHandlers() {
	reportFinishHandlers = []models.UserMessageHandler{}
	reportFinishHandlers = append(reportFinishHandlers, GetPhotoOrFile)
	reportFinishHandlers = append(reportFinishHandlers, GetKm)
	reportFinishHandlers = append(reportFinishHandlers, GetAnswerToWriteTripInfoFinish)
}
