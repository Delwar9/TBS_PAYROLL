package model

type Monthlydeduct struct {
	Id                       int     `json:"id"`
	Empcode                  int     `json:"empcode"`
	Month                    int     `json:"month"`
	Year                     int     `json:"year"`
	Incometax                float64 `json:"incometax"`
	Cpfloan                  float64 `json:"cpfloan"`
	Otheradv                 float64 `json:"otheradv"`
	Salaryadv                float64 `json:"salaryadv"`
	Loan_refno               int     `json:"loan_refno"`
	Loan_type                int     `json:"loan_type"`
	Investment               float64 `json:"investment"`
	Stuff_accid              int     `json:"stuff_accid"`
	Stuff_bankid             int     `json:"stuff_bankid"`
	Monthly_deduct_interest  float64 `json:"monthly_deduct_interest"`
	Monthly_deduct_principal float64 `json:"monthly_deduct_principal"`
	Due_principal            float64 `json:"due_principal"`
}
