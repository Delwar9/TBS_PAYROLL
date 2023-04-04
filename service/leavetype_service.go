package service

import (
	"github.com/gin-gonic/gin"
	"github.com/tools/payroll/model"
	"github.com/tools/payroll/repository"
)

type LeaveTypeService struct{}

func (l *LeaveTypeService) GetAllLeaveTypesService(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	repository := new(repository.LeaveTypes)
	output := repository.GetAllLeaveTypes()

	c.JSON(output.StatusCode, output)
}

func (l *LeaveTypeService) AddNewLeaveTypeService(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	var input model.Leavetype
	c.BindJSON(&input)
	repository := new(repository.LeaveTypes)
	output := repository.AddNewLeaveType(input)

	c.JSON(output.StatusCode, output)
}

func (l *LeaveTypeService) UpdateLeaveTypeService(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	var input model.Leavetype
	c.BindJSON(&input)
	repository := new(repository.LeaveTypes)
	output := repository.UpdateLeaveTypes(input)

	c.JSON(output.StatusCode, output)
}

func (l *LeaveTypeService) DeleteLeaveTypeService(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	var input model.Leavetype
	c.BindJSON(&input)
	repository := new(repository.LeaveTypes)
	output := repository.DeleteLeaveType(input)

	c.JSON(output.StatusCode, output)
}

func (l *LeaveTypeService) GetSingleLeaveTypeService(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	var input model.Leavetype
	c.BindJSON(&input)
	repository := new(repository.LeaveTypes)
	output := repository.GetSingleLeaveType(input)

	c.JSON(output.StatusCode, output)
}

func (l *LeaveTypeService) AddRouters(router *gin.Engine) {
	router.GET("/api/v1/leave/leavetypes", l.GetAllLeaveTypesService)
	router.POST("/api/v1/leave/addnewleavetypes", l.AddNewLeaveTypeService)
	router.PUT("/api/v1/leave/updateleavetypes", l.UpdateLeaveTypeService)
	router.DELETE("/api/v1/leave/deleteleavetypes", l.DeleteLeaveTypeService)
	router.POST("/api/v1/leave/singleleavetypes", l.GetSingleLeaveTypeService)
}
