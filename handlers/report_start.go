package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/LevOrlov5404/trip-data-receiver-bot/infrastructure"
	"github.com/LevOrlov5404/trip-data-receiver-bot/models"
	"github.com/LevOrlov5404/trip-data-receiver-bot/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var reportStartHandlers []models.UserMessageHandler

func GetReportStartHandlers() []models.UserMessageHandler {
	if reportStartHandlers == nil {
		initReportStartHandlers()
	}
	return reportStartHandlers
}
func GetPhotoOrFile(bot *tgbotapi.BotAPI, message *tgbotapi.Message, user *models.User) (string, error) {
	if (message.Document == nil && message.Photo == nil) || message.Text != "" {
		return handleFail(user)
	}

	if message.Document != nil {
		user.TripInfo.TelegramFileID = message.Document.FileID
	} else {
		user.TripInfo.TelegramFileID = (*message.Photo)[0].FileID
	}

	return "Принял ваше фото. Теперь введите пробег (км, целое число).", nil
}
func GetKm(bot *tgbotapi.BotAPI, message *tgbotapi.Message, user *models.User) (string, error) {
	if message.Text == "" {
		return handleFail(user)
	}

	km, ok := infrastructure.GetIntFromString(message.Text)
	if !ok || km < 0 {
		return handleFail(user)
	}

	user.TripInfo.Km = km

	return fmt.Sprintf("%s, принял ваш пробег: %d. Записать данные? /yes  |  /no", user.FullName, km), nil
}
func GetAnswerToWriteTripInfoStart(bot *tgbotapi.BotAPI, message *tgbotapi.Message, user *models.User) (string, error) {
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

	filePath, err := infrastructure.NewUserFile(fileBytes, user.FullName+"/Начало поездки", time.Now().Format("02_01_2006_15_04_05"))
	if err != nil {
		return "", err
	}

	db, err := repository.ConnectToDB()
	if err != nil {
		return "", err
	}
	defer db.Close()

	err = repository.AddTripInfoStart(db, user.ID, time.Now(), user.TripInfo.Km, filePath)
	if err != nil {
		return "", err
	}

	return "Записал начальные данные о поездке. Доброго пути!", nil
}
func initReportStartHandlers() {
	reportStartHandlers = []models.UserMessageHandler{}
	reportStartHandlers = append(reportStartHandlers, GetPhotoOrFile)
	reportStartHandlers = append(reportStartHandlers, GetKm)
	reportStartHandlers = append(reportStartHandlers, GetAnswerToWriteTripInfoStart)
}
