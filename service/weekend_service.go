package service

import (
	"github.com/gin-gonic/gin"
	"github.com/tools/payroll/model"
	"github.com/tools/payroll/repository"
)

type WeekendService struct{}

func (weekendService *WeekendService) GetWeekendService(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var weekend model.Weekend
	c.ShouldBind(&weekend)
	loanallotmentRepository := new(repository.WeekendRepository)
	myOutput := loanallotmentRepository.GetWeekendRepository(weekend)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (weekendService *WeekendService) InsertWeekendService(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var weekend []model.Weekend
	c.ShouldBind(&weekend)
	loanallotmentRepository := new(repository.WeekendRepository)
	myOutput := loanallotmentRepository.InsertWeekendRepository(weekend)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (weekendService *WeekendService) AddRouters(router *gin.Engine) {
	router.POST("/api/v1/weekend-get-weekend-by-year", weekendService.GetWeekendService)
	router.POST("/api/v1/weekend-save-weekend-by-year", weekendService.InsertWeekendService)
}
