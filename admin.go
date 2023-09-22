package server

type TemplateData struct {
	StatItems []StatItem
}

type StatItem struct {
	QuizType          string
	NumQuizzes        int
	NumParticipations int
	AvgTime           float64
}
