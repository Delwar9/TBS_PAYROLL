package repository

import (
	"fmt"
	"net/http"

	"github.com/tools/payroll/dto"
	"github.com/tools/payroll/model"
	"github.com/tools/payroll/util"
)

type MonthlydeductRepository struct{}

func (monthlydeductRepository *MonthlydeductRepository) InsertMonthyDeductReposioty(input model.Monthlydeduct) dto.ResponseDto {
	var res dto.ResponseDto

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	tx := db.Begin()
	tx.SavePoint("savepoint")

	var output1 []model.Loanpayschedule

	// result1 := tx.Where(&model.Loanpayschedule{Month: input.Month, Year: input.Year, Pause_flag: 0}).Find(&output1)
	result1 := tx.Model(&model.Loanpayschedule{}).Where(map[string]interface{}{"month": input.Month, "year": input.Year, "pause_flag": 0}).Find(&output1)
	if result1.RowsAffected == 0 {
		res.Message = "No Data Found"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotFound
		tx.RollbackTo("savepoint")
		return res
	}

	for _, v := range output1 {
		var check5 model.Loanallot
		_ = tx.Raw("select l.loan_type from payroll.loanallot l where refno = ?", v.Refno).First(&check5)
		if check5.Loan_type == 1 {
			fmt.Println("----CPF----")
			var check1 dto.CPF_total_loan_with_insterest
			var check2 dto.CPF_total_loan_with_insterest
			_ = tx.Raw("select ((select l.loan_amount from payroll.loanallot l where l.empcode = ? and l.loan_type = 1) + sum(l.monthly_deduct_interest)) as total_loan_with_interest  from payroll.loanpayschedule l where l.refno = (select max(refno) from payroll.loanallot l where l.empcode = ? and l.loan_type = 1)", v.Empcode, v.Empcode).First(&check1)
			_ = tx.Raw("select sum(m.cpfloan) as total_loan_with_interest from payroll.monthlydeduct m where loan_refno =  (select max(refno) from payroll.loanallot l where l.empcode = ? and l.loan_type = 1)", v.Empcode).First(&check2)
			fmt.Println("Check1: ", check1.Total_loan_with_interest)
			fmt.Println("Check2: ", check2.Total_loan_with_interest)
			if check1.Total_loan_with_interest == check2.Total_loan_with_interest {
				res.Message = "Already Completed CPF Deduct"
				res.IsSuccess = false
				res.Payload = nil
				res.StatusCode = http.StatusNotAcceptable
				tx.RollbackTo("savepoint")
				return res
			}
		} else if check5.Loan_type == 0 {
			fmt.Println("----Advance----")
			var check3 dto.Advance_total_loan
			var check4 dto.Advance_total_loan
			_ = tx.Raw("select sum(l.monthly_deduct_principal) as total_loan from payroll.loanpayschedule l where l.refno = (select max(refno) from payroll.loanallot l where l.empcode = ? and l.loan_type = 0)", v.Empcode).First(&check3)
			_ = tx.Raw("select sum(m.cpfloan) as total_loan from payroll.monthlydeduct m where loan_refno =  (select max(refno) from payroll.loanallot l where l.empcode = ? and l.loan_type = 0)", v.Empcode).First(&check4)
			fmt.Println("Check3: ", check3.Total_loan)
			fmt.Println("Check4: ", check4.Total_loan)
			if check3.Total_loan == check4.Total_loan {
				res.Message = "Already Completed Advance from Salary Deduct"
				res.IsSuccess = false
				res.Payload = nil
				res.StatusCode = http.StatusNotAcceptable
				tx.RollbackTo("savepoint")
				return res
			}
		} else {
			res.Message = "No Loan Type Found"
			res.IsSuccess = false
			res.Payload = nil
			res.StatusCode = http.StatusNotFound
			tx.RollbackTo("savepoint")
			return res
		}
	}

	var input3 model.Loanallot

	fmt.Println("TEST 1: ", input.Id)
	// var output3 []model.Monthlydeduct
	var input2 model.Monthlydeduct
	check2 := tx.Raw("select month, year from payroll.monthlydeduct where month = ? and year = ?", input.Month, input.Year).First(&input2)
	if check2.RowsAffected > 0 {
		if input2.Month == input.Month && input2.Year == input.Year {
			res.Message = "Already Inserted"
			res.IsSuccess = false
			res.Payload = nil
			res.StatusCode = http.StatusConflict
			tx.RollbackTo("savepoint")
			return res
		}
	}
	for _, v := range output1 {
		o := 1
		fmt.Println("lenth of O:", o)
		// FIXME: Need to check the condition after confirm with sir(duplicate data)

		input.Empcode = v.Empcode
		input.Loan_refno = v.Refno

		_ = tx.Raw("select * from payroll.loanallot l where l.refno = ?", v.Refno).First(&input3)
		input.Stuff_accid = input3.Stuff_accid
		input.Stuff_bankid = input3.Cash_bankid
		input.Monthly_deduct_principal = v.Monthly_deduct_principal
		input.Monthly_deduct_interest = v.Monthly_deduct_interest
		input.Due_principal = v.Due_principal
		if input3.Loan_type == 0 {
			input.Loan_type = 0
			input.Salaryadv = v.Installment_amount
			input.Cpfloan = 0
		} else if input3.Loan_type == 1 {
			input.Loan_type = 1
			input.Cpfloan = v.Installment_amount
			input.Salaryadv = 0
		}
		_ = tx.Raw("select coalesce ((max(id) + 1), 1) from payroll.monthlydeduct").First(&input.Id)
		fmt.Println("TEST 2: ", input.Id)
		input.Id--
		fmt.Println("TEST 3: ", input.Id)
		input.Id = input.Id + o
		fmt.Println("TEST 4: ", input.Id)
		_ = tx.Create(&input)
		fmt.Println("TEST 5: ", input.Id)
		// output3 = append(output3, input)
		o++
	}

	var output2 []model.Monthlydeduct
	_ = tx.Raw("select * from payroll.monthlydeduct where month = ? and year = ?", input.Month, input.Year).Find(&output2)
	// result1 := tx.Raw("select *  from payroll.loanpayschedule l2 where l2."month" = 4  and l2."year" = 2023")

	// TODO: this sql will be using in future
	// select ((select l.loan_amount from payroll.loanallot l where l.empcode = 2307 and l.loan_type = 1) + sum(l.monthly_deduct_interest)) as total_loan_with_interest  from payroll.loanpayschedule l where l.refno = (select max(refno) from payroll.loanallot l where l.empcode = 2307 and l.loan_type = 1)

	// select sum(l.monthly_deduct_principal) as total_loan from payroll.loanpayschedule l where l.refno = (select max(refno) from payroll.loanallot l where l.empcode = 2307 and l.loan_type = 0)
	tx.Commit()
	res.Message = "Successfully Inserted"
	res.IsSuccess = true
	res.Payload = output2
	res.StatusCode = http.StatusOK

	return res
}
