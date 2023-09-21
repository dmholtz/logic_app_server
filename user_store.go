package server

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"log"
	"math"

	"database/sql"
	"errors"
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

type UserStore interface {
	Signup(c Credentials) error
	Login(c Credentials) (string, error)
	Logout(token string) error
	UserIdFromToken(token string) (int, error)
	ResetUser(userId int) error

	Achievements(user_id int) ([]Achievement, error)
	Leaderboard() ([]LeaderbordEntry, error)

	GetCompetitionQuiz(user_id int) (Quiz, error)
	FindQuiz(user_id int, qc QuizProperties) (Quiz, error)
	SolveQuiz(user_id int, ss SolveSubmission) (SolveSubmissionResponse, error)
}

type MyUserStore struct {
	DB *sql.DB

	// cached query strings
	queryAchievements    string
	queryFindQuiz        string
	queryCompetitionQuiz string
	queryLeaderboard     string

	// cached query results
	numAchievements int
}

func NewMyUserStore(db *sql.DB) *MyUserStore {
	// cache query strings
	queryAchievements := ReadQueryFile("db/queries/achievements.sql")
	queryCompetitionQuiz := ReadQueryFile("db/queries/competition_quiz.sql")
	queryFindQuiz := ReadQueryFile("db/queries/find_quiz.sql")
	queryLeaderboard := ReadQueryFile("db/queries/leaderboard.sql")

	// cache number of achievements
	var numAchievements int
	q_err := db.QueryRow("SELECT COUNT(*) FROM achievement").Scan(&numAchievements)
	if q_err != nil {
		log.Fatal(q_err)
	}

	return &MyUserStore{
		DB:                   db,
		queryAchievements:    queryAchievements,
		queryCompetitionQuiz: queryCompetitionQuiz,
		queryFindQuiz:        queryFindQuiz,
		queryLeaderboard:     queryLeaderboard,
		numAchievements:      numAchievements,
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

func (us *MyUserStore) ResetUser(userId int) error {
	// delete all quiz_participations of user
	_, err := us.DB.Exec("DELETE FROM quiz_participation WHERE user_id = ?", userId)
	if err != nil {
		return err
	}

	// delete all achieved achievements of user
	_, err = us.DB.Exec("DELETE FROM achieved WHERE user_id = ?", userId)
	if err != nil {
		return err
	}

	return nil
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

func (us *MyUserStore) GetCompetitionQuiz(user_id int) (Quiz, error) {
	var quiz_id int
	var quiz_type string
	var time_limit float64
	var question string
	var answer_str string
	var solution_str string
	row := us.DB.QueryRow(us.queryCompetitionQuiz, user_id)
	err := row.Scan(&quiz_id, &quiz_type, &time_limit, &question, &answer_str, &solution_str)
	if err == sql.ErrNoRows {
		// generate a quiz
		randomQuizProperties := RandomQuizProperties()
		quiz, err := us.GenerateQuiz(randomQuizProperties, true)
		return quiz, err
	} else if err != nil {
		return Quiz{}, err
	}

	// extract list of answers and solutions from JSON representation in answer_str and solution_str
	var answers []string
	if err := json.Unmarshal([]byte(answer_str), &answers); err != nil {
		return Quiz{}, err
	}
	var solutions []bool
	if err := json.Unmarshal([]byte(solution_str), &solutions); err != nil {
		return Quiz{}, err
	}
	return Quiz{QuizId: quiz_id, Type: quiz_type, TimeLimit: time_limit, Question: question, Answers: answers, Solutions: solutions}, nil
}

func (us *MyUserStore) FindQuiz(user_id int, qc QuizProperties) (Quiz, error) {
	var quiz_id int
	var quiz_type string
	var question string
	var answer_str string
	var solution_str string
	row := us.DB.QueryRow(us.queryFindQuiz, qc.Type, qc.Difficulty, qc.NumVars, user_id)
	err := row.Scan(&quiz_id, &quiz_type, &question, &answer_str, &solution_str)
	if err == sql.ErrNoRows {
		// generate a quiz
		quiz, err := us.GenerateQuiz(qc, false)
		return quiz, err
	} else if err != nil {
		return Quiz{}, err
	}

	// extract list of answers and solutions from JSON representation in answer_str and solution_str
	var answers []string
	if err := json.Unmarshal([]byte(answer_str), &answers); err != nil {
		return Quiz{}, err
	}
	var solutions []bool
	if err := json.Unmarshal([]byte(solution_str), &solutions); err != nil {
		return Quiz{}, err
	}
	return Quiz{QuizId: quiz_id, Type: quiz_type, TimeLimit: float64(qc.TimeLimit), Question: question, Answers: answers, Solutions: solutions}, nil
}

func (us *MyUserStore) SolveQuiz(userId int, ss SolveSubmission) (SolveSubmissionResponse, error) {
	points, err := us.CalculatePoints(ss)
	if err != nil {
		return SolveSubmissionResponse{}, err
	}
	_, err = us.DB.Exec("INSERT INTO quiz_participation (quiz_id, user_id, correct, time, points) VALUES (?, ?, ?, ?, ?)", ss.QuizId, userId, ss.IsCorrect, ss.Time, points)
	if err != nil {
		return SolveSubmissionResponse{}, err
	}
	err = us.UpdateAchievements(userId)
	if err != nil {
		return SolveSubmissionResponse{}, err
	}
	return SolveSubmissionResponse{Points: points}, nil
}

func (us *MyUserStore) CalculatePoints(ss SolveSubmission) (int, error) {
	if !ss.IsCorrect {
		return 0, nil
	}
	var difficulty string
	var numVars int
	var timeLimit int
	var isCompetitionMode bool
	row := us.DB.QueryRow("SELECT difficulty, num_vars, time_limit, is_competition_mode FROM quiz WHERE id = ?", ss.QuizId)
	if err := row.Err(); err == sql.ErrNoRows {
		return 0, errors.New("Invalid quiz id.")
	} else if err != nil {
		return 0, err
	}
	row.Scan(&difficulty, &numVars, &timeLimit, &isCompetitionMode)
	if !isCompetitionMode {
		return 0, nil
	}
	var points int
	if difficulty == "easy" {
		points = 1
	} else if difficulty == "medium" {
		points = 2
	} else if difficulty == "hard" {
		points = 3
	}
	points *= numVars
	points *= int(math.Max((float64(timeLimit)-ss.Time)/5, 1))

	return points, nil
}

func (us *MyUserStore) UpdateAchievements(userId int) error {
	// collect all achievements that are not yet achieved
	rows, err := us.DB.Query("SELECT id, sql FROM achievement WHERE id not in (SELECT achievement_id FROM achieved WHERE user_id = ?)", userId)
	if err != nil {
		return err
	}

	achievementIds := make([]int, 0)
	achievementQueries := make([]string, 0)
	for rows.Next() {
		var achievement_id int
		var achievementQuery string
		if err := rows.Scan(&achievement_id, &achievementQuery); err != nil {
			return err
		}
		achievementIds = append(achievementIds, achievement_id)
		achievementQueries = append(achievementQueries, achievementQuery)
	}
	rows.Close()

	for i := 0; i < len(achievementIds); i++ {
		// check if achievement is achieved
		row := us.DB.QueryRow(achievementQueries[i], userId)
		var count int
		if err := row.Scan(&count); err != nil {
			// achievement not achieved
			continue
		}
		// achievement achieved
		_, err := us.DB.Exec("INSERT INTO achieved (user_id, achievement_id) VALUES (?, ?)", userId, achievementIds[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (us *MyUserStore) GenerateQuiz(qc QuizProperties, isCompetition bool) (Quiz, error) {
	quiz := GenerateTask(qc)

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.Encode(quiz.Answers)
	answerStr := string(buf.Bytes())

	buf.Reset()
	enc.Encode(quiz.Solutions)
	solutionStr := string(buf.Bytes())

	var res sql.Result
	var err error
	if isCompetition {
		// competition mode
		timeLimit := TimeFromQuizProperties(qc)
		quiz.TimeLimit = float64(timeLimit)
		res, err = us.DB.Exec("INSERT INTO quiz (type, difficulty, num_vars, time_limit, is_competition_mode, question, answers, solutions) VALUES (?, ?, ?, ?, 1, ?, ?, ?)",
			quiz.Type, qc.Difficulty, qc.NumVars, timeLimit, quiz.Question, answerStr, solutionStr)
	} else {
		// practice mode
		quiz.TimeLimit = float64(qc.TimeLimit)
		res, err = us.DB.Exec("INSERT INTO quiz (type, difficulty, num_vars, time_limit, is_competition_mode, question, answers, solutions) VALUES (?, ?, ?, NULL, 0, ?, ?, ?)",
			quiz.Type, qc.Difficulty, qc.NumVars, quiz.Question, answerStr, solutionStr)
	}

	if err != nil {
		return Quiz{}, err
	}

	quizId, err := res.LastInsertId()
	if err != nil {
		return Quiz{}, err
	}
	quiz.QuizId = int(quizId)

	return quiz, nil
}
