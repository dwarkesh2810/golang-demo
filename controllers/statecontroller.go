package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/dwarkesh2810/golang-demo/dto"
	"github.com/dwarkesh2810/golang-demo/models"
	"github.com/dwarkesh2810/golang-demo/pkg/helpers"
	"github.com/dwarkesh2810/golang-demo/pkg/validations"
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
// @Title Country wise states
// @Description Country wise states
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
	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&search); !isValid {
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, validations.ValidationErrorResponse(u.Controller, valid.Errors))
		return
	}
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

// Filter State
// @Title Filter states
// @Description Fetch Filter States
// @Param search formData string true "Search State"
// @Param open_page formData int false "if you want to open specific page than give page number"
// @Param page_size formData int false "how much data you want to show at a time default it will give 10 records"
// @Success 200 {object} object
// @Failure 403
// @router /search_state [post]
func (u *StateController) FilterStates() {
	var bodyData dto.SearchRequest
	if err := u.ParseForm(&bodyData); err != nil {
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "parsing"))
		return
	}
	json.Unmarshal(u.Ctx.Input.RequestBody, &bodyData)
	if len(bodyData.Search) < 2 {
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "serchfield"))
		return
	}
	search := helpers.CapitalizeWords(bodyData.Search)
	result, pagination_data, _ := models.FilterStates(search, bodyData.OpenPage, bodyData.PageSize)
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
	helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "searchnotfound"))
}

// Get State ...
// @Title get car
// @Desciption Get all car
// @Param state_id formData int true "get perticuler states using state_id"
// @Success 201 {object} string
// @Failure 403
// @router /get_state [post]
func (c *StateController) GetState() {
	var bodyData dto.GetStateRequest
	if err := c.ParseForm(&bodyData); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &bodyData)

	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&bodyData); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		return
	}

	Data, err := models.GetState(bodyData.Id)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, err.Error())
	}
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, Data, helpers.TranslateMessage(c.Ctx, "success", "read"), "")
}
