package controllers

import (
	"encoding/json"
	"fmt"

	"strings"

	"github.com/beego/beego/v2/core/validation"
	"github.com/dwarkesh2810/golang-demo/conf"
	"github.com/dwarkesh2810/golang-demo/dto"
	"github.com/dwarkesh2810/golang-demo/models"
	"github.com/dwarkesh2810/golang-demo/pkg/helpers"
	"github.com/dwarkesh2810/golang-demo/pkg/logger"
	"github.com/dwarkesh2810/golang-demo/pkg/validations"
	_ "github.com/lib/pq"

	beego "github.com/beego/beego/v2/server/web"
)

type HomeSettingController struct {
	beego.Controller
}

// RegisterSettings
// @Title After Login User Can Register Home Page settings
// @Description In this function after login can register Home page settings
// @Param	setting_data   formData 	file	false		"body for file"
// @Param	setting_data   formData 	string	false		"body for file"
// @Param	data_type   formData 	string	false		"body for html text or other "
// @Param	section   formData 	string	false		"body for file"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 200 {object} models.HomePagesSettingTable
// @Failure 403
// @router /register_settings [post]
func (c *HomeSettingController) RegisterSettings() {

	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))

	var settings dto.HomeSeetingInsert
	var filePath string

	if err := c.ParseForm(&settings); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &settings)

	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&settings); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		logger.InsertAuditLogs(c.Ctx, "Error : Validation error", userId)
		return
	}

	data_types := strings.ToUpper(settings.DataType)
	uploadDir := conf.ConfigMaps["basepath"] + "Home/files/images"
	if data_types == "LOGO" {
		uploadDir = conf.ConfigMaps["basepath"] + "Home/files/logo"
	} else if data_types != "BANNER" {
		filePath = ""
	}

	if data_types == "LOGO" || data_types == "BANNER" {
		file, fileHeader, err := c.GetFile("setting_data")

		ok := validations.ValidImageType(fileHeader.Filename)
		if !ok {
			helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "validation", "ValidImage"))
			logger.InsertAuditLogs(c.Ctx, "Error :"+logger.LogMessage(c.Ctx, "validation.ValidImage"), userId)
			return
		}
		if err != nil {

			section_failed_msg := "file_failed"
			section := "home_page_setting_failed_message_section"
			message_failed := helpers.TranslateMessage(c.Ctx, section, section_failed_msg)

			helpers.ApiFailedResponse(c.Ctx.ResponseWriter, message_failed)
			logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
			return
		}

		filePath, err = helpers.UploadFile(file, fileHeader, uploadDir)
		if err != nil {
			helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "upload"))
			logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
			return
		}
	}

	tokenData := helpers.GetTokenClaims(c.Ctx)
	userID := tokenData["User_id"]
	result, err := models.RegisterSetting(settings, userID.(float64), filePath)
	if result != 0 {
		section_success_msg := "create"
		section := "success"
		message_success := helpers.TranslateMessage(c.Ctx, section, section_success_msg)
		helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, "", message_success, "")
		logger.InsertAuditLogs(c.Ctx, "Error :"+logger.LogMessage(c.Ctx, "success.create"), userId)
		return
	}

	helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
	logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
}

// UpdateSettings
// @Title After Login User Can Update Home Page settings
// @Description In this function after login user  can update Home page settings
// @Param	setting_data   formData 	file	false		"body for file"
// @Param	data_type   formData 	string	false		"body for file"
// @Param	section   formData 	string	false		"body for file"
// @Param	setting_id   formData 	int		false		"body for file"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 200 {object} models.HomePagesSettingTable
// @Failure 403
// @router /update_settings [post]
func (c *HomeSettingController) UpdateSettings() {

	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))

	var settings dto.HomeSeetingUpdate
	var filePath string

	if err := c.ParseForm(&settings); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
		return
	}

	section_message := "update"
	section := "success"
	message := helpers.TranslateMessage(c.Ctx, section, section_message)

	json.Unmarshal(c.Ctx.Input.RequestBody, &settings)

	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&settings); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		logger.InsertAuditLogs(c.Ctx, "Error : Validation error", userId)
		return
	}

	data_types := strings.ToUpper(settings.DataType)

	uploadDir := conf.ConfigMaps["basepath"] + "Home/files/images"

	if data_types == "LOGO" {
		uploadDir = conf.ConfigMaps["basepath"] + "Home/files/logo"

	} else if data_types != "BANNER" {
		filePath = ""
	}

	if data_types == "LOGO" || data_types == "BANNER" {
		file, fileHeader, err := c.GetFile("setting_data")

		ok := validations.ValidImageType(fileHeader.Filename)
		if !ok {
			helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "validation", "ValidImage"))
			logger.InsertAuditLogs(c.Ctx, "Error :"+logger.LogMessage(c.Ctx, "validation.ValidImage"), userId)
			return
		}
		if err != nil {
			helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "filenotfound"))
			logger.InsertAuditLogs(c.Ctx, "Error :"+logger.LogMessage(c.Ctx, "error.filenotfound"), userId)
			return
		}

		filePath, err = helpers.UploadFile(file, fileHeader, uploadDir)
		if err != nil {
			helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "upload"))
			logger.InsertAuditLogs(c.Ctx, "Error :"+logger.LogMessage(c.Ctx, "error.upload"), userId)
			return
		}
	}

	tokenData := helpers.GetTokenClaims(c.Ctx)
	userID := tokenData["User_id"]
	result, err := models.UpdateSetting(settings, filePath, userID.(float64))

	if result != 0 {
		helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, "", message, "")
		logger.InsertAuditLogs(c.Ctx, logger.LogMessage(c.Ctx, "success.update"), userId)
		return
	}

	helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
	logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
}

// FetchSettings
// @Title After Login User Can Fetch Data Home Page settings
// @Description In this function after login user  can Fetch Data Home page settings
// @Param open_page formData int false "if you want to open specific page than give page number"
// @Param page_size formData int false "how much data you want to show at a time default it will give 10 records"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 200 {object} models.HomePagesSettingTable
// @Failure 403
// @router /fetch_settings [post]
func (c *HomeSettingController) FetchSettings() {

	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))

	var search dto.HomeSeetingSearch
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

	result, pagination_data, err := models.FetchSettingPaginations(search.OpenPage, search.PageSize)

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
	}
	helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
	logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
}

// @Title DeleteSetting
// @Description delete Setting From Database by setting_id
// @Param setting_id formData int true "User can delete setting after login by setting_id"
// @Param lang query string false "use en-US or hi-IN"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 200 {string} string
// @Failure 403
// @router /delete_settings [post]
func (c *HomeSettingController) DeleteSetting() {

	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))

	section_message := "delete"
	section := "success"
	message := helpers.TranslateMessage(c.Ctx, section, section_message)

	var home_settings dto.HomeSeetingDelete
	if err := c.ParseForm(&home_settings); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &home_settings)

	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&home_settings); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		logger.InsertAuditLogs(c.Ctx, "Error : Validation error", userId)
		return
	}

	result := models.HomePageSettingExistsDelete(home_settings)
	if result == nil {
		helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, "", message, "")
		logger.InsertAuditLogs(c.Ctx, logger.LogMessage(c.Ctx, "success.delete"), userId)
		return
	}

	helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
	logger.InsertAuditLogs(c.Ctx, "Error :"+result.Error(), userId)
}

// ExportFile
// @Title After Login User Can Export File in Home Page settings
// @Description In this function after login user  can Export File in Home page settings
// @Param file_type  formData string true "Here only select file within [XLSX,CSV,PDF]"
// @Param starting_from  formData int true "you want data from row first or id 1 to 20 so you can pass starting_from as 1"
// @Param limit  formData int true "How Much you want to export data Ex.10"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 200 {object} models.HomePagesSettingTable
// @Failure 403
// @router /export [post]
func (c *HomeSettingController) ExportFile() {

	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))

	var fileTypes dto.FileType
	if err := c.ParseForm(&fileTypes); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &fileTypes)

	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&fileTypes); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		logger.InsertAuditLogs(c.Ctx, "Error : Validation error", userId)
		return
	}

	create_file_type := strings.ToUpper(fileTypes.FileType)

	if create_file_type == "" {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "type"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+logger.LogMessage(c.Ctx, "error.type"), userId)
		return
	}

	if create_file_type == "XLSX" || create_file_type == "PDF" || create_file_type == "CSV" {
		res_data, _ := models.ExportData(fileTypes.Limit, fileTypes.SratingFrom)
		res_s, _ := helpers.TransformToKeyValuePairs(res_data)
		header := helpers.ExtractKeys(res_s)

		res_result, _ := helpers.CreateFile(res_s, header, "", "apps", create_file_type)
		if res_result == "" {
			helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "create"))
			logger.InsertAuditLogs(c.Ctx, "Error :"+logger.LogMessage(c.Ctx, "error.create"), userId)
			return
		}
		helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, res_result, helpers.TranslateMessage(c.Ctx, "success", "filecreate"), "")
		logger.InsertAuditLogs(c.Ctx, "Error :"+logger.LogMessage(c.Ctx, "success.filecreate"), userId)
		return
	}

	helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "create"))
	logger.InsertAuditLogs(c.Ctx, "Error :"+logger.LogMessage(c.Ctx, "error.create"), userId)
}

// ImportFile
// @Title After Login User Can Import File in Home Page settings
// @Description In this function after login user  can Import File in Home page settings
// @Param import_type  formData file true "Here only select file within [XLSX,CSV]"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 200 {object} models.HomePagesSettingTable
// @Failure 403
// @router /import [post]
func (c *HomeSettingController) ImportFile() {

	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))

	file, fileHeader, err := c.GetFile("import_type")
	if err != nil {
		c.Ctx.WriteString("Error uploading file")
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
		return
	}

	ok := validations.ImportValidFileType(fileHeader.Filename)
	if !ok {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "validation", "ValidFile"))
		logger.InsertAuditLogs(c.Ctx, logger.LogMessage(c.Ctx, "validation.ValidFile"), userId)
		return
	}

	uploadDir := conf.ConfigMaps["basepath"] + "FILES/IMPORT"
	filePath, err := helpers.UploadFile(file, fileHeader, uploadDir)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "create"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+logger.LogMessage(c.Ctx, "error.crete"), userId)
		return
	}
	defer helpers.RemoveFileByPath(filePath)

	var allRows []map[string]interface{}

	switch {
	case strings.HasSuffix(filePath, ".xlsx"):
		allRows, err = helpers.ReadXLSXFile(filePath)
		if err != nil {
			helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "notread"))
			logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
			return
		}
		result, update_id, err := models.RegisterSettingBatchsss(dto.HomeSeetingInsert{}, float64(userId), filePath, allRows)
		if (len(result) > 0 && len(update_id) == 0) || (len(result) > 0 && len(update_id) > 0) || (len(result) == 0 && len(update_id) > 0) {
			helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, "", helpers.TranslateMessage(c.Ctx, "success", "upload"), "")
			logger.InsertAuditLogs(c.Ctx, logger.LogMessage(c.Ctx, "success.upload"), userId)
			return
		}

		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "create"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)

	case strings.HasSuffix(filePath, ".csv"):

		allRows, err = helpers.ReadCSVFile(filePath)

		if err != nil {

			helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "notread"))
			logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
			return
		}

		result, update_id, err := models.RegisterSettingBatchsss(dto.HomeSeetingInsert{}, 100, filePath, allRows)
		if result != nil || update_id != nil {
			helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, "", helpers.TranslateMessage(c.Ctx, "success", "upload"), "")
			logger.InsertAuditLogs(c.Ctx, logger.LogMessage(c.Ctx, "success.upload"), userId)
			return
		}
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "create"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)

	default:

		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "create"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+logger.LogMessage(c.Ctx, "error.create"), userId)
		return
	}
}

// FiltersFetchSettings
// @Title After Login User Can Filter Data Home Page settings
// @Description In this function after login user  can FilterData with pagination and Check Count Of Match Data from  Home page settings
// @Param open_page formData int false "if you want to open specific page than give page number"
// @Param page_size formData int false "how much data you want to show at a time default it will give 10 records"
// @Param setting_data formData string false "it filter in database and give match"
// @Param data_type formData string false "it filter in database and give match"
// @Param unique_code formData string false "it filter in database and give match"
// @Param section formData string false "it filter in database and give match"
// @Param apply_position formData string false "if you apply_position pass start than it will match record with starting of a string or if you  apply_position not pass it will search in perticular/allcolumns  all string"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 200 {object} models.HomePagesSettingTable
// @Failure 403
// @router /filter_hpst [post]
func (c *HomeSettingController) FiltersFetchSettings() {

	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))

	var search dto.HomeSeetingSearchFilter
	if err := c.ParseForm(&search); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &search)

	searchFields := map[string]string{
		"setting_data": search.SettingData,
		"section":      search.Section,
		"data_type":    search.DataType,
		"unique_code":  search.UniqueCode,
	}
	if search.Section != "" || search.SettingData != "" || search.DataType != "" || search.UniqueCode != "" {
		result, pagination_data, _ := models.FilterWithPaginationFetchSettings(search.OpenPage, search.PageSize, search.ApplySearchPosition, searchFields)
		if result == nil && pagination_data["matchCount"] == 0 {
			helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "searchnotfound"))
			logger.InsertAuditLogs(c.Ctx, "Error :"+logger.LogMessage(c.Ctx, "error.searchnotfound"), userId)
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
		}
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "datanotfound"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+logger.LogMessage(c.Ctx, "error.datanotfound"), userId)
	}
	helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "search"))
	logger.InsertAuditLogs(c.Ctx, "Error :"+logger.LogMessage(c.Ctx, "error.search"), userId)
}
