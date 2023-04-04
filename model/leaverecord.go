package model

type Leaverecord struct {
	Id                 int `json:"id"`
	Empcode            int `json:"empcode"`
	Earnleave          int `json:"earnleave"`
	Year               int `json:"year"`
	Earnconsume        int `json:"earnconsume"`
	Rest_of_leave      int `json:"rest_of_leave"`
	El_earned_nextyear int `json:"el_earned_nextyear"`
}
