package infrastructure

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const pathPrefix = "/home/lev/Documents/ImageReceiverBotDocuments/"

func GetFileFromTelegramByFileID(bot *tgbotapi.BotAPI, fileID string) ([]byte, error) {
	getFileURL, err := bot.GetFileDirectURL(fileID)
	if err != nil {
		return nil, err
	}
	fmt.Println(getFileURL)

	resp, err := http.Get(getFileURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

func NewUserFile(fileBytes []byte, fileFolder string, fileName string) (string, error) {
	folderPath := pathPrefix+fileFolder
	
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err = os.MkdirAll(folderPath, 0777)
		if err != nil {
			return "", err
		}
	}
	
	filePath := folderPath+"/"+fileName
	newFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer newFile.Close()

	_, err = newFile.Write(fileBytes)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func GetIntFromString(s string) (int, bool) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}

	return i, true
}
