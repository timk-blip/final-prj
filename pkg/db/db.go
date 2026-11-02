package db

import (
	"database/sql"
	_ "embed"

	_ "modernc.org/sqlite"
)

//go:embed scheduler.sql
var schema string

var DB *sql.DB
var err error

func Init(dbFile string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	DB = db
	return db, nil
}

func AddTask(date, title, comment, repeat string) (int64, error) {
	res, err := DB.Exec(`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		date, title, comment, repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
