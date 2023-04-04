package model

type Leavetype struct {
	Lcode    int    `json:"lcode"`
	Lname    string `json:"lname"`
	Earndays int    `json:"earndays"`
}

type LeaveCount2 struct {
	Leavetype
	Leave
}
