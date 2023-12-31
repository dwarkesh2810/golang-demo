package middleware

import (
	"strings"

	"github.com/beego/beego/v2/server/web/context"
	"github.com/dgrijalva/jwt-go"
	"github.com/dwarkesh2810/golang-demo/conf"
	"github.com/dwarkesh2810/golang-demo/models"
	"github.com/dwarkesh2810/golang-demo/pkg/helpers"
)

var key = conf.Env.JwtSecret
var jwtKey = []byte(key)

func JWTMiddleware(ctx *context.Context) {
	tokenString := ctx.Input.Header("Authorization")
	if tokenString == "" {
		helpers.ApiFailedResponse(ctx.ResponseWriter, helpers.TranslateMessage(ctx, "error", "token"))
		return
	}
	bearer := ContainsBearer(tokenString)
	if bearer {
		tokenString = tokenString[7:]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		helpers.ApiFailedResponse(ctx.ResponseWriter, helpers.TranslateMessage(ctx, "error", "token"))
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	userId := claims["ID"].(float64)
	data, err := models.GetUserDetails(userId)
	if err != nil {
		helpers.ApiFailedResponse(ctx.ResponseWriter, helpers.TranslateMessage(ctx, "error", "token"))
		return
	}
	if data.UserId == 0 {
		helpers.ApiFailedResponse(ctx.ResponseWriter, helpers.TranslateMessage(ctx, "error", "token"))
		return
	}
	ctx.Input.SetData("user", token.Claims.(jwt.MapClaims))
}

func ContainsBearer(token string) bool {
	// Convert the token to lowercase to make the comparison case-insensitive
	lowerToken := strings.ToLower(token)
	// Check if the token starts with "bearer "
	return strings.HasPrefix(lowerToken, "bearer ")
}
