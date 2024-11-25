package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error){ 
	connStr := "host=localhost port=5433 user=admin password='Zxczxczxc1' dbname=postgres sslmode=disable"
    return sql.Open("postgres", connStr)
    
}