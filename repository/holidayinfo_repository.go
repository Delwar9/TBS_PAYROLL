package repository

import (
	"net/http"

	"github.com/tools/payroll/dto"
	"github.com/tools/payroll/model"
	"github.com/tools/payroll/util"
)

type HolidayinfoRepository struct{}

func (holidayinfoRepository *HolidayinfoRepository) GetHolidayinfoRepository(input model.Holidayinfo) dto.ResponseDto {
	var res dto.ResponseDto

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var output []model.Holidayinfo
	result := db.Where("year = ?", input.Year).Order("holidate").Find(&output)
	if result.RowsAffected == 0 {
		res.Message = "No Data Found"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusNotFound
		return res
	}

	for i := 0; i < len(output); i++ { // NOTE: remove extra value from date field
		output[i].Holidate = output[i].Holidate[0:10]
	}

	res.Message = "Successfully Retrieved Data"
	res.IsSuccess = true
	res.Payload = output
	res.StatusCode = http.StatusOK

	return res
}

func (holidayinfoRepository *HolidayinfoRepository) InsertHolidayinfoRepository(input []model.Holidayinfo) dto.ResponseDto {
	var res dto.ResponseDto

	db := util.CreatePayrollConnectionUsingGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	if len(input) == 0 {
		res.Message = "Please provide input data"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusBadRequest
		return res
	}

	var year int = input[0].Year

	_ = db.Exec("delete from payroll.holidayinfo where year = ? ", year)
	result1 := db.Create(&input)
	if result1.RowsAffected == 0 {
		res.Message = "Failed to Insert Data"
		res.IsSuccess = false
		res.Payload = nil
		res.StatusCode = http.StatusInternalServerError
		return res
	}

	res.Message = "Successfully Inserted Data"
	res.IsSuccess = true
	res.Payload = input
	res.StatusCode = http.StatusOK

	return res
}
