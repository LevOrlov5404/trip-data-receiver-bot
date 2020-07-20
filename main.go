package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/LevOrlov5404/trip-data-receiver-bot/handlers"
	"github.com/LevOrlov5404/trip-data-receiver-bot/models"

	// "github.com/LevOrlov5404/trip-data-receiver-bot/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
)

var (
	telegramBotToken string       = "1273007508:AAH_hqU_qFimVVv1OW5jYdJmgbfqhCr-2-g"
	users            models.Users = models.Users{}
)

func sendMessageToChatIdByBot(message string, chatID int64, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("не удалось отправить сообщение по причине: %v", err)
	}
}

func main() {
	// db, err := repository.ConnectToDB()
	// defer db.Close()

	// err = repository.AddUser(db, 0, "levchik")
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

		if update.Message.Text != "" {
			replyMsg := ""

			user := users[update.Message.From.ID]
			if user.MessageHandlersArray != nil {
				replyMsg, err = user.MessageHandlersArray[user.MessageHandlerNum](update.Message, user)
				if err != nil {
					if clientErr, ok := err.(models.ClientError); ok {
						replyMsg = clientErr.Error()
					} else {
						log.Printf("не удалось обработать сообщение пользователя по причине: %v", err)
						replyMsg = "Внутренняя ошибка. Попробуйте позже."
					}
				} else {
					if user.MessageHandlerNum < len(user.MessageHandlersArray) - 1 {
						user.MessageHandlerNum++
					} else {
						user.MessageHandlersArray = nil
						user.MessageHandlerNum = 0
					}
					user.CurrentFail = 0
				}
			} else {
				replyMsg, err = handlers.HandleTextMessage(update.Message, user)
				if err != nil {
					log.Printf("не удалось обработать сообщение пользователя по причине: %v", err)
					replyMsg = "Внутренняя ошибка. Попробуйте позже."
				}
			}

			if replyMsg != "" {
				sendMessageToChatIdByBot(replyMsg, update.Message.Chat.ID, bot)
			}
		} else if update.Message.Document != nil {
			getFileURL, err := bot.GetFileDirectURL(update.Message.Document.FileID)
			if err != nil {
				log.Printf("не удалось получить ссылку для скачивания файла по причине: %v", err)
				return
			}
			fmt.Println(getFileURL)

			resp, err := http.Get(getFileURL)
			if err != nil {
				log.Printf("не удалось получить запросом файл по причине: %v", err)
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("не удалось считать файл по причине: %v", err)
				return
			}

			newFile, err := os.Create("/home/lev/Documents/ImageReceiverBotDocuments/" + update.Message.Document.FileName)
			if err != nil {
				log.Printf("не удалось создать новый файл по причине: %v", err)
				return
			}
			defer func() {
				if err := newFile.Close(); err != nil {
					log.Printf("не удалось закрыть новый файл по причине: %v", err)
				}
			}()

			_, err = newFile.Write(body)
			if err != nil {
				log.Printf("не удалось записать в файл по причине: %v", err)
				return
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "принял ваш файл")
			_, err = bot.Send(msg)
			if err != nil {
				log.Printf("не удалось отправить сообщение по причине: %v", err)
				return
			}
		}
	}
}
