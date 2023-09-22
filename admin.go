package server

type TemplateData struct {
	StatItems []StatItem
}

type StatItem struct {
	QuizType          string
	NumParticipations int
	AvgTime           float64
}
