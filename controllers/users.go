package controllers

import (
	"encoding/json"
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
	_ = c.Ctx.Input.GetData("Lang").(string)
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
