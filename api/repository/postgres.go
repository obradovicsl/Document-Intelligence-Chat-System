package repository

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
    connStr := os.Getenv("DB_URL")

    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }

    if err = DB.Ping(); err != nil {
        panic(err)
    }

    fmt.Println("Connected to NeonDB")
}