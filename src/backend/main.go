package main

import (
	"backend/db"

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

	/*password := "secret"
	hash, _ := hashPassword(password)

	fmt.Println("Password:", password)
    fmt.Println("Hash:    ", hash)

    match := checkPasswordHash(password, hash)
    fmt.Println("Match:   ", match)*/

}