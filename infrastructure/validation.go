package infrastructure

import (
	"log"
	"regexp"
)

const (
	onlyRussianSymbolsPattern = `^[-А-Яа-яё]*$`
)

func CheckStrHasOnlyRuSymbols(str string) bool {
	matched, err := regexp.Match(onlyRussianSymbolsPattern, []byte(str))
	if err != nil {
		log.Printf("ошибка при проверке на русские символы: %v", err)
		return false
	}
	if !matched {
		return false
	}
	return true
}
