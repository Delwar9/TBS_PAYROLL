package model

import (
	"encoding/json"

	"github.com/tools/payroll/util"
)

const _AuthInfoTable = "hrm.authentication_info"

// const _ACLInfoTable = "hrm.acl_info"

var _ValidRoles = map[string]bool{
	"HR":          true,
	"DEPT_HEAD":   true,
	"EMP_MANAGER": true,
	"USER":        true,
	"SUPER_ADMIN": true,
}

// AuthenticationInfo reprents the row in hrm.authentication_info table
type AuthenticationInfo struct {
	Empcode        string      `json:"empcode"` // Mandatory field
	EmaiID         string      `json:"email,"`  // Key
	Password       string      `json:"password"`
	Role           interface{} `json:"role"`
	ExtraInfo      interface{} `json:"extrainfo,omitempty"`
	PasswordExpiry string      `json:"pwdexpiry"`
	PasswordStatus string      `json:"pwdstat"`
	Fname          string      `json:"empname,omitempty"`
	Lname          string      `json:"lastname,omitempty"`
}

// GetInsertStatement returns insert sql statement
func (e *AuthenticationInfo) GetInsertStatement() (string, []interface{}) {
	var authInfoType AuthenticationInfo
	fmap := util.NewSQLFieldMap(_AuthInfoTable, authInfoType)
	return fmap.GenerateInsertScript(e)
}

// IsValidRole return true if the input role map is valid
func IsValidRole(roleMap map[string]interface{}) bool {
	if len(roleMap) == 0 {
		return false
	}
	for role := range roleMap {
		if _, isFound := _ValidRoles[role]; !isFound {
			return false
		}
	}
	return true
}

// BuildAuthInfo converts db results
func BuildAuthInfo(objs []interface{}) *AuthenticationInfo {
	jb, _ := json.Marshal(objs)
	authData := make([]AuthenticationInfo, 0)
	err := json.Unmarshal(jb, &authData)
	if err != nil {
		return nil
	}
	return &authData[0]
}

// ACLInfo contains entry of acl_info table
type ACLInfo struct {
	Action  string `json:"action"`
	EmpCode string `json:"empcode"`
	Role    string `json:"roles"`
}

// BuildACLInfo converts db results
func BuildACLInfo(objs []interface{}) []ACLInfo {
	jb, _ := json.Marshal(objs)
	authData := make([]ACLInfo, 0)
	err := json.Unmarshal(jb, &authData)
	if err != nil {
		return nil
	}
	return authData
}
