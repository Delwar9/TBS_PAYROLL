package model

type Leave struct {
	Refno     int    `json:"refno"`
	Empcode   int    `json:"empcode"`
	Leavedate string `json:"leavedate"`
	Leavefrom string `json:"leavefrom"`
	Leaveto   string `json:"leaveto"`
	Casual    int    `json:"casual"`
	Medical   int    `json:"medical"`
	Special   int    `json:"special"`
	Earn      int    `json:"earn"`
	Cause     string `json:"cause"`
	Study     int    `json:"study"`
	Maternity int    `json:"maternity"`
	// Leavetype  int    `json:"leavetype"`
	Leavewopay int `json:"leavewopay"`
	Extraordi  int `json:"extraordi"`
	Festival   int `json:"festival"`
	Sick       int `json:"sick"`
}
