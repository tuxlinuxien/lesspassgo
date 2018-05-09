package server

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/mattn/go-sqlite3" // db sqlite3
)

const tables = `
    CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        email TEXT NOT NULL UNIQUE,
        password TEXT
    );

    CREATE TABLE IF NOT EXISTS passwords (
        id TEXT UNIQUE PRIMARY KEY,
        user_id TEXT NOT NULL,
        login TEXT NOT NULL,
        site TEXT NOT NULL,
        uppercase BOOLEAN DEFAULT TRUE,
        symbols BOOLEAN DEFAULT TRUE,
        lowercase BOOLEAN DEFAULT TRUE,
        numbers BOOLEAN DEFAULT TRUE,
        counter INTEGER DEFAULT 1,
        version INTEGER DEFAULT 2,
        length INTEGER DEFAULT 16,
        created TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        modified TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );
    CREATE UNIQUE INDEX IF NOT EXISTS passwords_unique ON passwords (user_id, login, site);
    CREATE INDEX IF NOT EXISTS passwords_user_id ON passwords (user_id);
`

var (
	conn *sql.DB
	// ErrUserNotFound .
	ErrUserNotFound = errors.New("user not found")
	// ErrPasswordNotFound .
	ErrPasswordNotFound = errors.New("password not found")
)

func openDB(path string) {
	_conn, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalln(err)
	}
	conn = _conn
}

func createTable() {
	_, err := conn.Exec(tables)
	if err != nil {
		log.Fatalln(err)
	}
}
