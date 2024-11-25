package db

import (
	"backend/utils"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error){ 
	connStr := "host=localhost port=5433 user=admin password='Zxczxczxc1' dbname=postgres sslmode=disable"
    return sql.Open("postgres", connStr)
    
}

func InsertUsers(db *sql.DB, login, password string)  error { 
	hashedPassword, err := utils.HashPassword(password)
	
	if err != nil { 
		return fmt.Errorf("Ошибка хэширования пароля: %w", err)
	}

	query := `INSERT INTO product(login, password) VALUES ($1, $2)`
	_, err = db.Exec(query, login, hashedPassword)

	if err != nil { 
		return fmt.Errorf("Ошибка добавления пользователя: %w", err)
	}
	return nil
}