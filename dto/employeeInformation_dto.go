package dto

import (
	"github.com/tools/payroll/model"
)

type EmployeeInformation struct {
	StaffInformation model.Staffinformation `json:"staffinformation"`
	SalaryStructure  model.Salarystructure  `json:"salarystructure"`
}

type EmployeeInformation_archive struct {
	StaffInformation_archive model.Staffinformation_archive `json:"staffinformation_archive"`
	SalaryStructure_archive  model.Salarystructure_archive  `json:"salarystructure_archive"`
}

type UpdateEmployeeInformation struct {
	EmployeeInformation         EmployeeInformation         `json:"employeeinformation"`
	EmployeeInformation_archive EmployeeInformation_archive `json:"employeeinformation_archive"`
}

type EmployeeInformationUpdate struct {
	StaffInformation         model.Staffinformation         `json:"staffinformation"`
	SalaryStructure          model.Salarystructure          `json:"salarystructure"`
	StaffInformation_archive model.Staffinformation_archive `json:"staffinformation_archive"`
	SalaryStructure_archive  model.Salarystructure_archive  `json:"salarystructure_archive"`
}

type EmployeeInformationDTO struct {
	Senioriy_serial int    `json:"senioriy_serial"`
	Empname         string `json:"empname"`
	Deptname        string `json:"deptname"`
	Designame       string `json:"designame"`
}

type MaxEmployeeInformationDTO struct {
	Empcode int `json:"empcode"`
}
type CustomEmpSaveDTO struct {
	Pempcode int    `json:"pempcode"`
	Pempname string `json:"pempname"`
}
