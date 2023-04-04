package dto

type Staffinformations struct {
	Empcode         string `json:"empcode"`
	Empname         string `json:"empname"`
	Deptname        string `json:"deptname"`
	Designame       string `json:"designame"`
	Grade           string `json:"grade"`
	Senioriy_serial int    `json:"senioriy_serial"`
	Bloodgroup      string `json:"bloodgroup"`
}
