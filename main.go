package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/LevOrlov5404/trip-data-receiver-bot/repository"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	_ "github.com/lib/pq"
)

type (
	UserMessageHandler func(message *tgbotapi.Message, user *User) (replyMsg string, err error)
	User struct{
		Id int
		CurrentMessageHandler UserMessageHandler
	}
	Users map[int]*User
)

var (
	telegramBotToken string = "1273007508:AAH_hqU_qFimVVv1OW5jYdJmgbfqhCr-2-g"
	defaultMessageHandler = func(message *tgbotapi.Message, user *User) (replyMsg string, err error) {return}
	users Users = Users{}
)

func CreateUser(id int) *User {
	return &User{
		Id: id,
		CurrentMessageHandler: defaultMessageHandler,
	}
}

func handleTextMessage(message *tgbotapi.Message, users Users) (replyMsg string, err error) {
	requestMsg := message.Text
	userInfo := message.From
	
	if strings.Contains(requestMsg, "Привет") || strings.Contains(requestMsg, "привет") || strings.Contains(requestMsg, "/start") {
		replyMsg = "Привет. Я бот, принимающий данные о командировке. Напиши /help"
	} else if requestMsg == "/help" {
		replyMsg = "Вначале зарегистрируйся, если еще не сделал этого. Затем можно отправлять данные о командировке.\n"+
			"Зарегистрироваться: /registrate\n"+
			"Отправить отчет: /report"
	} else if requestMsg == "/registrate" {
		if _, ok := users[userInfo.ID]; !ok {
			users[userInfo.ID] = CreateUser(userInfo.ID)
		}
	} else if requestMsg == "/report" {
		replyMsg = "smth with report"
	} else {
		replyMsg = "Не знаю, как ответить"
	}
	return
}

func main() {
	db := repository.ConnectToDB()
	defer db.Close()

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

		if update.Message.Text != "" {
			replyMsg, err := handleTextMessage(update.Message, users)
			if err != nil {
				log.Printf("не удалось обработать сообщение пользователя по причине: %v", err)
				return
			}

			// replyMsg = fmt.Sprintf("%s %d %s", update.Message.From.UserName, update.Message.From.ID, replyMsg)
			if replyMsg != "" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, replyMsg)
				_, err := bot.Send(msg)
				if err != nil {
					log.Printf("не удалось отправить сообщение по причине: %v", err)
				}
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
