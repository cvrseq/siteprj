package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func queries(){ 
	connStr := "host=localhost port=5433 user=admin password='Zxczxczxc1' dbname=postgres sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil { 
        fmt.Println("Ошибка при открытии подключения: ", err)
        return 
    }

    defer db.Close()


    err = db.Ping()
    if err != nil { 
        fmt.Println("Ошибка при проверке соединения: ", err)
        return 
    }

    fmt.Println("Успешно подключились к базе данных!")
}