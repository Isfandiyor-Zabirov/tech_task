package main

import (
	"log"
	"tech_task/internal/sber"
)

func main() {
	err := sber.SendRequest("./inputs.json")
	if err != nil {
		log.Println("Error:", err.Error())
	}
}
