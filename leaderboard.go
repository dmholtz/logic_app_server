package server

type LeaderbordEntry struct {
	Username   string  `json:"username"`
	Experience float64 `json:"experience"`
	Points     int     `json:"points"`
}
