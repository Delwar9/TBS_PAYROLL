package service

import (
	"github.com/gin-gonic/gin"
	"github.com/tools/payroll/dto"
	"github.com/tools/payroll/model"
	"github.com/tools/payroll/repository"
)

type HrService struct{}

func (hrrestservice *HrService) GetAllHRDeptName(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	rep := new(repository.HrdeptRepository)
	res := rep.GetAllHrDept()
	c.JSON(res.StatusCode, res)
}

// GetOne returns all info of one
func (hrrestservice *HrService) GetHRDeptByDeptCode(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var code model.Hrdept
	c.ShouldBind(&code)
	rep := new(repository.HrdeptRepository)
	res := rep.GetHrDeptByDeptCode(code)
	c.JSON(res.StatusCode, res)
}

func (hrrestservice *HrService) InsertHRDeptService(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	var a model.Hrdept
	c.ShouldBind(&a)

	repo := new(repository.HrdeptRepository)
	response := repo.AddHrDept(a)
	c.JSON(response.StatusCode, response)
}

func (hrrestservice *HrService) UpdateHRDeptService(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	var a model.Hrdept
	c.ShouldBind(&a)

	repo := new(repository.HrdeptRepository)
	response := repo.UpdateHrDeptNamebyDeptcode(a)
	c.JSON(response.StatusCode, response)
}

func (hrrestservice *HrService) DeleteHRDeptService(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	var a model.Hrdept
	c.ShouldBind(&a)

	repo := new(repository.HrdeptRepository)
	response := repo.DeleteHrDeptbyDeptcode(a)
	c.JSON(response.StatusCode, response)
}

func (hrrestservice *HrService) GetMaxDeptCode(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var input model.Hrdept
	c.ShouldBind(&input)

	rep := new(repository.HrdeptRepository)
	res := rep.MaxDeptCode(input)
	c.JSON(res.StatusCode, res)
}

func (hrrestservice *HrService) GetAllHRDeptName2(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	rep := new(repository.HrdeptRepository)
	var k dto.Abcd
	c.ShouldBind(&k)
	res := rep.GetAllHrDept2(k)
	c.JSON(res.StatusCode, res)
}

// GetAllHrDept2

func (hrrestservice *HrService) AddRouters(router *gin.Engine) {
	router.POST("/api/v1/hrdept/add", hrrestservice.InsertHRDeptService)
	router.GET("/api/v1/hrdept/getalldept", hrrestservice.GetAllHRDeptName)
	router.POST("/api/v1/hrdept/getdeptnamebydeptcode", hrrestservice.GetHRDeptByDeptCode)
	router.PATCH("/api/v1/hrdept/updatedeptname", hrrestservice.UpdateHRDeptService)
	router.DELETE("/api/v1/hrdept/deletedept", hrrestservice.DeleteHRDeptService)
	router.GET("/api/v1/hrdept/maxdeptcode", hrrestservice.GetMaxDeptCode)
	router.POST("/api/v1/hrdept/getalldept2", hrrestservice.GetAllHRDeptName2)
}
