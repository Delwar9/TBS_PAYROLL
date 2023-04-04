package service

import (
	"github.com/gin-gonic/gin"
	"github.com/tools/payroll/dto"
	"github.com/tools/payroll/model"
	"github.com/tools/payroll/repository"
)

type LoanallotmentService struct{}

func (loanallotmentService *LoanallotmentService) GetMaxRefNoService(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	loanallotmentRepository := new(repository.LoanallotmentRepository)
	myOutput := loanallotmentRepository.MaxRefNo()
	c.JSON(myOutput.StatusCode, myOutput)
}

func (loanallotmentService *LoanallotmentService) GetUserInfoService(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var empname dto.GetEmployeeInfoWithDeptAndDesigInputDto
	c.ShouldBind(&empname)
	loanallotmentRepository := new(repository.LoanallotmentRepository)
	myOutput := loanallotmentRepository.GetAllEmployeeWithDesigAndDept(empname)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (loanallotmentService *LoanallotmentService) GetEmployeeLoanInfo(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var empcode model.Loanallot
	c.ShouldBind(&empcode)
	loanallotmentRepository := new(repository.LoanallotmentRepository)
	myOutput := loanallotmentRepository.GetEmployeeLoanInfo(empcode)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (loanallotmentService *LoanallotmentService) ValidateForwardingDate(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var loan model.Loanallot
	c.ShouldBind(&loan)
	loanallotmentRepository := new(repository.LoanallotmentRepository)
	myOutput := loanallotmentRepository.ValidateForwardingDate(loan)
	c.JSON(myOutput.StatusCode, myOutput)
}

// TODO: CPF Loan
func (loanallotmentService *LoanallotmentService) TerminateCPFloan(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var loan model.Loanallot
	c.ShouldBind(&loan)
	loanallotmentRepository := new(repository.LoanallotmentRepository)
	myOutput := loanallotmentRepository.TerminateCPFloan(loan)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (loanallotmentService *LoanallotmentService) GetInstallmentAmountPerMonthForCPF(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var loanallotment model.Loanallot
	c.ShouldBind(&loanallotment)
	loanallotmentRepository := new(repository.LoanallotmentRepository)
	myOutput := loanallotmentRepository.GetInstallmentAmountPerMonthForCPF(loanallotment)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (loanallotmentService *LoanallotmentService) GetNoOfInstallmentForCPF(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var loanallotment model.Loanallot
	c.ShouldBind(&loanallotment)
	loanallotmentRepository := new(repository.LoanallotmentRepository)
	myOutput := loanallotmentRepository.GetNoOfInstallmentForCPF(loanallotment)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (loanallotmentService *LoanallotmentService) MaxNoOfInstallmentForCPF(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var loanallotment model.Loanallot
	c.ShouldBind(&loanallotment)
	loanallotmentRepository := new(repository.LoanallotmentRepository)
	myOutput := loanallotmentRepository.MaxNoOfInstallmentForCPF(loanallotment)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (loanallotmentService *LoanallotmentService) BankLoadForLoanForCPF(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	bankloanRepositoryForCPF := new(repository.LoanallotmentRepository)
	myOutput := bankloanRepositoryForCPF.BankLoadForLoanForCPF()
	c.JSON(myOutput.StatusCode, myOutput)
}

func (loanallotmentService *LoanallotmentService) GetABankForCPF(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var accid dto.Accid_Name_List_dto
	c.ShouldBind(&accid)
	loanallotmentRepository := new(repository.LoanallotmentRepository)
	myOutput := loanallotmentRepository.GetABankForLoanForCPF(accid)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (loanallotmentService *LoanallotmentService) InsertLoanAllotmentForCPF(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var loanallotment model.Loanallot
	c.ShouldBind(&loanallotment)
	loanallotmentRepository := new(repository.LoanallotmentRepository)
	myOutput := loanallotmentRepository.InsertLoanAllotmentForCPF(loanallotment)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (loanallotmentService *LoanallotmentService) GetLimitofCPFLoanAmount(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var limit dto.GetLimitofCPFLoanAmountInputDto
	c.ShouldBind(&limit)
	loanallotmentRepository := new(repository.LoanallotmentRepository)
	myOutput := loanallotmentRepository.GetLimitofCPFLoanAmount(limit)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (loanallotmentService *LoanallotmentService) GetLoanPayScheduleByEmpcodeForCPF(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var loanallotment model.Loanallot
	c.ShouldBind(&loanallotment)
	loanallotmentRepository := new(repository.LoanallotmentRepository)
	myOutput := loanallotmentRepository.GetLoanPayScheduleByEmpcodeForCPF(loanallotment)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (loanallotmentService *LoanallotmentService) GetDataEditLoanAllotmentByRefNoForCPF(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var loanallotment model.Loanallot
	// var loanallotment dto.GetDataToEditLoanAllotmentOutputDto
	c.ShouldBind(&loanallotment)
	loanallotmentRepository := new(repository.LoanallotmentRepository)
	myOutput := loanallotmentRepository.GetDataEditLoanAllotmentByRefNoForCPF(loanallotment)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (loanallotmentService *LoanallotmentService) EditLoanAllotmentForCPF(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var loanallotment model.Loanallot
	c.ShouldBind(&loanallotment)
	loanallotmentRepository := new(repository.LoanallotmentRepository)
	myOutput := loanallotmentRepository.EditLoanAllotmentForCPF(loanallotment)
	c.JSON(myOutput.StatusCode, myOutput)
}

// TODO: Advance
func (loanallotmentService *LoanallotmentService) TerminateAdvanceAgainstSalary(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var loan model.Loanallot
	c.ShouldBind(&loan)
	loanallotmentRepository := new(repository.LoanallotmentRepository)
	myOutput := loanallotmentRepository.TerminateAdvanceAgainstSalary(loan)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (loanallotmentService *LoanallotmentService) GetAllEmployeeForLoanServiceForAdvance(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	bankloanRepository := new(repository.LoanallotmentRepository)
	myOutput := bankloanRepository.GetAllEmployeeForAdvance()
	c.JSON(myOutput.StatusCode, myOutput)
}

func (loanallotmentService *LoanallotmentService) GetSingleEmployeeForAdvance(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var loanallotment model.Loanallot
	c.ShouldBind(&loanallotment)
	loanallotmentRepository := new(repository.LoanallotmentRepository)
	myOutput := loanallotmentRepository.GetSingleEmployeeForAdvance(loanallotment)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (loanallotmentService *LoanallotmentService) GetAllBankForLoanServiceForAdvance(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	bankloanRepository := new(repository.LoanallotmentRepository)
	myOutput := bankloanRepository.BankLoadForLoan()
	c.JSON(myOutput.StatusCode, myOutput)
}

func (loanallotmentService *LoanallotmentService) GetABankRepositoryForAdvance(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var accid dto.Accid_Name_List_dto
	c.ShouldBind(&accid)
	loanallotmentRepository := new(repository.LoanallotmentRepository)
	myOutput := loanallotmentRepository.GetABankForLoan(accid)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (loanallotmentService *LoanallotmentService) GetNoOfInstallmentForAdvance(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var loanallotment model.Loanallot
	c.ShouldBind(&loanallotment)
	loanallotmentRepository := new(repository.LoanallotmentRepository)
	myOutput := loanallotmentRepository.GetNoOfInstallmentWithModForAdvance(loanallotment)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (loanallotmentService *LoanallotmentService) GeInstallmentAmountForAdvance(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var loanallotment model.Loanallot
	c.ShouldBind(&loanallotment)
	loanallotmentRepository := new(repository.LoanallotmentRepository)
	myOutput := loanallotmentRepository.GeInstallmentAmountForAdvance(loanallotment)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (loanallotmentService *LoanallotmentService) InsertLoanAllotmentForAdvance(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var loanallotment model.Loanallot
	c.ShouldBind(&loanallotment)
	loanallotmentRepository := new(repository.LoanallotmentRepository)
	myOutput := loanallotmentRepository.InsertLoanAllotmentForAdvance(loanallotment)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (loanallotmentService *LoanallotmentService) GetLoanPayScheduleByEmpcodeForAdvance(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var loanallotment model.Loanallot
	c.ShouldBind(&loanallotment)
	loanallotmentRepository := new(repository.LoanallotmentRepository)
	myOutput := loanallotmentRepository.GetLoanPayScheduleByEmpcodeForAdvance(loanallotment)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (loanallotmentService *LoanallotmentService) GetDataEditLoanAllotmentByRefNoForAdvance(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var loanallotment model.Loanallot
	// var loanallotment dto.GetDataToEditLoanAllotmentOutputDto
	c.ShouldBind(&loanallotment)
	loanallotmentRepository := new(repository.LoanallotmentRepository)
	myOutput := loanallotmentRepository.GetDataEditLoanAllotmentByRefNoForAdvance(loanallotment)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (loanallotmentService *LoanallotmentService) EditLoanAllotmentForAdvance(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var loanallotment model.Loanallot
	c.ShouldBind(&loanallotment)
	loanallotmentRepository := new(repository.LoanallotmentRepository)
	myOutput := loanallotmentRepository.EditLoanAllotmentForAdvance(loanallotment)
	c.JSON(myOutput.StatusCode, myOutput)
}

// FIXME: This function is not used anywhere
func (loanallotmentService *LoanallotmentService) NextLoanCycle(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var a model.Loanallot
	c.ShouldBind(&a)
	loanallotmentRepository := new(repository.LoanallotmentRepository)
	myOutput := loanallotmentRepository.NextLoanCycle(a)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (loanallotmentService *LoanallotmentService) AddRouters(router *gin.Engine) {
	router.GET("/api/v1/loanallotment-get-max-refno", loanallotmentService.GetMaxRefNoService)
	router.POST("/api/v1/loanallotment-get-user-info", loanallotmentService.GetUserInfoService)
	router.POST("/api/v1/loanallotment-get-loan-info-by-empcode", loanallotmentService.GetEmployeeLoanInfo)
	router.POST("/api/v1/loanallotment-validate-fowarding-date", loanallotmentService.ValidateForwardingDate)
	// TODO: CPF
	router.POST("/api/v1/loanallotment-terminate-cpf-loan", loanallotmentService.TerminateCPFloan)
	router.POST("/api/v1/loanallotment-get-installment-amount-per-month-for-cpf", loanallotmentService.GetInstallmentAmountPerMonthForCPF)
	router.POST("/api/v1/loanallotment-get-no-of-installment-for-cpf", loanallotmentService.GetNoOfInstallmentForCPF)
	router.POST("/api/v1/loanallotment-max-no-of-installment-for-cpf", loanallotmentService.MaxNoOfInstallmentForCPF)
	router.GET("/api/v1/loanallotment-get-all-banklist-for-loan-for-cpf", loanallotmentService.BankLoadForLoanForCPF)
	router.POST("/api/v1/loanallotment-get-a-bank-for-cpf", loanallotmentService.GetABankForCPF)
	router.POST("/api/v1/loanallotment-insert-loanallotment-for-cpf", loanallotmentService.InsertLoanAllotmentForCPF)
	router.POST("/api/v1/loanallotment-get-maxlimit-of-loan-amount-for-cpf", loanallotmentService.GetLimitofCPFLoanAmount)
	router.POST("/api/v1/loanallotment-get-loanpayschedule-by-empcode-for-cpf", loanallotmentService.GetLoanPayScheduleByEmpcodeForCPF)
	router.POST("/api/v1/loanallotment-get-data-to-edit-loanallot-loanpayschedule-by-refno-for-cpf", loanallotmentService.GetDataEditLoanAllotmentByRefNoForCPF)
	router.PATCH("/api/v1/loanallotment-edit-loanallotment-for-cpf", loanallotmentService.EditLoanAllotmentForCPF)
	// TODO: Advance
	router.POST("/api/v1/loanallotment-terminate-for-advance", loanallotmentService.TerminateAdvanceAgainstSalary)
	router.GET("/api/v1/loanallotment-get-all-employeelist-for-advance", loanallotmentService.GetAllEmployeeForLoanServiceForAdvance)
	router.POST("/api/v1/loanallotment-get-single-employee-for-advance", loanallotmentService.GetSingleEmployeeForAdvance)
	router.GET("/api/v1/loanallotment-get-all-banklist-for-advance", loanallotmentService.GetAllBankForLoanServiceForAdvance)
	router.POST("/api/v1/loanallotment-get-a-bank-for-advance", loanallotmentService.GetABankRepositoryForAdvance)
	router.POST("/api/v1/loanallotment-get-no-of-installment-for-advance", loanallotmentService.GetNoOfInstallmentForAdvance)
	router.POST("/api/v1/loanallotment-get-installment-amount-for-advance", loanallotmentService.GeInstallmentAmountForAdvance)
	router.POST("/api/v1/loanallotment-insert-loanallotment-for-advance", loanallotmentService.InsertLoanAllotmentForAdvance)
	router.POST("/api/v1/loanallotment-get-loanpayschedule-by-empcode-for-advance", loanallotmentService.GetLoanPayScheduleByEmpcodeForAdvance)
	router.POST("/api/v1/loanallotment-get-data-to-edit-loanallot-loanpayschedule-by-refno-for-advance", loanallotmentService.GetDataEditLoanAllotmentByRefNoForAdvance)
	router.PATCH("/api/v1/loanallotment-edit-loanallotment-for-advance", loanallotmentService.EditLoanAllotmentForAdvance)
	// router.PATCH("/api/v1/loanallotment-get-next-loan-cycle", loanallotmentService.NextLoanCycle)
}
