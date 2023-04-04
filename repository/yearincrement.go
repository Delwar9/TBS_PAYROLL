package repository

import (
	"fmt"
	"math"

	"github.com/tools/payroll/dto"
	"github.com/tools/payroll/model"
	"github.com/tools/payroll/util"
)

type YearIncrementRepository struct{}

func (yearincrement *YearIncrementRepository) GetAnEmployeeInfo(anemp model.Staffinformation) dto.ResponseDto {
	var op dto.ResponseDto

	var output1 dto.Employeedto
	if anemp.Empcode == 0 {
		op.IsSuccess = false
		op.Message = "Valid Empcode is required!"
		op.StatusCode = 404
		op.Payload = nil
		return op
	}

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	tx := db.Begin()
	tx.SavePoint("savepoint")

	result := tx.Raw(`SELECT a.id,b.refno,a.empcode,a.empname,b.desigcode,b.sal_scale ,	
	a.accno ,a.joindate ,a.probationperiod ,a.confirmdate ,	b.deptcode ,a.active_salary,
	b.incrementno,b.incrementamount ,b.consolited ,b.basic ,b.houserent ,b.dareness ,b.pfund ,b.specialallow ,b.carallow ,b.cpf ,
	b.incometax ,b.specialda ,b.incentive ,b.conveyance ,b.medical ,b.otherallow ,b.extraallow ,b.technical ,b.mobile ,b.pubsalary ,
	b.business ,b.charge ,b.eyeallow ,b.cosecretary,b.grosssalary  FROM payroll.staffinformation a FULL OUTER JOIN payroll.salarystructure b 
	ON a.empcode  = b.empcode where a.empcode =? and b.empcode =? and b.refno in (select max(refno) as refno from payroll.salarystructure)`, anemp.Empcode, anemp.Empcode).First(&output1)
	if result.RowsAffected == 0 {
		tx.RollbackTo("savepoint")
		op.IsSuccess = false
		op.Message = "no Yearly Increments found!"
		op.StatusCode = 404
		op.Payload = nil
	}

	op.IsSuccess = true
	op.Message = "Yearly Increments found!"
	op.StatusCode = 200
	op.Payload = output1

	return op
}

func (yearincrement *YearIncrementRepository) AddYearIncrement(yearIncrement model.Yearincrement) dto.ResponseDto {
	var op dto.ResponseDto

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	// var op1 model.Yearincrement
	// var op2 model.Yearincrement
	tx := db.Begin()
	tx.SavePoint("savepoint")

	// codes

	var input model.Salarystructure
	fmt.Println("initial desigcode: ", yearIncrement.Desigcode)
	fmt.Println("initial deptcode: ", yearIncrement.Deptcode)
	fmt.Println("initial incrementno: ", yearIncrement.Incrementno)
	fmt.Println("initial incrementamount: ", yearIncrement.Incrementamount)
	fmt.Println("initial empcode:", yearIncrement.Empcode)

	if yearIncrement.Basic != 0 {
		yearIncrement.Consolited = 0
		inc_amt := yearIncrement.Incrementamount * float64(yearIncrement.Incrementno)
		fmt.Println("increment: ", inc_amt)
		fmt.Println("Pre basic: ", yearIncrement.Basic)
		yearIncrement.Basic = yearIncrement.Basic + inc_amt
		fmt.Println("Post basic: ", yearIncrement.Basic)

		total := yearIncrement.Basic + yearIncrement.Houserent + yearIncrement.Dareness + yearIncrement.Specialallow + yearIncrement.Specialda +
			yearIncrement.Arrear + yearIncrement.Binder_wf + yearIncrement.Incentive + yearIncrement.Conveyance + yearIncrement.Medical +
			yearIncrement.Otherallow + yearIncrement.Extraallow + yearIncrement.Technical + yearIncrement.Mobile + yearIncrement.Pubsalary +
			yearIncrement.Business + yearIncrement.Charge + yearIncrement.Eyeallow + yearIncrement.Cosecretary + yearIncrement.Carallow +
			yearIncrement.Leaserent
		yearIncrement.Grosssalary = total
		yearIncrement.Pfund = math.Round(yearIncrement.Basic * 0.1)
		yearIncrement.Cpf = math.Round(yearIncrement.Pfund * 2)

		fmt.Println("gross: ", yearIncrement.Grosssalary)

		input.Empcode = yearIncrement.Empcode
		input.Desigcode = yearIncrement.Desigcode
		input.Deptcode = yearIncrement.Deptcode

		input.Incrementno = yearIncrement.Incrementno
		input.Incrementamount = yearIncrement.Incrementamount
		input.Branchcode = yearIncrement.Branchcode
		input.Bankid = yearIncrement.Bankid
		input.Pfbankid = yearIncrement.Pfbankid
		input.Groupins = yearIncrement.Groupins
		input.Cpfloan = yearIncrement.Cpfloan
		input.Stamp = yearIncrement.Stamp
		input.Seniority_serial = yearIncrement.Seniority_serial
		input.Sal_scale = yearIncrement.Sal_scale
		input.Overtime = yearIncrement.Overtime
		input.Food = yearIncrement.Food
		input.Salaryadv = yearIncrement.Salaryadv
		input.Otheradv = yearIncrement.Otheradv
		input.Specialallow1 = yearIncrement.Specialallow1
		input.Seniority_serial = yearIncrement.Seniority_serial
		input.Telephone = yearIncrement.Telephone

		input.Basic = yearIncrement.Basic
		input.Consolited = yearIncrement.Consolited
		input.Houserent = yearIncrement.Houserent
		input.Entertainment = yearIncrement.Entertainment
		input.Housemaint = yearIncrement.Housemaint
		input.Dareness = yearIncrement.Dareness
		input.Specialallow = yearIncrement.Specialallow
		input.Specialda = yearIncrement.Specialda
		input.Arrear = yearIncrement.Arrear
		input.Binder_wf = yearIncrement.Binder_wf
		input.Incentive = yearIncrement.Incentive
		input.Conveyance = yearIncrement.Conveyance
		input.Incometax = yearIncrement.Incometax
		input.Bonusrate = yearIncrement.Bonusrate
		input.Medical = yearIncrement.Medical
		input.Otherallow = yearIncrement.Otherallow
		input.Extraallow = yearIncrement.Extraallow
		input.Technical = yearIncrement.Technical
		input.Mobile = yearIncrement.Mobile
		input.Pubsalary = yearIncrement.Pubsalary
		input.Business = yearIncrement.Business
		input.Charge = yearIncrement.Charge
		input.Eyeallow = yearIncrement.Eyeallow
		input.Cosecretary = yearIncrement.Cosecretary
		input.Carallow = yearIncrement.Carallow
		input.Leaserent = yearIncrement.Leaserent
		input.Grosssalary = yearIncrement.Grosssalary
		input.Pfund = yearIncrement.Pfund
		input.Cpf = yearIncrement.Cpf

		fmt.Println("input gross: ", input.Grosssalary)
		fmt.Println("input pfund: ", input.Pfund)
		fmt.Println("input cpf: ", input.Cpf)

		_ = tx.Raw("select coalesce ((max(id) + 1), 1) from payroll.salarystructure").First(&input.Id)
		_ = tx.Raw("select coalesce ((max(refno) + 1), 1) from payroll.salarystructure").First(&input.Refno)
		res := tx.Create(&input)
		if res.RowsAffected == 0 {
			tx.RollbackTo("savepoint")
			op.IsSuccess = false
			op.Message = "Salary Increment in Salary Structure not added!"
			op.StatusCode = 404
			op.Payload = nil
		}
		yearIncrement.Refno = input.Refno
		yearIncrement.Id = input.Id
		yearIncrement.Empcode = input.Empcode
		// _ = tx.Raw("select coalesce ((max(id) + 1), 1) from payroll.yearincrement").First(&yearIncrement.Id)
		// _ = tx.Raw("select coalesce ((max(refno) + 1), 1) from payroll.yearincrement").First(&yearIncrement.Refno)
		fmt.Println("refno: ", yearIncrement.Refno)
		fmt.Println("id: ", yearIncrement.Id)
		result := tx.Create(&yearIncrement)
		if result.RowsAffected == 0 {
			tx.RollbackTo("savepoint")
			op.IsSuccess = false
			op.Message = "Yearly Increment not added!"
			op.StatusCode = 404
			op.Payload = nil
		}

	}
	if yearIncrement.Consolited != 0 {
		yearIncrement.Basic = 0
		inc_amt := yearIncrement.Incrementamount * float64(yearIncrement.Incrementno)
		fmt.Println("increment: ", inc_amt)
		fmt.Println("Pre Consolidate: ", yearIncrement.Consolited)
		yearIncrement.Consolited = yearIncrement.Consolited + inc_amt
		fmt.Println("Post Consolidate: ", yearIncrement.Consolited)

		total := yearIncrement.Consolited + yearIncrement.Houserent + yearIncrement.Dareness + yearIncrement.Specialallow + yearIncrement.Specialda +
			yearIncrement.Arrear + yearIncrement.Binder_wf + yearIncrement.Incentive + yearIncrement.Conveyance + yearIncrement.Medical +
			yearIncrement.Otherallow + yearIncrement.Extraallow + yearIncrement.Technical + yearIncrement.Mobile + yearIncrement.Pubsalary +
			yearIncrement.Business + yearIncrement.Charge + yearIncrement.Eyeallow + yearIncrement.Cosecretary + yearIncrement.Carallow +
			yearIncrement.Leaserent
		yearIncrement.Grosssalary = total
		yearIncrement.Pfund = math.Round(yearIncrement.Consolited * 0.1)
		yearIncrement.Cpf = math.Round(yearIncrement.Pfund * 2)

		fmt.Println("gross: ", yearIncrement.Grosssalary)

		input.Empcode = yearIncrement.Empcode
		input.Desigcode = yearIncrement.Desigcode
		input.Deptcode = yearIncrement.Deptcode

		input.Incrementno = yearIncrement.Incrementno
		input.Incrementamount = yearIncrement.Incrementamount
		input.Branchcode = yearIncrement.Branchcode
		input.Bankid = yearIncrement.Bankid
		input.Pfbankid = yearIncrement.Pfbankid
		input.Groupins = yearIncrement.Groupins
		input.Cpfloan = yearIncrement.Cpfloan
		input.Stamp = yearIncrement.Stamp
		input.Seniority_serial = yearIncrement.Seniority_serial
		input.Sal_scale = yearIncrement.Sal_scale
		input.Overtime = yearIncrement.Overtime
		input.Food = yearIncrement.Food
		input.Salaryadv = yearIncrement.Salaryadv
		input.Otheradv = yearIncrement.Otheradv
		input.Specialallow1 = yearIncrement.Specialallow1
		input.Seniority_serial = yearIncrement.Seniority_serial
		input.Telephone = yearIncrement.Telephone

		input.Basic = yearIncrement.Basic
		input.Houserent = yearIncrement.Houserent
		input.Entertainment = yearIncrement.Entertainment
		input.Housemaint = yearIncrement.Housemaint
		input.Dareness = yearIncrement.Dareness
		input.Specialallow = yearIncrement.Specialallow
		input.Specialda = yearIncrement.Specialda
		input.Arrear = yearIncrement.Arrear
		input.Binder_wf = yearIncrement.Binder_wf
		input.Incentive = yearIncrement.Incentive
		input.Conveyance = yearIncrement.Conveyance
		input.Incometax = yearIncrement.Incometax
		input.Bonusrate = yearIncrement.Bonusrate
		input.Medical = yearIncrement.Medical
		input.Otherallow = yearIncrement.Otherallow
		input.Extraallow = yearIncrement.Extraallow
		input.Technical = yearIncrement.Technical
		input.Mobile = yearIncrement.Mobile
		input.Pubsalary = yearIncrement.Pubsalary
		input.Business = yearIncrement.Business
		input.Charge = yearIncrement.Charge
		input.Eyeallow = yearIncrement.Eyeallow
		input.Cosecretary = yearIncrement.Cosecretary
		input.Carallow = yearIncrement.Carallow
		input.Leaserent = yearIncrement.Leaserent
		input.Grosssalary = yearIncrement.Grosssalary
		input.Pfund = yearIncrement.Pfund
		input.Cpf = yearIncrement.Cpf

		fmt.Println("input gross: ", input.Grosssalary)
		fmt.Println("input pfund: ", input.Pfund)
		fmt.Println("input cpf: ", input.Cpf)

		_ = tx.Raw("select coalesce ((max(id) + 1), 1) from payroll.salarystructure").First(&input.Id)
		_ = tx.Raw("select coalesce ((max(refno) + 1), 1) from payroll.salarystructure").First(&input.Refno)
		res := tx.Create(&input)
		if res.RowsAffected == 0 {
			tx.RollbackTo("savepoint")
			op.IsSuccess = false
			op.Message = "Salary Increment in Salary Structure not added!"
			op.StatusCode = 404
			op.Payload = nil
		}
		yearIncrement.Refno = input.Refno
		yearIncrement.Id = input.Id
		yearIncrement.Empcode = input.Empcode

		// _ = tx.Raw("select coalesce ((max(id) + 1), 1) from payroll.yearincrement").First(&yearIncrement.Id)
		// _ = tx.Raw("select coalesce ((max(refno) + 1), 1) from payroll.yearincrement").First(&yearIncrement.Refno)
		fmt.Println("refno: ", yearIncrement.Refno)
		fmt.Println("id: ", yearIncrement.Id)
		fmt.Println("Desicode increment: ", yearIncrement.Desigcode)
		fmt.Println("Deptcode increment: ", yearIncrement.Deptcode)
		fmt.Println("Desicode salary: ", input.Desigcode)

		fmt.Println("Deptcode salary: ", input.Deptcode)
		result := tx.Create(&yearIncrement)
		if result.RowsAffected == 0 {
			tx.RollbackTo("savepoint")
			op.IsSuccess = false
			op.Message = "Yearly Increment not added!"
			op.StatusCode = 404
			op.Payload = nil
		}
	}

	tx.Commit()
	op.IsSuccess = true
	op.Message = "Yearly Increment added!"
	op.StatusCode = 200
	op.Payload = yearIncrement

	return op
}
