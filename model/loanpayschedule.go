package model

type Loanpayschedule struct {
	Sl_id                    int     `json:"sl_id"`
	Refno                    int     `json:"refno"`
	Installment_id           int     `json:"installment_id"`
	Empcode                  int     `json:"empcode"`
	Empname                  string  `json:"empname"`
	Month                    int     `json:"month"`
	Year                     int     `json:"year"`
	Amount                   float64 `json:"amount"`
	Monthly_deduct_interest  float64 `json:"monthly_deduct_interest"`
	Monthly_deduct_principal float64 `json:"monthly_deduct_principal"`
	Interest_rate            float64 `json:"interest_rate"`
	Installment_amount       float64 `json:"installment_amount"`
	Due_principal            float64 `json:"due_principal"`
	Pause_flag               int     `json:"pause_flag"`
}
