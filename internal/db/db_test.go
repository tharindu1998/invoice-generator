package db

import (
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestInitDB_Success(t *testing.T) {
	os.Setenv("MYSQL_DSN", "root:root@tcp(127.0.0.1:3306)/invoice?parseTime=true")

	db, err := InitDB()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if db == nil {
		t.Fatal("expected db not to be nil")
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("expected db.Ping() to succeed, got %v", err)
	}
}

func TestInitDB_InvalidDSN(t *testing.T) {
	os.Setenv("MYSQL_DSN", "invalid_dsn")

	db, err := InitDB()
	if err == nil {
		t.Fatal("expected error for invalid DSN, got nil")
	}

	if db != nil {
		t.Fatal("expected db to be nil when DSN invalid")
	}
}
