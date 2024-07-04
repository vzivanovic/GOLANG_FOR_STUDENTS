package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitLocationDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./locations.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS user_locations (
		username TEXT PRIMARY KEY,
		latitude REAL,
		longitude REAL
	);
	`
	_, err = DB.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

func InitLocationHistoryDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./location_history.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS location_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		latitude REAL,
		longitude REAL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = DB.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

func CloseDB() {
	if err := DB.Close(); err != nil {
		log.Fatalf("Failed to close database: %v", err)
	}
}
