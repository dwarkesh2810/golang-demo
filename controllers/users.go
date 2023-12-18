package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/dgrijalva/jwt-go"
	"github.com/dwarkesh2810/golang-demo/dto"
	"github.com/dwarkesh2810/golang-demo/helpers"
	"github.com/dwarkesh2810/golang-demo/models"
)

type UserController struct {
	beego.Controller
}

var key, _ = beego.AppConfig.String("JWT_SEC_KEY")
var jwtKey = []byte(key)

// Login ...
// @Title login User
// @Desciption login
// @Param body body models.UserLoginRequest true "login User"
// @Param lang query string false "use en-US or hi-IN"
// @Param   Accept-Languages   header  string  false  "Bearer YourAccessToken"
// @Success 201 {object} object
// @Failure 403
// @router /login [post]
func (c *UserController) Login() {
	_ = c.Ctx.Input.GetData("lang")
	var user dto.UserLoginRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	HashPassWord, err := models.GetUserByEmail(user.Username)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, err.Error())
		return
	}
	err = helpers.VerifyHashedData(HashPassWord.Password, user.Password)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, err.Error())
		return
	}
	userData, err := models.LoginUser(user.Username, HashPassWord.Password)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, err.Error())
		return
	}
	if userData.Email == "" && userData.FirstName == "" {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, "Unauthorized")
		return
	}
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &dto.JwtClaim{Email: userData.Email, ID: int(userData.UserId), StandardClaims: jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
	}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, err.Error())
		return
	}
	data := map[string]interface{}{"User_Data": token.Claims, "Tokan": tokenString}
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, data, "Login successfully", "")
}

// RegisterNewUser ...
// @Title Insert New User
// @Desciption new users
// @Param lang query string false "use en-US or hi-IN"
// @Param body body models.NewUserRequest true "Insert New User"
// @Success 201 {object} models.Users
// @Failure 403
// @router /register [post]
func (c *UserController) RegisterNewUser() {
	_ = c.Ctx.Input.GetData("Lang").(string)
	var bodyData dto.NewUserRequest
	if err := c.ParseForm(&bodyData); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, err.Error())
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &bodyData)
	data, _ := models.GetUserByEmail(bodyData.Email)
	if data.Email == bodyData.Email {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, "email already exist")
		return
	}
	output, err := models.InsertNewUser(bodyData)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, err.Error())
		return
	}
	go models.VerifyEmail(output.Email, output.FirstName)
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, output, "user Register Successfullty", "")
}

// GetAll ...
// @Title Get All
// @Description get Users
// @Param lang query string false "use en-US or hi-IN"
// @Param page query int false "For pagination"
// @Param limit query int false "per page user"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 200 {object} models.Users
// @Failure 403
// @router /secure/users/ [get]
func (c *UserController) GetAllUsers() {
	var search dto.HomeSeetingSearch
	if err := c.ParseForm(&search); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, "Parsing Data Error")
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &search)

	tableName := "users"
	query := `SELECT u.first_name , u.last_name, u.email, u.phone_number
	FROM users as u
	ORDER BY u.user_id
	LIMIT ? OFFSET ?
	
`
	result, pagination_data, _ := models.FetchSettingPaginations(search.OpenPage, search.PageSize, tableName, query)
	if pagination_data["pageOpen_error"] == 1 {
		current := pagination_data["current_page"]
		last := pagination_data["last_page"]
		message := fmt.Sprintf("PAGE NUMBER %d IS NOT EXISTS , LAST PAGE NUMBER IS %d", current, last)
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, message)
		return
	}

	if result != nil {
		section_message := ""
		section := ""
		message := helpers.TranslateMessage(c.Ctx, section, section_message)
		helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, result, message, pagination_data)
		return
	}
	helpers.ApiFailedResponse(c.Ctx.ResponseWriter, "Not Found Data Please Try Again")
}

// VerifyEmailOTP ...
// @Title verify otp for email
// @Desciption otp verification for eamil
// @Param body body models.VerifyEmailOTPRequest true "otp verification for email"
// @Param lang query string false "use en-US or hi-IN"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} string
// @Failure 403
// @router /secure/verify_email_otp [post]
func (c *UserController) VerifyEmailOTP() {
	// lang := c.Ctx.Input.GetData("Lang").(string)
	var bodyData dto.VerifyEmailOTPRequest
	if err := c.ParseForm(&bodyData); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, err.Error())
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &bodyData)
	data, err := models.GetEmailOTP(bodyData.Username, bodyData.Otp)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, err.Error())
		return
	}
	if data.OtpCode != bodyData.Otp {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, "Please Enter valid OTP")
		return
	}
	err = models.UpdateIsVerified(data.UserId)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, err.Error())
		return
	}
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, "Verified", "user verified Successfullty", "")
}

// UpdateUser .
// @Title update User
// @Desciption update users
// @Param body body models.UpdateUserRequest true "update New User"
// @Param lang query string false "use en-US or hi-IN"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} models.Users
// @Failure 403
// @router /secure/update [put]
func (c *UserController) UpdateUser() {
	// lang := c.Ctx.Input.GetData("Lang").(string)
	var bodyData dto.UpdateUserRequest
	if err := c.ParseForm(&bodyData); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, "Parsing Data Error")
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &bodyData)
	data, err := models.GetUserDetails(bodyData.Id)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, err.Error())
		return
	}
	if bodyData.Email != data.Email {
		res, _ := models.GetUserByEmail(bodyData.Email)
		if res.Email == bodyData.Email {
			helpers.ApiFailedResponse(c.Ctx.ResponseWriter, "Email Already exist")
			return
		}
	}
	output, err := models.UpdateUser(bodyData)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, err.Error())
		return
	}
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, output, "Updated successfully", "")
}

// ResetPassword ...
// @Title Reset password
// @Desciption Reset password
// @Param body body models.ResetUserPassword true "reset password"
// @Param lang query string false "use en-US or hi-IN"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} models.Users
// @Failure 403
// @router /secure/reset_pass [post]
func (c *UserController) ResetPassword() {
	// lang := c.Ctx.Input.GetData("Lang").(string)
	claims := helpers.GetTokenClaims(c.Ctx)
	id := claims["User_id"].(float64)
	output, err := models.GetUserDetails(id)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, err.Error())
		return
	}
	var bodyData dto.ResetUserPassword
	if err := c.ParseForm(&bodyData); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, "Parsing Data Error")
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &bodyData)
	err = helpers.VerifyHashedData(output.Password, bodyData.CurrentPass)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, err.Error())
		return
	}
	if bodyData.ConfirmPass != bodyData.NewPass {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, "New Password and confirm password not match")
		return
	}
	uppass, err := models.ResetPassword(bodyData.NewPass, id)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, err.Error())
		return
	}
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, uppass, "Password Reset successFully", "")
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
	// lang := c.Ctx.Input.GetData("Lang").(string)
	idString := c.Ctx.Input.Params()
	id, err := strconv.Atoi(idString["1"])
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, err.Error())
		return
	}
	data, err := models.DeleteUser(id)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, err.Error())
		return
	}
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, data, "User Remove successFully", "")
}
