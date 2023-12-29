package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/dgrijalva/jwt-go"
	"github.com/dwarkesh2810/golang-demo/conf"
	"github.com/dwarkesh2810/golang-demo/dto"
	"github.com/dwarkesh2810/golang-demo/models"
	"github.com/dwarkesh2810/golang-demo/pkg/helpers"
	"github.com/dwarkesh2810/golang-demo/pkg/logger"
	"github.com/dwarkesh2810/golang-demo/pkg/validations"
)

type UserController struct {
	beego.Controller
}

var key = conf.Env.JwtSecret
var jwtKey = []byte(key)

// Login ...
// @Title login User
// @Desciption login
// @Param body body dto.UserLoginRequest true "login User"
// @Param lang query string false "use en-US or hi-IN"
// @Param   Accept-Languages   header  string  false  "Bearer YourAccessToken"
// @Success 201 {object} object
// @Failure 403
// @router /login [post]
func (c *UserController) Login() {
	var userReq dto.UserLoginRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &userReq)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		logger.InsertAuditLogs(c.Ctx, err.Error(), 0)
		return
	}

	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&userReq); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		logger.InsertAuditLogs(c.Ctx, "Error : Validation error", 0)
		return
	}

	dbUser, err := models.GetUserByEmail(userReq.Username)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), 0)
		return
	}
	if dbUser.Email == "" {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "login"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+logger.LogMessage(c.Ctx, "error.login"), 0)
		return
	}
	err = helpers.VerifyHashedData(dbUser.Password, userReq.Password)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "credential"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), uint(dbUser.UserId))
		return
	}
	userData, err := models.LoginUser(userReq.Username, dbUser.Password)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), uint(dbUser.UserId))
		return
	}
	if userData.Email == "" && userData.FirstName == "" {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "login"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+logger.LogMessage(c.Ctx, "error.login"), uint(userData.UserId))
		return
	}
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &dto.JwtClaim{Email: userData.Email, ID: int(userData.UserId), StandardClaims: jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
	}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "accesstoken"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+logger.LogMessage(c.Ctx, "error.accesstoken"), uint(userData.UserId))
		return
	}
	data := map[string]interface{}{"User_Data": token.Claims, "Tokan": tokenString}
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, data, helpers.TranslateMessage(c.Ctx, "success", "login"), "")

	logger.InsertAuditLogs(c.Ctx, fmt.Sprintf("Login by user : %d", userData.UserId), uint(userData.UserId))
}

// RegisterNewUser ...
// @Title Insert New User
// @Desciption new users
// @Param lang query string false "use en-US or hi-IN"
// @Param body body dto.NewUserRequest true "Insert New User"
// @Success 201 {object} models.Users
// @Failure 403
// @router /register [post]
func (c *UserController) RegisterNewUser() {
	bodyData := dto.NewUserRequest{}
	if err := c.ParseForm(&bodyData); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), 0)
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &bodyData)

	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&bodyData); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		logger.InsertAuditLogs(c.Ctx, "Error : Validation error", 0)
		return
	}

	data, _ := models.GetUserByEmail(bodyData.Email)
	if data.Email == bodyData.Email || data.PhoneNumber == bodyData.PhoneNumber {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "userexist"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+logger.LogMessage(c.Ctx, "error.userexist"), 0)
		return
	}
	output, err := models.InsertNewUser(bodyData)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), 0)
		return
	}
	go models.VerifyEmail(output.Email, output.FirstName)
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, output, helpers.TranslateMessage(c.Ctx, "success", "user"), "")

	logger.InsertAuditLogs(c.Ctx, fmt.Sprintf("New user registered: %d", output.UserId), uint(output.UserId))
}

// GetAll ...
// @Title Get All
// @Description get Users
// @Param lang query string false "use en-US or hi-IN"
// @Param body body dto.PaginationReq false "Insert New User"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 200 {object} models.Users
// @Failure 403
// @router /secure/users/ [post]
func (c *UserController) GetAllUsers() {
	var search dto.PaginationReq
	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))

	if err := c.ParseForm(&search); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &search)

	result, pagination_data, err := models.FetchUsers(search.OpenPage, search.PageSize)

	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "datanotfound"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
	}
	if pagination_data["pageOpen_error"] == 1 {
		current := pagination_data["current_page"]
		last := pagination_data["last_page"]
		message := fmt.Sprintf(helpers.TranslateMessage(c.Ctx, "error", "page"), current, last)
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, message)
		logger.InsertAuditLogs(c.Ctx, "Error :"+fmt.Sprintf(logger.LogMessage(c.Ctx, "error.page"), current, last), userId)
		return
	}

	if result != nil {
		section_message := ""
		section := ""
		message := helpers.TranslateMessage(c.Ctx, section, section_message)
		helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, result, message, pagination_data)
		logger.InsertAuditLogs(c.Ctx, logger.LogMessage(c.Ctx, "success.read"), userId)
		return
	} else {
		helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, result, helpers.TranslateMessage(c.Ctx, "success", "data"), nil)
		logger.InsertAuditLogs(c.Ctx, logger.LogMessage(c.Ctx, "success.data"), userId)
		return
	}

}

// VerifyEmailOTP ...
// @Title verify otp for email
// @Desciption otp verification for eamil
// @Param body body dto.VerifyEmailOTPRequest true "otp verification for email"
// @Param lang query string false "use en-US or hi-IN"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} string
// @Failure 403
// @router /secure/verify_email_otp [post]
func (c *UserController) VerifyEmailOTP() {
	var bodyData dto.VerifyEmailOTPRequest

	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))
	if err := c.ParseForm(&bodyData); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &bodyData)

	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&bodyData); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		logger.InsertAuditLogs(c.Ctx, "Error : Validation error", userId)
		return
	}
	data, err := models.GetEmailOTP(bodyData.Username, bodyData.Otp)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
		return
	}
	if data.OtpCode != bodyData.Otp {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "otp"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+logger.LogMessage(c.Ctx, "error.otp"), userId)
		return
	}
	err = models.UpdateIsVerified(data.UserId)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
		return
	}
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, "Verified", helpers.TranslateMessage(c.Ctx, "success", "verify"), "")
	logger.InsertAuditLogs(c.Ctx, logger.LogMessage(c.Ctx, "success.verify"), userId)
}

// UpdateUser .
// @Title update User
// @Desciption update users
// @Param body body dto.UpdateUserRequest true "update New User"
// @Param lang query string false "use en-US or hi-IN"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} models.Users
// @Failure 403
// @router /secure/update [put]
func (c *UserController) UpdateUser() {
	var bodyData dto.UpdateUserRequest

	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))

	if err := c.ParseForm(&bodyData); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &bodyData)

	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&bodyData); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		logger.InsertAuditLogs(c.Ctx, "Error : Validation error", userId)
		return
	}
	data, err := models.GetUserDetails(bodyData.Id)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
		return
	}
	if bodyData.Email != data.Email {
		res, _ := models.GetUserByEmail(bodyData.Email)
		if res.Email == bodyData.Email {
			helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "userexist"))
			logger.InsertAuditLogs(c.Ctx, "Error : "+logger.LogMessage(c.Ctx, "error.userexist"), userId)
			return
		}
		if data.Email == "" {
			helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "invalidid"))
			logger.InsertAuditLogs(c.Ctx, "Error : "+logger.LogMessage(c.Ctx, "error.invalidid"), userId)
			return
		}
	}
	user, err := models.UpdateUser(bodyData)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
		return
	}

	userDetails := dto.UserDetailsResponse{Id: user.UserId, FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, Country: user.CountryId}
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, userDetails, helpers.TranslateMessage(c.Ctx, "success", "update"), "")
	logger.InsertAuditLogs(c.Ctx, logger.LogMessage(c.Ctx, "success.update"), userId)
}

// ResetPassword ...
// @Title Reset password
// @Desciption Reset password
// @Param body body dto.ResetUserPassword true "reset password"
// @Param lang query string false "use en-US or hi-IN"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} models.Users
// @Failure 403
// @router /secure/reset_pass [post]
func (c *UserController) ResetPassword() {
	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))
	output, err := models.GetUserDetails(userId)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
		return
	}
	var bodyData dto.ResetUserPassword
	if err := c.ParseForm(&bodyData); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &bodyData)

	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&bodyData); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		logger.InsertAuditLogs(c.Ctx, "Error : Validation error", userId)
		return
	}
	err = helpers.VerifyHashedData(output.Password, bodyData.CurrentPass)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "hashpassword"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+logger.LogMessage(c.Ctx, "error.hashpassword"), userId)
		return
	}
	if bodyData.ConfirmPass != bodyData.NewPass {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "passwordnotmatched"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+logger.LogMessage(c.Ctx, "error.passwordnotmatched"), userId)

		return
	}
	uppass, err := models.ResetPassword(bodyData.NewPass, int(userId))
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
		return
	}
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, uppass, helpers.TranslateMessage(c.Ctx, "success", "passwordreset"), "")
	logger.InsertAuditLogs(c.Ctx, logger.LogMessage(c.Ctx, "success.passwordreset"), userId)
}

// @Title delete user
// @Description delete user
// @Param id path int true "user id to delete recode"
// @Param lang query string false "use en-US or hi-IN"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 200 {string} string
// @Failure 403
// @router /secure/delete/:id([0-9]+) [delete]
func (c *UserController) DeleteUser() {

	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))

	idString := c.Ctx.Input.Params()
	id, err := strconv.Atoi(idString["1"])
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "convert", "atoi"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
		return
	}
	data, err := models.DeleteUser(id)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
		return
	}
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, data, helpers.TranslateMessage(c.Ctx, "success", "delete"), "")
	logger.InsertAuditLogs(c.Ctx, logger.LogMessage(c.Ctx, "success.delete"), userId)
}

// SendOtp ...
// @Title forgot password
// @Desciption forgot password
// @Param body body dto.SendOtpData true "forgot password this is send otp on mobile and email"
// @Param lang query string false "use en-US or hi-IN"
// @Success 201 {object} string
// @Failure 403
// @router /secure/forgot_pass [post]
func (c *UserController) ForgotPassword() {

	var bodyData dto.SendOtpData
	if err := c.ParseForm(&bodyData); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), 0)
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &bodyData)

	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&bodyData); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		logger.InsertAuditLogs(c.Ctx, "Error : "+"Validation error", 0)
		return
	}
	output, err := models.GetUserByEmail(bodyData.Username)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), 0)
		return
	}
	if output.Email == "" {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "login"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+logger.LogMessage(c.Ctx, "error.login"), 0)
		return
	}
	res, err := models.VerifyEmail(output.Email, output.FirstName)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), 0)
		return
	}
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "success", res), helpers.TranslateMessage(c.Ctx, "success", "otpsent"), "")
	logger.InsertAuditLogs(c.Ctx, logger.LogMessage(c.Ctx, "success.otpsent"), 0)
}

// VerifyOtpResetpassword ...
// @Title verify otp
// @Desciption otp verification for forgot password
// @Param body body dto.ResetUserPasswordOtp true "otp verification for forgot password"
// @Param lang query string false "use en-US or hi-IN"
// @Success 201 {object} string
// @Failure 403
// @router /secure/reset_pass_otp [post]
func (c *UserController) VerifyOtpResetpassword() {
	var bodyData dto.ResetUserPasswordOtp
	if err := c.ParseForm(&bodyData); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), 0)
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &bodyData)

	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&bodyData); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		logger.InsertAuditLogs(c.Ctx, "Error : "+"Validation error", 0)
		return
	}
	output, err := models.GetUserByEmail(bodyData.Username)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), 0)
		return
	}
	if output.Email == "" {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "login"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+logger.LogMessage(c.Ctx, "error.login"), 0)
		return
	}
	data, err := models.GetEmailOTP(bodyData.Username, bodyData.Otp)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), 0)
		return
	}
	if data.OtpCode != bodyData.Otp {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "otp"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+logger.LogMessage(c.Ctx, "error.otp"), 0)
		return
	}
	err = models.UpdateIsVerified(data.UserId)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), 0)
		return
	}
	uppass, err := models.ResetPassword(bodyData.NewPass, output.UserId)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), 0)
		return
	}
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, uppass, helpers.TranslateMessage(c.Ctx, "success", "passwordreset"), "")
	logger.InsertAuditLogs(c.Ctx, logger.LogMessage(c.Ctx, "success.passwordreset"), 0)
}

// SearchUser ...
// @Title Search User
// @Desciption SearchUser
// @Param search formData strings true "enter search"
// @Param open_page formData int false "if you want to open specific page than give page number"
// @Param page_size formData int false "how much data you want to show at a time default it will give 10 records"
// @Param lang query string false "use en-US or hi-IN"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} string
// @Failure 403
// @router /secure/search [post]
func (c *UserController) SearchUser() {

	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))

	var bodyData dto.SearchRequest
	if err := c.ParseForm(&bodyData); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &bodyData)

	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&bodyData); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		logger.InsertAuditLogs(c.Ctx, "Error : "+"Validation error", userId)
		return
	}
	if len(bodyData.Search) < 3 {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "serchfield"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+logger.LogMessage(c.Ctx, "error.serchfield"), userId)
		return
	}
	user, pagination_data, err := models.SearchUser(bodyData.Search, bodyData.PageSize, bodyData.OpenPage)

	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "datanotfound"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
	}

	if pagination_data["pageOpen_error"] == 1 {
		current := pagination_data["current_page"]
		last := pagination_data["last_page"]
		message := fmt.Sprintf(helpers.TranslateMessage(c.Ctx, "error", "page"), current, last)
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, message)
		logger.InsertAuditLogs(c.Ctx, "Error :"+fmt.Sprintf(logger.LogMessage(c.Ctx, "error.page"), current, last), userId)
		return
	}

	if user != nil {
		section_message := "read"
		section := "success"
		message := helpers.TranslateMessage(c.Ctx, section, section_message)
		helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, user, message, pagination_data)
		logger.InsertAuditLogs(c.Ctx, logger.LogMessage(c.Ctx, "success.read"), userId)
		return
	} else {
		helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, user, helpers.TranslateMessage(c.Ctx, "success", "data"), nil)
		logger.InsertAuditLogs(c.Ctx, logger.LogMessage(c.Ctx, "success.data"), userId)
		return
	}

}
