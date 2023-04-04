package model

type Staffinformation struct {
	Id              int    `json:"id"`
	Empcode         int    `json:"empcode"`
	Senioriy_serial int    `json:"senioriy_serial"`
	Empname         string `json:"empname"`
	Grade           string `json:"grade"`
	Accno           string `json:"accno"`
	Tin             string `json:"tin"`
	Joindate        string `json:"joindate"`
	Probationperiod int    `json:"probationperiod"`
	Confirmdate     string `json:"confirmdate"`
	Eduquali        string `json:"eduquali"`
	Mailadd         string `json:"mailadd"`
	Permanentadd    string `json:"permanentadd"`
	Maritalstatus   string `json:"maritalsatus"`
	Gender          string `json:"gender"`
	Religion        string `json:"religion"`
	Bloodgroup      string `json:"bloodgroup"`
	Active_salary   int    `json:"active_salary"`
	Entry_user      string `json:"entry_user"`
}

type Staffinformation_archive struct {
	Staffinformation

	Changedate   string `json:"changedate"`
	Changeuserid string `json:"changeuserid"`
	Flag_ed_del  string `json:"flag_ed_del"`
	Trackid      int    `json:"trackid"`
}
