package controllers

import (
	"encoding/json"
	"fmt"
	"log"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/dwarkesh2810/golang-demo/dto"
	"github.com/dwarkesh2810/golang-demo/models"
	"github.com/dwarkesh2810/golang-demo/pkg/helpers"
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

// FetchStates
// @Title Fetch Data Country
// @Description Fetch Data Country
// @Param open_page formData int false "if you want to open specific page than give page number"
// @Param page_size formData int false "how much data you want to show at a time default it will give 10 records"
// @Success 200 {object} object
// @Failure 403
// @router /list_states [post]
func (u *CountryController) FetchStates() {
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
func (u *CountryController) CountryWiseState() {
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
func (u *CountryController) FilterStates() {
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
