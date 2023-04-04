package service

import (
	"github.com/gin-gonic/gin"
	"github.com/tools/payroll/model"
	"github.com/tools/payroll/repository"
)

type DesignationRestService struct{}

func (designationRestService *DesignationRestService) GetAllDesignation(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	rep := new(repository.DesignationRepository)
	res := rep.GetAllDesignation()
	c.JSON(res.StatusCode, res)
}

func (designationRestService *DesignationRestService) GetDesignationByCode(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	var id model.Designation
	c.ShouldBind(&id)
	rep1 := new(repository.DesignationRepository)
	res1 := rep1.GetById(id)
	c.JSON(res1.StatusCode, res1)
}

func (designationRestService *DesignationRestService) AddDesignation(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	var obj model.Designation
	c.ShouldBind(&obj)
	rep := new(repository.DesignationRepository)
	res := rep.AddDesignation(obj)
	c.JSON(res.StatusCode, res)
}

func (designationRestService *DesignationRestService) UpdateDesignation(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	var obj1 model.Designation
	c.ShouldBind(&obj1)
	repo := new(repository.DesignationRepository)
	result2 := repo.UpdateDesignation(obj1)
	c.JSON(result2.StatusCode, result2)
}

func (designationRestService *DesignationRestService) DeleteDesignation(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	var a model.Designation
	c.ShouldBind(&a)

	repo := new(repository.DesignationRepository)
	response := repo.DeleteDesignation(a)
	c.JSON(response.StatusCode, response)
}

func (designationRestService *DesignationRestService) GetMaxDeptCode(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var input model.Designation
	c.ShouldBind(&input)

	rep := new(repository.DesignationRepository)
	res := rep.MaxDeptCode(input)
	c.JSON(res.StatusCode, res)
}

func (designationRestService *DesignationRestService) AddRouters(router *gin.Engine) {
	router.GET("/api/v1/designation/getalldesignation", designationRestService.GetAllDesignation)
	router.POST("/api/v1/designation/getdesigbycode", designationRestService.GetDesignationByCode)
	router.POST("/api/v1/designation/adddasignation", designationRestService.AddDesignation)
	router.PATCH("/api/v1/designation/updatedesigntion", designationRestService.UpdateDesignation)
	router.DELETE("/api/v1/designation/deletedesignation", designationRestService.DeleteDesignation)
	router.GET("/api/v1/designation/maxdesigcode", designationRestService.GetMaxDeptCode)
}
