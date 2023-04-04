package repository

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/tools/payroll/dto"
	"github.com/tools/payroll/model"
	"github.com/tools/payroll/util"
)

type HrdeptRepository struct{}

func (hrdeptrepo *HrdeptRepository) GetAllHrDept() dto.ResponseDto {
	var output dto.ResponseDto
	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var ab []model.Hrdept
	result := db.Order("deptcode").Find(&ab)
	if result.RowsAffected == 0 {
		output.Message = "No Department info found"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusNotFound
		return output
	}
	type tempOutPut struct {
		Output      []model.Hrdept `json:"output"`
		OutputCount int            `json:"outputCount"`
	}
	var tOutput tempOutPut
	tOutput.Output = ab
	tOutput.OutputCount = len(ab)
	output.Count = len(ab)
	output.Message = "List of departments"
	output.IsSuccess = true
	output.Payload = tOutput
	output.StatusCode = http.StatusOK

	return output
}

func (hrdeptrepo *HrdeptRepository) GetAllHrDept2(x dto.Abcd) dto.ResponseDto {
	var output dto.ResponseDto
	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var ab []model.Hrdept
	result := db.Order("deptcode").Offset(x.Offset).Limit(x.Limit).Find(&ab)
	if result.RowsAffected == 0 {
		output.Message = "No Department info found"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusNotFound
		return output
	}
	type tempOutPut struct {
		Output      []model.Hrdept `json:"output"`
		OutputCount int            `json:"outputCount"`
	}
	var tOutput tempOutPut
	tOutput.Output = ab
	tOutput.OutputCount = len(ab)
	_ = db.Raw("SELECT COUNT (*) FROM payroll.hrdept").First(&output.Count)
	// output.Count = len(ab)
	output.Message = "List of departments"
	output.IsSuccess = true
	output.Payload = tOutput
	output.StatusCode = http.StatusOK

	return output
}

func (hrdeptrepo *HrdeptRepository) GetHrDeptByDeptCode(c model.Hrdept) dto.ResponseDto {
	var output dto.ResponseDto
	if c.Deptcode <= 0 {
		output.Message = "Invalid code"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusBadRequest
		return output

	}

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	result := db.Raw("select * from payroll.hrdept where deptcode = ?", c.Deptcode).First(&c)
	if result.RowsAffected == 0 {
		output.Message = "No Dept Employee info found"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusNotFound
		return output
	}
	type tempOutput struct {
		Output model.Hrdept `json:"output"`
	}
	var tOutput tempOutput
	tOutput.Output = c
	output.Message = "Dept Employee info details found for given criteria"
	output.IsSuccess = true
	output.Payload = tOutput
	output.StatusCode = http.StatusOK

	return output
}

func (hrdeptrepo *HrdeptRepository) AddHrDept(c model.Hrdept) dto.ResponseDto {
	var output dto.ResponseDto
	// if c.Code <= 0 {
	// 	output.Message = "Invalid code"
	// 	output.IsSuccess = false
	// 	output.Payload = nil
	// 	output.StatusCode = http.StatusBadRequest
	// 	return output

	// }
	if c.Deptname == "" {
		output.Message = "Dept can't be null"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusBadRequest
		return output
	}
	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	tx := db.Begin()
	tx.SavePoint("savepoint")

	result := tx.Raw("Select * from payroll.hrdept where deptcode =?", c.Deptcode).First(&c)
	if result.RowsAffected != 0 {
		tx.RollbackTo("savepoint")
		output.Message = "Department Code is already exist"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusBadRequest
		return output
	}
	res := tx.Raw("select * from payroll.hrdept where lower (deptname)=? ", strings.ToLower(c.Deptname)).First(&c)
	if res.RowsAffected != 0 {
		tx.RollbackTo("savepoint")
		output.Message = "Deptartment name alread exist"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusBadRequest
		return output
	}
	//_ = tx.Raw("select coalesce ((max(deptcode)+1),1) from payroll.hrdept").First(&c.Deptcode)
	fmt.Println(c.Deptcode)
	fmt.Println(c.Deptname)
	result1 := tx.Create(&c)
	if result1.RowsAffected == 0 {
		tx.RollbackTo("savepoint")
		output.Message = "Department creation failed"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusNotFound
		return output
	}

	tx.Commit()
	type abc struct {
		Output model.Hrdept `json:"output"`
	}
	var a abc
	a.Output = c
	output.Message = "Employee create succesfully"
	output.IsSuccess = true
	output.Payload = a
	output.StatusCode = http.StatusCreated
	return output
}

func (hrdeptrepo *HrdeptRepository) UpdateHrDeptNamebyDeptcode(input model.Hrdept) dto.ResponseDto {
	var response dto.ResponseDto
	if input.Deptcode <= 0 {
		response.Message = " Code can't be null"
		response.IsSuccess = false
		response.Payload = nil
		response.StatusCode = http.StatusBadRequest
		return response
	}
	if input.Deptname == "" {
		response.Message = "Dept can't be null"
		response.IsSuccess = false
		response.Payload = nil
		response.StatusCode = http.StatusBadRequest
		return response
	}
	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var output model.Hrdept
	result := db.Where(&model.Hrdept{Deptcode: input.Deptcode}).First(&output)
	if result.RowsAffected == 0 {
		response.Message = "this code doesnot exists"
		response.IsSuccess = false
		response.Payload = nil
		response.StatusCode = http.StatusNotFound
		return response
	}
	output.Deptname = input.Deptname
	result1 := db.Where(&model.Hrdept{Deptcode: input.Deptcode}).Updates(&output)
	if result1.RowsAffected == 0 {
		response.Message = "No Employee info found for given criteria"
		response.IsSuccess = false
		response.Payload = nil
		response.StatusCode = http.StatusInternalServerError
		return response
	}
	response.Message = "Employee info updated successfully"
	response.IsSuccess = true
	response.Payload = output
	response.StatusCode = http.StatusOK

	return response
}

func (hrdeptrepo *HrdeptRepository) DeleteHrDeptbyDeptcode(c model.Hrdept) dto.ResponseDto {
	var output dto.ResponseDto
	if c.Deptcode <= 0 {
		output.Message = "Invalid code"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusBadRequest
		return output
	}
	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	result := db.Where("deptcode = ?", c.Deptcode).Delete(&c)
	if result.RowsAffected == 0 {
		output.Message = "No info found for given criteria"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusNotFound
		return output
	}
	output.Message = "Deleted successfully"
	output.IsSuccess = true
	output.Payload = nil
	output.StatusCode = http.StatusOK
	return output
}

func (hrdeptrepo *HrdeptRepository) MaxDeptCode(c model.Hrdept) dto.ResponseDto {
	var output dto.ResponseDto
	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	tx := db.Begin()
	tx.SavePoint("savepoint")
	var outputcode model.Hrdept
	result := tx.Raw("select max(deptcode)+1 from payroll.hrdept").First(&outputcode.Deptcode)
	if result.RowsAffected == 0 {
		tx.RollbackTo("savepoint")
		output.Message = "No info found for given criteria"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusNotFound
		return output
	} else if result.Error != nil {
		outputcode.Deptcode = 1
	}

	tx.Commit()

	output.IsSuccess = true
	output.StatusCode = 200
	output.Message = "Max code for new dept entry"
	output.Payload = outputcode

	return output
}
