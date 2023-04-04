package repository

import (
	"net/http"
	"strings"

	"github.com/tools/payroll/dto"
	"github.com/tools/payroll/model"
	"github.com/tools/payroll/util"
)

type DesignationRepository struct{}

func (designationrepo *DesignationRepository) GetAllDesignation() dto.ResponseDto {
	var output dto.ResponseDto
	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var obj []model.Designation
	result := db.Order("desigcode").Find(&obj)
	if result.RowsAffected == 0 {
		output.Message = "No designation info found!"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusNotFound
		return output
	}
	type tempOutPut struct {
		T_output    []model.Designation `json:"output"`
		OutputCount int                 `json:"outputCount"`
	}
	var tOutput tempOutPut
	tOutput.T_output = obj
	tOutput.OutputCount = len(obj)
	output.Message = "List of designations"
	output.IsSuccess = true
	output.Count = len(obj)
	output.Payload = tOutput
	output.StatusCode = http.StatusOK

	return output
}

func (designationrepo *DesignationRepository) GetById(id model.Designation) dto.ResponseDto {
	var output dto.ResponseDto
	if id.Desigcode <= 0 {
		output.Message = "Designation can't be null"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusBadRequest
		return output
	}

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	result := db.Raw("select * from payroll.Designation where desigcode = ?", id.Desigcode).First(&id)
	if result.RowsAffected == 0 {
		output.Message = "No country info found"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusNotFound
		return output
	}
	type tempOutput struct {
		Output model.Designation `json:"output"`
	}
	var tOutput tempOutput
	tOutput.Output = id
	output.Message = "Designation info details found for given criteria"
	output.IsSuccess = true
	output.Payload = tOutput
	output.StatusCode = http.StatusOK
	return output
}

func (designationrepo *DesignationRepository) AddDesignation(c model.Designation) dto.ResponseDto {
	var output dto.ResponseDto
	// if c.Desigcode <= 0 {
	// 	output.Message = "Invalid code"
	// 	output.IsSuccess = false
	// 	output.Payload = nil
	// 	output.StatusCode = http.StatusBadRequest
	// 	return output

	// }
	if c.Designame == "" {
		output.Message = "Name can't be null"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusBadRequest
		return output

	}
	// if c.Sdesig == "" {
	// 	output.Message = "Dept can't be null"
	// 	output.IsSuccess = false
	// 	output.Payload = nil
	// 	output.StatusCode = http.StatusBadRequest
	// 	return output
	// }
	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	tx := db.Begin()
	tx.SavePoint("savepoint1")

	// result := db.Raw("Select * from public.desig where code =?", c.Code).First(&c)
	// result := db.Where(&model.Designation{Desigcode: c.Desigcode}).First(&c)
	// if result.RowsAffected != 0 {
	// 	tx.RollbackTo("savepoint1")
	// 	output.Message = "Department Code is already exist"
	// 	output.IsSuccess = false
	// 	output.Payload = nil
	// 	output.StatusCode = http.StatusConflict
	// 	return output
	// }

	// result2 := db.Raw("Select * from payroll.designation where lower(designame) =?", strings.ToLower(c.Designame)).First(&c)
	// if result2.RowsAffected != 0 {
	// 	tx.RollbackTo("savepoint1")
	// 	output.Message = "Designation is alread exist"
	// 	output.IsSuccess = false
	// 	output.Payload = nil
	// 	output.StatusCode = http.StatusConflict
	// 	return output
	// }

	result1 := db.Raw("Select * from payroll.designation where lower(designame) =? or lower(sdesig) = ?", strings.ToLower(c.Designame), strings.ToLower(c.Sdesig)).First(&c)
	if result1.RowsAffected != 0 {
		tx.RollbackTo("savepoint1")
		output.Message = "Designation Or Short_Designation is alread exist"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusConflict
		return output
	}
	// result2 := db.Raw("Select * from public.desig where sdesignation =?", c.Sdesignation).First(&c)
	// if result2.RowsAffected !=0{
	// 	output.Message ="Designation is alread exist"
	// 	output.IsSuccess = false
	// 	output.Payload=nil
	// 	output.StatusCode = http.StatusBadRequest
	// 	return output

	//_ = tx.Raw("select coalesce ((max(desigcode)+1),1) from payroll.designation").First(&c.Desigcode)

	if c.Sdesig == "" {
		c.Sdesig = "0"
	}
	result3 := tx.Create(&c)
	if result3.RowsAffected == 0 {
		tx.RollbackTo("savepoint1")
		output.Message = "Designation creation failed"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusInternalServerError
		return output
	}

	tx.Commit()
	type abc struct {
		Output model.Designation `json:"output"`
	}
	var a abc
	a.Output = c
	output.Message = "Designation create succesfully"
	output.IsSuccess = true
	output.Payload = a
	output.StatusCode = http.StatusOK
	return output
}

func (designationrepo *DesignationRepository) UpdateDesignation(input model.Designation) dto.ResponseDto {
	var response dto.ResponseDto
	if input.Desigcode <= 0 {
		response.Message = " Code can't be null"
		response.IsSuccess = false
		response.Payload = nil
		response.StatusCode = http.StatusBadRequest
		return response
	}
	if input.Designame == "" {
		response.Message = "Designation can't be null"
		response.IsSuccess = false
		response.Payload = nil
		response.StatusCode = http.StatusBadRequest
		return response

	}
	if input.Sdesig == "" {
		response.Message = "ShortDesig can't be null"
		response.IsSuccess = false
		response.Payload = nil
		response.StatusCode = http.StatusBadRequest
		return response
	}
	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	tx := db.Begin()
	tx.SavePoint("savepoint1")

	var output model.Designation
	result := db.Where(&model.Designation{Desigcode: input.Desigcode}).First(&output)
	if result.RowsAffected == 0 {
		tx.RollbackTo("savepoint1")
		response.Message = "this code doesnot exists"
		response.IsSuccess = false
		response.Payload = nil
		response.StatusCode = http.StatusNotFound
		return response
	}
	output.Designame = input.Designame
	output.Sdesig = input.Sdesig
	result1 := db.Where(&model.Designation{Desigcode: input.Desigcode}).Updates(&output)
	if result1.RowsAffected == 0 {
		response.Message = "No Employee info found for given criteria"
		response.IsSuccess = false
		response.Payload = nil
		response.StatusCode = http.StatusInternalServerError
		return response
	}
	tx.Commit()
	response.Message = "Employee info updated successfully"
	response.IsSuccess = true
	response.Payload = output
	response.StatusCode = http.StatusOK

	return response
}

func (designationrepo *DesignationRepository) DeleteDesignation(c model.Designation) dto.ResponseDto {
	// var output dto.ResponseDto
	// if code.Desigcode <= 0 {
	// 	output.Message = "Invalid code"
	// 	output.IsSuccess = false
	// 	output.Payload = nil
	// 	output.StatusCode = http.StatusBadRequest
	// 	return output
	// }
	// db := util.CreatePayrollConnectionUsingGorm()
	// sqlDB, _ := db.DB()
	// defer sqlDB.Close()

	// var op model.Designation

	// // result := db.Where("desigcode = ?", code.Desigcode).Delete(&model.Designation{})
	// result1 := db.Where(&model.Designation{Desigcode: code.Desigcode}).First(&op)

	// if result1.RowsAffected == 0 {
	// 	output.Message = "No info found for given criteria"
	// 	output.IsSuccess = false
	// 	output.Payload = nil
	// 	output.StatusCode = http.StatusNotFound
	// 	return output
	// }
	// result := db.Where("desigcode = ?",code.Desigcode).Delete(&op)
	// if result.RowsAffected == 0 {
	// 	output.Message = "Data is not Deleted Cause of some error"
	// 	output.IsSuccess = false
	// 	output.Payload = nil
	// 	output.StatusCode = http.StatusInternalServerError
	// 	return output
	// }
	// output.Message = "Deleted successfully"
	// output.IsSuccess = true
	// output.Payload = nil
	// output.StatusCode = http.StatusOK
	// return output

	var output dto.ResponseDto
	if c.Desigcode <= 0 {
		output.Message = "Invalid code"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusBadRequest
		return output
	}
	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// result1 := db.Where(&model.Designation{Desigcode: c.Desigcode}).First(&output)
	result1 := db.Raw("select * from payroll.designation where desigcode = ?", c.Desigcode).First(&output)
	if result1.RowsAffected == 0 {

		output.Message = "this code doesnot exists"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusNotFound
		return output
	}

	// result := db.Where("desigcode =?", c.Desigcode).Delete(&c)
	result := db.Where(&model.Designation{Desigcode: c.Desigcode}).Delete(&c)
	if result.RowsAffected == 0 {
		output.Message = "No info found for given criteria"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusInternalServerError
		return output
	}
	output.Message = "Deleted successfully"
	output.IsSuccess = true
	output.Payload = nil
	output.StatusCode = http.StatusOK
	return output
}

func (designationrepo *DesignationRepository) MaxDeptCode(c model.Designation) dto.ResponseDto {
	var output dto.ResponseDto
	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	tx := db.Begin()
	tx.SavePoint("savepoint")
	var outputcode model.Designation
	result := tx.Raw("select max(desigcode)+1 from payroll.designation").First(&outputcode.Desigcode)
	if result.RowsAffected == 0 {
		tx.RollbackTo("savepoint")
		output.IsSuccess = false
		output.StatusCode = http.StatusNotFound
		output.Message = "Internal Server error!"
		output.Payload = nil
		return output
	} else if result.Error != nil {
		outputcode.Desigcode = 1
	}

	tx.Commit()

	output.IsSuccess = true
	output.StatusCode = 200
	output.Message = "Max code for new dept entry"
	output.Payload = outputcode

	return output
}
