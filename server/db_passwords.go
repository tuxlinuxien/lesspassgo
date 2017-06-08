package server

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

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

// GetPasswordsByUserID .
func GetPasswordsByUserID(userID string) []PasswordModel {
	rows, err := conn.Query("SELECT id, user_id, login, site, uppercase, lowercase, symbols, numbers, counter, version, length FROM passwords WHERE user_id = $1", userID)
	if err != nil {
		return []PasswordModel{}
	}
	defer rows.Close()
	var out = []PasswordModel{}
	for rows.Next() {
		p := PasswordModel{}
		err = rows.Scan(&p.ID, &p.UserID, &p.Login, &p.Site, &p.Uppercase, &p.Lowercase, &p.Symbols, &p.Numbers, &p.Counter, &p.Version, &p.Length)
		if err != nil {
			continue
		}
		out = append(out, p)
	}

	return out
}

// DeletePasswordByIDAndUserID .
func DeletePasswordByIDAndUserID(id, userID string) error {
	_, err := conn.Exec("DELETE FROM passwords WHERE id = $1 AND user_id = $2", id, userID)
	return err
}
