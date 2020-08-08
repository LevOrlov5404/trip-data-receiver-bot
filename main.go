package main

import (
	"log"

	"github.com/LevOrlov5404/trip-data-receiver-bot/domain/models"
	"github.com/LevOrlov5404/trip-data-receiver-bot/domain/repository"
	"github.com/LevOrlov5404/trip-data-receiver-bot/handlers"
	"github.com/LevOrlov5404/trip-data-receiver-bot/infrastructure"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
)

func main() {
	config := infrastructure.ReadConfig()
	users := models.Users{}

	err := repository.InitDB(config.PostgresDb)
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(config.TelegramBotToken)
	if err != nil {
		log.Fatal(err)
	}

	//Устанавливаем время обновления
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	//Получаем обновления от бота
	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if _, ok := users[update.Message.From.ID]; !ok {
			users[update.Message.From.ID] = models.CreateUser(update.Message.From.ID)
		}

		if update.Message.Text != "" || update.Message.Document != nil || update.Message.Photo != nil {
			go func() {
				user := users[update.Message.From.ID]
				user.Mux.Lock()

				replyMsg := ""

				if user.MessageHandlersArray != nil {
					replyMsg, err = user.MessageHandlersArray[user.MessageHandlerNum](bot, update.Message, user)
					if err != nil {
						if clientErr, ok := err.(models.ClientError); ok {
							replyMsg = clientErr.Error()
						} else {
							log.Printf("не удалось обработать сообщение пользователя по причине: %v", err)
							replyMsg = "Внутренняя ошибка. Попробуйте позже."
						}
					} else {
						if user.MessageHandlerNum < len(user.MessageHandlersArray)-1 {
							user.MessageHandlerNum++
						} else {
							user.MessageHandlersArray = nil
							user.MessageHandlerNum = 0
						}
						user.CurrentFail = 0
					}
				} else if update.Message.Text != "" {
					replyMsg, err = handlers.HandleTextMessage(update.Message, user)
					if err != nil {
						log.Printf("не удалось обработать сообщение пользователя по причине: %v", err)
						replyMsg = "Внутренняя ошибка. Попробуйте позже."
					}
				}

				if replyMsg != "" {
					infrastructure.SendMessageToChatID(bot, replyMsg, update.Message.Chat.ID)
				}

				user.Mux.Unlock()
			}()
		}
	}
}
