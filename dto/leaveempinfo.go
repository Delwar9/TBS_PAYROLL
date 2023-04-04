package dto

type LeaveEmpInfo struct {
	Empcode     int    `json:"empcode"`
	Empname     string `json:"empname"`
	Joindate    string `json:"joindate"`
	Confirmdate string `json:"confirmdate"`
	Deptname    string `json:"deptname"`
	Designame   string `json:"designame"`
}
