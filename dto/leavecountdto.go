package dto

type LeaveCount struct {
	// Leavetype model.Leavetype
	// Leave     model.Leave
	Lcode     int    `json:"lcode"`
	Lname     string `json:"lname"`
	Earndays  int    `json:"earndays"`
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
	Festival  int    `json:"festival"`
	Sick      int    `json:"sick"`
	// Leavetype  int    `json:"leavetype"`
	Leavewopay int `json:"leavewopay"`
	Extraordi  int `json:"extraordi"`
}
