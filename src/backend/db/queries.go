package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error){ 
	connStr := "host=localhost port=5433 user=admin password='Zxczxczxc1' dbname=postgres sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil { 
        return nil, fmt.Errorf("ошибка при открытии подключения: %w", err)
    }

    


    err = db.Ping()
    if err != nil { 
        
        return nil, fmt.Errorf("ошибка при проверке соединения: %w", err)
    }

    fmt.Println("Успешно подключились к базе данных!")
	return db, nil
}