package service

import (
	"github.com/gin-gonic/gin"
	"github.com/tools/payroll/dto"
	"github.com/tools/payroll/model"
	"github.com/tools/payroll/repository"
)

type EmployeeInformationService struct{}

func (e *EmployeeInformationService) AddEmployeeInformationService(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	var input dto.EmployeeInformation

	c.ShouldBind(&input)
	repository := new(repository.EmployeeInformation)
	output := repository.AddEmployeeInformation(input)

	c.JSON(output.StatusCode, output)
}

func (e *EmployeeInformationService) UpdateEmployeeInformationService(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	var input dto.EmployeeInformationUpdate
	c.ShouldBind(&input)
	repository := new(repository.EmployeeInformation)
	output := repository.UpdateEmployeeInformation(input)

	c.JSON(output.StatusCode, output)
}

func (e *EmployeeInformationService) DeleteEmployeeInformationService(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	var input dto.EmployeeInformationUpdate
	c.ShouldBind(&input)
	repository := new(repository.EmployeeInformation)
	output := repository.DeleteEmployeeInformation(input)

	c.JSON(output.StatusCode, output)
}

func (e *EmployeeInformationService) GetAnEmployeeInformationService(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	var input dto.EmployeeInformationUpdate
	c.ShouldBind(&input)
	repository := new(repository.EmployeeInformation)
	output := repository.GetAnEmployeeInformation(input)

	c.JSON(output.StatusCode, output)
}

func (e *EmployeeInformationService) GetMaxEmpCodeEmployeeInformationService(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	var input dto.EmployeeInformationUpdate
	c.ShouldBind(&input)
	repository := new(repository.EmployeeInformation)
	output := repository.GetMaxEmpcodeEmployeeInformation(input)

	c.JSON(output.StatusCode, output)
}

func (e *EmployeeInformationService) GetSerialListServie(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	var input model.Salarystructure
	c.ShouldBind(&input)
	repository := new(repository.EmployeeInformation)
	output := repository.GetEmployeeSerialList()

	c.JSON(output.StatusCode, output)
}

func (e *EmployeeInformationService) GetAllStaffinformationService(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	// var input model.Staffinformation
	// c.ShouldBind(&input)
	repository := new(repository.EmployeeInformation)
	output := repository.GetAllStaffinformation()

	c.JSON(output.StatusCode, output)
}

func (e *EmployeeInformationService) GetABankInfoService(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	var input model.Acchead
	c.ShouldBind(&input)
	repository := new(repository.EmployeeInformation)
	output := repository.GetAnDepositeBank(input)

	c.JSON(output.StatusCode, output)
}

func (e *EmployeeInformationService) GetAllBankInfoService(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	repository := new(repository.EmployeeInformation)
	output := repository.GetAllDepositeBank()

	c.JSON(output.StatusCode, output)
}

func (e *EmployeeInformationService) GetBankNameToUpdateService(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	var input model.Salarystructure
	c.ShouldBind(&input)
	repository := new(repository.EmployeeInformation)
	output := repository.GetBankNameToUpdate(input)

	c.JSON(output.StatusCode, output)
}

func (e *EmployeeInformationService) AddRouters(c *gin.Engine) {
	c.POST("/api/v1/addEmployeeInformation", e.AddEmployeeInformationService)
	c.PATCH("/api/v1/updateEmployeeInformation", e.UpdateEmployeeInformationService)
	c.DELETE("/api/v1/deleteEmployeeInformation", e.DeleteEmployeeInformationService)
	c.POST("/api/v1/getAnEmployeeInformation", e.GetAnEmployeeInformationService)
	c.GET("/api/v1/getMaxEmpCodeEmployeeInformation", e.GetMaxEmpCodeEmployeeInformationService)
	c.GET("/api/v1/getSerialList", e.GetSerialListServie)
	c.GET("/api/v1/getStaffinformation", e.GetAllStaffinformationService)
	c.POST("/api/v1/getABankInfo", e.GetABankInfoService)
	c.GET("/api/v1/getallbankinfo", e.GetAllBankInfoService)
	c.POST("/api/v1/getsalarybanknameforupdate", e.GetBankNameToUpdateService)
}
