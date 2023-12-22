package controllers

import (
	"encoding/json"
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/dwarkesh2810/golang-demo/dto"
	"github.com/dwarkesh2810/golang-demo/helpers"
	"github.com/dwarkesh2810/golang-demo/models"
)

type StateController struct {
	beego.Controller
}

// FetchStates
// @Title Fetch Data Country
// @Description Fetch Data Country
// @Param open_page formData int false "if you want to open specific page than give page number"
// @Param page_size formData int false "how much data you want to show at a time default it will give 10 records"
// @Success 200 {object} object
// @Failure 403
// @router /list_states [post]
func (u *StateController) FetchStates() {
	var search dto.PaginationReq
	if err := u.ParseForm(&search); err != nil {
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "parsing"))
		return
	}
	json.Unmarshal(u.Ctx.Input.RequestBody, &search)
	result, pagination_data, _ := models.FetchStateList(search.OpenPage, search.PageSize)
	if pagination_data["pageOpen_error"] == 1 {
		current := pagination_data["current_page"]
		last := pagination_data["last_page"]
		message := fmt.Sprintf("PAGE NUMBER %d IS NOT EXISTS , LAST PAGE NUMBER IS %d", current, last)
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, message)
		return
	}
	if result != nil {
		section_message := "read"
		section := "success"
		message := helpers.TranslateMessage(u.Ctx, section, section_message)
		helpers.ApiSuccessResponse(u.Ctx.ResponseWriter, result, message, pagination_data)
		return
	}
	helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "datanotfound"))
}

// Country Wise State
// @Title Fetch Data Country
// @Description Fetch Data Country
// @Param country_id formData int true "Enter Country Id To Search State"
// @Param open_page formData int false "if you want to open specific page than give page number"
// @Param page_size formData int false "how much data you want to show at a time default it will give 10 records"
// @Success 200 {object} object
// @Failure 403
// @router /country_wise_state [post]
func (u *StateController) CountryWiseState() {
	var search dto.CountryWiseState
	if err := u.ParseForm(&search); err != nil {
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "parsing"))
		return
	}
	json.Unmarshal(u.Ctx.Input.RequestBody, &search)
	result, pagination_data, _ := models.CountryWiseState(search.OpenPage, search.PageSize, search.CountryId)
	if pagination_data["pageOpen_error"] == 1 {
		current := pagination_data["current_page"]
		last := pagination_data["last_page"]
		message := fmt.Sprintf("PAGE NUMBER %d IS NOT EXISTS , LAST PAGE NUMBER IS %d", current, last)
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, message)
		return
	}
	if result != nil {
		section_message := "read"
		section := "success"
		message := helpers.TranslateMessage(u.Ctx, section, section_message)
		helpers.ApiSuccessResponse(u.Ctx.ResponseWriter, result, message, pagination_data)
		return
	}
	helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "datanotfound"))
}

// Country Wise State
// @Title Fetch Data Country
// @Description Fetch Data Country
// @Param search formData string true "Search State"
// @Param open_page formData int false "if you want to open specific page than give page number"
// @Param page_size formData int false "how much data you want to show at a time default it will give 10 records"
// @Success 200 {object} object
// @Failure 403
// @router /search_state [post]
func (u *StateController) FilterStates() {
	var search dto.SearchRequest
	if err := u.ParseForm(&search); err != nil {
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "parsing"))
		return
	}
	json.Unmarshal(u.Ctx.Input.RequestBody, &search)
	if len(search.Search) < 2 {
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "serchfield"))
		return
	}
	result, pagination_data, _ := models.FilterStates(search.Search, search.OpenPage, search.PageSize)
	if pagination_data["pageOpen_error"] == 1 {
		current := pagination_data["current_page"]
		last := pagination_data["last_page"]
		message := fmt.Sprintf("PAGE NUMBER %d IS NOT EXISTS , LAST PAGE NUMBER IS %d", current, last)
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, message)
		return
	}
	if result != nil {
		section_message := "read"
		section := "success"
		message := helpers.TranslateMessage(u.Ctx, section, section_message)
		helpers.ApiSuccessResponse(u.Ctx.ResponseWriter, result, message, pagination_data)
		return
	}
	helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "datanotfound"))
}
