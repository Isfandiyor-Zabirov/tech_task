package utils

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"
)

func ReadFile(filePath string) ([]string, error) {
	var (
		data []string
		err  error
	)

	splitFilePath := strings.SplitN(filePath, ".", -1)
	fileExtension := splitFilePath[len(splitFilePath)-1]

	if fileExtension != "json" {
		return []string{}, errors.New("Неверный формат файла. Формат должен быть .json ")
	}

	dataBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("readFile func os.ReadFile error:", err.Error())
		return []string{}, errors.New("Файл не найден ")
	}

	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		log.Println("Неверная структура данных в файле:", err.Error())
		return []string{}, errors.New("Неверная структура данных в файле ")
	}

	if len(data) == 0 {
		return []string{}, errors.New("Файл пустой ")
	}

	return data, nil
}
