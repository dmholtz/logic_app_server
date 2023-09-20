package server

import (
	"math/rand"

	l "github.com/dmholtz/logo"
	b "github.com/dmholtz/logo/builder"
	s "github.com/dmholtz/logo/scrambler"
)

// Calculate the time limit for a competition mode quiz based on its properties.
func TimeFromQuizProperties(qp QuizProperties) int {
	if qp.Type == "EQUIV" {
		switch qp.Difficulty {
		case "easy":
			return 60
		case "medium":
			return 45
		case "hard":
			return 30
		default:
			return 60
		}
	} else {
		switch qp.Difficulty {
		case "easy":
			return 30
		case "medium":
			return 20
		case "hard":
			return 15
		default:
			return 30
		}
	}
}

func RandomQuizProperties() QuizProperties {
	var quizType string = "SAT"
	switch rand.Intn(3) {
	case 0:
		quizType = "SAT"
	case 1:
		quizType = "TAUT"
	case 2:
		quizType = "EQUIV"
	}

	var difficulty string = "easy"
	switch rand.Intn(3) {
	case 0:
		difficulty = "easy"
	case 1:
		difficulty = "medium"
	case 2:
		difficulty = "hard"
	}

	numVars := rand.Intn(4) + 2
	timeLimit := rand.Intn(50) + 10

	return QuizProperties{
		Type:       quizType,
		Difficulty: difficulty,
		NumVars:    numVars,
		TimeLimit:  timeLimit}
}

// Generate a quiz based on its properties.
func GenerateTask(qp QuizProperties) Quiz {
	switch qp.Type {
	case "SAT":
		return generateSat(qp)
	case "TAUT":
		return generateTaut(qp)
	case "EQUIV":
		return generateEquiv(qp)
	default:
		panic("Invalid quiz type.")
	}
}

func generateSat(qp QuizProperties) Quiz {
	dnfBuilder := b.NewDnfBuilder(qp.NumVars, 4, 3)
	var formula l.LogicNode = l.Bottom()
	var solution []bool
	// enforce the correct scope
	for len(formula.Scope()) != qp.NumVars {
		if rand.Intn(2) == 0 {
			// not a tautology
			unsat := dnfBuilder.BuildUnsat()
			formula = &unsat
			solution = []bool{false, true}
		} else {
			// satisfiable
			sat := dnfBuilder.BuildSat()
			formula = &sat
			solution = []bool{true, false}
		}
	}
	if qp.Difficulty == "medium" || qp.Difficulty == "hard" {
		// scramble formula
		formula = s.DeMorganIteration(formula)
	}
	if qp.Difficulty == "hard" {
		// try to introduce -> and <-> operators
		formula = s.SubstituteArrows(formula)
	}
	formula = s.Simplify(formula)

	return Quiz{Type: "SAT", Question: formula.String(), Answers: []string{"yes", "no"}, Solutions: solution}
}

func generateTaut(qp QuizProperties) Quiz {
	dnfBuilder := b.NewDnfBuilder(qp.NumVars, 4, 3)
	var formula l.LogicNode = l.Bottom()
	var solution []bool
	// enforce the correct scope
	for len(formula.Scope()) != qp.NumVars {
		if rand.Intn(2) == 0 {
			// not a tautology
			sat := dnfBuilder.BuildSat()
			formula = l.Not(&sat)
			solution = []bool{false, true}
		} else {
			// satisfiable
			unsat := dnfBuilder.BuildUnsat()
			formula = l.Not(&unsat)
			solution = []bool{true, false}
		}
	}
	if qp.Difficulty == "medium" || qp.Difficulty == "hard" {
		// scramble formula
		formula = s.DeMorganIteration(formula)
	}
	if qp.Difficulty == "hard" {
		// try to introduce -> and <-> operators
		formula = s.SubstituteArrows(formula)
	}
	formula = s.Simplify(formula)

	return Quiz{Type: "TAUT", Question: formula.String(), Answers: []string{"yes", "no"}, Solutions: solution}
}

func generateEquiv(qp QuizProperties) Quiz {
	return Quiz{}
}
