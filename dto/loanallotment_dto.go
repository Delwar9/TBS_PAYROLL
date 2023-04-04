package dto

type LoanallotmentRefno struct {
	Refno int `json:"refno"`
}

type Accid_Name_List_dto struct {
	Accid int    `json:"accid"`
	Name  string `json:"name"`
}

type GetEmployeeInfoWithDeptAndDesigInputDto struct {
	Empname string `json:"empname"`
}

type GetEmployeeInfoWithDeptAndDesigOutputDto struct {
	Empcode   int    `json:"empcode"`
	Empname   string `json:"empname"`
	Deptname  string `json:"deptname"`
	Designame string `json:"designame"`
}

type GetEmployeeInfoWithEmpName_CurrentDue_DueInstallment struct {
	Empcode int `json:"empcode"`
}

type GetEmployeeInfoWithEmpName_CurrentDue_DueInstallmentOutput struct {
	Empname            string  `json:"empname"`
	Due_principal      float64 `json:"due_principal"`
	Installment_amount float64 `json:"installment_amount"`
	Loan_type          string  `json:"loan_type"`
	Empcode            int     `json:"empcode"`
}

type EmpNoCurrentDueDto struct {
	Empname            string `json:"empname"`
	Due_principal      string `json:"due_principal"`
	Installment_amount string `json:"installment_amount"`
	Loan_type          string `json:"loan_type"`
	Empcode            int    `json:"empcode"`
}

type NoOfInstallmentOutputDto struct {
	No_of_installment int `json:"no_of_installment"`
}

type InstallmentAmountOutputDto struct {
	Installment_amount float64 `json:"installment_amount"`
}

type GetLimitofCPFLoanAmountInputDto struct {
	Empcode     int     `json:"empcode"`
	Loan_amount float64 `json:"loan_amount"`
}

type GetLimitofCPFLoanAmountOutput1Dto struct {
	CPF float64 `json:"cpf"`
}

type GetLimitofCPFLoanAmountOutputDto struct {
	Loan_amount float64 `json:"loan_amount"`
}

type GetSingleEmployeeOutputDto struct {
	Stuff_accid int `json:"stuff_accid"`
}

type LoanpayscheduleforCPFOutputDto struct {
	Refno                      int     `json:"refno"`
	Installment_id             int     `json:"installment_id"`
	Empcode                    int     `json:"empcode"`
	Empname                    string  `json:"empname"`
	Month                      int     `json:"month"`
	Year                       int     `json:"year"`
	Amount                     float64 `json:"amount"`
	Monthly_deduct_interest    float64 `json:"monthly_deduct_interest"`
	Monthly_deduct_principal   float64 `json:"monthly_deduct_principal"`
	Interest_rate              float64 `json:"interest_rate"`
	Installment_amount         float64 `json:"installment_amount"`
	Due_principal              float64 `json:"due_principal"`
	Total_interest             float64 `json:"total_interest"`
	Total_amount_with_interest float64 `json:"total_amount_with_interest"`
}

type LoanpayscheduleforAdvanceOutputDto struct {
	Refno                      int     `json:"refno"`
	Installment_id             int     `json:"installment_id"`
	Empcode                    int     `json:"empcode"`
	Empname                    string  `json:"empname"`
	Month                      int     `json:"month"`
	Year                       int     `json:"year"`
	Amount                     float64 `json:"amount"`
	Monthly_deduct_principal   float64 `json:"monthly_deduct_principal"`
	Installment_amount         float64 `json:"installment_amount"`
	Due_principal              float64 `json:"due_principal"`
}

type GetDataToEditLoanAllotmentOutputDto struct {
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
	Flag_value              int     `json:"flag_value"`
}