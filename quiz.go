package server

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

type PossibleAnswer struct {
	Answer   string `json:"answer"`
	Solution bool   `json:"solution"`
}

type Quiz struct {
	QuizId          int              `json:"quiz_id"`
	Type            string           `json:"type"`
	TimeLimit       float64          `json:"time_limit"`
	Question        string           `json:"question"`
	PossibleAnswers []PossibleAnswer `json:"possible_answers"`
}

type QuizProperties struct {
	Type       string
	Difficulty string
	NumVars    int
	TimeLimit  int
}

func QuizPropertiesFromUrlQuery(queryParams url.Values) (QuizProperties, error) {
	qp := QuizProperties{}
	qp.Type = strings.ToUpper(queryParams.Get("type"))
	if qp.Type != "SAT" && qp.Type != "TAUT" && qp.Type != "EQUIV" {
		return qp, errors.New("Invalid quiz type.")
	}

	qp.Difficulty = strings.ToLower(queryParams.Get("difficulty"))
	if qp.Difficulty != "easy" && qp.Difficulty != "medium" && qp.Difficulty != "hard" {
		return qp, errors.New("Invalid quiz difficulty.")
	}

	// parse numVars query parameter to int
	numVars, err := strconv.ParseInt(queryParams.Get("numVars"), 10, 64)
	if err != nil || numVars < 1 || numVars > 10 {
		return qp, errors.New("Invalid number of variables.")
	}
	qp.NumVars = int(numVars)

	// parse timeLimit query parameter to float64
	timeLimit, err := strconv.ParseInt(queryParams.Get("timeLimit"), 10, 64)
	if err != nil || timeLimit < 0 {
		return qp, errors.New("Invalid time limit.")
	}
	qp.TimeLimit = int(timeLimit)
	return qp, nil
}