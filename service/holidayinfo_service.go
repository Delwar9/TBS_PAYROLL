package service

import (
	"github.com/gin-gonic/gin"
	"github.com/tools/payroll/model"
	"github.com/tools/payroll/repository"
)

type HolidayinfoService struct{}

func (holidayinfoService *HolidayinfoService) GetHolidayinfoService(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var holidayinfo model.Holidayinfo
	c.ShouldBind(&holidayinfo)
	loanallotmentRepository := new(repository.HolidayinfoRepository)
	myOutput := loanallotmentRepository.GetHolidayinfoRepository(holidayinfo)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (holidayinfoService *HolidayinfoService) InsertHolidayinfoService(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var holidayinfo []model.Holidayinfo
	c.ShouldBind(&holidayinfo)
	loanallotmentRepository := new(repository.HolidayinfoRepository)
	myOutput := loanallotmentRepository.InsertHolidayinfoRepository(holidayinfo)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (holidayinfoService *HolidayinfoService) AddRouters(router *gin.Engine) {
	router.POST("/api/v1/holidayinfo-get-holidayinfo-by-year", holidayinfoService.GetHolidayinfoService)
	router.POST("/api/v1/holidayinfo-save-holidayinfo-by-year", holidayinfoService.InsertHolidayinfoService)
}
