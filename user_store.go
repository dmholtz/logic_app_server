package server

import (
	b64 "encoding/base64"

	"crypto/rand"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserStore interface {
	Signup(c Credentials) error
	Login(c Credentials) (string, error)
	Logout(token string) error
}

type MyUserStore struct {
	db *sql.DB
}

func (us *MyUserStore) Signup(credential Credentials) error {
	if us.db.QueryRow("SELECT * FROM users WHERE username = ?", credential.Username).Scan() != sql.ErrNoRows {
		return errors.New("Username already exists")
	}

	pwdBytes := []byte(credential.Password)
	// generate random salt
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return err
	}
	// append salt to password
	pwdBytes = append(pwdBytes, salt...)
	// hash password
	hashedPwd, err := bcrypt.GenerateFromPassword(pwdBytes, bcrypt.DefaultCost)

	// store username, hashed password, and salt in database
	_, err = us.db.Exec("INSERT INTO users (username, hashed_password, salt) VALUES (?, ?, ?)", credential.Username, hashedPwd, salt)
	if err != nil {
		return err
	}

	return nil
}

func (us *MyUserStore) Login(credentials Credentials) (string, error) {
	row := us.db.QueryRow("SELECT * FROM users WHERE username = ?", credentials.Username)

	var id int
	var username string
	var hashedPwd []byte
	var salt []byte

	if err := row.Scan(&id, &username, &hashedPwd, &salt); err == sql.ErrNoRows {
		// user does not exist
		return "", errors.New("Invalid username or password1")
	} else if err != nil {
		return "", err
	}

	// append salt to password
	pwdBytes := []byte(credentials.Password)
	pwdBytes = append(pwdBytes, salt...)
	// hash password
	err := bcrypt.CompareHashAndPassword(hashedPwd, pwdBytes)
	if err != nil {
		// password does not match
		return "", errors.New("Invalid username or password3")
	}

	// generate token
	tokenBytes := make([]byte, 32)
	_, err = rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	// convert token to base64 string
	token := b64.StdEncoding.EncodeToString([]byte(tokenBytes))

	// store token in database
	_, err = us.db.Exec("INSERT INTO sessions (user_id, token) VALUES (?, ?)", id, token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (us *MyUserStore) Logout(token string) error {
	// delete token from database
	_, err := us.db.Exec("DELETE FROM sessions WHERE token = ?", token)
	if err != nil {
		return err
	}
	return nil
}
