package controllers

import (
	"encoding/json"
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/dwarkesh2810/golang-demo/dto"
	"github.com/dwarkesh2810/golang-demo/helpers"
	"github.com/dwarkesh2810/golang-demo/models"
)

type CountryController struct {
	beego.Controller
}

// FetchCountries
// @Title After Login User Can Fetch Data Country
// @Description In this function after login user  can Fetch Data Country
// @Param open_page formData int false "if you want to open specific page than give page number"
// @Param page_size formData int false "how much data you want to show at a time default it will give 10 records"
// @Success 200 {object} object
// @Failure 403
// @router /list_countries [post]
func (u *CountryController) FetchCountries() {
	var search dto.CountrySearch
	if err := u.ParseForm(&search); err != nil {
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "parsing"))
		return
	}
	json.Unmarshal(u.Ctx.Input.RequestBody, &search)

	result, pagination_data, _ := models.FetchCountriesList(search.OpenPage, search.PageSize)

	if pagination_data["pageOpen_error"] == 1 {
		current := pagination_data["current_page"]
		last := pagination_data["last_page"]
		message := fmt.Sprintf(helpers.TranslateMessage(u.Ctx, "error", "page"), current, last)
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, message)
		return
	}

	if result != nil {
		section_message := "found"
		section := "country_message"
		message := helpers.TranslateMessage(u.Ctx, section, section_message)

		helpers.ApiSuccessResponse(u.Ctx.ResponseWriter, result, message, pagination_data)
		return
	}
	helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "db"))
}
