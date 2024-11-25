package main

import (
	"backend/db"
	"backend/utils"
	"fmt"
)

func main() {
	dbConn, err := db.ConnectDB()
	if err != nil {
		fmt.Println("Ошибка подключения к базе данных:", err)
		return
	}
	defer dbConn.Close()

	fmt.Println("Сервер успешно запущен!")

	/* <-------------------------------------------------------------- */

	password := "secret"
	hash, _ := utils.HashPassword(password)

	fmt.Println("Password:", password)
    fmt.Println("Hash:    ", hash)

    match := utils.CheckPasswordHash(password, hash)
    fmt.Println("Match:   ", match)

}