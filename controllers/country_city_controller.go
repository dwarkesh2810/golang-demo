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
func (c *CountryController) FetchCountries() {

	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))

	var search dto.PaginationReq
	if err := c.ParseForm(&search); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &search)

	result, pagination_data, err := models.FetchCountriesList(search.OpenPage, search.PageSize)

	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
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

// FilterCountries
// @Title Filter Countries
// @Description it give country_name and country_id
// @Param open_page formData int false "if you want to open specific page than give page number"
// @Param page_size formData int false "how much data you want to show at a time default it will give 10 records"
// @Param search_param formData string false "it filter in database and give match"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 200 {object} object
// @Failure 403
// @router /filter_country [post]
func (c *CountryController) FilterCountries() {

	// claims := helpers.GetTokenClaims(c.Ctx)
	// userId := uint(claims["User_id"].(float64))
	userId := uint(1)

	var search dto.CountrySearch
	if err := c.ParseForm(&search); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, "Parsing Data Error")
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &search)
	searchFields := map[string]string{
		"country_name": search.SearchParam,
	}

	result, pagination_data, err := models.CountryFilter(search.OpenPage, search.PageSize, "start", searchFields)

	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
		return
	}

	if result == nil && pagination_data["matchCount"] == 0 {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, "Search Country Not Found")
		logger.InsertAuditLogs(c.Ctx, "Error :"+logger.LogMessage(c.Ctx, "error.datanotfound"), userId)
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

// filter_city
// @Title Filter Cities
// @Description it give city_names and country_id
// @Param open_page formData int false "if you want to open specific page than give page number"
// @Param page_size formData int false "how much data you want to show at a time default it will give 10 records"
// @Param state_id formData int false "state_id if you provide country_id than use it"
// @Param country_id formData int false "country_id if you provide state_id than use it"
// @Param country_name formData string false "if you provide proper country name than it returns all cities of that country"
// @Param state_name formData string false "if you provide proper state  name than it returns all cities of that states"
// @Param search_param formData string false "it filter in database and give match"
// @Success 200 {object} object
// @Failure 403
// @router /filter_city [post]
func (c *CountryController) FilterCity() {

	// claims := helpers.GetTokenClaims(c.Ctx)
	// userId := uint(claims["User_id"].(float64))
	userId := uint(1)

	var search dto.CitySearch
	if err := c.ParseForm(&search); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, "Parsing Data Error")
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &search)
	searchFields := map[string]string{
		"city_name": search.SearchParam,
	}

	otherFieldSCount := 2
	if search.CountryId > 0 && search.StateId > 0 {
		otherFieldSCount = 2
	}

	result, pagination_data, err := models.CityFilter(search.OpenPage, search.PageSize, search.CountryId, search.StateId, otherFieldSCount, "start", search.CountryName, search.StateName, searchFields)

	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
		return
	}

	if result == nil && pagination_data["matchCount"] == 0 {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, "Search City Not Found")
		logger.InsertAuditLogs(c.Ctx, "Error :"+logger.LogMessage(c.Ctx, "error.datanotfound"), userId)
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

// Filter Countries
// @Title Fetch Data Country
// @Description Fetch Data Country
// @Param search formData string true "Search Country"
// @Param open_page formData int false "if you want to open specific page than give page number"
// @Param page_size formData int false "how much data you want to show at a time default it will give 10 records"
// @Success 200 {object} object
// @Failure 403
// @router /search_country [post]
func (c *CountryController) FilterCountry() {

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
	search := helpers.CapitalizeWords(bodyData.Search)
	result, pagination_data, err := models.FilterCountries(search, bodyData.OpenPage, bodyData.PageSize)

	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
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

// Get country
// @Title get country
// @Desciption Get country
// @Param country_id formData int true "get perticuler states using country_id"
// @Success 201 {object} string
// @Failure 403
// @router /get_country [post]
func (c *CountryController) GetCountry() {

	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))

	var bodyData dto.GetCountryRequest
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

	Data, err := models.GetCountry(bodyData.Id)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, err.Error())
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
		return
	}
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, Data, helpers.TranslateMessage(c.Ctx, "success", "read"), "")
	logger.InsertAuditLogs(c.Ctx, logger.LogMessage(c.Ctx, "success.read"), userId)

}
