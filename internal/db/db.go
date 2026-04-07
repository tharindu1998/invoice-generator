package db

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func InitDB() (*sqlx.DB, error) {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "root:root@tcp(127.0.0.1:3306)/invoice?parseTime=true"
	}

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}

	DB = db
	return DB, nil
}
