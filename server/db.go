package server

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3" // db sqlite3
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
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
    CREATE UNIQUE INDEX passwords_unique ON passwords (user_id, login, site);
    CREATE INDEX passwords_user_id ON passwords (user_id);
`

var (
	conn *sql.DB
	// ErrUserNotFound .
	ErrUserNotFound = errors.New("user not found")
	// ErrPasswordNotFound .
	ErrPasswordNotFound = errors.New("password not found")
)

// UserModel .
type UserModel struct {
	ID       string    `json:"id"`
	Email    string    `json:"email"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

// CreateUser .
func CreateUser(email, passord string) error {
	passordBytes, err := bcrypt.GenerateFromPassword([]byte(passord), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = conn.Exec("INSERT INTO users (id, email, password) VALUES ($1, $2, $3);", uuid.NewV4().String(), email, string(passordBytes))
	return err
}

// AuthUser .
func AuthUser(email, password string) (*UserModel, error) {
	rows, err := conn.Query("SELECT id, password FROM users WHERE email = $1 LIMIT 1", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, ErrUserNotFound
	}
	var passwordBytes []byte
	var id string
	if err := rows.Scan(&id, &passwordBytes); err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword(passwordBytes, []byte(password)); err != nil {
		return nil, err
	}
	return &UserModel{ID: id, Email: email}, nil
}

// GetUserByID .
func GetUserByID(id string) (*UserModel, error) {
	rows, err := conn.Query("SELECT id, email FROM users WHERE id = $1 LIMIT 1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, ErrUserNotFound
	}
	u := &UserModel{}
	err = rows.Scan(&u.ID, &u.Email)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// PasswordModel .
type PasswordModel struct {
	ID        string `json:"id"`
	UserID    string `json:"-"`
	Login     string `json:"login"`
	Site      string `json:"site"`
	Uppercase bool   `json:"uppercase"`
	Symbols   bool   `json:"symbols"`
	Lowercase bool   `json:"lowercase"`
	Numbers   bool   `json:"numbers"`
	Counter   int    `json:"counter"`
	Version   int    `json:"version"`
	Length    int    `json:"length"`
}

// Update .
func (p *PasswordModel) Update() error {
	res, err := conn.Exec(
		"UPDATE passwords SET login=$1, site=$2, uppercase=$3, lowercase=$4, symbols=$5, numbers=$6, counter=$7, version=$8, length=$9 WHERE id=$10",
		p.Login, p.Site, p.Uppercase, p.Lowercase, p.Symbols, p.Numbers, p.Counter, p.Version, p.Length, p.ID,
	)
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}

// CreatePassword .
func CreatePassword(userID, login, site string, up, low, sym, num bool, c, v, l int) (*PasswordModel, error) {
	p := &PasswordModel{
		ID:        uuid.NewV4().String(),
		UserID:    userID,
		Login:     login,
		Site:      site,
		Uppercase: up,
		Lowercase: low,
		Symbols:   sym,
		Numbers:   num,
		Counter:   c,
		Version:   v,
		Length:    l,
	}
	_, err := conn.Exec(
		"INSERT INTO passwords (id, user_id, login, site, uppercase, lowercase, symbols, numbers, counter, version, length) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)",
		p.ID, p.UserID, p.Login, p.Site, p.Uppercase, p.Lowercase, p.Symbols, p.Numbers, p.Counter, p.Version, p.Length,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// GetPasswordByID .
func GetPasswordByID(id string) (*PasswordModel, error) {
	rows, err := conn.Query("SELECT id, user_id, login, site, uppercase, lowercase, symbols, numbers, counter, version, length FROM passwords WHERE id = $1 LIMIT 1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, ErrPasswordNotFound
	}
	p := &PasswordModel{}
	err = rows.Scan(&p.ID, &p.UserID, &p.Login, &p.Site, &p.Uppercase, &p.Lowercase, &p.Symbols, &p.Numbers, &p.Counter, &p.Version, &p.Length)
	if err != nil {
		return nil, err
	}
	return p, nil
}

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
