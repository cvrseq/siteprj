package main

import (
	"fmt"
)

func main() {
	// Подключаемся к базе данных
	dbConn, err := db.ConnectDB()
	if err != nil {
		fmt.Println("Ошибка подключения к базе данных:", err)
		return
	}
	defer dbConn.Close()

	fmt.Println("Сервер успешно запущен!")
	// Здесь можно вызвать другие функции для обработки запросов
}