package main

import (
	"log"

	"github.com/LevOrlov5404/trip-data-receiver-bot/handlers"
	"github.com/LevOrlov5404/trip-data-receiver-bot/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
)

var (
	telegramBotToken string       = "1273007508:AAH_hqU_qFimVVv1OW5jYdJmgbfqhCr-2-g"
	users            models.Users = models.Users{}
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
	// defer db.Close()

	// err = repository.TripInfoStartAdd(db, 0, time.Now(), 123, "lol")
	// fmt.Println(err)
	// err = repository.TripInfoFinishAdd(db, 5, time.Now(), 123, "lol")
	// fmt.Println(err)
	// tripInfoID, err := repository.GetNotFininishedTripInfoID(db, 1)
	// fmt.Println(err)
	// fmt.Println(tripInfoID)
	// err = repository.SetFinishedToTripInfo(db, 9223372036854775807)
	// fmt.Println(err)
	// kmDifference, err := repository.GetKmDifferenceByTripInfoID(db, 9223372036854775807)
	// fmt.Println(err)
	// fmt.Println(kmDifference)

	// err = repository.AddUser(db, 0, "levchik")
	// fmt.Println(err)
	// err = repository.SetUserName(db, 2, "test2")
	// fmt.Println(err)
	// dbUser, err := repository.GetUser(db, 0)
	// fmt.Println(err)
	// fmt.Println(dbUser)
	// return

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
			replyMsg := ""

			user := users[update.Message.From.ID]
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
			// } else if update.Message.Document != nil {
			// 	fmt.Println("have document")
			// 	fileBytes, err := inftastracture.GetFileFromTelegramByFileID(bot, update.Message.Document.FileID)
			// 	if err != nil {
			// 		// handle error
			// 		return
			// 	}

			// 	err = inftastracture.NewFileWithPath(fileBytes, "/home/lev/Documents/ImageReceiverBotDocuments/"+time.Now().Format("01_02_2006_15_04_05"))
			// 	if err != nil {
			// 		// handle error
			// 		return
			// 	}

			// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "принял ваш файл")
			// 	_, err = bot.Send(msg)
			// 	if err != nil {
			// 		log.Printf("не удалось отправить сообщение по причине: %v", err)
			// 		return
			// 	}
			// } else if update.Message.Photo != nil {
			// 	fmt.Println("have photo")
			// 	fileBytes, err := inftastracture.GetFileFromTelegramByFileID(bot, (*update.Message.Photo)[0].FileID)
			// 	if err != nil {
			// 		// handle error
			// 		return
			// 	}

			// 	err = inftastracture.NewFileWithPath(fileBytes, "/home/lev/Documents/ImageReceiverBotDocuments/"+time.Now().Format("01_02_2006_15_04_05"))
			// 	if err != nil {
			// 		// handle error
			// 		return
			// 	}

			// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "принял ваше фото")
			// 	_, err = bot.Send(msg)
			// 	if err != nil {
			// 		log.Printf("не удалось отправить сообщение по причине: %v", err)
			// 		return
			// 	}
		}
	}
}
