package server

import (
	b64 "encoding/base64"
	"log"

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
	UserIdFromToken(token string) (int, error)

	Achievements(user_id int) ([]Achievement, error)
	Leaderboard() ([]LeaderbordEntry, error)
}

type MyUserStore struct {
	DB *sql.DB

	// cached query strings
	queryAchievements string
	queryLeaderboard  string

	// cached query results
	numAchievements int
}

func NewMyUserStore(db *sql.DB) *MyUserStore {
	// cache query strings
	queryAchievements := ReadQueryFile("db/queries/achievements.sql")
	queryLeaderboard := ReadQueryFile("db/queries/leaderboard.sql")

	// cache number of achievements
	var numAchievements int
	q_err := db.QueryRow("SELECT COUNT(*) FROM achievement").Scan(&numAchievements)
	if q_err != nil {
		log.Fatal(q_err)
	}

	return &MyUserStore{
		DB:                db,
		queryAchievements: queryAchievements,
		queryLeaderboard:  queryLeaderboard,
		numAchievements:   numAchievements,
	}
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

func (us *MyUserStore) UserIdFromToken(token string) (int, error) {
	// check if token belongs to a session
	row := us.DB.QueryRow("SELECT user_id FROM sessions WHERE token = ?", token)
	if err := row.Err(); err == sql.ErrNoRows {
		return -1, errors.New("Invalid token")
	} else if err != nil {
		return -1, err
	}

	// get user_id from session
	var user_id int
	err := row.Scan(&user_id)
	if err != nil {
		return -1, err
	}

	return user_id, nil
}

func (us *MyUserStore) Achievements(user_id int) ([]Achievement, error) {
	achievements := make([]Achievement, 0)

	rows, err := us.DB.Query(us.queryAchievements, user_id)
	if err != nil {
		return achievements, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var description string
		var level string
		var isAchieved bool
		if err := rows.Scan(&name, &description, &level, &isAchieved); err != nil {
			return achievements, err
		}
		achievement := Achievement{Name: name, Description: description, Level: level, Achieved: isAchieved}
		achievements = append(achievements, achievement)
	}

	return achievements, nil
}

func (us *MyUserStore) Leaderboard() ([]LeaderbordEntry, error) {
	leaderboard := make([]LeaderbordEntry, 0)

	rows, err := us.DB.Query(us.queryLeaderboard)
	if err != nil {
		return leaderboard, err
	}
	defer rows.Close()

	for rows.Next() {
		var username string
		var xp float64
		var score int
		if err := rows.Scan(&username, &xp, &score); err != nil {
			return leaderboard, err
		}
		leaderboardEntry := LeaderbordEntry{Username: username, Experience: xp / float64(us.numAchievements), Points: score}
		leaderboard = append(leaderboard, leaderboardEntry)
	}

	return leaderboard, nil
}
