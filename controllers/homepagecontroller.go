package controllers

import (
	"encoding/json"
	"fmt"

	"strings"

	"github.com/dwarkesh2810/golang-demo/dto"
	"github.com/dwarkesh2810/golang-demo/helpers"
	"github.com/dwarkesh2810/golang-demo/models"
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
// @Param	data_type   formData 	string	false		"body for file"
// @Param	section   formData 	string	false		"body for file"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 200 {object} models.HomePagesSettingTable
// @Failure 403
// @router /register_settings [post]
func (u *HomeSettingController) RegisterSettings() {
	// logedIN := u.GetSession("user_login")
	// if logedIN == "" {
	// 	helpers.ApiFailedResponse(u.Ctx.ResponseWriter, "Please Login ")
	// 	return
	// }

	var settings dto.HomeSeetingInsert
	var filePath string

	if err := u.ParseForm(&settings); err != nil {
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "parsing"))
		return
	}

	json.Unmarshal(u.Ctx.Input.RequestBody, &settings)

	data_types := strings.ToUpper(settings.DataType)
	// uploadDir := os.Getenv("uploadHomePageImages")
	uploadDir := "uploads/Home/files/images"
	if data_types == "LOGO" {
		// uploadDir = os.Getenv("uploadHomePageLogos")
		uploadDir = "uploads/Home/files/logo"
	} else if data_types != "BANNER" {
		filePath = ""
	}
	if data_types == "LOGO" || data_types == "BANNER" {
		file, fileHeader, err := u.GetFile("setting_data")
		if err != nil {

			section_failed_msg := "file_failed"
			section := "home_page_setting_failed_message_section"
			message_failed := helpers.TranslateMessage(u.Ctx, section, section_failed_msg)

			helpers.ApiFailedResponse(u.Ctx.ResponseWriter, message_failed)
			return
		}

		filePath, err = helpers.UploadFile(file, fileHeader, uploadDir)
		if err != nil {
			helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "upload"))
			return
		}
	}

	tokenData := helpers.GetTokenClaims(u.Ctx)
	// log.Print(tokenData, "=================")
	userID := tokenData["User_id"]
	result, _ := models.RegisterSetting(settings, userID.(float64), filePath)
	if result != 0 {

		section_success_msg := "register"
		section := "home_page_setting_success_message_section"
		message_success := helpers.TranslateMessage(u.Ctx, section, section_success_msg)
		helpers.ApiSuccessResponse(u.Ctx.ResponseWriter, "", message_success, "")
		return
	}

	helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "db"))
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
func (u *HomeSettingController) UpdateSettings() {
	// logedIN := u.GetSession("user_login")
	// if logedIN == "" {
	// 	helpers.ApiFailedResponse(u.Ctx.ResponseWriter, "Please Login ")
	// 	return
	// }
	var settings dto.HomeSeetingUpdate
	var filePath string

	if err := u.ParseForm(&settings); err != nil {
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "parsing"))
		return
	}

	section_message := "update"
	section := "home_page_setting_success_message_section"
	message := helpers.TranslateMessage(u.Ctx, section, section_message)

	json.Unmarshal(u.Ctx.Input.RequestBody, &settings)
	data_types := strings.ToUpper(settings.DataType)

	// uploadDir := os.Getenv("uploadHomePageImages")
	uploadDir := "uploads/Home/files/images"

	if data_types == "LOGO" {
		// uploadDir = os.Getenv("uploadHomePageLogos")
		uploadDir = "uploads/Home/files/logo"

	} else if data_types != "BANNER" {
		filePath = ""
	}

	if data_types == "LOGO" || data_types == "BANNER" {
		file, fileHeader, err := u.GetFile("setting_data")
		if err != nil {
			helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "filenotfound"))
			return
		}

		filePath, err = helpers.UploadFile(file, fileHeader, uploadDir)
		if err != nil {
			helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "upload"))
			return
		}
	}

	tokenData := helpers.GetTokenClaims(u.Ctx)
	userID := tokenData["User_id"]
	result, _ := models.UpdateSetting(settings, filePath, userID.(float64))

	if result != 0 {
		helpers.ApiSuccessResponse(u.Ctx.ResponseWriter, "", message, "")
		return
	}

	helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "db"))
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
func (u *HomeSettingController) FetchSettings() {
	var search dto.HomeSeetingSearch
	if err := u.ParseForm(&search); err != nil {
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "parsing"))
		return
	}
	json.Unmarshal(u.Ctx.Input.RequestBody, &search)

	tableName := "home_pages_setting_table"
	query := `
	SELECT hpst.section, hpst.data_type, hpst.setting_data, hpst.created_date, hpst.updated_date,
	concat(umt.first_name,' ',umt.last_name) as created_by
	FROM home_pages_setting_table as hpst
	LEFT JOIN users as umt ON umt.user_id = hpst.created_by
	ORDER BY hpst.created_date DESC
	LIMIT ? OFFSET ?
`
	result, pagination_data, _ := models.FetchSettingPaginations(search.OpenPage, search.PageSize, tableName, query)
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
	helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "db"))
}

func (u *HomeSettingController) DeleteSetting() {

	// logedIN := u.GetSession("user_login")
	// if logedIN == "" {
	// 	helpers.ApiFailedResponse(u.Ctx.ResponseWriter, "Please Login ")
	// 	return
	// }

	section_message := "delete"
	section := "home_page_setting_success_message_section"
	message := helpers.TranslateMessage(u.Ctx, section, section_message)

	var home_settings dto.HomeSeetingDelete
	if err := u.ParseForm(&home_settings); err != nil {
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "parsing"))
		return
	}
	json.Unmarshal(u.Ctx.Input.RequestBody, &home_settings)
	result := models.HomePageSettingExistsDelete(home_settings)
	if result != 0 {
		helpers.ApiSuccessResponse(u.Ctx.ResponseWriter, "", message, "")
		return
	}

	helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "db"))
}

// ExportFile
// @Title After Login User Can Export File in Home Page settings
// @Description In this function after login user  can Export File in Home page settings
// @Param file_type  formData string true "Here only select file within [XLSX,CSV,PDF]"
// @Param limit  formData int true "How Much you want to export data Ex.10"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 200 {object} models.HomePagesSettingTable
// @Failure 403
// @router /export [post]
func (c *HomeSettingController) ExportFile() {
	// logedIN := c.GetSession("user_login")
	// if logedIN == "" {
	// 	helpers.ApiFailedResponse(c.Ctx.ResponseWriter, "Please Login ")
	// 	return
	// }
	var fileTypes dto.FileType
	if err := c.ParseForm(&fileTypes); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &fileTypes)

	create_file_type := strings.ToUpper(fileTypes.FileType)

	if create_file_type == "" {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "type"))
		return
	}

	if create_file_type == "XLSX" || create_file_type == "PDF" || create_file_type == "CSV" {
		res_data, _ := models.ExportData(fileTypes.Limit, fileTypes.SratingFrom)
		res_s, _ := helpers.TransformToKeyValuePairs(res_data)
		header := helpers.ExtractKeys(res_s)

		res_result, _ := helpers.CreateFile(res_s, header, "", "apps", create_file_type)
		if res_result == "" {
			helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "create"))
			return
		}
		helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, res_result, helpers.TranslateMessage(c.Ctx, "success", "filecreate"), "")
		return
	}

	helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "create"))
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
	// logedIN := c.GetSession("user_login")
	// if logedIN == "" {
	// 	helpers.ApiFailedResponse(c.Ctx.ResponseWriter, "Please Login ")
	// 	return
	// }
	file, fileHeader, err := c.GetFile("import_type")
	if err != nil {
		c.Ctx.WriteString("Error uploading file")
		return
	}

	uploadDir := "uploads/FILES/IMPORT"
	filePath, err := helpers.UploadFile(file, fileHeader, uploadDir)
	if err != nil {
		// fmt.Println("Error uploading file:", err)
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "create"))
		return
	}
	defer helpers.RemoveFileByPath(filePath)

	var allRows []map[string]interface{}

	switch {
	case strings.HasSuffix(filePath, ".xlsx"):
		allRows, err = helpers.ReadXLSXFile(filePath)
		if err != nil {
			// fmt.Println("Error reading XLSX file:", err)
			helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "notread"))
			return
		}
		result, update_id, _ := models.RegisterSettingBatchsss(dto.HomeSeetingInsert{}, 35, filePath, allRows)
		if (len(result) > 0 && len(update_id) == 0) || (len(result) > 0 && len(update_id) > 0) || (len(result) == 0 && len(update_id) > 0) {
			helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, "", helpers.TranslateMessage(c.Ctx, "sucess", "upload"), "")
			return
		}

		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "create"))

	case strings.HasSuffix(filePath, ".csv"):

		allRows, err = helpers.ReadCSVFile(filePath)

		if err != nil {
			// fmt.Println("Error reading CSV file:", err)
			helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "notread"))
			return
		}

		result, update_id, _ := models.RegisterSettingBatchsss(dto.HomeSeetingInsert{}, 100, filePath, allRows)
		if result != nil || update_id != nil {
			helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, "", helpers.TranslateMessage(c.Ctx, "success", "upload"), "")
			return
		}

	default:
		// fmt.Println("Unsupported file format")
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "create"))
		return
	}
}

func (u *HomeSettingController) FiltersFetchSettings() {

	var search dto.HomeSeetingSearchFilter
	if err := u.ParseForm(&search); err != nil {
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "parsing"))
		return
	}

	json.Unmarshal(u.Ctx.Input.RequestBody, &search)

	result, pagination_data, _ := models.FilterFetchSettingPagination(search.OpenPage, search.PageSize, search.SearchParam)
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
	helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "db"))
}
