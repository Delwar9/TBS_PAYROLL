package repository

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tools/payroll/dto"
	"github.com/tools/payroll/model"
	"github.com/tools/payroll/util"
)

type LoanallotmentRepository struct{}

func (loanallotmentRepository *LoanallotmentRepository) MaxRefNo() dto.ResponseDto {
	var res dto.ResponseDto
	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	tx := db.Begin()
	tx.SavePoint("savepoint")
	// var output model.Loanallot
	var output dto.LoanallotmentRefno
	// result := tx.Raw("select max(refno)+1 as refno from payroll.loanallotbackup").First(&ouptut.Refno)
	result := tx.Raw("select max(refno)+1 as refno from payroll.loanallot").First(&output.Refno)
	if result.RowsAffected == 0 {
		tx.RollbackTo("savepoint")
		res.IsSuccess = false
		res.StatusCode = http.StatusInternalServerError
		res.Message = "Oops! Something went wrong. Please try again later."
		res.Payload = nil
		return res
	}

	if result.Error != nil {
		output.Refno = 1
	}

	tx.Commit()

	// var temp tempOutput

	res.IsSuccess = true
	res.StatusCode = http.StatusOK
	res.Message = "Max code for new Loanallotment"
	res.Payload = output

	return res
}

func (loanallotmentRepository *LoanallotmentRepository) GetEmployeeLoanInfo(input model.Loanallot) dto.ResponseDto {
	var res dto.ResponseDto

	if input.Empcode == 0 {
		res.Message = "Please Type Employee Code"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var outputFinal model.Loanallot

	var outputFinal1 []model.Loanallot

	var output dto.GetEmployeeInfoWithEmpName_CurrentDue_DueInstallmentOutput

	var Output3 dto.EmpNoCurrentDueDto

	result := db.Raw("select s.empname, l.due_principal, l.installment_amount, l.loan_type, s.empcode from payroll.loanallot l, payroll.staffinformation s where l.empcode = ? and s.empcode = l.empcode and l.refno = (select max(refno) from payroll.loanallot las where las.empcode = ?)", input.Empcode, input.Empcode).First(&outputFinal)
	if result.RowsAffected == 0 {
		var output1 model.Staffinformation
		result1 := db.Where(model.Staffinformation{Empcode: input.Empcode}).First(&output1)
		if result1.RowsAffected == 0 {
			res.Message = "Employee Not Found"
			res.IsSuccess = false
			res.Payload = nil
			res.StatusCode = http.StatusNotFound
			return res
		}
		var output2 dto.EmpNoCurrentDueDto
		output2.Empname = output1.Empname
		output2.Empcode = output1.Empcode
		output2.Due_principal = ""
		output2.Installment_amount = "no due"
		output2.Loan_type = "Not Assigned Yet"
		res.Message = output1.Empname + " have no loan."
		res.IsSuccess = true
		res.Payload = output2
		res.StatusCode = http.StatusOK
		return res
	}

	// repo1 := new(LoanallotmentRepository)
	// repo1Func := repo1.TerminateCPFloan(input)
	// if !repo1Func.IsSuccess {
	// 	output.Empname = outputFinal.Empname
	// 	output.Installment_amount = outputFinal.Installment_amount
	// 	output.Due_principal = outputFinal.Due_principal
	// 	output.Empcode = outputFinal.Empcode
	// 	output.Loan_type = "CPF Loan Allotment"
	// 	res.Message = repo1Func.Message
	// 	res.IsSuccess = true
	// 	res.Payload = output
	// 	res.StatusCode = http.StatusOK
	// 	return res
	// }

	// repo2 := new(LoanallotmentRepository)
	// repo2Func := repo2.TerminateAdvanceAgainstSalary(input)
	// if !repo2Func.IsSuccess {
	// 	output.Empname = outputFinal.Empname
	// 	output.Installment_amount = outputFinal.Installment_amount
	// 	output.Due_principal = outputFinal.Due_principal
	// 	output.Empcode = outputFinal.Empcode
	// 	output.Loan_type = "Advance against salary"
	// 	res.Message = repo2Func.Message
	// 	res.IsSuccess = true
	// 	res.Payload = output
	// 	res.StatusCode = http.StatusOK
	// 	return res
	// }

	output.Empname = outputFinal.Empname
	output.Installment_amount = outputFinal.Installment_amount
	output.Due_principal = outputFinal.Due_principal
	output.Empcode = outputFinal.Empcode
	if outputFinal.Loan_type == 0 {
		output.Loan_type = "Advance against salary"
	} else if outputFinal.Loan_type == 1 {
		output.Loan_type = "CPF Loan Allotment"
	}
	Repo2 := new(LoanallotmentRepository)
	Func1 := Repo2.TerminateCPFloan(input)
	Func2 := Repo2.TerminateAdvanceAgainstSalary(input)

	// result1 := db.Raw("SELECT * FROM payroll.loanallot l  WHERE l.empcode = ? AND loan_type IN (1, 0)", input.Empcode).Find(&outputFinal1)
	result1 := db.Raw(`SELECT * FROM payroll.loanallot l  WHERE l.empcode = ? AND loan_type IN (1, 0) and l.refno = (select max(refno) from payroll.loanallot las where las.empcode = ? and las.loan_type = 0)
	union 
	SELECT * FROM payroll.loanallot l  WHERE l.empcode = ? AND loan_type IN (1, 0) and l.refno = (select max(refno) from payroll.loanallot las where las.empcode = ? and las.loan_type = 1)`, input.Empcode, input.Empcode, input.Empcode, input.Empcode).Find(&outputFinal1)

	// if result1.Count() == 0 {
	if result1.RowsAffected == 2 {
		if outputFinal1[0].Loan_type == 0 && outputFinal1[1].Loan_type == 1 {
			// if outputFinal1[0].Loan_type == 0 {
			// 	repo2Func := Repo2.TerminateAdvanceAgainstSalary(input)
			// 	if !repo2Func.IsSuccess {
			// 		output.Empname = outputFinal.Empname
			// 		output.Installment_amount = outputFinal.Installment_amount
			// 		output.Due_principal = outputFinal.Due_principal
			// 		output.Empcode = outputFinal.Empcode
			// 		output.Loan_type = "CPF Loan Allotment"
			// 		res.Message = repo2Func.Message
			// 		res.IsSuccess = false
			// 		res.Payload = output
			// 		res.StatusCode = http.StatusOK
			// 		return res
			// 	}
			// } else if outputFinal1[1].Loan_type == 1 {
			// 	repo2Func := Repo2.TerminateCPFloan(input)
			// 	if !repo2Func.IsSuccess {
			// 		output.Empname = outputFinal.Empname
			// 		output.Installment_amount = outputFinal.Installment_amount
			// 		output.Due_principal = outputFinal.Due_principal
			// 		output.Empcode = outputFinal.Empcode
			// 		output.Loan_type = "Advance against salary"
			// 		res.Message = repo2Func.Message
			// 		res.IsSuccess = false
			// 		res.Payload = output
			// 		res.StatusCode = http.StatusOK
			// 		return res
			// 	}
			// }
			fmt.Println("Part 1")
			if !Func1.IsSuccess && !Func2.IsSuccess {
				fmt.Println("Both Loan Type Running")
				output.Empname = outputFinal.Empname
				output.Installment_amount = outputFinal1[0].Installment_amount + outputFinal1[1].Installment_amount
				output.Due_principal = outputFinal.Due_principal
				output.Empcode = outputFinal.Empcode
				output.Loan_type = "CPF Loan and Advance against Salary"
				res.Message = "Both Loan Type Assigned"
				res.IsSuccess = true
				res.Payload = output
				res.StatusCode = http.StatusOK
				return res

			} else if Func1.IsSuccess && !Func2.IsSuccess {
				fmt.Println("Advance Loan Type Running")
				output.Empname = outputFinal.Empname
				// output.Installment_amount = outputFinal1[0].Installment_amount + outputFinal1[1].Installment_amount
				output.Due_principal = outputFinal.Due_principal
				output.Empcode = outputFinal.Empcode
				output.Loan_type = "Advance against Salary"
				res.Message = "Advance Loan Type Running"
				res.IsSuccess = true
				res.Payload = output
				res.StatusCode = http.StatusOK
				return res

			} else if !Func1.IsSuccess && Func2.IsSuccess {
				fmt.Println("CPF Loan Type Running")
				output.Empname = outputFinal.Empname
				// output.Installment_amount = outputFinal1[0].Installment_amount + outputFinal1[1].Installment_amount
				output.Due_principal = outputFinal.Due_principal
				output.Empcode = outputFinal.Empcode
				output.Loan_type = "CPF Loan"
				res.Message = "CPF Loan Type Running"
				res.IsSuccess = true
				res.Payload = output
				res.StatusCode = http.StatusOK
				return res

			} else {
				fmt.Println("Both Loan Type previously running but now terminated")
				Output3.Empname = outputFinal.Empname
				Output3.Installment_amount = ""
				Output3.Due_principal = "no due"
				Output3.Empcode = outputFinal.Empcode
				Output3.Loan_type = "Not Assigned Yet"
				res.Message = "Can Take any loan"
				res.IsSuccess = true
				res.Payload = Output3
				res.StatusCode = http.StatusOK
				return res
			}
		} else if outputFinal1[0].Loan_type == 1 && outputFinal1[1].Loan_type == 0 {
			// if outputFinal1[1].Loan_type == 0 {
			// 	repo2Func := Repo2.TerminateAdvanceAgainstSalary(input)
			// 	if !repo2Func.IsSuccess {
			// 		output.Empname = outputFinal.Empname
			// 		output.Installment_amount = outputFinal.Installment_amount
			// 		output.Due_principal = outputFinal.Due_principal
			// 		output.Empcode = outputFinal.Empcode
			// 		output.Loan_type = "CPF Loan Allotment"
			// 		res.Message = repo2Func.Message
			// 		res.IsSuccess = false
			// 		res.Payload = output
			// 		res.StatusCode = http.StatusOK
			// 		return res
			// 	}
			// } else if outputFinal1[0].Loan_type == 1 {
			// 	repo2Func := Repo2.TerminateCPFloan(input)
			// 	if !repo2Func.IsSuccess {
			// 		output.Empname = outputFinal.Empname
			// 		output.Installment_amount = outputFinal.Installment_amount
			// 		output.Due_principal = outputFinal.Due_principal
			// 		output.Empcode = outputFinal.Empcode
			// 		output.Loan_type = "Advance against salary"
			// 		res.Message = repo2Func.Message
			// 		res.IsSuccess = false
			// 		res.Payload = output
			// 		res.StatusCode = http.StatusOK
			// 		return res
			// 	}
			// }
			fmt.Println("Part 2")
			if !Func1.IsSuccess && !Func2.IsSuccess {
				fmt.Println("Both Loan Type Running")
				output.Empname = outputFinal.Empname
				output.Installment_amount = outputFinal1[0].Installment_amount + outputFinal1[1].Installment_amount
				output.Due_principal = outputFinal.Due_principal
				output.Empcode = outputFinal.Empcode
				output.Loan_type = "CPF Loan and Advance against Salary"
				res.Message = "Both Loan Type Assigned"
				res.IsSuccess = true
				res.Payload = output
				res.StatusCode = http.StatusOK
				return res
			} else if Func1.IsSuccess && !Func2.IsSuccess {
				fmt.Println("Advance Loan Type Running")
				output.Empname = outputFinal.Empname
				// output.Installment_amount = outputFinal1[0].Installment_amount + outputFinal1[1].Installment_amount
				output.Due_principal = outputFinal.Due_principal
				output.Empcode = outputFinal.Empcode
				output.Loan_type = "Advance against Salary"
				res.Message = "Advance Loan Type Running"
				res.IsSuccess = true
				res.Payload = output
				res.StatusCode = http.StatusOK
				return res

			} else if !Func1.IsSuccess && Func2.IsSuccess {
				fmt.Println("CPF Loan Type Running")
				output.Empname = outputFinal.Empname
				// output.Installment_amount = outputFinal1[0].Installment_amount + outputFinal1[1].Installment_amount
				output.Due_principal = outputFinal.Due_principal
				output.Empcode = outputFinal.Empcode
				output.Loan_type = "CPF Loan"
				res.Message = "CPF Loan Type Running"
				res.IsSuccess = true
				res.Payload = output
				res.StatusCode = http.StatusOK
				return res

			} else {

				fmt.Println("Both Loan Type previously running but now terminated")
				Output3.Empname = outputFinal.Empname
				Output3.Installment_amount = ""
				Output3.Due_principal = "no due"
				Output3.Empcode = outputFinal.Empcode
				Output3.Loan_type = "Not Assigned Yet"
				res.Message = "Can Take any loan"
				res.IsSuccess = true
				res.Payload = Output3
				res.StatusCode = http.StatusOK
				return res
			}
		}
	}
	if result1.RowsAffected == 1 {
		if outputFinal1[0].Loan_type == 0 {
			fmt.Println("Advance Loan Running/Previous Running")
			// func1 := Repo2.TerminateAdvanceAgainstSalary(input)
			if !Func2.IsSuccess {
				output.Empname = outputFinal.Empname
				output.Installment_amount = outputFinal.Installment_amount
				output.Due_principal = outputFinal.Due_principal
				output.Empcode = outputFinal.Empcode
				output.Loan_type = "Advance against salary"
				res.Message = Func2.Message
				res.IsSuccess = true
				res.Payload = output
				res.StatusCode = http.StatusOK
				return res
			} else {
				Output3.Empname = outputFinal.Empname
				Output3.Installment_amount = ""
				Output3.Due_principal = "no due"
				Output3.Empcode = outputFinal.Empcode
				Output3.Loan_type = "Not Assigned Yet"
				res.Message = "Can Take any loan"
				res.IsSuccess = true
				res.Payload = Output3
				res.StatusCode = http.StatusOK
				return res
			}
		} else if outputFinal1[0].Loan_type == 1 {
			fmt.Println("CPF Loan Running/Previous Running")
			if !Func1.IsSuccess {
				output.Empname = outputFinal.Empname
				output.Installment_amount = outputFinal.Installment_amount
				output.Due_principal = outputFinal.Due_principal
				output.Empcode = outputFinal.Empcode
				output.Loan_type = "CPF Loan"
				res.Message = Func1.Message
				res.IsSuccess = true
				res.Payload = output
				res.StatusCode = http.StatusOK
				return res
			} else {
				Output3.Empname = outputFinal.Empname
				Output3.Installment_amount = ""
				Output3.Due_principal = "no due"
				Output3.Empcode = outputFinal.Empcode
				Output3.Loan_type = "Not Assigned Yet"
				res.Message = "Can Take any loan"
				res.IsSuccess = true
				res.Payload = Output3
				res.StatusCode = http.StatusOK
				return res
			}
		}
	}
	if result1.RowsAffected == 0 {
		Output3.Empname = outputFinal.Empname
		Output3.Installment_amount = ""
		Output3.Due_principal = "no due"
		Output3.Empcode = outputFinal.Empcode
		Output3.Loan_type = "Not Assigned Yet"
		res.Message = "Can Take any loan"
		res.IsSuccess = true
		res.Payload = Output3
		res.StatusCode = http.StatusOK
		return res
	}
	res.Message = "Successfully Get Employee Loan Info"
	res.IsSuccess = true
	res.Payload = output
	res.StatusCode = http.StatusOK

	return res
}

func (loanallotmentRepository *LoanallotmentRepository) GetAllEmployeeWithDesigAndDept(input dto.GetEmployeeInfoWithDeptAndDesigInputDto) dto.ResponseDto {
	var res dto.ResponseDto

	if input.Empname == "" {
		res.Message = "Please Type Employee Name"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	var output []dto.GetEmployeeInfoWithDeptAndDesigOutputDto
	result := db.Raw("select sa.empcode, s.empname, h2.deptname, d.designame from payroll.designation d, payroll.hrdept h2, payroll.staffinformation s , payroll.salarystructure sa where lower(s.empname) like ? and sa.empcode = s.empcode and h2.deptcode = sa.deptcode and d.desigcode = sa.desigcode order by s.empname", "%"+strings.ToLower(input.Empname)+"%").Find(&output)
	if result.RowsAffected == 0 {
		res.Message = "Not Found"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotFound
		return res
	}

	res.Message = "Successfully Get All Employee List"
	res.IsSuccess = true
	res.Payload = output
	res.Count = len(output)
	res.StatusCode = http.StatusOK

	return res
}

func (loanallotmentRepository *LoanallotmentRepository) BankLoadForLoan() dto.ResponseDto {
	var res dto.ResponseDto

	db := util.CreateConnectionToAccountsSchemaUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var output []dto.Accid_Name_List_dto
	result := db.Raw("select a.accid, a.name from accounts.acchead a where a.parent in(10101100,10101150) or a.accid= 10100101 and a.lr = 'L' order by name").Find(&output)
	if result.RowsAffected == 0 {
		res.Message = "No Data Found"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotFound
		return res
	}

	// type tempOutput struct {
	// 	Output      []dto.BankLoanListDto `json:"output"`
	// 	OutputCount int                   `json:"output_count"`
	// }

	// var tOutput tempOutput
	// tOutput.Output = output
	// tOutput.OutputCount = len(output)
	res.Message = "Successfully Get All Bank Loan List"
	res.IsSuccess = true
	res.Payload = output
	res.Count = len(output)
	res.StatusCode = http.StatusOK

	return res
}

func (loanallotmentRepository *LoanallotmentRepository) GetABankForLoan(input dto.Accid_Name_List_dto) dto.ResponseDto {
	var res dto.ResponseDto

	if input.Accid == 0 {
		res.Message = "Please Select Bank"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	db := util.CreateConnectionToAccountsSchemaUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var output dto.Accid_Name_List_dto
	result := db.Raw(`(select a.accid, a.name from accounts.acchead a where a.parent in(10101100,10101150) or a.accid= 10100101 and a.lr = 'L' order by name)
	INTERSECT
	(select a.accid, a.name from  accounts.acchead a  where a.accid = ?)`, input.Accid).First(&output)
	if result.RowsAffected == 0 {
		res.Message = "No Data Found"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotFound
		return res
	}

	res.Message = "Successfully Get Bank"
	res.IsSuccess = true
	res.Payload = output
	res.StatusCode = http.StatusOK

	return res
}

func (loanallotmentRepository *LoanallotmentRepository) GetLimitofCPFLoanAmount(input dto.GetLimitofCPFLoanAmountInputDto) dto.ResponseDto {
	var res dto.ResponseDto

	if input.Empcode == 0 {
		res.Message = "Please Select Employee"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Loan_amount == 0 {
		res.Message = "Please Enter Loan Amount"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	var output1 dto.GetLimitofCPFLoanAmountOutput1Dto

	result1 := db.Raw("select s.cpf from payroll.salarystructure s where s.empcode = ?", input.Empcode).First(&output1)
	if result1.RowsAffected == 0 {
		res.Message = "No Data Found"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotFound
		return res
	}

	validLoanAmount := output1.CPF * 0.8
	validLoanAmountstr := strconv.Itoa(int(validLoanAmount))

	if input.Loan_amount > validLoanAmount {
		res.Message = "he/she cant take CPF loan over : " + validLoanAmountstr + "BDT"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotAcceptable
		return res
	}

	var output dto.GetLimitofCPFLoanAmountOutputDto
	output.Loan_amount = input.Loan_amount
	res.Message = "Successfully Get Limit of CPF Loan Amount"
	res.IsSuccess = true
	res.Payload = output
	res.StatusCode = http.StatusOK

	return res
}

func (loanallotmentRepository *LoanallotmentRepository) ValidateForwardingDate(input model.Loanallot) dto.ResponseDto {
	var res dto.ResponseDto

	if input.Effective_year == 0 || input.Effective_month > 12 {
		res.Message = "Please Select Effective_month"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Effective_year == 0 {
		res.Message = "Please Select Effective_year"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	dt := time.Now()

	fulldate := dt.Format("2006-01-02")
	fmt.Println("fulldate: ", fulldate)

	y := fulldate[0:4]
	outputyear := y
	fmt.Println("After Slice", outputyear)
	outputyearint, err := strconv.Atoi(outputyear) // Convert String to Int
	if err != nil {                                // Error Handling
		panic(err)
	}

	if input.Effective_year < outputyearint {
		res.Message = "Effective date couldn't be previous year month againts current year month"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotAcceptable
		return res
	}

	m := fulldate[5:7]
	outputmonth := m
	fmt.Println("After Slice", outputmonth)
	outputmonthint, err := strconv.Atoi(outputmonth) // Convert String to Int
	if err != nil {                                  // Error Handling
		panic(err)
	}

	if input.Effective_year == outputyearint && input.Effective_month < outputmonthint {
		res.Message = "Effective date couldn't be previous year month againts current year month"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotAcceptable
		return res
	}

	res.Message = "Successfully Validate Forwarding Date"
	res.IsSuccess = true
	res.Payload = nil
	res.StatusCode = http.StatusOK
	return res
}

// TODO: CPF Loan
func (loanallotmentRepository *LoanallotmentRepository) TerminateCPFloan(input model.Loanallot) dto.ResponseDto {
	var res dto.ResponseDto

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	tx := db.Begin()
	tx.SavePoint("savepoint")
	var output1 model.Loanpayschedule
	// result1 := tx.Raw("select * from payroll.staffinformation s where s.empcode = ?", input.Empcode).First(&input)
	// if result1.RowsAffected == 0 {
	// 	res.Message = "There is no employee for this empcode"
	// 	res.IsSuccess = false
	// 	res.Payload = nil
	// 	res.StatusCode = http.StatusNotFound
	// 	tx.RollbackTo("savepoint")
	// 	return res
	// }

	result2 := tx.Raw("select lp.empcode , lp.month , lp.year from payroll.loanpayschedule lp where lp.refno = (select max(refno) from payroll.loanallot las where las.empcode = ? and las.loan_type = 1) order by lp.installment_id desc limit 1", input.Empcode).First(&output1)
	if result2.RowsAffected == 0 {
		res.Message = "No previous CPF loan found"
		res.IsSuccess = true
		res.Payload = nil
		res.StatusCode = http.StatusOK
		tx.RollbackTo("savepoint")
		return res
	}

	dt := time.Now()
	fulldate := dt.Format("2006-01-02")
	fmt.Println("fulldate: ", fulldate)

	y := fulldate[0:4]
	outputyear := y
	fmt.Println("After Slice Year", outputyear)
	outputyearint, err := strconv.Atoi(outputyear) // Convert String to Int
	if err != nil {                                // Error Handling
		panic(err)
	}
	fmt.Println("Year", output1.Year)
	if output1.Year > outputyearint {
		res.Message = "Already CPF loan running."
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotAcceptable
		tx.RollbackTo("savepoint")
		return res
	} else if output1.Year < outputyearint {
		res.Message = "There is no CPF loan running for this empcode."
		res.IsSuccess = true
		res.Payload = nil
		res.StatusCode = http.StatusOK
		return res
	}

	m := fulldate[5:7]
	outputmonth := m
	fmt.Println("After Slice month", outputmonth)
	outputmonthint, err := strconv.Atoi(outputmonth) // Convert String to Int
	if err != nil {                                  // Error Handling
		panic(err)
	}
	fmt.Println("Month", output1.Month)

	if output1.Year == outputyearint && output1.Month >= outputmonthint {
		res.Message = "Already CPF loan running."
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotAcceptable
		tx.RollbackTo("savepoint")
		return res
	} else if output1.Year == outputyearint && output1.Month < outputmonthint {
		res.Message = "There is no CPF loan running for this empcode."
		res.IsSuccess = true
		res.Payload = nil
		res.StatusCode = http.StatusOK
		return res
	}

	tx.Commit()

	return res
}

func (loanallotmentRepository *LoanallotmentRepository) GetInstallmentAmountPerMonthForCPF(input model.Loanallot) dto.ResponseDto {
	var res dto.ResponseDto

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var output dto.InstallmentAmountOutputDto

	input.Interest_rate = input.Interest_rate / 12 / 100

	if input.No_of_installment <= 0 {
		res.Message = "No of Installment can't be 0 or less than 0"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotAcceptable
		return res
	}

	if input.No_of_installment > 60 {
		res.Message = "No of Installment can't be greater than 60"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotAcceptable
		return res
	}

	Installment_amount := (input.Loan_amount * input.Interest_rate * math.Pow(1+input.Interest_rate, float64(input.No_of_installment))) / (math.Pow(1+input.Interest_rate, float64(input.No_of_installment)) - 1) // installment amount per month using EMI formula
	output.Installment_amount = math.Round(Installment_amount)
	//output.Installment_amount = math.Round(output.Installment_amount*100) / 100
	// Installment_amount_int_type, Installment_amount_float_type := math.Modf(output.Installment_amount) // separate integer and float part
	// fmt.Println("float type: ", Installment_amount_float_type)
	// fmt.Println("int type: ", Installment_amount_int_type)
	// if Installment_amount_float_type >= 0.5 { // if float part is greater than 0.5 then add 1 to integer part
	// Installment_amount_int_type = Installment_amount_int_type + 1
	// }
	// output.Installment_amount = float64(Installment_amount_int_type) // convert integer to float64

	res.Message = "Successfully Get Installment Amount Per Month"
	res.IsSuccess = true
	res.Payload = output
	res.StatusCode = http.StatusOK

	return res
}

func (loanallotmentRepository *LoanallotmentRepository) GetNoOfInstallmentForCPF(input model.Loanallot) dto.ResponseDto {
	var res dto.ResponseDto

	var output dto.NoOfInstallmentOutputDto
	input.Interest_rate = input.Interest_rate / 12 / 100

	float64NoOfInstallment := (math.Log(input.Installment_amount) - math.Log(input.Installment_amount-input.Loan_amount*input.Interest_rate)) / math.Log(1+input.Interest_rate)
	fmt.Println("No of Installment Float type: ", float64NoOfInstallment)
	No_of_installment := math.Round(float64NoOfInstallment)
	output.No_of_installment = int(No_of_installment)
	fmt.Println("No of Installment Int Type After Ceil: ", output.No_of_installment)
	if output.No_of_installment > 60 {
		res.Message = "Can't be greater than 60 months or 5 years. Please increase monthly installment amount."
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotAcceptable
		return res
	} else if output.No_of_installment < 0 {
		res.Message = "Please increase monthly installment amount."
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotAcceptable
		return res
	}
	res.Message = "Successfully Get No Of Installment Amount"
	res.IsSuccess = true
	res.Payload = output
	res.StatusCode = http.StatusOK

	return res
}

func (loanallotmentRepository *LoanallotmentRepository) BankLoadForLoanForCPF() dto.ResponseDto {
	var res dto.ResponseDto

	// db := util.CreateConnectionToAccountsSchemaUsingGorm()
	// sqlDB, _ := db.DB()
	// defer sqlDB.Close()

	db := util.CreateConnectionToPFAccountsSchemaUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var output []dto.Accid_Name_List_dto
	result := db.Raw("select accid,name from pfaccounts.pfacchead where parent in(10101000) order by name").Find(&output)
	if result.RowsAffected == 0 {
		res.Message = "No Data Found"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotFound
		return res
	}

	res.Message = "Successfully Get All cash & bankid for CPF"
	res.IsSuccess = true
	res.Payload = output
	res.Count = len(output)
	res.StatusCode = http.StatusOK

	return res
}

func (loanallotmentRepository *LoanallotmentRepository) GetABankForLoanForCPF(input dto.Accid_Name_List_dto) dto.ResponseDto {
	var res dto.ResponseDto

	if input.Accid == 0 {
		res.Message = "Please Select Cash/Bank ID"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	db := util.CreateConnectionToPFAccountsSchemaUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var output dto.Accid_Name_List_dto
	result := db.Raw(`(select accid,name from pfaccounts.pfacchead where parent in(10101000) or accid = 10100101)
	INTERSECT
	(select a.accid, a.name from  pfaccounts.pfacchead a  where a.accid = ?)`, input.Accid).First(&output)
	if result.RowsAffected == 0 {
		res.Message = "No Data Found"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotFound
		return res
	}

	res.Message = "Successfully Get Bank"
	res.IsSuccess = true
	res.Payload = output
	res.StatusCode = http.StatusOK

	return res
}

func (loanallotmentRepository *LoanallotmentRepository) MaxNoOfInstallmentForCPF(input model.Loanallot) dto.ResponseDto {
	var res dto.ResponseDto
	if input.No_of_installment == 0 {
		res.Message = "Please Enter No Of Installment"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.No_of_installment > 60 {
		res.Message = "Maximum No Of Installment is 60"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotAcceptable
		return res
	}
	res.Message = "Its a valid no of installment"
	res.IsSuccess = true
	res.Payload = nil
	res.StatusCode = http.StatusOK
	return res
}

// TODO: CPF Loan Insert with Loanpayschedule
func (loanallotmentRepository *LoanallotmentRepository) InsertLoanAllotmentForCPF(input model.Loanallot) dto.ResponseDto {
	var res dto.ResponseDto

	input.Loan_type = 1 // CPF Loan
	input.Branchcode = 0

	if input.Allotment_date == "" {
		res.Message = "Please Select Date"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Empcode == 0 {
		res.Message = "Please Select Employee"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Loan_amount == 0 {
		res.Message = "Please Enter Loan Amount"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.No_of_installment == 0 {
		res.Message = "Please Enter No of Installment"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Interest_rate == 0 {
		res.Message = "Please Enter Interest Rate"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Installment_amount == 0 {
		res.Message = "Please Enter Installment Amount"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Effective_month == 0 {
		res.Message = "Please Select Effective Month"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Effective_year == 0 {
		res.Message = "Please Select Effective Year"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Remarks == "" {
		res.Message = "Please Enter Remarks"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Entry_user == "" {
		res.Message = "Please Enter Entry User"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Compid == 0 {
		res.Message = "Please Select Company Name"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Compyearid == 0 {
		res.Message = "Please Select Company Year"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	tx := db.Begin()
	tx.SavePoint("savepoint")

	forwardingdateRepo := new(LoanallotmentRepository)
	forwardingdateFunc := forwardingdateRepo.ValidateForwardingDate(input)
	if !forwardingdateFunc.IsSuccess {
		res.Message = "Effective date couldn't be previous year month againts current year month"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotAcceptable
		tx.RollbackTo("savepoint")
		return res
	}

	result2 := tx.Raw("select sa.empcode, s.empname, h2.deptcode , d.desigcode from payroll.designation d, payroll.hrdept h2, payroll.staffinformation s , payroll.salarystructure sa where sa.empcode = ? and sa.empcode = s.empcode and h2.deptcode = sa.deptcode and d.desigcode = sa.desigcode", input.Empcode).First(&input)
	if result2.Error != nil {
		res.Message = "Error Occured"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusInternalServerError
		tx.RollbackTo("savepoint")
		return res
	}

	if input.Empname == "" {
		res.Message = "Employee Name is empty"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		tx.RollbackTo("savepoint")
		return res
	}

	if input.Deptcode == 0 {
		res.Message = "Department Name is empty"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		tx.RollbackTo("savepoint")
		return res
	}

	if input.Desigcode == 0 {
		res.Message = "Designation Name is empty"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		tx.RollbackTo("savepoint")
		return res
	}

	var output model.Loanallot
	result1 := tx.Where("refno = ?", input.Refno).First(&output)
	if result1.RowsAffected != 0 {
		res.Message = "Loan Allotment Already Exist"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusConflict
		tx.RollbackTo("savepoint")
		return res
	}
	// ID Autoincrement
	_ = tx.Raw("select coalesce ((max(refno) + 1), 1) from payroll.loanallot").First(&input.Refno)

	input.Due_month = 0
	input.Due_year = 0
	input.Due_principal = 0
	input.Monthly_deduct_interest = 0
	input.Monthly_deduct_pricipal = 0

	// TODO: loanpayschedule table
	var output2 model.Loanpayschedule
	var output3 []model.Loanpayschedule
	output2.Refno = input.Refno
	output2.Empcode = input.Empcode
	output2.Empname = input.Empname
	output2.Interest_rate = input.Interest_rate
	effective_year := input.Effective_year
	effective_month := input.Effective_month

	for i := 1; i <= input.No_of_installment; i++ {

		output2.Installment_id = i
		output2.Year = effective_year + (effective_month+i-2)/12
		output2.Month = (effective_month + i - 1) % 12
		if output2.Month == 0 {
			output2.Month = 12
		}

		_ = tx.Raw("select coalesce ((max(sl_id) + 1), 1) from payroll.loanpayschedule").First(&output2.Sl_id)
		output2.Sl_id--
		output2.Sl_id = output2.Sl_id + output2.Installment_id

		if i == 1 {
			output2.Amount = input.Loan_amount

			Monthly_deduct_interest := output2.Amount * (input.Interest_rate / 1200)
			output2.Monthly_deduct_interest = math.Round(Monthly_deduct_interest)
			output2.Monthly_deduct_principal = input.Installment_amount - output2.Monthly_deduct_interest
			output2.Due_principal = input.Loan_amount - output2.Monthly_deduct_principal
			output2.Installment_amount = input.Installment_amount
		} else if output2.Installment_id > 1 && output2.Installment_id < input.No_of_installment {
			// } else if output2.Installment_id > 1 {
			output2.Amount = output2.Due_principal
			Monthly_deduct_interest := output2.Due_principal * (input.Interest_rate / 1200)
			// Monthly_deduct_interest_int_type, Monthly_deduct_interest_float_type := math.Modf(Monthly_deduct_interest)	// in case of 0.50 or 0.51
			// if Monthly_deduct_interest_float_type >= 0.50 {
			// 	// if Monthly_deduct_interest_float_type >= 0.51 {
			// 	b := Monthly_deduct_interest_int_type + 1
			// 	output2.Monthly_deduct_interest = b
			// } else {
			// 	a := Monthly_deduct_interest_int_type
			// 	output2.Monthly_deduct_interest = a
			// }

			output2.Monthly_deduct_interest = math.Round(Monthly_deduct_interest)
			output2.Monthly_deduct_principal = input.Installment_amount - output2.Monthly_deduct_interest
			output2.Due_principal = output2.Due_principal - output2.Monthly_deduct_principal
			output2.Installment_amount = input.Installment_amount

		} else if output2.Installment_id == input.No_of_installment {
			// } else if output2.Due_principal < input.Installment_amount {
			output2.Amount = output2.Due_principal
			// fmt.Println("output2.Amount", output2.Amount)
			Monthly_deduct_interest := output2.Due_principal * (input.Interest_rate / 1200)
			output2.Monthly_deduct_interest = math.Round(Monthly_deduct_interest)
			output2.Installment_amount = output2.Monthly_deduct_interest + output2.Due_principal
			output2.Monthly_deduct_principal = output2.Installment_amount - output2.Monthly_deduct_interest
			output2.Due_principal = 0
		}

		output3 = append(output3, output2)

		input.Total_interest += output2.Monthly_deduct_interest

	}

	if input.Total_interest == 0 {
		res.Message = "Total Interest is empty"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		tx.RollbackTo("savepoint")
		return res
	}

	// TODO: loanallot insert
	result := tx.Create(&input)
	if result.RowsAffected == 0 {
		res.Message = "Failed to Insert Loan Allotment"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusInternalServerError
		tx.RollbackTo("savepoint")
		return res
	}

	// TODO: loanpayschedule insert
	// _ = tx.Raw("select coalesce ((max(sl_id) + 1), 1) from payroll.loanpayschedule").First(&output2.Sl_id)

	result5 := tx.Create(&output3)
	if result5.RowsAffected == 0 {
		res.Message = "Failed to Insert Loan Allotment in sc"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusInternalServerError
		tx.RollbackTo("savepoint")
		return res
	}

	// for organized console print
	// for _, v := range output3 {
	// 	output3 = append(output3, v)
	// 	fmt.Println("This is V Ulala", v)
	// }
	db1 := util.CreateConnectionToPFAccountsSchemaUsingGorm()
	sqlDB1, _ := db1.DB()
	defer sqlDB1.Close()

	_ = db1.Exec("CALL pfaccounts.cpfloanallotvoucher(? , ? , ? , ? , ? , ? , ? , ? , ?)", input.Empcode, input.Empname, input.Loan_amount, input.Total_interest, input.Compid, input.Compyearid, input.Allotment_date, input.Entry_user, 0)

	tx.Commit()
	res.Message = "Loan Allotment Inserted Successfully"
	res.IsSuccess = true
	res.Payload = input
	res.StatusCode = http.StatusOK

	return res
}

func (loanallotmentRepository *LoanallotmentRepository) GetLoanPayScheduleByEmpcodeForCPF(input model.Loanallot) dto.ResponseDto {
	var res dto.ResponseDto

	if input.Empcode <= 0 {
		res.Message = "Please select your desire employee name from Checking CPF Loan List"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	var output []dto.LoanpayscheduleforCPFOutputDto

	result := db.Raw(`SELECT l2.empcode, l2.empname, l2.total_interest, l.refno, l.installment_id, l."month", l."year", l.amount, l.monthly_deduct_principal, l.monthly_deduct_interest, l.due_principal, l.interest_rate , l.installment_amount
	FROM payroll.loanpayschedule l
	INNER JOIN (
			SELECT las.empcode, las.empname, las.total_interest, MAX(las.refno) AS max_refno
			FROM payroll.loanallot las
			WHERE las.loan_type = 1 AND las.empcode = ?
			GROUP BY las.empcode, las.empname, las.total_interest
	) l2 ON l.refno = l2.max_refno and l2.max_refno = (select max(refno) from payroll.loanallot las where las.empcode = ? and las.loan_type = 1)
	ORDER BY l.installment_id;
	`, input.Empcode, input.Empcode).Find(&output)
	if result.RowsAffected == 0 {
		res.Message = "You have no CPF loan"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	type tempOutput struct {
		Refno                      int                                  `json:"refno"`
		Empname                    string                               `json:"empname"`
		Loan_Amount                float64                              `json:"loan_amount"`
		Total_interest             float64                              `json:"total_interest"`
		Total_amount_with_interest float64                              `json:"total_amount_with_interest"`
		Output                     []dto.LoanpayscheduleforCPFOutputDto `json:"output"`
	}

	var temp tempOutput
	temp.Refno = output[0].Refno
	temp.Empname = output[0].Empname
	temp.Loan_Amount = output[0].Amount
	temp.Total_interest = output[0].Total_interest
	temp.Total_amount_with_interest = output[0].Total_interest + output[0].Amount
	temp.Output = output
	res.Message = "CPF Loan Pay Schedule Found"
	res.IsSuccess = true
	res.Payload = temp
	res.StatusCode = http.StatusOK

	return res
}

func (loanallotmentRepository *LoanallotmentRepository) GetDataEditLoanAllotmentByRefNoForCPF(input model.Loanallot) dto.ResponseDto {
	var res dto.ResponseDto
	if input.Refno <= 0 {
		res.Message = "Plase Enter Your Ref No"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	tx := db.Begin()
	tx.SavePoint("savepoint")
	result1 := tx.Where(&model.Loanallot{Refno: input.Refno, Loan_type: 1}).First(&input)
	if result1.RowsAffected == 0 {
		res.Message = "There is no CPF Loan Allotment against this Ref No"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotFound
		tx.RollbackTo("savepoint")
		return res
	}
	// fmt.Println("test1: ", input)
	result2 := tx.Where(&model.Loanallot{Refno: input.Refno, Loan_type: 1, Empcode: input.Empcode}).First(&input) // 1st option
	// var Output2 dto.GetDataToEditLoanAllotmentOutputDto	 // 2nd option
	// result2 := tx.Raw("select * from payroll.loanallot where refno = ? and loan_type = 1 and empcode = ?", input.Refno, input.Empcode).Find(&Output2) // 2nd option
	if result2.RowsAffected == 0 {
		res.Message = "There is no CPF Loan Allotment against this Ref No"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotFound
		tx.RollbackTo("savepoint")
		return res
	}
	// fmt.Println("test2: ", input)

	type tempOutput struct {
		Flag_value int             `json:"flag_value"`
		Output     model.Loanallot `json:"output"`
	}

	var temp tempOutput
	temp.Flag_value = 1
	temp.Output = input
	// TODO: there is another option that will implete in future (if the refno already exist in monthly deduct table then you can't edit the loan allotment)
	var check1 model.Monthlydeduct
	result3 := tx.Where(&model.Monthlydeduct{Loan_refno: input.Refno}).First(&check1)
	if result3.RowsAffected > 0 {
		res.Message = "You can't edit this loan allotment because it is already in monthly deduct table"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotAcceptable
		tx.RollbackTo("savepoint")
		return res
	}

	tx.Commit()
	res.Message = "Get CPF Loanallot Successfully"
	res.IsSuccess = true
	res.Payload = temp // 1st option
	// Output2.Flag_value = 1 // 2nd option
	// res.Payload = Output2 // 2nd option
	res.StatusCode = http.StatusOK

	return res
}

func (loanallotmentRepository *LoanallotmentRepository) EditLoanAllotmentForCPF(input model.Loanallot) dto.ResponseDto {
	var res dto.ResponseDto

	input.Branchcode = 0

	if input.Allotment_date == "" {
		res.Message = "Please Enter Allotment Date"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Loan_amount <= 0 {
		res.Message = "Please Enter Loan Amount"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Installment_amount <= 0 {
		res.Message = "Please Enter Installment Amount"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.No_of_installment <= 0 {
		res.Message = "Please Enter No Of Installment"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Interest_rate <= 0 {
		res.Message = "Please Enter Interest Rate"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Effective_month == 0 {
		res.Message = "Please Select Effective Month"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Effective_year == 0 {
		res.Message = "Please Select Effective Year"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Remarks == "" {
		res.Message = "Please Enter Remarks"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Entry_user == "" {
		res.Message = "Please Enter Entry User"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Compid == 0 {
		res.Message = "Please Select Company Name"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Compyearid == 0 {
		res.Message = "Please Select Company Year"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	tx := db.Begin()
	tx.SavePoint("savepoint")

	forwardingdateRepo := new(LoanallotmentRepository)
	forwardingdateFunc := forwardingdateRepo.ValidateForwardingDate(input)
	if !forwardingdateFunc.IsSuccess {
		res.Message = "Effective date couldn't be previous year month againts current year month"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotAcceptable
		tx.RollbackTo("savepoint")
		return res
	}

	result2 := tx.Raw("select sa.empcode, s.empname, h2.deptcode , d.desigcode from payroll.designation d, payroll.hrdept h2, payroll.staffinformation s , payroll.salarystructure sa where sa.empcode = ? and sa.empcode = s.empcode and h2.deptcode = sa.deptcode and d.desigcode = sa.desigcode", input.Empcode).First(&input)
	if result2.Error != nil {
		res.Message = "Error Occured"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusInternalServerError
		tx.RollbackTo("savepoint")
		return res
	}

	if input.Empname == "" {
		res.Message = "Employee Name is empty"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		tx.RollbackTo("savepoint")
		return res
	}

	if input.Deptcode == 0 {
		res.Message = "Department Name is empty"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		tx.RollbackTo("savepoint")
		return res
	}

	if input.Desigcode == 0 {
		res.Message = "Designation Name is empty"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		tx.RollbackTo("savepoint")
		return res
	}

	var input2 model.Loanpayschedule

	result3 := tx.Where(&model.Loanpayschedule{Refno: input.Refno}).Delete(&input2)
	if result3.RowsAffected == 0 {
		res.Message = "Something went wrong. Please try again later."
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusInternalServerError
		tx.RollbackTo("savepoint")
		return res
	}

	var output2 model.Loanpayschedule
	var output3 []model.Loanpayschedule
	output2.Refno = input.Refno
	output2.Empcode = input.Empcode
	output2.Empname = input.Empname
	output2.Interest_rate = input.Interest_rate
	effective_year := input.Effective_year
	effective_month := input.Effective_month

	for i := 1; i <= input.No_of_installment; i++ {

		output2.Installment_id = i
		output2.Year = effective_year + (effective_month+i-2)/12
		output2.Month = (effective_month + i - 1) % 12
		if output2.Month == 0 {
			output2.Month = 12
		}

		_ = tx.Raw("select coalesce ((max(sl_id) + 1), 1) from payroll.loanpayschedule").First(&output2.Sl_id)
		output2.Sl_id--
		output2.Sl_id = output2.Sl_id + output2.Installment_id

		if i == 1 {
			output2.Amount = input.Loan_amount

			Monthly_deduct_interest := output2.Amount * (input.Interest_rate / 1200)
			output2.Monthly_deduct_interest = math.Round(Monthly_deduct_interest)
			output2.Monthly_deduct_principal = input.Installment_amount - output2.Monthly_deduct_interest
			output2.Due_principal = input.Loan_amount - output2.Monthly_deduct_principal
			output2.Installment_amount = input.Installment_amount
		} else if output2.Installment_id > 1 && output2.Installment_id < input.No_of_installment {
			// } else if output2.Installment_id > 1 {
			output2.Amount = output2.Due_principal
			Monthly_deduct_interest := output2.Due_principal * (input.Interest_rate / 1200)
			// Monthly_deduct_interest_int_type, Monthly_deduct_interest_float_type := math.Modf(Monthly_deduct_interest)	// in case of 0.50 or 0.51
			// if Monthly_deduct_interest_float_type >= 0.50 {
			// 	// if Monthly_deduct_interest_float_type >= 0.51 {
			// 	b := Monthly_deduct_interest_int_type + 1
			// 	output2.Monthly_deduct_interest = b
			// } else {
			// 	a := Monthly_deduct_interest_int_type
			// 	output2.Monthly_deduct_interest = a
			// }

			output2.Monthly_deduct_interest = math.Round(Monthly_deduct_interest)
			output2.Monthly_deduct_principal = input.Installment_amount - output2.Monthly_deduct_interest
			output2.Due_principal = output2.Due_principal - output2.Monthly_deduct_principal
			output2.Installment_amount = input.Installment_amount

		} else if output2.Installment_id == input.No_of_installment {
			// } else if output2.Due_principal < input.Installment_amount {
			output2.Amount = output2.Due_principal
			// fmt.Println("output2.Amount", output2.Amount)
			Monthly_deduct_interest := output2.Due_principal * (input.Interest_rate / 1200)
			output2.Monthly_deduct_interest = math.Round(Monthly_deduct_interest)
			output2.Installment_amount = output2.Monthly_deduct_interest + output2.Due_principal
			output2.Monthly_deduct_principal = output2.Installment_amount - output2.Monthly_deduct_interest
			output2.Due_principal = 0
		}

		output3 = append(output3, output2)

		input.Total_interest += output2.Monthly_deduct_interest

	}

	if input.Total_interest == 0 {
		res.Message = "Total Interest is empty"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		tx.RollbackTo("savepoint")
		return res
	}

	result1 := tx.Where(&model.Loanallot{Refno: input.Refno, Empcode: input.Empcode, Loan_type: 1}).Updates(&input)
	if result1.RowsAffected == 0 {
		res.Message = "Something went wrong. Please try again later."
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusInternalServerError
		tx.RollbackTo("savepoint")
		return res
	}

	result5 := tx.Create(&output3)
	if result5.RowsAffected == 0 {
		res.Message = "Failed to Insert Loan Allotment in sc"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusInternalServerError
		tx.RollbackTo("savepoint")
		return res
	}

	db1 := util.CreateConnectionToPFAccountsSchemaUsingGorm()
	sqlDB1, _ := db1.DB()
	defer sqlDB1.Close()

	_ = db1.Exec("CALL pfaccounts.cpfloanallotvoucher(? , ? , ? , ? , ? , ? , ? , ? , ?)", input.Empcode, input.Empname, input.Loan_amount, input.Total_interest, input.Compid, input.Compyearid, input.Allotment_date, input.Entry_user, 1)

	tx.Commit()
	res.Message = "Edit CPC Loan Successfully"
	res.IsSuccess = true
	res.Payload = input
	res.StatusCode = http.StatusOK

	return res
}

// TODO: Advance Loan
// FIXME: this function is not using
func (loanallotmentRepository *LoanallotmentRepository) GetNoOfInstallmentForAdvance(input model.Loanallot) dto.ResponseDto {
	var res dto.ResponseDto

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var output dto.NoOfInstallmentOutputDto

	remainder := int(input.Loan_amount) % int(input.Installment_amount)

	fracNo_of_Installment := input.Loan_amount / input.Installment_amount

	output.No_of_installment = int(input.Loan_amount) / int(input.Installment_amount)
	roundstr := strconv.Itoa(output.No_of_installment)

	if remainder != 0 {
		res.Message = "Installment Amount is not valid. Installment amount should be around  " + roundstr
		res.IsSuccess = false
		res.Payload = fracNo_of_Installment
		res.StatusCode = http.StatusBadRequest
		return res
	}

	res.Message = "Successfully Get No Of Installment"
	res.IsSuccess = true
	res.Payload = output
	res.StatusCode = http.StatusOK
	return res
}

func (loanallotmentRepository *LoanallotmentRepository) TerminateAdvanceAgainstSalary(input model.Loanallot) dto.ResponseDto {
	var res dto.ResponseDto

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	tx := db.Begin()
	tx.SavePoint("savepoint")
	var output1 model.Loanpayschedule
	// result1 := tx.Raw("select * from payroll.staffinformation s where s.empcode = ?", input.Empcode).First(&input)
	// if result1.RowsAffected == 0 {
	// 	res.Message = "There is no employee for this empcode."
	// 	res.IsSuccess = false
	// 	res.Payload = nil
	// 	res.StatusCode = http.StatusNotFound
	// 	tx.RollbackTo("savepoint")
	// 	return res
	// }

	result2 := tx.Raw("select lp.empcode , lp.month , lp.year  from payroll.loanpayschedule lp  where lp.refno = (select max(refno) from payroll.loanallot las where las.empcode = ? and las.loan_type = 0) order by lp.installment_id desc limit 1", input.Empcode).First(&output1)
	if result2.RowsAffected == 0 {
		res.Message = "No previous Advance against salary found."
		res.IsSuccess = true
		res.Payload = nil
		res.StatusCode = http.StatusOK
		tx.RollbackTo("savepoint")
		return res
	}

	dt := time.Now()
	fulldate := dt.Format("2006-01-02")
	fmt.Println("fulldate: ", fulldate)

	y := fulldate[0:4]
	outputyear := y
	fmt.Println("After Slice Year", outputyear)
	outputyearint, err := strconv.Atoi(outputyear) // Convert String to Int
	if err != nil {                                // Error Handling
		panic(err)
	}
	fmt.Println("Year", output1.Year)
	if output1.Year > outputyearint {
		res.Message = "Already Advance against salary running."
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotAcceptable
		tx.RollbackTo("savepoint")
		return res
	} else if output1.Year < outputyearint {
		res.Message = "There is no Advance against salary running for this empcode."
		res.IsSuccess = true
		res.Payload = nil
		res.StatusCode = http.StatusOK
		return res
	}

	m := fulldate[5:7]
	outputmonth := m
	fmt.Println("After Slice month", outputmonth)
	outputmonthint, err := strconv.Atoi(outputmonth) // Convert String to Int
	if err != nil {                                  // Error Handling
		panic(err)
	}
	fmt.Println("Month", output1.Month)

	if output1.Year == outputyearint && output1.Month >= outputmonthint {
		res.Message = "Already Advance against salary running."
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotAcceptable
		tx.RollbackTo("savepoint")
		return res
	} else if output1.Year == outputyearint && output1.Month < outputmonthint {
		res.Message = "There is no Advance against salary running for this empcode."
		res.IsSuccess = true
		res.Payload = nil
		res.StatusCode = http.StatusOK
		return res
	}

	tx.Commit()

	return res
}

func (loanallotmentRepository *LoanallotmentRepository) GetAllEmployeeForAdvance() dto.ResponseDto {
	var res dto.ResponseDto

	db := util.CreateConnectionToAccountsSchemaUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var output []dto.Accid_Name_List_dto
	result := db.Raw("select a.accid, a.name from accounts.acchead a where a.parent in(11102000) and a.lr = 'L' order by name").Find(&output)
	if result.RowsAffected == 0 {
		res.Message = "No Data Found"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotFound
		return res
	}

	res.Message = "Successfully Get All Employee List"
	res.IsSuccess = true
	res.Payload = output
	res.Count = len(output)
	res.StatusCode = http.StatusOK

	return res
}

func (loanallotmentRepository *LoanallotmentRepository) GetSingleEmployeeForAdvance(input model.Loanallot) dto.ResponseDto {
	var res dto.ResponseDto

	var output dto.GetSingleEmployeeOutputDto

	db := util.CreateConnectionToAccountsSchemaUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	result := db.Raw("select a.accid as stuff_accid from accounts.acchead a where a.parent in(11102000) and a.lr = 'L' and a.accid = ?", input.Stuff_accid).First(&output)
	if result.RowsAffected == 0 {
		res.Message = "No Data Found"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotFound
		return res
	}

	res.Message = "Successfully Get Single Employee"
	res.IsSuccess = true
	res.Payload = output
	res.StatusCode = http.StatusOK

	return res
}

func (loanallotmentRepository *LoanallotmentRepository) GeInstallmentAmountForAdvance(input model.Loanallot) dto.ResponseDto {
	var res dto.ResponseDto

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var output dto.InstallmentAmountOutputDto

	//remainder := int(input.Loan_amount) % int(input.No_of_installment)

	//fracInstallment_amount := input.Loan_amount / float64(input.No_of_installment)

	output.Installment_amount = input.Loan_amount / float64(input.No_of_installment)
	//roundstr := strconv.Itoa(int(output.Installment_amount))

	// if remainder != 0 {
	// 	res.Message = "Installment Amount is not valid. Installment amount should be around  " + roundstr
	// 	res.IsSuccess = false
	// 	res.Payload = fracInstallment_amount
	// 	res.StatusCode = http.StatusBadRequest
	// 	return res
	// }

	res.Message = "Successfully Get No Of Installment"
	res.IsSuccess = true
	res.Payload = output
	res.StatusCode = http.StatusOK
	return res
}

func (loanallotmentRepository *LoanallotmentRepository) GetNoOfInstallmentWithModForAdvance(input model.Loanallot) dto.ResponseDto {
	var res dto.ResponseDto

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var output dto.NoOfInstallmentOutputDto

	No_of_installment := input.Loan_amount / input.Installment_amount

	No_of_installment = math.Round(No_of_installment)

	output.No_of_installment = int(No_of_installment)

	res.Message = "Successfully Get No Of Installment"
	res.IsSuccess = true
	res.Payload = output
	res.StatusCode = http.StatusOK
	return res
}

// TODO: Advance Loan Insert with LoanPaySchedule
func (loanallotmentRepository *LoanallotmentRepository) InsertLoanAllotmentForAdvance(input model.Loanallot) dto.ResponseDto {
	var res dto.ResponseDto

	input.Loan_type = 0 // Advance against salary
	input.Branchcode = 0

	if input.Allotment_date == "" {
		res.Message = "Please Select Date"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Empcode == 0 {
		res.Message = "Please Select Employee"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Loan_amount == 0 {
		res.Message = "Please Enter Loan Amount"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.No_of_installment == 0 {
		res.Message = "Please Enter No of Installment"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Installment_amount == 0 {
		res.Message = "Please Enter Installment Amount"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Effective_month == 0 {
		res.Message = "Please Select Effective Month"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Effective_year == 0 {
		res.Message = "Please Select Effective Year"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Remarks == "" {
		res.Message = "Please Enter Remarks"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Entry_user == "" {
		res.Message = "Please Enter Entry User"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Compid == 0 {
		res.Message = "Please Select Company Name"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Compyearid == 0 {
		res.Message = "Please Select Company Year"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Stuff_accid == 0 {
		res.Message = "Please Select Stuff Account"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	tx := db.Begin()
	tx.SavePoint("savepoint")

	forwardingdateRepo := new(LoanallotmentRepository)
	forwardingdateFunc := forwardingdateRepo.ValidateForwardingDate(input)
	if !forwardingdateFunc.IsSuccess {
		res.Message = "Effective date couldn't be previous year month againts current year month"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotAcceptable
		tx.RollbackTo("savepoint")
		return res
	}

	result2 := tx.Raw("select sa.empcode, s.empname, h2.deptcode , d.desigcode from payroll.designation d, payroll.hrdept h2, payroll.staffinformation s , payroll.salarystructure sa where sa.empcode = ? and sa.empcode = s.empcode and h2.deptcode = sa.deptcode and d.desigcode = sa.desigcode", input.Empcode).First(&input)
	if result2.Error != nil {
		res.Message = "Error Occured"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusInternalServerError
		tx.RollbackTo("savepoint")
		return res
	}

	if input.Empname == "" {
		res.Message = "Employee Name is empty"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		tx.RollbackTo("savepoint")
		return res
	}

	if input.Deptcode == 0 {
		res.Message = "Department Name is empty"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		tx.RollbackTo("savepoint")
		return res
	}

	if input.Desigcode == 0 {
		res.Message = "Designation Name is empty"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		tx.RollbackTo("savepoint")
		return res
	}

	var output model.Loanallot
	result1 := tx.Where("refno = ?", input.Refno).First(&output)
	if result1.RowsAffected != 0 {
		res.Message = "Loan Allotment Already Exist"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusConflict
		tx.RollbackTo("savepoint")
		return res
	}
	// ID Autoincrement
	_ = tx.Raw("select coalesce ((max(refno) + 1), 1) from payroll.loanallot").First(&input.Refno)

	input.Due_month = 0
	input.Due_year = 0
	input.Due_principal = 0
	input.Monthly_deduct_interest = 0
	input.Monthly_deduct_pricipal = 0
	input.Interest_rate = 0
	input.Total_interest = 0

	// TODO: loanpayschedule table
	var output2 model.Loanpayschedule
	var output3 []model.Loanpayschedule
	output2.Refno = input.Refno
	output2.Empcode = input.Empcode
	output2.Empname = input.Empname
	output2.Interest_rate = input.Interest_rate
	effective_year := input.Effective_year
	effective_month := input.Effective_month

	for i := 1; i <= input.No_of_installment; i++ {

		output2.Installment_id = i
		output2.Year = effective_year + (effective_month+i-2)/12
		output2.Month = (effective_month + i - 1) % 12
		if output2.Month == 0 {
			output2.Month = 12
		}

		_ = tx.Raw("select coalesce ((max(sl_id) + 1), 1) from payroll.loanpayschedule").First(&output2.Sl_id)
		output2.Sl_id--
		output2.Sl_id = output2.Sl_id + output2.Installment_id
		if i == 1 {
			output2.Amount = input.Loan_amount
			output2.Installment_amount = input.Installment_amount
			output2.Monthly_deduct_principal = input.Installment_amount
			output2.Due_principal = output2.Amount - input.Installment_amount
		} else if output2.Installment_id > 1 && output2.Installment_id < input.No_of_installment {
			// output2.Amount = output3[i-1].Due_principal
			output2.Amount = output2.Due_principal
			output2.Installment_amount = input.Installment_amount
			output2.Monthly_deduct_principal = input.Installment_amount
			output2.Due_principal = output2.Amount - input.Installment_amount
		} else if output2.Installment_id == input.No_of_installment {
			output2.Amount = output2.Due_principal
			output2.Installment_amount = output2.Due_principal
			output2.Monthly_deduct_principal = output2.Installment_amount
			output2.Due_principal = output2.Due_principal - output2.Installment_amount
		}

		output3 = append(output3, output2)

	}

	// TODO: loanallot insert
	result := tx.Create(&input)
	if result.RowsAffected == 0 {
		res.Message = "Failed to Insert Loan Allotment"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusInternalServerError
		tx.RollbackTo("savepoint")
		return res
	}

	// TODO: loanpayschedule insert
	result5 := tx.Create(&output3)
	if result5.RowsAffected == 0 {
		res.Message = "Failed to Insert Loan Allotment in sc"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusInternalServerError
		tx.RollbackTo("savepoint")
		return res
	}

	db1 := util.CreateConnectionToPFAccountsSchemaUsingGorm()
	sqlDB1, _ := db1.DB()
	defer sqlDB1.Close()

	_ = db1.Exec("CALL pfaccounts.saladvallotvoucher(? , ? , ? , ? , ? , ? , ? , ?)", input.Stuff_accid, input.Empname, input.Loan_amount, input.Compid, input.Compyearid, input.Allotment_date, input.Entry_user, 0)

	tx.Commit()
	res.Message = "Loan Allotment Inserted Successfully"
	res.IsSuccess = true
	res.Payload = input
	res.StatusCode = http.StatusOK
	return res
}

func (loanallotmentRepository *LoanallotmentRepository) GetLoanPayScheduleByEmpcodeForAdvance(input model.Loanallot) dto.ResponseDto {
	var res dto.ResponseDto

	if input.Empcode <= 0 {
		res.Message = "Please select your desire employee name from Checking Advance Loan List."
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	var output []dto.LoanpayscheduleforAdvanceOutputDto

	result := db.Raw(`SELECT l2.empcode, l2.empname, l.refno, l.installment_id, l."month", l."year", l.amount, l.monthly_deduct_principal, l.monthly_deduct_interest, l.due_principal , l.installment_amount
	FROM payroll.loanpayschedule l
	INNER JOIN (
		SELECT las.empcode, las.empname, MAX(las.refno) AS max_refno
		FROM payroll.loanallot las
		WHERE las.loan_type = 0 AND las.empcode = ?
		GROUP BY las.empcode, las.empname
	) l2 ON l.refno = l2.max_refno and l2.max_refno = (select max(refno) from payroll.loanallot las where las.empcode = ? and las.loan_type = 0)
	ORDER BY l.installment_id`, input.Empcode, input.Empcode).Find(&output)
	if result.RowsAffected == 0 {
		res.Message = "You have no Advance from salary."
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	type tempOutput struct {
		Refno             int                                      `json:"refno"`
		Empname           string                                   `json:"empname"`
		Total_interest    float64                                  `json:"total_interest,omitempty"`
		Total_loan_amount float64                                  `json:"total_loan_amount,omitempty"`
		Output            []dto.LoanpayscheduleforAdvanceOutputDto `json:"output"`
	}

	var temp tempOutput
	temp.Refno = output[0].Refno
	temp.Empname = output[0].Empname
	// temp.Total_interest = output[0].Total_interest
	temp.Total_loan_amount = output[0].Amount
	temp.Output = output
	res.Message = "Advance from salary Schedule Found."
	res.IsSuccess = true
	res.Payload = temp
	res.StatusCode = http.StatusOK

	return res
}

func (loanallotmentRepository *LoanallotmentRepository) GetDataEditLoanAllotmentByRefNoForAdvance(input model.Loanallot) dto.ResponseDto {
	var res dto.ResponseDto
	if input.Refno <= 0 {
		res.Message = "Plase Enter Your Ref No"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	tx := db.Begin()
	tx.SavePoint("savepoint")
	// input.Loan_type = 0
	result1 := tx.Model(&model.Loanallot{}).Where(map[string]interface{}{"refno": input.Refno, "loan_type": 0}).First(&input)
	if result1.RowsAffected == 0 {
		res.Message = "There is no Advance against Salary for this Ref No"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotFound
		tx.RollbackTo("savepoint")
		return res
	}
	// fmt.Println("test1: ", input)
	result2 := tx.Model(&model.Loanallot{}).Where(map[string]interface{}{"refno": input.Refno, "loan_type": 0, "empcode": input.Empcode}).First(&input) // 1st option
	// var Output2 dto.GetDataToEditLoanAllotmentOutputDto	 // 2nd option
	// result2 := tx.Raw("select * from payroll.loanallot where refno = ? and loan_type = 0 and empcode = ?", input.Refno, input.Empcode).Find(&Output2) // 2nd option
	if result2.RowsAffected == 0 {
		res.Message = "There is no Advance against Salary for this Ref No"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotFound
		tx.RollbackTo("savepoint")
		return res
	}
	// fmt.Println("test2: ", input)

	type tempOutput struct {
		Flag_value int             `json:"flag_value"`
		Output     model.Loanallot `json:"output"`
	}

	var temp tempOutput
	temp.Flag_value = 1
	temp.Output = input
	// TODO: there is another option that will implete in future (if the refno already exist in monthly deduct table then you can't edit the loan allotment)
	var check1 model.Monthlydeduct
	result3 := tx.Where(&model.Monthlydeduct{Loan_refno: input.Refno}).First(&check1)
	if result3.RowsAffected > 0 {
		res.Message = "You can't edit this loan allotment because it is already in monthly deduct table"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotAcceptable
		tx.RollbackTo("savepoint")
		return res
	}

	tx.Commit()
	res.Message = "Get Advance against Salary Successfully"
	res.IsSuccess = true
	res.Payload = temp // 1st option
	// Output2.Flag_value = 1 // 2nd option
	// res.Payload = Output2 // 2nd option
	res.StatusCode = http.StatusOK

	return res
}

func (loanallotmentRepository *LoanallotmentRepository) EditLoanAllotmentForAdvance(input model.Loanallot) dto.ResponseDto {
	var res dto.ResponseDto

	if input.Allotment_date == "" {
		res.Message = "Please Select Allotment Date"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Loan_amount <= 0 {
		res.Message = "Please Enter Loan Amount"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Loan_amount <= 0 {
		res.Message = "Please Enter Loan Amount"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.No_of_installment <= 0 {
		res.Message = "Please Enter No of Installment"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Effective_month <= 0 {
		res.Message = "Please Select Effective Month"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Effective_year <= 0 {
		res.Message = "Please Select Effective Year"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Remarks == "" {
		res.Message = "Please Enter Remarks"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Entry_user == "" {
		res.Message = "Please Enter Entry User"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Compid <= 0 {
		res.Message = "Please Select Company Name"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Compyearid <= 0 {
		res.Message = "Please Select Company Year"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Stuff_accid <= 0 {
		res.Message = "Please Select Stuff Account"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	tx := db.Begin()
	tx.SavePoint("savepoint")

	forwardingdateRepo := new(LoanallotmentRepository)
	forwardingdateFunc := forwardingdateRepo.ValidateForwardingDate(input)
	if !forwardingdateFunc.IsSuccess {
		res.Message = "Effective date couldn't be previous year month againts current year month"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotAcceptable
		tx.RollbackTo("savepoint")
		return res
	}

	result2 := tx.Raw("select sa.empcode, s.empname, h2.deptcode , d.desigcode from payroll.designation d, payroll.hrdept h2, payroll.staffinformation s , payroll.salarystructure sa where sa.empcode = ? and sa.empcode = s.empcode and h2.deptcode = sa.deptcode and d.desigcode = sa.desigcode", input.Empcode).First(&input)
	if result2.Error != nil {
		res.Message = "Error Occured"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusInternalServerError
		tx.RollbackTo("savepoint")
		return res
	}

	if input.Empname == "" {
		res.Message = "Employee Name is empty"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		tx.RollbackTo("savepoint")
		return res
	}

	if input.Deptcode == 0 {
		res.Message = "Department Name is empty"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		tx.RollbackTo("savepoint")
		return res
	}

	if input.Desigcode == 0 {
		res.Message = "Designation Name is empty"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		tx.RollbackTo("savepoint")
		return res
	}

	result1 := tx.Model(&model.Loanallot{}).Where(map[string]interface{}{"refno": input.Refno, "empcode": input.Empcode, "loan_type": 0}).Updates(&input)
	if result1.RowsAffected == 0 {
		res.Message = "Something went wrong. Please try again later."
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusInternalServerError
		tx.RollbackTo("savepoint")
		return res
	}

	var input2 model.Loanpayschedule

	result3 := tx.Where(&model.Loanpayschedule{Refno: input.Refno}).Delete(&input2)
	if result3.RowsAffected == 0 {
		res.Message = "Something went wrong. Please try again later."
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusInternalServerError
		tx.RollbackTo("savepoint")
		return res
	}

	// TODO: loanpayschedule table
	var output2 model.Loanpayschedule
	var output3 []model.Loanpayschedule
	output2.Refno = input.Refno
	output2.Empcode = input.Empcode
	output2.Empname = input.Empname
	output2.Interest_rate = input.Interest_rate
	effective_year := input.Effective_year
	effective_month := input.Effective_month

	for i := 1; i <= input.No_of_installment; i++ {

		output2.Installment_id = i
		output2.Year = effective_year + (effective_month+i-2)/12
		output2.Month = (effective_month + i - 1) % 12
		if output2.Month == 0 {
			output2.Month = 12
		}

		_ = tx.Raw("select coalesce ((max(sl_id) + 1), 1) from payroll.loanpayschedule").First(&output2.Sl_id)
		output2.Sl_id--
		output2.Sl_id = output2.Sl_id + output2.Installment_id
		if i == 1 {
			output2.Amount = input.Loan_amount
			output2.Installment_amount = input.Installment_amount
			output2.Monthly_deduct_principal = input.Installment_amount
			output2.Due_principal = output2.Amount - input.Installment_amount
		} else if output2.Installment_id > 1 && output2.Installment_id < input.No_of_installment {
			// output2.Amount = output3[i-1].Due_principal
			output2.Amount = output2.Due_principal
			output2.Installment_amount = input.Installment_amount
			output2.Monthly_deduct_principal = input.Installment_amount
			output2.Due_principal = output2.Amount - input.Installment_amount
		} else if output2.Installment_id == input.No_of_installment {
			output2.Amount = output2.Due_principal
			output2.Installment_amount = output2.Due_principal
			output2.Monthly_deduct_principal = output2.Installment_amount
			output2.Due_principal = output2.Due_principal - output2.Installment_amount
		}

		output3 = append(output3, output2)

	}

	result5 := tx.Create(&output3)
	if result5.RowsAffected == 0 {
		res.Message = "Failed to Insert Loan Allotment Schedule"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusInternalServerError
		tx.RollbackTo("savepoint")
		return res
	}

	db1 := util.CreateConnectionToPFAccountsSchemaUsingGorm()
	sqlDB1, _ := db1.DB()
	defer sqlDB1.Close()
	_ = db1.Exec("CALL pfaccounts.saladvallotvoucher(? , ? , ? , ? , ? , ? , ? , ?)", input.Stuff_accid, input.Empname, input.Loan_amount, input.Compid, input.Compyearid, input.Allotment_date, input.Entry_user, 1)

	tx.Commit()
	res.Message = "Edit Advance against Salary Successfully"
	res.IsSuccess = true
	res.Payload = input
	res.StatusCode = http.StatusOK

	return res
}

// FIXME: not ready yet
func (loanallotmentRepository *LoanallotmentRepository) InsertLoanAllotment(input model.Loanallot) dto.ResponseDto {
	var res dto.ResponseDto

	dt := time.Now()

	input.Allotment_date = dt.Format("01-02-2006")

	if input.Allotment_date == "" {
		res.Message = "Please Select Date"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Loan_type <= -1 || input.Loan_type >= 2 {
		res.Message = "Please Select Loan Type"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Loan_type == 0 {
		if input.Interest_rate != 0 {
			res.Message = "Interest Rate remain 0"
			res.IsSuccess = false
			res.Payload = nil
			res.StatusCode = http.StatusBadRequest
			return res
		}

		if input.Total_interest != 0 {
			res.Message = "Total Interest remain 0"
			res.IsSuccess = false
			res.Payload = nil
			res.StatusCode = http.StatusBadRequest
			return res
		}

		if input.Monthly_deduct_interest != 0 {
			res.Message = "Monthly Deduct Interest remain 0"
			res.IsSuccess = false
			res.Payload = nil
			res.StatusCode = http.StatusBadRequest
			return res
		}
	}

	if input.Loan_type == 1 {
		if input.Interest_rate == 0 {
			res.Message = "Please Enter Interest Rate"
			res.IsSuccess = false
			res.Payload = nil
			res.StatusCode = http.StatusBadRequest
			return res
		}

		if input.Total_interest == 0 {
			res.Message = "Please Enter Total Interest"
			res.IsSuccess = false
			res.Payload = nil
			res.StatusCode = http.StatusBadRequest
			return res
		}

		if input.Monthly_deduct_interest == 0 {
			res.Message = "Please Enter Monthly Deduct Interest"
			res.IsSuccess = false
			res.Payload = nil
			res.StatusCode = http.StatusBadRequest
			return res
		}
	}

	if input.Loan_amount == 0 {
		res.Message = "Please Enter Loan Amount"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.No_of_installment == 0 {
		res.Message = "Please Enter No of Installment"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Installment_amount == 0 {
		res.Message = "Please Enter Installment Amount"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Effective_month == 0 {
		res.Message = "Please Select Effective Month"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Effective_year == 0 {
		res.Message = "Please Select Effective Year"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Empname == "" {
		res.Message = "Please Select Employee Name"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Branchcode == 0 {
		res.Message = "Please Select Branch Name"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Remarks == "" {
		res.Message = "Please Enter Remarks"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Deptcode == 0 {
		res.Message = "Please Select Department Name"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Desigcode == 0 {
		res.Message = "Please Select Designation Name"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Due_month == 0 {
		res.Message = "Please Select Due Month"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Due_year == 0 {
		res.Message = "Please Select Due Year"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Due_principal == 0 {
		res.Message = "Please Enter Due Pricipal"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	if input.Entry_user == "" {
		res.Message = "Please Enter Entry User"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var output model.Loanallot
	result1 := db.Where("refno = ?", input.Refno).First(&output)
	if result1.RowsAffected != 0 {
		res.Message = "Loan Allotment Already Exist"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusConflict
		return res
	}
	// ID Autoincrement
	_ = db.Raw("select coalesce ((max(refno) + 1), 1) from payroll.loanallot").First(&input.Refno)

	result := db.Create(&input)
	if result.RowsAffected == 0 {
		res.Message = "Failed to Insert Loan Allotment"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusInternalServerError
		return res
	}

	res.Message = "Loan Allotment Inserted Successfully"
	res.IsSuccess = true
	res.Payload = input
	res.StatusCode = http.StatusOK

	return res
}

func (loanallotmentRepository *LoanallotmentRepository) NextLoanCycle(input model.Loanallot) dto.ResponseDto {
	var res dto.ResponseDto

	if input.Empcode == 0 {
		res.Message = "Please Select Employee"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var input1 model.Loanallot
	// result1 := db.Where("empcode = ?", input.Empcode).First(&output)
	// result1 := db.Raw(" SELECT * FROM payroll.loanallot WHERE empcode = ?", input.Empcode).First(&output)
	result1 := db.Where(&model.Loanallot{Empcode: input.Empcode}).First(&input1)
	if result1.RowsAffected == 0 {
		res.Message = "No Loan Allotment Found"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotFound
		return res
	}

	// fmt.Println("Input Loan Amount: ", input1.Loan_amount)

	var output2 model.Loanallot
	var output model.Loanallot
	output2.Loan_amount = input1.Loan_amount - input1.Installment_amount
	// fmt.Println("Test1 : ", output2.Loan_amount)
	// output.Installment_amount = input1.Installment_amount
	if output2.Loan_amount+1 <= input1.Installment_amount {
		// fmt.Println("Test2 : ", output2.Loan_amount)
		// output.Loan_amount = input.Installment_amount
		// result2 := db.Raw(`UPDATE payroll.loanallot
		// SET loan_amount =0
		// WHERE empcode = ?`, input.Empcode).Updates(&output2)
		if output2.Loan_amount < 0 {
			output2.Loan_amount = 0
		}
		result2 := db.Model(&model.Loanallot{}).Where(&model.Loanallot{Empcode: input.Empcode}).Updates(map[string]interface{}{"loan_amount": output2.Loan_amount})
		if result2.RowsAffected == 0 {
			res.Message = "Server Error. Please try again later."
			res.IsSuccess = false
			res.Payload = nil
			res.StatusCode = http.StatusInternalServerError
			return res
		}
		// fmt.Println("Test3 : ", output2.Loan_amount)
		res.Message = "Successfully Get Loan Allotment"
		res.IsSuccess = true
		res.Payload = output2
		res.StatusCode = http.StatusOK
		return res

	}

	output.Loan_amount = output2.Loan_amount
	// fmt.Println("Output Loan Amount vvvv: ", output.Loan_amount)
	fmt.Println("Test4 : ", output.Loan_amount)
	result := db.Where(&model.Loanallot{Empcode: input.Empcode}).Updates(&output)
	if result.RowsAffected == 0 {
		res.Message = "Server Error. Please try again later."
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusInternalServerError
		return res
	}

	res.Message = "Successfully Get Loan Allotment"
	res.IsSuccess = true
	res.Payload = output
	res.StatusCode = http.StatusOK

	return res
}
