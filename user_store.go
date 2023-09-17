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

	Leaderboard() ([]LeaderbordEntry, error)
}

type MyUserStore struct {
	DB *sql.DB
}

func (us *MyUserStore) Signup(credential Credentials) error {
	if us.DB.QueryRow("SELECT * FROM users WHERE username = ?", credential.Username).Scan() != sql.ErrNoRows {
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
	_, err = us.DB.Exec("INSERT INTO users (username, hashed_password, salt) VALUES (?, ?, ?)", credential.Username, hashedPwd, salt)
	if err != nil {
		return err
	}

	return nil
}

func (us *MyUserStore) Login(credentials Credentials) (string, error) {
	row := us.DB.QueryRow("SELECT * FROM users WHERE username = ?", credentials.Username)

	var id int
	var username string
	var hashedPwd []byte
	var salt []byte

	if err := row.Scan(&id, &username, &hashedPwd, &salt); err == sql.ErrNoRows {
		// user does not exist
		return "", errors.New("Invalid username or password")
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
		return "", errors.New("Invalid username or password")
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
	_, err = us.DB.Exec("INSERT INTO sessions (user_id, token) VALUES (?, ?)", id, token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (us *MyUserStore) Logout(token string) error {
	// check if token belongs to a session
	row := us.DB.QueryRow("SELECT * FROM sessions WHERE token = ?", token)
	if err := row.Err(); err == sql.ErrNoRows {
		return errors.New("Invalid token")
	} else if err != nil {
		return err
	}

	// delete token from database
	_, err := us.DB.Exec("DELETE FROM sessions WHERE token = ?", token)
	if err != nil {
		return err
	}
	return nil
}

func (us *MyUserStore) Leaderboard() ([]LeaderbordEntry, error) {
	leaderboard := make([]LeaderbordEntry, 0)

	rows, err := us.DB.Query("SELECT username, SUM(points) as score FROM users, quiz_participation qp WHERE users.id = qp.user_id GROUP BY username ORDER BY score DESC")
	if err != nil {
		return leaderboard, err
	}
	defer rows.Close()

	for rows.Next() {
		var username string
		var score int
		if err := rows.Scan(&username, &score); err != nil {
			return leaderboard, err
		}
		leaderboardEntry := LeaderbordEntry{Username: username, Experience: 0.5, Points: score}
		leaderboard = append(leaderboard, leaderboardEntry)
	}

	return leaderboard, nil
}
