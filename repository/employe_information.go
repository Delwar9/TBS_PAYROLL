package repository

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tools/payroll/dto"
	"github.com/tools/payroll/model"
	"github.com/tools/payroll/util"
)

type EmployeeInformation struct{}

func (e *EmployeeInformation) AddEmployeeInformation(input dto.EmployeeInformation) dto.ResponseDto {
	var res dto.ResponseDto
	var latest dto.EmployeeInformation

	fmt.Println("Hello")
	if input.StaffInformation.Empcode == 0 {
		res.Message = "Emp Code is required"
		res.IsSuccess = false
		res.StatusCode = http.StatusBadRequest
		res.Payload = nil

		return res
	}

	if input.StaffInformation.Accno == "" {
		res.Message = "Account Number is required"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}
	if input.StaffInformation.Senioriy_serial <= 0 {
		res.Message = "Serial number is required"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}

	if input.StaffInformation.Empname == "" {
		res.Message = "Name is required"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}

	if input.StaffInformation.Joindate == "" {
		res.Message = "Joining date is required"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}

	// if input.StaffInformation.Mailadd == "" {
	// 	res.Message = "Mail address is required"
	// 	res.StatusCode = http.StatusBadRequest
	// 	res.IsSuccess = false
	// 	res.Payload = nil
	// 	return res
	// }
	// if input.StaffInformation.Permanentadd == "" {
	// 	res.Message = "Permanent address is required"
	// 	res.StatusCode = http.StatusBadRequest
	// 	res.IsSuccess = false
	// 	res.Payload = nil
	// 	return res
	// }
	if input.StaffInformation.Eduquali == "" {
		res.Message = "Education qualification is required"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}
	if input.StaffInformation.Gender == "" {
		res.Message = "Gender is required"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}
	if input.StaffInformation.Grade == "" {
		res.Message = "Grade is required"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}
	if input.StaffInformation.Maritalstatus == "" {
		res.Message = "Marital status is required"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}
	if input.StaffInformation.Religion == "" {
		res.Message = "Religion is required"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}
	if input.StaffInformation.Bloodgroup == "" {
		res.Message = "Salary is required"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}
	if input.StaffInformation.Active_salary < 0 && input.StaffInformation.Active_salary > 1 {
		res.Message = "Active status is either 0 or 1"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}

	// validation for salary structure
	if input.SalaryStructure.Bankid <= 0 {
		res.Message = "Bank id is required"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	tx := db.Begin()
	tx.SavePoint("savePoint1")
	db1 := util.CreateConnectionToPFAccountsSchemaUsingGorm()
	sqlDB1, _ := db1.DB()
	defer sqlDB1.Close()

	tx1 := db1.Begin()
	tx1.SavePoint("savePoint1")

	// var input model.EmployeeInformation

	// validation for staff information
	var output dto.EmployeeInformation

	res1 := tx.Raw("select * from payroll.staffinformation where empcode=?", input.StaffInformation.Empcode).First(&output.StaffInformation)

	if res1.RowsAffected != 0 {
		tx.RollbackTo("savePoint1")
		res.Message = "Employee code already exists"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}
	res2 := tx.Raw("select * from payroll.staffinformation where senioriy_serial=?", input.StaffInformation.Senioriy_serial).First(&output.StaffInformation)

	if res2.RowsAffected != 0 {
		tx.RollbackTo("savePoint1")
		res.Message = "Serial already exists"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}

	if input.SalaryStructure.Basic > 0 {
		input.SalaryStructure.Consolited = 0
	}
	if input.SalaryStructure.Consolited > 0 {
		input.SalaryStructure.Basic = 0
	}
	if input.SalaryStructure.Basic > 0 && input.SalaryStructure.Consolited > 0 {

		res.Message = "Either Basic or Consolited salary is required"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}
	if input.SalaryStructure.Basic <= 0 && input.SalaryStructure.Consolited <= 0 {
		res.IsSuccess = false
		res.Message = "Neither basic nor consolited can't be zero"
		res.StatusCode = http.StatusBadRequest
		res.Payload = nil
		return res
	}
	if input.SalaryStructure.Basic > 0 && input.SalaryStructure.Consolited <= 0 {
		input.SalaryStructure.Consolited = 0
		total := input.SalaryStructure.Basic + input.SalaryStructure.Houserent + input.SalaryStructure.Dareness + input.SalaryStructure.Specialallow + input.SalaryStructure.Specialda + input.SalaryStructure.Arrear + input.SalaryStructure.Binder_wf + input.SalaryStructure.Incentive +
			input.SalaryStructure.Conveyance + input.SalaryStructure.Medical + input.SalaryStructure.Otherallow + input.SalaryStructure.Extraallow + input.SalaryStructure.Technical + input.SalaryStructure.Mobile + input.SalaryStructure.Pubsalary + input.SalaryStructure.Business + input.SalaryStructure.Charge +
			input.SalaryStructure.Eyeallow + input.SalaryStructure.Cosecretary + input.SalaryStructure.Carallow + input.SalaryStructure.Leaserent

		input.SalaryStructure.Grosssalary = total
	}
	if input.SalaryStructure.Consolited > 0 && input.SalaryStructure.Basic <= 0 {
		input.SalaryStructure.Basic = 0
		total := input.SalaryStructure.Consolited + input.SalaryStructure.Houserent + input.SalaryStructure.Dareness + input.SalaryStructure.Specialallow + input.SalaryStructure.Specialda + input.SalaryStructure.Arrear + input.SalaryStructure.Binder_wf + input.SalaryStructure.Incentive +
			input.SalaryStructure.Conveyance + input.SalaryStructure.Medical + input.SalaryStructure.Otherallow + input.SalaryStructure.Extraallow + input.SalaryStructure.Technical + input.SalaryStructure.Mobile + input.SalaryStructure.Pubsalary + input.SalaryStructure.Business + input.SalaryStructure.Charge +
			input.SalaryStructure.Eyeallow + input.SalaryStructure.Cosecretary + input.SalaryStructure.Carallow + input.SalaryStructure.Leaserent

		input.SalaryStructure.Grosssalary = total
	}
	input.SalaryStructure.Empcode = input.StaffInformation.Empcode
	input.SalaryStructure.Seniority_serial = input.StaffInformation.Senioriy_serial
	input.SalaryStructure.Refno = 1
	input.SalaryStructure.Incrementno = 0
	input.SalaryStructure.Incrementamount = 0
	input.StaffInformation.Active_salary = 1

	_ = tx.Raw("select coalesce ((max(id) + 1), 1) from payroll.staffinformation").First(&input.StaffInformation.Id)
	result1 := tx.Create(&input.StaffInformation) // u are se
	if result1.RowsAffected == 0 {
		tx.RollbackTo("savePoint1")
		res.Message = "Oops! Something went wrong. please try again later."
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusInternalServerError
		return res
	}

	_ = tx.Raw("select coalesce ((max(id) + 1), 1) from payroll.salarystructure").First(&input.SalaryStructure.Id)
	result := tx.Create(&input.SalaryStructure)
	if result.RowsAffected == 0 {
		tx.RollbackTo("savePoint1")
		res.Message = "Oops! Something went wrong. please try again later."
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusInternalServerError
		return res
	}
	var input2 dto.CustomEmpSaveDTO
	input2.Pempcode = input.StaffInformation.Empcode
	input2.Pempname = input.StaffInformation.Empname
	fmt.Println("test1", input2.Pempcode)
	fmt.Println("test2", input2.Pempname)

	_ = tx.Raw("select * from payroll.staffinformation where empcode = ?", input.SalaryStructure.Empcode).First(&latest.StaffInformation)
	_ = tx.Raw("select * from payroll.salarystructure where empcode = ?", input.SalaryStructure.Empcode).First(&latest.SalaryStructure)
	// Procedure called below
	_ = db1.Exec("CALL pfaccounts.addemploanmember(?, ?)", input2.Pempcode, input2.Pempname)
	output = latest

	tx.Commit()
	res.Message = "Successfully added"
	res.IsSuccess = true
	res.Payload = output
	res.StatusCode = http.StatusOK
	// res.Count = len(output)
	return res
}

func (e *EmployeeInformation) UpdateEmployeeInformation(input dto.EmployeeInformationUpdate) dto.ResponseDto {
	var res dto.ResponseDto
	var input1 dto.UpdateEmployeeInformation
	var latest dto.EmployeeInformation
	var op dto.EmployeeInformation

	db1 := util.CreateConnectionToPFAccountsSchemaUsingGorm()
	sqlDB1, _ := db1.DB()
	defer sqlDB1.Close()

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	tx := db.Begin()
	tx.SavePoint("savePoint1")
	var output dto.EmployeeInformation
	var output1 dto.EmployeeInformation_archive
	_ = tx.Raw("select * from payroll.staffinformation where empcode = ?", input.StaffInformation.Empcode).First(&output.StaffInformation)
	prevname := output.StaffInformation.Empname
	_ = tx.Raw("select * from payroll.salarystructure where empcode = ?", output.StaffInformation.Empcode).First(&output.SalaryStructure)
	//_ = tx.Raw("select * from payroll.salarystructure where id = ?", output.SalaryStructure.Id).First(&output.SalaryStructure)

	input1.EmployeeInformation_archive.StaffInformation_archive.Id = output.StaffInformation.Id
	input1.EmployeeInformation_archive.StaffInformation_archive.Empcode = output.StaffInformation.Empcode
	input1.EmployeeInformation_archive.StaffInformation_archive.Senioriy_serial = output.StaffInformation.Senioriy_serial
	input1.EmployeeInformation_archive.StaffInformation_archive.Empname = output.StaffInformation.Empname
	input1.EmployeeInformation_archive.StaffInformation_archive.Grade = output.StaffInformation.Grade
	input1.EmployeeInformation_archive.StaffInformation_archive.Accno = output.StaffInformation.Accno
	input1.EmployeeInformation_archive.StaffInformation_archive.Tin = output.StaffInformation.Tin
	input1.EmployeeInformation_archive.StaffInformation_archive.Joindate = output.StaffInformation.Joindate
	input1.EmployeeInformation_archive.StaffInformation_archive.Probationperiod = output.StaffInformation.Probationperiod
	input1.EmployeeInformation_archive.StaffInformation_archive.Confirmdate = output.StaffInformation.Confirmdate
	input1.EmployeeInformation_archive.StaffInformation_archive.Eduquali = output.StaffInformation.Eduquali
	input1.EmployeeInformation_archive.StaffInformation_archive.Mailadd = output.StaffInformation.Mailadd
	input1.EmployeeInformation_archive.StaffInformation_archive.Permanentadd = output.StaffInformation.Permanentadd
	input1.EmployeeInformation_archive.StaffInformation_archive.Maritalstatus = output.StaffInformation.Maritalstatus
	input1.EmployeeInformation_archive.StaffInformation_archive.Gender = output.StaffInformation.Gender
	input1.EmployeeInformation_archive.StaffInformation_archive.Religion = output.StaffInformation.Religion
	input1.EmployeeInformation_archive.StaffInformation_archive.Bloodgroup = output.StaffInformation.Bloodgroup
	input1.EmployeeInformation_archive.StaffInformation_archive.Active_salary = output.StaffInformation.Active_salary
	input1.EmployeeInformation_archive.StaffInformation_archive.Entry_user = output.StaffInformation.Entry_user

	dt := time.Now()

	input1.EmployeeInformation_archive.StaffInformation_archive.Changedate = dt.Format("2006-01-02 15:04:05")
	input1.EmployeeInformation_archive.StaffInformation_archive.Changeuserid = input.StaffInformation.Entry_user
	input1.EmployeeInformation_archive.StaffInformation_archive.Flag_ed_del = "Update"

	result := tx.Raw("select coalesce ((max(trackid) + 1), 1) from payroll.staffinformation_archive").First(&output1.StaffInformation_archive.Trackid)
	if result.Error != nil {
		tx.RollbackTo("savePoint1")
		res.Message = "Error in getting track id"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}

	input1.EmployeeInformation_archive.StaffInformation_archive.Trackid = output1.StaffInformation_archive.Trackid

	input1.EmployeeInformation_archive.SalaryStructure_archive.Id = output.SalaryStructure.Id
	input1.EmployeeInformation_archive.SalaryStructure_archive.Empcode = output.SalaryStructure.Empcode
	input1.EmployeeInformation_archive.SalaryStructure_archive.Refno = output.SalaryStructure.Refno
	input1.EmployeeInformation_archive.SalaryStructure_archive.Incrementno = output.SalaryStructure.Incrementno
	input1.EmployeeInformation_archive.SalaryStructure_archive.Incrementamount = output.SalaryStructure.Incrementamount
	input1.EmployeeInformation_archive.SalaryStructure_archive.Desigcode = output.SalaryStructure.Desigcode
	input1.EmployeeInformation_archive.SalaryStructure_archive.Branchcode = output.SalaryStructure.Branchcode
	input1.EmployeeInformation_archive.SalaryStructure_archive.Deptcode = output.SalaryStructure.Deptcode
	input1.EmployeeInformation_archive.SalaryStructure_archive.Bankid = output.SalaryStructure.Bankid
	input1.EmployeeInformation_archive.SalaryStructure_archive.Pfbankid = output.SalaryStructure.Pfbankid
	input1.EmployeeInformation_archive.SalaryStructure_archive.Consolited = output.SalaryStructure.Consolited
	input1.EmployeeInformation_archive.SalaryStructure_archive.Basic = output.SalaryStructure.Basic
	input1.EmployeeInformation_archive.SalaryStructure_archive.Houserent = output.SalaryStructure.Houserent
	input1.EmployeeInformation_archive.SalaryStructure_archive.Conveyance = output.SalaryStructure.Conveyance
	input1.EmployeeInformation_archive.SalaryStructure_archive.Medical = output.SalaryStructure.Medical
	input1.EmployeeInformation_archive.SalaryStructure_archive.Entertainment = output.SalaryStructure.Entertainment
	input1.EmployeeInformation_archive.SalaryStructure_archive.Housemaint = output.SalaryStructure.Housemaint
	input1.EmployeeInformation_archive.SalaryStructure_archive.Incometax = output.SalaryStructure.Incometax
	input1.EmployeeInformation_archive.SalaryStructure_archive.Bonusrate = output.SalaryStructure.Bonusrate
	input1.EmployeeInformation_archive.SalaryStructure_archive.Arrear = output.SalaryStructure.Arrear
	input1.EmployeeInformation_archive.SalaryStructure_archive.Cpf = output.SalaryStructure.Cpf
	input1.EmployeeInformation_archive.SalaryStructure_archive.Groupins = output.SalaryStructure.Groupins
	input1.EmployeeInformation_archive.SalaryStructure_archive.Cpfloan = output.SalaryStructure.Cpfloan
	input1.EmployeeInformation_archive.SalaryStructure_archive.Stamp = output.SalaryStructure.Stamp
	input1.EmployeeInformation_archive.SalaryStructure_archive.Pfund = output.SalaryStructure.Pfund
	input1.EmployeeInformation_archive.SalaryStructure_archive.Sal_scale = output.SalaryStructure.Sal_scale
	input1.EmployeeInformation_archive.SalaryStructure_archive.Seniority_serial = output.SalaryStructure.Seniority_serial
	input1.EmployeeInformation_archive.SalaryStructure_archive.Telephone = output.SalaryStructure.Telephone
	input1.EmployeeInformation_archive.SalaryStructure_archive.Incentive = output.SalaryStructure.Incentive
	input1.EmployeeInformation_archive.SalaryStructure_archive.Specialallow = output.SalaryStructure.Specialallow
	input1.EmployeeInformation_archive.SalaryStructure_archive.Overtime = output.SalaryStructure.Overtime
	input1.EmployeeInformation_archive.SalaryStructure_archive.Food = output.SalaryStructure.Food
	input1.EmployeeInformation_archive.SalaryStructure_archive.Salaryadv = output.SalaryStructure.Salaryadv
	input1.EmployeeInformation_archive.SalaryStructure_archive.Otherallow = output.SalaryStructure.Otherallow
	input1.EmployeeInformation_archive.SalaryStructure_archive.Otheradv = output.SalaryStructure.Otheradv
	input1.EmployeeInformation_archive.SalaryStructure_archive.Carallow = output.SalaryStructure.Carallow
	input1.EmployeeInformation_archive.SalaryStructure_archive.Specialallow1 = output.SalaryStructure.Specialallow1
	input1.EmployeeInformation_archive.SalaryStructure_archive.Binder_wf = output.SalaryStructure.Binder_wf
	input1.EmployeeInformation_archive.SalaryStructure_archive.Leaserent = output.SalaryStructure.Leaserent
	input1.EmployeeInformation_archive.SalaryStructure_archive.Dareness = output.SalaryStructure.Dareness
	input1.EmployeeInformation_archive.SalaryStructure_archive.Specialda = output.SalaryStructure.Specialda
	input1.EmployeeInformation_archive.SalaryStructure_archive.Extraallow = output.SalaryStructure.Extraallow
	input1.EmployeeInformation_archive.SalaryStructure_archive.Technical = output.SalaryStructure.Technical
	input1.EmployeeInformation_archive.SalaryStructure_archive.Mobile = output.SalaryStructure.Mobile
	input1.EmployeeInformation_archive.SalaryStructure_archive.Pubsalary = output.SalaryStructure.Pubsalary
	input1.EmployeeInformation_archive.SalaryStructure_archive.Business = output.SalaryStructure.Business
	input1.EmployeeInformation_archive.SalaryStructure_archive.Charge = output.SalaryStructure.Charge
	input1.EmployeeInformation_archive.SalaryStructure_archive.Eyeallow = output.SalaryStructure.Eyeallow
	input1.EmployeeInformation_archive.SalaryStructure_archive.Cosecretary = output.SalaryStructure.Cosecretary
	input1.EmployeeInformation_archive.SalaryStructure_archive.Grosssalary = output.SalaryStructure.Grosssalary

	input1.EmployeeInformation_archive.SalaryStructure_archive.Changedate = input1.EmployeeInformation_archive.StaffInformation_archive.Changedate
	input1.EmployeeInformation_archive.SalaryStructure_archive.Changeuserid = input1.EmployeeInformation_archive.StaffInformation_archive.Changeuserid
	input1.EmployeeInformation_archive.SalaryStructure_archive.Flag_ed_del = input1.EmployeeInformation_archive.StaffInformation_archive.Flag_ed_del
	input1.EmployeeInformation_archive.SalaryStructure_archive.Trackid = input1.EmployeeInformation_archive.StaffInformation_archive.Trackid

	result = tx.Create(&input1.EmployeeInformation_archive.StaffInformation_archive)
	if result.RowsAffected == 0 {
		tx.RollbackTo("savepoint")
		res.IsSuccess = false
		res.StatusCode = http.StatusInternalServerError
		res.Message = "Staff Information_archive insert failed"
		res.Payload = nil
		return res

	}
	result = tx.Create(&input1.EmployeeInformation_archive.SalaryStructure_archive)
	if result.RowsAffected == 0 {
		tx.RollbackTo("savepoint")
		res.IsSuccess = false
		res.StatusCode = http.StatusInternalServerError
		res.Message = "Salary Structure_archive insert failed"
		res.Payload = nil
		return res
	}

	input.StaffInformation.Id = output.StaffInformation.Id
	input.SalaryStructure.Id = output.StaffInformation.Id
	input.StaffInformation.Empcode = output.StaffInformation.Empcode
	if input.StaffInformation.Empcode == 0 {
		res.Message = "EMP Code is required"
		res.IsSuccess = false
		res.StatusCode = http.StatusBadRequest
		res.Payload = nil

		return res
	}
	if input.StaffInformation.Accno == "" {
		res.Message = "Account Number is required"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}
	if input.StaffInformation.Senioriy_serial <= 0 {
		res.Message = "Serial number is required"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}

	if input.StaffInformation.Empname == "" {
		res.Message = "Name is required"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}

	if input.StaffInformation.Joindate == "" {
		res.Message = "Joining date is required"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}

	// if input.StaffInformation.Mailadd == "" {
	// 	res.Message = "Mail address is required"
	// 	res.StatusCode = http.StatusBadRequest
	// 	res.IsSuccess = false
	// 	res.Payload = nil
	// 	return res
	// }
	// if input.StaffInformation.Permanentadd == "" {
	// 	res.Message = "Permanent address is required"
	// 	res.StatusCode = http.StatusBadRequest
	// 	res.IsSuccess = false
	// 	res.Payload = nil
	// 	return res
	// }
	if input.StaffInformation.Eduquali == "" {
		res.Message = "Education qualification is required"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}
	if input.StaffInformation.Gender == "" {
		res.Message = "Gender is required"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}
	if input.StaffInformation.Grade == "" {
		res.Message = "Grade is required"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}
	if input.StaffInformation.Maritalstatus == "" {
		res.Message = "Marital status is required"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}
	if input.StaffInformation.Religion == "" {
		res.Message = "Religion is required"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}
	if input.StaffInformation.Bloodgroup == "" {
		res.Message = "Salary is required"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}
	if input.StaffInformation.Active_salary < 0 && input.StaffInformation.Active_salary > 1 {
		res.Message = "Active status is either 0 or 1"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}

	// validation for salary structure
	if input.SalaryStructure.Bankid <= 0 {
		res.Message = "Bank id is required"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}
	_ = tx.Raw("select * from payroll.salarystructure where empcode = ?", input.StaffInformation.Empcode).First(output.SalaryStructure.Empcode)

	output.SalaryStructure.Basic = 0
	output.SalaryStructure.Consolited = 0
	if input.SalaryStructure.Basic > 0 && input.SalaryStructure.Consolited <= 0 {
		input.SalaryStructure.Consolited = 0
		total := input.SalaryStructure.Basic + input.SalaryStructure.Houserent + input.SalaryStructure.Dareness + input.SalaryStructure.Specialallow + input.SalaryStructure.Specialda + input.SalaryStructure.Arrear + input.SalaryStructure.Binder_wf + input.SalaryStructure.Incentive +
			input.SalaryStructure.Conveyance + input.SalaryStructure.Medical + input.SalaryStructure.Otherallow + input.SalaryStructure.Extraallow + input.SalaryStructure.Technical + input.SalaryStructure.Mobile + input.SalaryStructure.Pubsalary + input.SalaryStructure.Business + input.SalaryStructure.Charge +
			input.SalaryStructure.Eyeallow + input.SalaryStructure.Cosecretary + input.SalaryStructure.Carallow + input.SalaryStructure.Leaserent

		input.SalaryStructure.Grosssalary = total
		fmt.Println("basic", input.SalaryStructure.Basic)
		fmt.Println("gross salary for basic", input.SalaryStructure.Grosssalary)
	}

	if input.SalaryStructure.Consolited > 0 && input.SalaryStructure.Basic <= 0 {
		input.SalaryStructure.Basic = 0
		total := input.SalaryStructure.Consolited + input.SalaryStructure.Houserent + input.SalaryStructure.Dareness + input.SalaryStructure.Specialallow + input.SalaryStructure.Specialda + input.SalaryStructure.Arrear + input.SalaryStructure.Binder_wf + input.SalaryStructure.Incentive +
			input.SalaryStructure.Conveyance + input.SalaryStructure.Medical + input.SalaryStructure.Otherallow + input.SalaryStructure.Extraallow + input.SalaryStructure.Technical + input.SalaryStructure.Mobile + input.SalaryStructure.Pubsalary + input.SalaryStructure.Business + input.SalaryStructure.Charge +
			input.SalaryStructure.Eyeallow + input.SalaryStructure.Cosecretary + input.SalaryStructure.Carallow + input.SalaryStructure.Leaserent

		input.SalaryStructure.Grosssalary = total
		fmt.Println("gross salary for consolited", input.SalaryStructure.Consolited)
		fmt.Println("gross salary for consolited", input.SalaryStructure.Grosssalary)
	}
	input.SalaryStructure.Empcode = input.StaffInformation.Empcode
	input.SalaryStructure.Seniority_serial = input.StaffInformation.Senioriy_serial

	//_ = tx.Raw("select coalesce ((max(id) + 1), 1) from payroll.staffinformation").First(&input.StaffInformation.Id)
	result1 := tx.Model(&input.StaffInformation).Where("empcode = ?", input.StaffInformation.Empcode).Updates(input.StaffInformation)
	if result1.RowsAffected == 0 {
		tx.RollbackTo("savePoint1")
		res.Message = "Staff information update failed!"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusInternalServerError
		return res
	}

	//_ = tx.Raw("select coalesce ((max(id) + 1), 1) from payroll.salarystructure").First(&input.SalaryStructure.Id)

	result3 := tx.Model(&input.SalaryStructure).Where("empcode = ?", input.StaffInformation.Empcode).Updates(input.SalaryStructure)
	if result3.RowsAffected == 0 {
		tx.RollbackTo("savePoint1")
		res.Message = "Salary structure update failed!"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusInternalServerError
		return res
	}
	fmt.Println("prev emp name ->", prevname)
	var input2 dto.CustomEmpSaveDTO
	input2.Pempcode = input.StaffInformation.Empcode
	input2.Pempname = input.StaffInformation.Empname
	fmt.Println("test1", input2.Pempcode)
	fmt.Println("test2", input2.Pempname)

	_ = tx.Raw("select * from payroll.staffinformation where empcode = ?", input.StaffInformation.Empcode).First(&latest.StaffInformation)
	_ = tx.Raw("select * from payroll.salarystructure where empcode = ?", input.StaffInformation.Empcode).First(&latest.SalaryStructure)
	// _ = tx.Raw("select * from payroll.staffinformation_archive where empcode = ?", input.StaffInformation.Empcode).First(&latest.EmployeeInformation_archive.StaffInformation_archive)
	// _ = tx.Raw("select * from payroll.salarystructure_archive where empcode = ?", input.StaffInformation.Empcode).First(&latest.EmployeeInformation_archive.SalaryStructure_archive)

	_ = db1.Exec("CALL pfaccounts.updateempname_loanmemberacc(?, ?, ?)", input2.Pempcode, input2.Pempname, prevname)
	op = latest

	tx.Commit()
	res.Message = "Successfully updated"
	res.IsSuccess = true
	res.Payload = op
	res.StatusCode = http.StatusOK

	return res
}

func (e *EmployeeInformation) DeleteEmployeeInformation(input dto.EmployeeInformationUpdate) dto.ResponseDto {
	var res dto.ResponseDto
	var input1 dto.UpdateEmployeeInformation

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	tx := db.Begin()
	tx.SavePoint("savePoint1")
	var output dto.EmployeeInformation
	var output1 dto.EmployeeInformation_archive

	result3 := tx.Raw("select * from payroll.staffinformation where empcode = ?", input.StaffInformation.Empcode).First(&output.StaffInformation)
	if result3.RowsAffected == 0 {
		tx.RollbackTo("savePoint1")
		res.Message = "Employee code doesn't exist in Staffinformation."
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusInternalServerError
		return res
	}

	result4 := tx.Raw("select * from payroll.salarystructure where empcode = ?", input.StaffInformation.Empcode).First(&output.SalaryStructure)
	if result4.RowsAffected != 0 {
		tx.RollbackTo("savePoint1")
		res.Message = "Data already exist in other module, you can't delete this employee."
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusInternalServerError
		return res
	}

	input1.EmployeeInformation_archive.StaffInformation_archive.Staffinformation = output.StaffInformation
	// fmt.Println("hi heelo halo wwiwiwiiwiwwiwiwiiwiwiwiwiwiwiw")
	// fmt.Println(input1.EmployeeInformation_archive.StaffInformation_archive)

	// input1.EmployeeInformation_archive.StaffInformation_archive.Id = output.StaffInformation.Id
	// input1.EmployeeInformation_archive.StaffInformation_archive.Empcode = output.StaffInformation.Empcode
	// input1.EmployeeInformation_archive.StaffInformation_archive.Serial = output.StaffInformation.Serial
	// input1.EmployeeInformation_archive.StaffInformation_archive.Name = output.StaffInformation.Name
	// input1.EmployeeInformation_archive.StaffInformation_archive.Grade = output.StaffInformation.Grade
	// input1.EmployeeInformation_archive.StaffInformation_archive.Accno = output.StaffInformation.Accno
	// input1.EmployeeInformation_archive.StaffInformation_archive.Tin = output.StaffInformation.Tin
	// input1.EmployeeInformation_archive.StaffInformation_archive.Joindate = output.StaffInformation.Joindate
	// input1.EmployeeInformation_archive.StaffInformation_archive.Probationperiod = output.StaffInformation.Probationperiod
	// input1.EmployeeInformation_archive.StaffInformation_archive.Confirmdate = output.StaffInformation.Confirmdate
	// input1.EmployeeInformation_archive.StaffInformation_archive.Eduquali = output.StaffInformation.Eduquali
	// input1.EmployeeInformation_archive.StaffInformation_archive.Mailadd = output.StaffInformation.Mailadd
	// input1.EmployeeInformation_archive.StaffInformation_archive.Permanentadd = output.StaffInformation.Permanentadd
	// input1.EmployeeInformation_archive.StaffInformation_archive.Maritalstatus = output.StaffInformation.Maritalstatus
	// input1.EmployeeInformation_archive.StaffInformation_archive.Gender = output.StaffInformation.Gender
	// input1.EmployeeInformation_archive.StaffInformation_archive.Religion = output.StaffInformation.Religion
	// input1.EmployeeInformation_archive.StaffInformation_archive.Bloodgroup = output.StaffInformation.Bloodgroup
	// input1.EmployeeInformation_archive.StaffInformation_archive.Active = output.StaffInformation.Active

	dt := time.Now()

	input1.EmployeeInformation_archive.StaffInformation_archive.Changedate = dt.Format("2006-01-02 15:04:05")
	input1.EmployeeInformation_archive.StaffInformation_archive.Changeuserid = strconv.Itoa(output.StaffInformation.Id)
	input1.EmployeeInformation_archive.StaffInformation_archive.Flag_ed_del = "Delete"

	result := tx.Raw("select coalesce ((max(trackid) + 1), 1) from payroll.staffinformation_archive").First(&output1.StaffInformation_archive.Trackid)
	if result.Error != nil {
		tx.RollbackTo("savePoint1")
		res.Message = "Error in getting track id"
		res.StatusCode = http.StatusBadRequest
		res.IsSuccess = false
		res.Payload = nil
		return res
	}

	input1.EmployeeInformation_archive.StaffInformation_archive.Trackid = output1.StaffInformation_archive.Trackid

	// input1.EmployeeInformation_archive.SalaryStructure_archive.Salarystructure = output.SalaryStructure

	// input1.EmployeeInformation_archive.SalaryStructure_archive.Changedate = input1.EmployeeInformation_archive.StaffInformation_archive.Changedate
	// input1.EmployeeInformation_archive.SalaryStructure_archive.Changeuserid = input1.EmployeeInformation_archive.StaffInformation_archive.Changeuserid
	// input1.EmployeeInformation_archive.SalaryStructure_archive.Flag_ed_del = input1.EmployeeInformation_archive.StaffInformation_archive.Flag_ed_del
	// input1.EmployeeInformation_archive.SalaryStructure_archive.Trackid = input1.EmployeeInformation_archive.StaffInformation_archive.Trackid

	result = tx.Create(&input1.EmployeeInformation_archive.StaffInformation_archive)
	if result.RowsAffected == 0 {
		tx.RollbackTo("savepoint")
		res.IsSuccess = false
		res.Message = "Staff Information_archive insert failed"
		res.Payload = nil
		return res

	}
	// result = tx.Create(&input1.EmployeeInformation_archive.SalaryStructure_archive)
	// if result.RowsAffected == 0 {
	// 	tx.RollbackTo("savepoint")
	// 	res.IsSuccess = false
	// 	res.Message = "SalaryStructure_archive insert failed"
	// 	res.Payload = nil
	// 	return res
	// }

	result = tx.Model(&input1.EmployeeInformation.StaffInformation).Where("empcode = ?", input.StaffInformation.Empcode).Delete(&input1.EmployeeInformation.StaffInformation)
	if result.RowsAffected == 0 {
		tx.RollbackTo("savepoint")
		res.IsSuccess = false
		res.Message = "Staff Information code delete failed"
		res.Payload = nil
		return res
	}
	// _ = tx.Raw("select * from payroll.staffinformation where empcode = ?", input.StaffInformation.Empcode).First(&output.StaffInformation)
	// result = tx.Model(&input1.EmployeeInformation.SalaryStructure).Where("id = ?", output.StaffInformation.Id).Delete(&input1.EmployeeInformation.SalaryStructure)
	// if result.RowsAffected == 0 {
	// 	tx.RollbackTo("savepoint")
	// 	res.IsSuccess = false
	// 	res.Message = "SalaryStructure delete failed"
	// 	res.Payload = nil
	// 	return res
	// }

	tx.Commit()
	res.Message = "Successfully Deleted"
	res.IsSuccess = true
	res.Payload = nil
	res.StatusCode = http.StatusAccepted
	return res
}

func (e *EmployeeInformation) GetAnEmployeeInformation(input dto.EmployeeInformationUpdate) dto.ResponseDto {
	var res dto.ResponseDto

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	tx := db.Begin()
	tx.SavePoint("savePoint1")
	var output dto.EmployeeInformationUpdate

	fmt.Println("input", input.StaffInformation.Empcode)

	result := tx.Raw("select * from payroll.staffinformation where empcode = ?", input.StaffInformation.Empcode).First(&output.StaffInformation)
	if result.RowsAffected == 0 {
		tx.RollbackTo("savePoint1")
		res.IsSuccess = false
		res.StatusCode = http.StatusNotFound
		res.Message = "Staff Information code not found"
		res.Payload = nil
		return res
	}
	fmt.Println("output", output.StaffInformation.Empcode)

	//_ = tx.Raw("select * from payroll.staffinformation where empcode = ?", input.StaffInformation.Empcode).First(&output.StaffInformation)
	result1 := tx.Raw("select * from payroll.salarystructure where empcode = ?", output.StaffInformation.Empcode).First(&output.SalaryStructure)
	if result1.RowsAffected == 0 {
		tx.RollbackTo("savePoint1")
		res.IsSuccess = false
		res.StatusCode = http.StatusNotFound
		res.Message = "Salary not found"
		res.Payload = nil
		return res
	}

	tx.Commit()
	res.Message = "Get Successful"
	res.IsSuccess = true
	res.Payload = output
	res.StatusCode = http.StatusOK

	return res
}

func (e *EmployeeInformation) GetMaxEmpcodeEmployeeInformation(input dto.EmployeeInformationUpdate) dto.ResponseDto {
	var res dto.ResponseDto

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	tx := db.Begin()
	tx.SavePoint("savePoint1")
	var output dto.MaxEmployeeInformationDTO

	result := tx.Raw("select max(empcode+1) from payroll.staffinformation").First(&output.Empcode)
	if result.RowsAffected == 0 {
		tx.RollbackTo("savePoint1")
		res.IsSuccess = false
		res.StatusCode = http.StatusNotFound
		res.Message = "Internal Server Error"
		res.Payload = nil
		return res
	} else if result.Error != nil {
		output.Empcode = 1
	}
	fmt.Println("output", output.Empcode)

	tx.Commit()
	res.Message = "Max Empcode for new Employee Insert"
	res.IsSuccess = true
	res.Payload = output
	res.StatusCode = http.StatusOK

	return res
}

func (e *EmployeeInformation) GetEmployeeSerialList() dto.ResponseDto {
	var res dto.ResponseDto

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	// var deptcode model.Hrdept
	tx := db.Begin()
	tx.SavePoint("savePoint1")
	var output []dto.EmployeeInformationDTO

	result := tx.Raw(`
	select a.senioriy_serial, a.empname ,c.deptname, d.designame
	FROM  payroll.staffinformation a
	INNER JOIN payroll.salarystructure b
	ON a.empcode=b.empcode
	INNER JOIN payroll.hrdept c
	ON c.deptcode =b.deptcode
	INNER JOIN payroll.designation d
	ON d.desigcode =b.desigcode
	order by c.deptcode , a.senioriy_serial `).Find(&output)
	// result := tx.Raw("select senioriy_serial from payroll.staffinformation order by senioriy_serial asc").Find(&output)
	if result.RowsAffected == 0 {
		tx.RollbackTo("savePoint1")
		res.IsSuccess = false
		res.StatusCode = http.StatusNotFound
		res.Message = "Serial list not found"
		res.Payload = nil
		return res
	}
	fmt.Println("output", output)

	tx.Commit()
	res.Message = "Get All Serial Successfully"
	res.IsSuccess = true
	res.Payload = output
	res.Count = len(output)
	res.StatusCode = http.StatusOK

	return res
}

// type Staffinformations interface {
// 	Staffinformation model.Staffinformation
// 	SalaryStructure  model.Salarystructure
// 	Designation      model.Designation
// 	Department       model.Hrdept

// }
func (e *EmployeeInformation) GetAllStaffinformation() dto.ResponseDto {
	var res dto.ResponseDto
	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	//var output1 []dto.Staffinformations
	var salarystructureop []model.Salarystructure
	tx := db.Begin()
	tx.SavePoint("savePoint1")

	// result := tx.Raw(` select a.empcode, a.empname ,c.deptname, d.designame, a.grade ,a.senioriy_serial, a.bloodgroup  
	// FROM  payroll.staffinformation a
	// INNER JOIN payroll.salarystructure b
	// ON a.empcode=b.empcode
	// INNER JOIN payroll.hrdept c
	// ON c.deptcode =b.deptcode
	// INNER JOIN payroll.designation d
	// ON d.desigcode =b.desigcode        
	// order by a.senioriy_serial `).Find(&output1)
	// if result.RowsAffected == 0 {
	// 	tx.RollbackTo("savePoint1")
	// 	res.IsSuccess = false
	// 	res.StatusCode = http.StatusNotFound
	// 	res.Message = "Staff Information code not found"
	// 	res.Payload = nil
	// 	return res
	// }

	result1 := tx.Raw(`select e.* from payroll.salarystructure e, (select empcode , max(refno) refno from payroll.salarystructure group by empcode) f 
	where e.refno = f.refno and e.empcode =f.empcode`).Find(&salarystructureop)
	if result1.RowsAffected == 0 {
		tx.RollbackTo("savePoint1")
		res.IsSuccess = false
		res.StatusCode = http.StatusNotFound
		res.Message = "Salary Structure not found"
		res.Payload = nil
		return res
	}

	tx.Commit()
	res.Message = "Get All Employee Successfully"
	res.IsSuccess = true
	res.Payload = salarystructureop
	res.StatusCode = http.StatusOK

	return res
}

func (e *EmployeeInformation) GetAnDepositeBank(input model.Acchead) dto.ResponseDto {
	var res dto.ResponseDto

	db := util.CreateConnectionToAccountsSchemaUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	tx := db.Begin()
	tx.SavePoint("savePoint1")
	var output dto.DepositeBankDTO

	result := tx.Raw("select accid,name from accounts.acchead where lower(name) =?", strings.ToLower(input.Name)).First(&output)
	if result.RowsAffected == 0 {
		tx.RollbackTo("savePoint1")
		res.IsSuccess = false
		res.StatusCode = http.StatusNotFound
		res.Message = "Bank not found"
		res.Payload = nil
		return res
	}
	fmt.Println("output", output)

	tx.Commit()
	res.Message = "Deposite bank fetch Successfully"
	res.IsSuccess = true
	res.Payload = output
	res.StatusCode = http.StatusOK

	return res
}

func (e *EmployeeInformation) GetAllDepositeBank() dto.ResponseDto {
	var res dto.ResponseDto

	db := util.CreateConnectionToAccountsSchemaUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	tx := db.Begin()
	tx.SavePoint("savePoint1")
	var output []dto.DepositeBankDTO

	result := tx.Raw("select a.accid, a.name from accounts.acchead a where a.parent in(10101100,10101150) and a.lr = 'L' order by name").Find(&output)
	if result.RowsAffected == 0 {
		tx.RollbackTo("savePoint1")
		res.IsSuccess = false
		res.StatusCode = http.StatusNotFound
		res.Message = "Bank not found"
		res.Payload = nil
		return res
	}

	fmt.Println("output", output)

	tx.Commit()
	res.Message = "All Deposite bank fetch Successfully"
	res.IsSuccess = true
	res.Count = len(output)
	res.Payload = output
	res.StatusCode = http.StatusOK

	return res
}

func (e *EmployeeInformation) GetBankNameToUpdate(input model.Salarystructure) dto.ResponseDto {
	var res dto.ResponseDto

	db := util.CreateConnectionToAccountsSchemaUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	tx := db.Begin()
	tx.SavePoint("savePoint1")
	var output dto.DepositeBankDTO

	result := tx.Raw("select a.accid, a.name from accounts.acchead a where a.accid =?", input.Bankid).First(&output)
	if result.RowsAffected == 0 {
		tx.RollbackTo("savePoint1")
		res.IsSuccess = false
		res.StatusCode = http.StatusNotFound
		res.Message = "Bank info is not found"
		res.Payload = nil
		return res
	}
	fmt.Println("output", output)

	tx.Commit()
	res.Message = "Bank info Successfully found"
	res.IsSuccess = true
	res.Payload = output
	res.StatusCode = http.StatusOK

	return res
}
