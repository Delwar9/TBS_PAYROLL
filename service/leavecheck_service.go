package service

import (
	"github.com/gin-gonic/gin"
	"github.com/tools/payroll/model"
	"github.com/tools/payroll/repository"
)

type Leavecheckservice struct{}

func (l *Leavecheckservice) AllEmpLeaveCheckService(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	var input model.Leaverecord
	c.ShouldBind(&input)
	repository := new(repository.Leavecheck)
	output := repository.CheckEarnLeave(input)

	c.JSON(output.StatusCode, output)
}

func (l *Leavecheckservice) AddRouters(router *gin.Engine) {
	router.POST("/api/v1/leave/allempleavecheck", l.AllEmpLeaveCheckService)

}
