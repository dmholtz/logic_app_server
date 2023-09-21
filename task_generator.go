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
			// not satisfiable
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
			// tautology
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
	equivBuilder := b.NewEquivalentFormulaBuilder(qp.NumVars, qp.NumVars+3)
	// enforce the correct scope
	for len(equivBuilder.BaseFormula.Scope())+1 != qp.NumVars {
		equivBuilder = b.NewEquivalentFormulaBuilder(qp.NumVars, 6)
	}

	answerFormulas := make([]l.LogicNode, 0)
	solutionsRaw := make([]bool, 0)

	numEquivalents := rand.Intn(4) + 1
	for i := 0; i < numEquivalents; i++ {
		answerFormulas = append(answerFormulas, equivBuilder.Equivalent())
		solutionsRaw = append(solutionsRaw, true)
	}
	for i := numEquivalents; i < 4; i++ {
		answerFormulas = append(answerFormulas, equivBuilder.NotEquivalent())
		solutionsRaw = append(solutionsRaw, false)
	}

	// remove arrows and apply commutativity to all formulas
	question := equivBuilder.Question()
	question = s.Traverse(question, s.RemoveIff)
	question = s.Traverse(question, s.RemoveImplies)
	for i := 0; i < 4; i++ {
		answerFormulas[i] = s.Traverse(answerFormulas[i], s.RemoveIff)
		answerFormulas[i] = s.Traverse(answerFormulas[i], s.RemoveImplies)
		answerFormulas[i] = s.TraverseProbabilistic(answerFormulas[i], s.Commute, 0.5)
	}

	if qp.Difficulty == "medium" || qp.Difficulty == "hard" {
		// scramble formula with DeMorgan
		for i := 0; i < 4; i++ {
			answerFormulas[i] = s.DeMorganIteration(answerFormulas[i])
		}
	}

	if qp.Difficulty == "hard" {
		// try to introduce -> and <-> operators
		for i := 0; i < 4; i++ {
			answerFormulas[i] = s.SubstituteArrows(answerFormulas[i])
		}
	}

	// simplify all formulas
	for i := 0; i < 4; i++ {
		answerFormulas[i] = s.Simplify(answerFormulas[i])
	}

	// permute answers and their corresponding solutions
	perm := rand.Perm(4)
	answers := make([]string, 4)
	solutions := make([]bool, 4)
	for i := 0; i < 4; i++ {
		answers[i] = answerFormulas[perm[i]].String()
		solutions[i] = solutionsRaw[perm[i]]
	}

	return Quiz{Type: "EQUIV", Question: question.String(), Answers: answers, Solutions: solutions}
}
