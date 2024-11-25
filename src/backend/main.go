package main

import (
	"backend/db"
	"backend/utils"
	"fmt"
	"log"
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

	
	/* <-------------------------------------------------------------- */

	err = db.InsertUsers(dbConn, "admin", "admin_password")
	if err != nil { 
		log.Fatalf("Ошибка при добавлении пользователя: %s", err)
	}

	fmt.Println("Пользователь успешно добавлен в базу данных!")
}