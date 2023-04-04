package service

import (
	"github.com/gin-gonic/gin"
	"github.com/tools/payroll/dto"
	"github.com/tools/payroll/model"
	"github.com/tools/payroll/repository"
)

type LeaveService struct{}

func (l *LeaveService) GetEmpInformationService(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	var input model.Leave
	c.ShouldBind(&input)
	repository := new(repository.Leave)
	output := repository.GetEmployInformation(input)

	c.JSON(output.StatusCode, output)
}

func (l *LeaveService) CheckRemainingLeaveServices(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	var input dto.LeaveCount
	c.ShouldBind(&input)
	repository := new(repository.Leave)
	output := repository.CheckConsumeLeave(input)

	c.JSON(output.StatusCode, output)
}

func (l *LeaveService) CreateLeave(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	var input model.Leave
	c.ShouldBind(&input)
	repository := new(repository.Leave)
	output := repository.EntryANewLeave(input)

	c.JSON(output.StatusCode, output)
}

func (l *LeaveService) leaveStatusService(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	var input model.Leave
	c.ShouldBind(&input)
	repository := new(repository.Leave)
	output := repository.LeaveStatus(input)

	c.JSON(output.StatusCode, output)
}

func (l *LeaveService) AddRouters(router *gin.Engine) {
	router.POST("/api/v1/leave/getempinfo", l.GetEmpInformationService)
	router.POST("api/v1/leave/checkleave", l.CheckRemainingLeaveServices)
	router.POST("/api/v1/leave/createleave", l.CreateLeave)
	router.POST("/api/v1/leave/leavestatus", l.leaveStatusService)
}
