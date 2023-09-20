package server

type SolveSubmission struct {
	QuizId    int     `json:"quizId"`
	IsCorrect bool    `json:"isCorrect"`
	Time      float64 `json:"time"`
}

type SolveSubmissionResponse struct {
	Points int `json:"points"`
}
