package repository

import (
	"github.com/tools/payroll/dto"
	"github.com/tools/payroll/model"
	"github.com/tools/payroll/util"
)

type LeaveTypes struct{}

func (l *LeaveTypes) GetAllLeaveTypes() dto.ResponseDto {
	var res dto.ResponseDto
	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var output []model.Leavetype
	result := db.Order("lcode").Find(&output)

	if result.RowsAffected == 0 {
		res.IsSuccess = false
		res.StatusCode = 404
		res.Message = "No Leave Types Found"
		res.Payload = nil
		return res
	}
	res.IsSuccess = true
	res.StatusCode = 200
	res.Message = "Leave Types Found"
	res.Payload = output
	return res
}

func (l *LeaveTypes) AddNewLeaveType(input model.Leavetype) dto.ResponseDto {
	var res dto.ResponseDto
	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	if input.Lname == "" {
		res.IsSuccess = false
		res.StatusCode = 400
		res.Message = "Valid Leave Type Name is required"
		res.Payload = nil
		return res
	}
	if input.Earndays <= 0 {
		res.IsSuccess = false
		res.StatusCode = 400
		res.Message = "Valid Earn Days is required"
		res.Payload = nil
		return res
	}
	result1 := db.Where("lname = ?", input.Lname).Find(&model.Leavetype{})
	if result1.RowsAffected > 0 {
		res.IsSuccess = false
		res.StatusCode = 400
		res.Message = "Leave Type Name already exists"
		res.Payload = nil
		return res
	}
	result2 := db.Where("lcode = ?", input.Lcode).Find(&model.Leavetype{})
	if result2.RowsAffected > 0 {
		res.IsSuccess = false
		res.StatusCode = 400
		res.Message = "Leave Type Code already exists"
		res.Payload = nil
		return res
	}
	_ = db.Raw("select coalesce ((max(lcode)+1),1) from payroll.leavetype").First(&input.Lcode)
	result := db.Create(&input)
	if result.RowsAffected == 0 {
		res.IsSuccess = false
		res.StatusCode = 500
		res.Message = "Unable to create Leave Type"
		res.Payload = nil
		return res
	}
	res.IsSuccess = true
	res.StatusCode = 200
	res.Message = "Leave Type Created"
	res.Payload = input
	return res
}

func (l *LeaveTypes) UpdateLeaveTypes(input model.Leavetype) dto.ResponseDto {
	var res dto.ResponseDto
	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	if input.Lcode <= 0 {
		res.IsSuccess = false
		res.StatusCode = 400
		res.Message = "Valid Leave Type Code is required"
		res.Payload = nil
		return res
	}
	result := db.Where("lcode = ?", input.Lcode).Find(&model.Leavetype{})
	if result.RowsAffected == 0 {
		res.IsSuccess = false
		res.StatusCode = 404
		res.Message = "Leave Type not found"
		res.Payload = nil
		return res
	}
	result1 := db.Model(&model.Leavetype{}).Where("lcode = ?", input.Lcode).Updates(input)
	if result1.RowsAffected == 0 {
		res.IsSuccess = false
		res.StatusCode = 500
		res.Message = "Unable to update Leave Type"
		res.Payload = nil
		return res
	}
	res.IsSuccess = true
	res.StatusCode = 200
	res.Message = "Leave Type Updated"
	res.Payload = input
	return res
}

func (l *LeaveTypes) DeleteLeaveType(input model.Leavetype) dto.ResponseDto {
	var res dto.ResponseDto
	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	var output model.Leavetype
	if input.Lcode <= 0 {
		res.IsSuccess = false
		res.StatusCode = 400
		res.Message = "Valid Leave Type Code is required"
		res.Payload = nil
		return res
	}
	_ = db.Raw("select * from payroll.leavetype where lcode = ?", input.Lcode).First(&output)
	result := db.Where("lcode = ?", input.Lcode).First(&model.Leavetype{})
	if result.RowsAffected == 0 {
		res.IsSuccess = false
		res.StatusCode = 404
		res.Message = "Leave Type not found"
		res.Payload = nil
		return res
	}
	result1 := db.Where("lcode = ?", input.Lcode).Delete(&model.Leavetype{})
	if result1.RowsAffected == 0 {
		res.IsSuccess = false
		res.StatusCode = 500
		res.Message = "Unable to delete Leave Type"
		res.Payload = nil
		return res
	}
	res.IsSuccess = true
	res.StatusCode = 200
	res.Message = "Leave Type Deleted"
	res.Payload = output
	return res
}

func (l *LeaveTypes) GetSingleLeaveType(input model.Leavetype) dto.ResponseDto {
	var res dto.ResponseDto
	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	var output model.Leavetype
	if input.Lcode <= 0 {
		res.IsSuccess = false
		res.StatusCode = 400
		res.Message = "Valid Leave Type Code is required"
		res.Payload = nil
		return res
	}

	_ = db.Raw("select * from payroll.leavetype where lcode = ?", input.Lcode).First(&output)
	result := db.Where("lcode = ?", input.Lcode).First(&model.Leavetype{})
	if result.RowsAffected == 0 {
		res.IsSuccess = false
		res.StatusCode = 404
		res.Message = "Leave Type not found"
		res.Payload = nil
		return res
	}

	res.IsSuccess = true
	res.StatusCode = 200
	res.Message = "Leave Type Found"
	res.Payload = output
	return res
}
