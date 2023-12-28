package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/dwarkesh2810/golang-demo/dto"
	"github.com/dwarkesh2810/golang-demo/models"
	"github.com/dwarkesh2810/golang-demo/pkg/helpers"
	"github.com/dwarkesh2810/golang-demo/pkg/logger"
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
func (c *StateController) FetchStates() {
	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))
	var search dto.PaginationReq
	if err := c.ParseForm(&search); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &search)
	result, pagination_data, err := models.FetchStateList(search.OpenPage, search.PageSize)

	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "datanotfound"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
		return
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
		section_message := "read"
		section := "success"
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

// Country Wise State
// @Title Country wise states
// @Description Country wise states
// @Param country_id formData int true "Enter Country Id To Search State"
// @Param open_page formData int false "if you want to open specific page than give page number"
// @Param page_size formData int false "how much data you want to show at a time default it will give 10 records"
// @Success 200 {object} object
// @Failure 403
// @router /country_wise_state [post]
func (c *StateController) CountryWiseState() {
	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))
	var search dto.CountryWiseState
	if err := c.ParseForm(&search); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &search)
	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&search); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		logger.InsertAuditLogs(c.Ctx, "Error : Validation error", userId)
		return
	}
	result, pagination_data, err := models.CountryWiseState(search.OpenPage, search.PageSize, search.CountryId)

	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "datanotfound"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
		return
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
		section_message := "read"
		section := "success"
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

// Filter State
// @Title Filter states
// @Description Fetch Filter States
// @Param search formData string true "Search State"
// @Param open_page formData int false "if you want to open specific page than give page number"
// @Param page_size formData int false "how much data you want to show at a time default it will give 10 records"
// @Success 200 {object} object
// @Failure 403
// @router /search_state [post]
func (c *StateController) FilterStates() {
	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))
	var bodyData dto.SearchRequest
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

	if len(bodyData.Search) < 2 {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "search"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+logger.LogMessage(c.Ctx, "error.search"), userId)
		return
	}
	search := helpers.CapitalizeWords(bodyData.Search)
	result, pagination_data, err := models.FilterStates(search, bodyData.OpenPage, bodyData.PageSize)

	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "datanotfound"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
		return
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
		section_message := "read"
		section := "success"
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

// Get State ...
// @Title get car
// @Desciption Get all car
// @Param state_id formData int true "get perticuler states using state_id"
// @Success 201 {object} string
// @Failure 403
// @router /get_state [post]
func (c *StateController) GetState() {
	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))
	var bodyData dto.GetStateRequest
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

	Data, err := models.GetState(bodyData.Id)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, err.Error())
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
		return
	}
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, Data, helpers.TranslateMessage(c.Ctx, "success", "read"), "")
	logger.InsertAuditLogs(c.Ctx, logger.LogMessage(c.Ctx, "success.read"), userId)
}
