package model

type Weekend struct {
	Month    int    `json:"month"`
	Noofweek int    `json:"noofweek"`
	Date     string `json:"date"`
	Year     int    `json:"year"`
}
