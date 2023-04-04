package service

import (
	"github.com/gin-gonic/gin"
	"github.com/tools/payroll/model"
	"github.com/tools/payroll/repository"
)

type YearIncrementService struct{}

func (yearincrement *YearIncrementService) GetAnEmpService(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	var input model.Staffinformation

	c.ShouldBind(&input)
	repository := new(repository.YearIncrementRepository)
	output := repository.GetAnEmployeeInfo(input)

	c.JSON(output.StatusCode, output)
}

func (yearincrement *YearIncrementService) AddYearIncrement(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	var yearIncrement model.Yearincrement
	c.ShouldBind(&yearIncrement)

	repository := new(repository.YearIncrementRepository)
	result := repository.AddYearIncrement(yearIncrement)
	c.JSON(result.StatusCode, result)
}

func (yearincrement *YearIncrementService) AddRouters(c *gin.Engine) {
	c.POST("/api/v1/getanemp", yearincrement.GetAnEmpService)
	c.POST("/api/v1/newincrement", yearincrement.AddYearIncrement)
}
