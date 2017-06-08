package server

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
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
