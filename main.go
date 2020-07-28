package main

import (
	"log"

	"github.com/LevOrlov5404/trip-data-receiver-bot/handlers"
	"github.com/LevOrlov5404/trip-data-receiver-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
)

const (
	telegramBotToken string = "1273007508:AAH_hqU_qFimVVv1OW5jYdJmgbfqhCr-2-g"
)

func sendMessageToChatId(bot *tgbotapi.BotAPI, message string, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("не удалось отправить сообщение по причине: %v", err)
	}
}

func main() {
	// db, err := repository.ConnectToDB()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer db.Close()

	// status, err := repository.GetUserIsBlockedStatus(db, 0)
	// fmt.Println(err)
	// fmt.Println(status)
	// return

	users := models.Users{}

	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
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
					sendMessageToChatId(bot, replyMsg, update.Message.Chat.ID)
				}

				user.Mux.Unlock()
			}()
		}
	}
}
