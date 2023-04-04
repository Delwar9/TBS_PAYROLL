package service

import (
	"github.com/gin-gonic/gin"
	"github.com/tools/payroll/model"
	"github.com/tools/payroll/repository"
)

type MonthlyDeductService struct{}

func (monthlyDeductService *MonthlyDeductService) GetMonthlyDeductService(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	var input model.Monthlydeduct
	c.ShouldBind(&input)
	repository := new(repository.MonthlydeductRepository)
	output := repository.InsertMonthyDeductReposioty(input)
	c.JSON(output.StatusCode, output)
}

func (m *MonthlyDeductService) AddRouters(router *gin.Engine) {
	router.POST("/api/v1/monthlydeduct-get-monthlydeduct-by-year", m.GetMonthlyDeductService)
	// router.POST("/api/v1/monthlydeduct-save-monthlydeduct-by-year", m.InsertMonthlyDeductService)
}
