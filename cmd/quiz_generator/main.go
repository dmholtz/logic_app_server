package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	s "github.com/dmholtz/logic_app_server"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, dbErr := sql.Open("sqlite3", "./db.sqlite3")
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	defer db.Close()

	us := s.NewMyUserStore(db)

	start := time.Now()
	quizzes := make([]s.Quiz, 0)
	for _, qt := range [3]string{"SAT", "TAUT", "EQUIV"} {
		for _, diff := range [3]string{"easy", "medium", "hard"} {
			for _, numVars := range [9]int{2, 3, 4, 5, 6, 7, 8, 9, 10} {
				for i := 0; i < numVars; i++ {
					qp := s.QuizProperties{Type: qt, Difficulty: diff, NumVars: numVars, TimeLimit: 60}
					quiz, _ := us.GenerateQuiz(qp, false)
					quizzes = append(quizzes, quiz)
				}
			}
		}
	}
	timeElapsed := time.Since(start)
	fmt.Printf("Time elapsed for generating %d quizzes: %s", len(quizzes), timeElapsed)

	// marshal quizzes
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.Encode(quizzes)

	// write to file
	_ = os.WriteFile("quizzes.json", buf.Bytes(), 0644)

}
