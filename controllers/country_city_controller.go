package controllers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/dwarkesh2810/golang-demo/dto"
	"github.com/dwarkesh2810/golang-demo/models"
	"github.com/dwarkesh2810/golang-demo/pkg/helpers"
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
func (u *CountryController) FetchCountries() {
	var search dto.PaginationReq
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

// FilterCountries
// @Title Filter Countries
// @Description it give country_name and country_id
// @Param open_page formData int false "if you want to open specific page than give page number"
// @Param page_size formData int false "how much data you want to show at a time default it will give 10 records"
// @Param search_param formData string false "it filter in database and give match"
// @Success 200 {object} object
// @Failure 403
// @router /filter_country [post]
func (u *CountryController) FilterCountries() {
	var search dto.CountrySearch
	if err := u.ParseForm(&search); err != nil {
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, "Parsing Data Error")
		return
	}

	json.Unmarshal(u.Ctx.Input.RequestBody, &search)
	searchFields := map[string]string{
		"country_name": search.SearchParam,
	}

	result, pagination_data, _ := models.CountryFilter(search.OpenPage, search.PageSize, "start", searchFields)
	log.Print(pagination_data)
	if result == nil && pagination_data["matchCount"] == 0 {
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, "Search Country Not Found")
		return
	}

	if pagination_data["pageOpen_error"] == 1 {
		current := pagination_data["current_page"]
		last := pagination_data["last_page"]
		message := fmt.Sprintf("PAGE NUMBER %d IS NOT EXISTS , LAST PAGE NUMBER IS %d", current, last)
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, message)
		return
	}

	if result != nil {
		section_message := "found"
		section := "home_page_setting_success_message_section"
		message := helpers.TranslateMessage(u.Ctx, section, section_message)
		helpers.ApiSuccessResponse(u.Ctx.ResponseWriter, result, message, pagination_data)
		return
	}

	helpers.ApiFailedResponse(u.Ctx.ResponseWriter, "Not Found Data Please Try Again")
}

// filter_city
// @Title Filter Countries
// @Description it give country_name and country_id
// @Param open_page formData int false "if you want to open specific page than give page number"
// @Param page_size formData int false "how much data you want to show at a time default it will give 10 records"
// @Param state_id formData int false "how much data you want to show at a time default it will give 10 records"
// @Param country_id formData int false "how much data you want to show at a time default it will give 10 records"
// @Param search_param formData string false "it filter in database and give match"
// @Success 200 {object} object
// @Failure 403
// @router /filter_city [post]
func (u *CountryController) FilterCity() {
	var search dto.CitySearch
	if err := u.ParseForm(&search); err != nil {
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, "Parsing Data Error")
		return
	}
	json.Unmarshal(u.Ctx.Input.RequestBody, &search)
	searchFields := map[string]string{
		"city_name": search.SearchParam,
	}

	otherFieldSCount := 0
	if search.CountryId > 0 && search.StateId > 0 {
		otherFieldSCount = 2
	}

	result, pagination_data, _ := models.CityFilter(search.OpenPage, search.PageSize, search.CountryId, search.StateId, otherFieldSCount, "start", searchFields)
	if result == nil && pagination_data["matchCount"] == 0 {
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, "Search City Not Found")
		return
	}

	if pagination_data["pageOpen_error"] == 1 {
		current := pagination_data["current_page"]
		last := pagination_data["last_page"]
		message := fmt.Sprintf("PAGE NUMBER %d IS NOT EXISTS , LAST PAGE NUMBER IS %d", current, last)
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, message)
		return
	}

	if result != nil {
		section_message := "found"
		section := "home_page_setting_success_message_section"
		message := helpers.TranslateMessage(u.Ctx, section, section_message)
		helpers.ApiSuccessResponse(u.Ctx.ResponseWriter, result, message, pagination_data)
		return
	}
	helpers.ApiFailedResponse(u.Ctx.ResponseWriter, "Not Found Data Please Try Again")

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
func (u *CountryController) FilterCountry() {
	var bodyData dto.SearchRequest
	if err := u.ParseForm(&bodyData); err != nil {
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "parsing"))
		return
	}
	json.Unmarshal(u.Ctx.Input.RequestBody, &bodyData)
	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&bodyData); !isValid {
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, validations.ValidationErrorResponse(u.Controller, valid.Errors))
		return
	}
	search := helpers.CapitalizeWords(bodyData.Search)
	result, pagination_data, _ := models.FilterCountries(search, bodyData.OpenPage, bodyData.PageSize)
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

// Get country
// @Title get country
// @Desciption Get country
// @Param country_id formData int true "get perticuler states using country_id"
// @Success 201 {object} string
// @Failure 403
// @router /get_country [post]
func (c *CountryController) GetCountry() {
	var bodyData dto.GetCountryRequest
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

	Data, err := models.GetCountry(bodyData.Id)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, err.Error())
	}
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, Data, helpers.TranslateMessage(c.Ctx, "success", "read"), "")
}