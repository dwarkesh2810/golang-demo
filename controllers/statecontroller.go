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

	var search dto.PaginationReq
	if err := c.ParseForm(&search); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), 0)
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &search)
	result, pagination_data, err := models.FetchStateList(search.OpenPage, search.PageSize)

	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "datanotfound"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), 0)
		return
	}
	if pagination_data["pageOpen_error"] == 1 {
		current := pagination_data["current_page"]
		last := pagination_data["last_page"]
		message := fmt.Sprintf(helpers.TranslateMessage(c.Ctx, "error", "page"), current, last)
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, message)
		logger.InsertAuditLogs(c.Ctx, "Error :"+fmt.Sprintf(logger.LogMessage(c.Ctx, "error.page"), current, last), 0)
		return
	}
	section_message := "read"
	section := "success"
	message := helpers.TranslateMessage(c.Ctx, section, section_message)
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, result, message, pagination_data)
	logger.InsertAuditLogs(c.Ctx, logger.LogMessage(c.Ctx, "success.read"), 0)
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

	var search dto.CountryWiseState
	if err := c.ParseForm(&search); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), 0)
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &search)
	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&search); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		logger.InsertAuditLogs(c.Ctx, "Error : Validation error", 0)
		return
	}
	result, pagination_data, err := models.CountryWiseState(search.OpenPage, search.PageSize, search.CountryId)

	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "datanotfound"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), 0)
		return
	}

	if pagination_data["pageOpen_error"] == 1 {
		current := pagination_data["current_page"]
		last := pagination_data["last_page"]
		message := fmt.Sprintf(helpers.TranslateMessage(c.Ctx, "error", "page"), current, last)
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, message)
		logger.InsertAuditLogs(c.Ctx, "Error :"+fmt.Sprintf(logger.LogMessage(c.Ctx, "error.page"), current, last), 0)
		return
	}
	section_message := "read"
	section := "success"
	message := helpers.TranslateMessage(c.Ctx, section, section_message)
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, result, message, pagination_data)
	logger.InsertAuditLogs(c.Ctx, logger.LogMessage(c.Ctx, "success.read"), 0)
}

// Filter State
// @Title Filter states
// @Description Fetch Filter States
// @Param search formData string false "Search State"
// @Param open_page formData int false "if you want to open specific page than give page number"
// @Param page_size formData int false "how much data you want to show at a time default it will give 10 records"
// @Success 200 {object} object
// @Failure 403
// @router /search_state [post]
func (c *StateController) FilterStates() {

	var bodyData dto.SearchRequest
	if err := c.ParseForm(&bodyData); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), 0)
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &bodyData)

	if bodyData.Search != "" {
		if len(bodyData.Search) < 2 {
			helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "serchfield"))
			logger.InsertAuditLogs(c.Ctx, "Error : "+logger.LogMessage(c.Ctx, "error.serchfield"), 0)
			return
		}
	}
	search := helpers.CapitalizeWords(bodyData.Search)
	result, pagination_data, err := models.FilterStates(search, bodyData.OpenPage, bodyData.PageSize)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "datanotfound"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), 0)
		return
	}

	if pagination_data["pageOpen_error"] == 1 {
		current := pagination_data["current_page"]
		last := pagination_data["last_page"]
		message := fmt.Sprintf(helpers.TranslateMessage(c.Ctx, "error", "page"), current, last)
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, message)
		logger.InsertAuditLogs(c.Ctx, "Error :"+fmt.Sprintf(logger.LogMessage(c.Ctx, "error.page"), current, last), 0)
		return
	}
	if len(result) == 0 {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "searchnotfound"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+logger.LogMessage(c.Ctx, "error.searchnotfound"), 0)
		return
	}
	section_message := "read"
	section := "success"
	message := helpers.TranslateMessage(c.Ctx, section, section_message)
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, result, message, pagination_data)
	logger.InsertAuditLogs(c.Ctx, logger.LogMessage(c.Ctx, "success.read"), 0)
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

	Data, err := models.GetState(bodyData.Id)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, err.Error())
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), 0)
		return
	}
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, Data, helpers.TranslateMessage(c.Ctx, "success", "read"), "")
	logger.InsertAuditLogs(c.Ctx, logger.LogMessage(c.Ctx, "success.read"), 0)
}
