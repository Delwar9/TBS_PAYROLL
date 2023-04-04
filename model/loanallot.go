package model

type Loanallot struct {
	Refno                   int     `json:"refno"`
	Empcode                 int     `json:"empcode"`
	Allotment_date          string  `json:"allotment_date"`
	Loan_amount             float64 `json:"loan_amount"`
	Installment_amount      float64 `json:"installment_amount"`
	No_of_installment       int     `json:"no_of_installment"`
	Interest_rate           float64 `json:"interest_rate"`
	Effective_month         int     `json:"effective_month"`
	Effective_year          int     `json:"effective_year"`
	Empname                 string  `json:"empname"`
	Branchcode              int     `json:"branchcode"`
	Remarks                 string  `json:"remarks"`
	Loan_type               int     `json:"loan_type"`
	Deptcode                int     `json:"deptcode"`
	Desigcode               int     `json:"desigcode"`
	Due_month               int     `json:"due_month"`
	Due_year                int     `json:"due_year"`
	Total_interest          float64 `json:"total_interest"`
	Monthly_deduct_interest float64 `json:"monthly_deduct_interest"`
	Monthly_deduct_pricipal float64 `json:"monthly_deduct_pricipal"`
	Due_principal           float64 `json:"due_principal"`
	Entry_user              string  `json:"entry_user"`
	Compid                  int     `json:"compid"`
	Compyearid              int     `json:"compyearid"`
	Cash_bankid             int     `json:"cash_bankid"`
	Cash_bank_name          string  `json:"cash_bank_name"`
	Stuff_accid             int     `json:"stuff_accid"`
}
