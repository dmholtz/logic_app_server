package server

type Achievement struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       string `json:"level"`
	Achieved    bool   `json:"achieved"`
}
