package service

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tools/payroll/model"
	"github.com/tools/payroll/util"
)

var _asLogger = logrus.New()

// AuthenticationRESTService provides reference data related rest services
type AuthenticationRESTService struct {
	dbUtil                 *util.PGSqlDBUtil
	jwtSigningKey          []byte
	bypassAuth             map[string]bool
	adminEmail             string
	adminPassword          string
	adminEmpCode           string
	cache                  map[string]UserRoleInfo
	acl                    map[string]bool
	aclEnabled             map[string]bool
	empSpecificActionCache map[string]map[string]bool
	// EmpCode-->Action->true/false
}

type authDataInput struct {
	Email          string      `json:"email"`
	Password       string      `json:"pwd"`
	NewPassword    string      `json:"newPwd"`
	Role           interface{} `json:"role,omitempty"`
	ExtraInfo      interface{} `json:"extrainfo,omitempty"`
	PasswordStatus *string     `json:"pwdstat,omitempty"`
}

// AuthorizationClaims JWTTokenClaims
type AuthorizationClaims struct {
	Empcode string `json:"empcode"`
	Email   string `json:"email"`
	jwt.StandardClaims
}
type authServiceConfig struct {
	JWTKey        *string  `json:"jwtKey"`
	BypassAuth    []string `json:"bypassAuth"`
	AdminEmail    string   `json:"adminEmailId"`
	AdminPassword string   `json:"adminPassword"`
	AdminEmpCode  string   `json:"adminEmpCode"`
}

// UserRoleInfo contains user role related inforation
type UserRoleInfo struct {
	RoleMap  map[string]bool
	EmpCode  string
	EmpEmail string
	// Other stuff
}

func buildUserRoleInfo(authInfo *model.AuthenticationInfo) UserRoleInfo {
	roles := make(map[string]bool)
	jb, _ := json.Marshal(authInfo.Role)
	json.Unmarshal(jb, &roles)
	if len(roles) == 0 {
		roles["USER"] = true
	}
	return UserRoleInfo{EmpCode: authInfo.Empcode, EmpEmail: authInfo.EmaiID, RoleMap: roles}
}

// NewAuthenticationRESTService retuens a new initialized version of the service
func NewAuthenticationRESTService(config []byte, dbUtil *util.PGSqlDBUtil, verbose bool) *AuthenticationRESTService {
	service := new(AuthenticationRESTService)
	if err := service.Init(config, dbUtil, verbose); err != nil {
		_asLogger.Errorf("unable to intialize service instance %v", err)
		return nil
	}
	return service
}

// Init intializes the service instance
func (s *AuthenticationRESTService) Init(config []byte, dbUtil *util.PGSqlDBUtil, verbose bool) error {
	if verbose {
		_asLogger.SetLevel(logrus.DebugLevel)
	}
	if dbUtil == nil {
		return fmt.Errorf("null DB Util reference passed ")
	}
	s.dbUtil = dbUtil
	var conf authServiceConfig
	err := json.Unmarshal(config, &conf)
	if err != nil {
		_asLogger.Errorf("unable to parse config json file %v", err)
		return err
	}
	if conf.JWTKey != nil && len(*conf.JWTKey) > 0 {
		s.jwtSigningKey = []byte(*conf.JWTKey)
	}
	s.bypassAuth = make(map[string]bool)
	s.bypassAuth["/"] = true
	if conf.BypassAuth != nil && len(conf.BypassAuth) > 0 {
		for _, url := range conf.BypassAuth {
			s.bypassAuth[url] = true
		}
	}
	s.adminEmail = conf.AdminEmail
	s.adminPassword = conf.AdminPassword
	s.adminEmpCode = conf.AdminEmpCode
	s.cache = make(map[string]UserRoleInfo)
	s.acl = make(map[string]bool)
	s.aclEnabled = make(map[string]bool)
	s.empSpecificActionCache = make(map[string]map[string]bool)
	s.loadACLInfo()
	_asLogger.Infof("successfully initialized AuthenticationRESTService")
	return nil
}

// AddRouters add api end points specific to this service
func (s *AuthenticationRESTService) AddRouters(router *gin.Engine) {
	router.POST("/api/auth/login", func(c *gin.Context) {
		resp := s.validateLogin(c)
		c.JSON(http.StatusOK, resp)
	})
	router.POST("/api/auth/updatepwd", func(c *gin.Context) {
		resp := s.updatePwd(c)
		c.JSON(http.StatusOK, resp)
	})
	// Allows to update role,extrainfo
	router.POST("/api/auth/updateinfo", func(c *gin.Context) {
		// Check the auth of the caller
		if s.HasPriviledge("UPDATE_AUTH_INFO", c) {
			resp := s.updateEntry(c)
			c.JSON(http.StatusOK, resp)

		} else {
			c.JSON(http.StatusUnauthorized, "Unauthorized")
		}
	})
	// Allows to reset password
	router.POST("/api/auth/resetpwd", func(c *gin.Context) {
		// Check the auth of the caller
		resp := s.resetPassword(c)
		c.JSON(http.StatusOK, resp)
	})
	// Get user roles only for admin
	router.GET("/api/auth/roles", func(c *gin.Context) {
		// Check the auth of the caller
		if s.HasPriviledge("VIEW_USER_ROLES", c) {
			resp := s.getUserRoles(c)
			c.JSON(http.StatusOK, resp)

		} else {
			c.JSON(http.StatusUnauthorized, "Unauthorized")
		}
	})
}

func (s *AuthenticationRESTService) validateLogin(c *gin.Context) APIResponse {
	var input authDataInput
	var finalResp APIResponse
	if !parseInput(c, &input) {
		return buildResponse(false, "Invalid input provided", nil)
	}
	// Check if super admin or not

	isOk, msg, authInfo := s.getAuthInfo(input)

	if !isOk {
		return buildResponse(false, msg, authInfo)
	}
	authInfo.Password = ""
	if authInfo.PasswordStatus == "DEFAULT" {
		authInfo.PasswordExpiry = ""
		finalResp = buildResponse(true, "Default password. Please change the password", authInfo)

	} else {
		today := s.getExpiryDate(0)
		if strings.Compare(today, authInfo.PasswordExpiry) > 0 {

			authInfo.PasswordStatus = "EXPIRED"
			finalResp = buildResponse(true, "Password expired. Please change the password", authInfo)
		} else {
			finalResp = buildResponse(true, "Login successful", authInfo)
		}
	}

	jwtToken := s.createJWTToken(authInfo.Empcode, authInfo.EmaiID)
	finalResp.Token = &jwtToken
	// Create the role info
	roleInfo := buildUserRoleInfo(authInfo)
	s.cache[authInfo.Empcode] = roleInfo
	return finalResp
}

func (s *AuthenticationRESTService) updatePwd(c *gin.Context) APIResponse {
	var input authDataInput
	if !parseInput(c, &input) {
		return buildResponse(false, "Invalid input provided", nil)
	}
	isOk, msg, authInfo := s.getAuthInfo(input)
	if !isOk {
		return buildResponse(false, msg, authInfo)
	}
	if strings.Compare(input.Password, input.NewPassword) == 0 {
		return buildResponse(false, "Old and new password can not be same", nil)
	}
	newPwdHash := s.getHashOf(input.NewPassword)
	newExpDate := s.getExpiryDate(90)
	sql := "update hrm.authentication_info set password=$1 , pwdexpiry=$2 , pwdstat='VALID' where email=$3"
	updCount, err := s.dbUtil.UpdateRecords(sql, []interface{}{newPwdHash, newExpDate, input.Email})
	if err != nil || updCount == 0 {
		_asLogger.Errorf("error in updatin %v", err)
		return buildResponse(false, "Password change failure", nil)
	}
	return buildResponse(true, "Password changed successfully", nil)
}

func (s *AuthenticationRESTService) getAuthInfo(input authDataInput) (bool, string, *model.AuthenticationInfo) {
	if len(input.Email) == 0 || len(input.Password) == 0 {
		return false, "Email or password missing", nil
	}
	if input.Email == s.adminEmail && input.Password == s.adminPassword {
		adminAuthInfo := model.AuthenticationInfo{
			Empcode:        s.adminEmpCode,
			EmaiID:         s.adminEmail,
			PasswordStatus: "VALID",
			PasswordExpiry: "20990101",
			Role: map[string]interface{}{
				"SUPER_ADMIN": true,
			},
			Lname: "Admin",
			Fname: "Default",
		}
		return true, "Login successful", &adminAuthInfo
	}
	records, err := s.dbUtil.QueryRecords("select A.*,B.empname,B.lastname from hrm.authentication_info A INNER JOIN hrm.employee_general_info B ON A.empcode = B.empcode WHERE A.email = $1 ", input.Email)
	if err != nil {
		_asLogger.Errorf("error in query hrm.authentication_info %v", err)
		return false, "Error in reading authentication records", nil
	}
	if len(records) == 0 {
		return false, "Login failed. Email or password is wrong.", nil
	}
	authInfo := model.BuildAuthInfo(records)
	if authInfo.Password != s.getHashOf(input.Password) {
		return false, "Login failed. Email or password is wrong", nil
	}
	return true, "", authInfo
}

func (s *AuthenticationRESTService) getUserRoles(c *gin.Context) APIResponse {
	empcode := c.Param("empcode")
	records, err := s.dbUtil.QueryRecords("select A.* from hrm.authentication_info A  WHERE A.empcode = $1 ", empcode)
	if err != nil {
		_asLogger.Errorf("error in query hrm.authentication_info %v", err)

		return buildResponse(false, "Error in reading authentication records", nil)
	}
	if len(records) == 0 {
		return buildResponse(false, "Login failed. Email or password is wrong.", nil)
	}
	authInfo := model.BuildAuthInfo(records)
	authInfo.Password = ""
	authInfo.PasswordExpiry = ""
	authInfo.PasswordStatus = ""
	return buildResponse(true, "User roles retrived successfully ", authInfo)
}

func (s *AuthenticationRESTService) resetPassword(c *gin.Context) APIResponse {
	var input authDataInput
	if !parseInput(c, &input) {
		return buildResponse(false, "Invalid input provided", nil)
	}
	records, err := s.dbUtil.QueryRecords("select * from hrm.authentication_info where email=$1 ", input.Email)
	if err != nil {
		_asLogger.Errorf("error in query hrm.authentication_info %v", err)
		return buildResponse(false, "Error in checking records", nil)
	}
	if len(records) == 0 {
		return buildResponse(false, "Email id not found", nil)
	}
	if input.PasswordStatus != nil && *input.PasswordStatus == "RESET" {
		// newRandomPwd := fmt.Sprintf("%x", time.Now().Unix())
		newPwdHash := s.getHashOf("passw0rd")
		newExpDate := s.getExpiryDate(3)
		sql := "update hrm.authentication_info set password=$1 , pwdexpiry=$2 , pwdstat='DEFAULT' where email=$3"
		updCount, err := s.dbUtil.UpdateRecords(sql, []interface{}{newPwdHash, newExpDate, input.Email})
		if err != nil || updCount == 0 {
			_asLogger.Errorf("error in updating %v", err)
			return buildResponse(false, "Password reset failure", nil)
		}
		// Need to send a email
		return buildResponse(true, "Password reset success", nil)
	}

	return buildResponse(false, "Invalid input given", input)
}

func (s *AuthenticationRESTService) updateEntry(c *gin.Context) APIResponse {
	var input authDataInput
	if !parseInput(c, &input) {
		return buildResponse(false, "Invalid input provided", nil)
	}
	records, err := s.dbUtil.QueryRecords("select * from hrm.authentication_info where email=$1 ", input.Email)
	if err != nil {
		_asLogger.Errorf("error in query hrm.authentication_info %v", err)
		return buildResponse(false, "Error in checking records", nil)
	}
	if len(records) == 0 {
		return buildResponse(false, "Email id not found", nil)
	}
	if input.PasswordStatus != nil && *input.PasswordStatus == "RESET" {
		newPwdHash := s.getHashOf("passw0rd")
		newExpDate := s.getExpiryDate(3)
		sql := "update hrm.authentication_info set password=$1 , pwdexpiry=$2 , pwdstat='DEFAULT' where email=$3"
		updCount, err := s.dbUtil.UpdateRecords(sql, []interface{}{newPwdHash, newExpDate, input.Email})
		if err != nil || updCount == 0 {
			_asLogger.Errorf("error in updating %v", err)
			return buildResponse(false, "Password reset failure", nil)
		}
		return buildResponse(true, "Password reset success", nil)
	}
	updParams := make([]interface{}, 0)
	sql := ""
	if input.Role != nil && input.ExtraInfo == nil {
		sql = "update hrm.authentication_info set role=$1  where email=$2"
		updParams = append(updParams, input.Role)
		updParams = append(updParams, input.Email)
	} else if input.Role != nil && input.ExtraInfo != nil {
		sql = "update hrm.authentication_info set role=$1 , extrainfo=$2  where email=$3 "
		updParams = append(updParams, input.Role)
		updParams = append(updParams, input.ExtraInfo)
		updParams = append(updParams, input.Email)
	} else if input.Role == nil && input.ExtraInfo != nil {
		sql = "update hrm.authentication_info set extrainfo=$1  where email=$2 "
		updParams = append(updParams, input.ExtraInfo)
		updParams = append(updParams, input.Email)
	} else {
		return buildResponse(false, "Invalid input", nil)
	}
	updCount, err := s.dbUtil.UpdateRecords(sql, updParams)
	if err != nil || updCount == 0 {
		_asLogger.Errorf("error in role/extn updating %v", err)
		return buildResponse(false, "Role/Extn update failure", nil)
	}
	return buildResponse(true, "Role/Extn update success", nil)
}

// CreateAuthInfo creates a authentication entry with default password
func (s *AuthenticationRESTService) CreateAuthInfo(emailID, empCode string, role map[string]interface{}, isMigration bool) (bool, string, string, []interface{}, string) {
	if len(emailID) == 0 || len(empCode) == 0 || role == nil || len(role) == 0 || !model.IsValidRole(role) {
		return false, "EmployeeID/Email/Role should not be blank or invalid", "", nil, ""
	}
	// Check for email id
	records, err := s.dbUtil.QueryRecords("select email from hrm.authentication_info where email=$1 ", emailID)
	if err != nil {
		_asLogger.Errorf("error in query hrm.authentication_info %v", err)
		return false, "Error in checking records", "", nil, ""
	}
	if len(records) > 0 {
		return false, "EmailID aready taken", "", nil, ""
	}
	var authRecord model.AuthenticationInfo
	defaultPassword := "passw0rd" // TOBE Randomly generated
	authRecord.Empcode = empCode
	authRecord.EmaiID = emailID
	authRecord.Password = s.getHashOf(defaultPassword)
	authRecord.Role = role
	authRecord.ExtraInfo = map[string]string{}
	authRecord.PasswordExpiry = s.getExpiryDate(90)
	authRecord.PasswordStatus = "DEFAULT"
	sql, params := authRecord.GetInsertStatement()

	return true, "", sql, params, defaultPassword
}

func (s *AuthenticationRESTService) getHashOf(password string) string {
	shaBytes := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", shaBytes)
}

func (s *AuthenticationRESTService) getExpiryDate(days int) string {
	return time.Now().AddDate(0, 0, days).Format("20060102")
}

// func (s *AuthenticationRESTService) checkAuth(c *gin.Context) bool {
// 	_asLogger.Infof("URL %s", c.Request.URL)

// 	url := c.Request.URL
// 	uri := url.RequestURI()
// 	if s.jwtSigningKey == nil || strings.EqualFold(uri, "/") {
// 		//jwt is not available
// 		return true
// 	}
// 	if _, isFound := s.bypassAuth[uri]; isFound {
// 		return true
// 	}
// 	authHeader := c.Request.Header.Get("Authorization")
// 	_asLogger.Infof("Authorization header received %s", authHeader)
// 	if len(authHeader) == 0 || !strings.HasPrefix(authHeader, "Bearer") {
// 		_asLogger.Infof("Invalid authotization header %s", authHeader)

// 		return false
// 	}
// 	tokenStrs := strings.Split(authHeader, " ")
// 	if len(tokenStrs) != 2 {
// 		_asLogger.Infof("Unable to parse authorization token %s", authHeader)
// 		return false
// 	}
// 	tokenStr := tokenStrs[1]
// 	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("UnexpectetokenStrd signing method: %v", token.Header["alg"])
// 		}
// 		return s.jwtSigningKey, nil
// 	})
// 	if err != nil {
// 		_asLogger.Errorf("Token parse error %v", err)
// 		return false
// 	}
// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		empCode := fmt.Sprintf("%v", claims["empcode"])
// 		empEmail := fmt.Sprintf("%v", claims["email"])
// 		//TODO:Verify empCode and emailID in the login cache
// 		//TODO:Also get the user role
// 		if role, isFound := s.cache[empCode]; isFound {

// 			c.Set("__ROLE_INFO__", role)
// 			_asLogger.Infof("Logged in user %s %s", empCode, empEmail)
// 			return true
// 		}
// 		_asLogger.Errorf("Role entry not found")
// 	}
// 	return false
// }

func (s *AuthenticationRESTService) createJWTToken(empCode, empEmail string) string {
	if s.jwtSigningKey == nil {
		return ""
	}
	expDate := time.Now().Add(1 * time.Hour).Unix()
	// stdClaim :=
	claim := AuthorizationClaims{
		empCode,
		empEmail,
		jwt.StandardClaims{
			ExpiresAt: expDate,
			Issuer:    "HRMS Systems",
			Id:        empCode,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := token.SignedString(s.jwtSigningKey)
	if err != nil {
		_asLogger.Errorf("error in generating token %v", err)
		return ""
	}
	_asLogger.Infof("Generated token %v", tokenStr)

	return tokenStr
}

// GetLoggedInUserEmpCode returns logged-in users empcode
func (s *AuthenticationRESTService) GetLoggedInUserEmpCode(c *gin.Context) string {
	if usrInfo, isExisting := c.Get("__ROLE_INFO__"); isExisting {
		if usrDetails, isOk := usrInfo.(UserRoleInfo); isOk {
			return usrDetails.EmpCode
		}
	}
	return "----"
}

// GetLoggedInUserRoleInfo returns logged-in users empcode
func (s *AuthenticationRESTService) GetLoggedInUserRoleInfo(c *gin.Context) (string, map[string]bool) {
	if usrInfo, isExisting := c.Get("__ROLE_INFO__"); isExisting {
		if usrDetails, isOk := usrInfo.(UserRoleInfo); isOk {
			return usrDetails.EmpCode, usrDetails.RoleMap
		}
	}
	return "", map[string]bool{}
}

// HasPriviledge checks if the logged in user has right priviledge for the input action
func (s *AuthenticationRESTService) HasPriviledge(action string, c *gin.Context) bool {
	if _, isFound := s.aclEnabled[action]; !isFound {
		return true
	}
	if usrInfo, isExisting := c.Get("__ROLE_INFO__"); isExisting {
		if usrDetails, isOk := usrInfo.(UserRoleInfo); isOk {
			return s.hasEntry(action, usrDetails.RoleMap)
		}
	}
	return false
}

func (s *AuthenticationRESTService) hasEntry(action string, roles map[string]bool) bool {
	for role := range roles {
		k := fmt.Sprintf("%s:%s", action, role)
		if _, isFound := s.acl[k]; isFound {
			return true
		}
	}
	return false
}

func (s *AuthenticationRESTService) loadACLInfo() {
	recs, err := s.dbUtil.QueryRecords("select * from hrm.acl_info")
	if err != nil || len(recs) == 0 {
		_asLogger.Warnf("unable to load acl_info records %v", err)
		return
	}
	aclRecords := model.BuildACLInfo(recs)
	for _, aclentry := range aclRecords {
		s.aclEnabled[aclentry.Action] = true
		roles := strings.Split(aclentry.Role, ",")
		if aclentry.EmpCode == "*" {
			for _, role := range roles {
				k := fmt.Sprintf("%s:%s", aclentry.Action, role)
				s.acl[k] = true
				_asLogger.Info("Loaded acl entry ", k)
			}
		}
		// else {
		// 	//This is employee specific entry
		// }
	}
}
