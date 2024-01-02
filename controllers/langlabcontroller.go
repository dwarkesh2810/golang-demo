package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/dwarkesh2810/golang-demo/conf"
	"github.com/dwarkesh2810/golang-demo/dto"
	"github.com/dwarkesh2810/golang-demo/models"
	"github.com/dwarkesh2810/golang-demo/pkg/helpers"
	"github.com/dwarkesh2810/golang-demo/pkg/logger"
	"github.com/dwarkesh2810/golang-demo/pkg/validations"
)

type LangLableController struct {
	beego.Controller
}

// InsertLanguageLables
// @Title After Login admin Can create language lable
// @Description In this function after login it will work
// @Param	lable_code      formData      string	      true		"lable code"
// @Param	section         formData      string	      true		"section like success or failed or errors"
// @Param	ENGlang_value   formData 	  string	      true		"here you pass original message value in english"
// @Param	lang_ini_code   formData 	  string	      true		"to use  for hindi [hi-IN],for gujarati [gu-IN]"
// @Param	otherlang_value formData 	  string	      true		"here you can pass ENGlanguage message  converted otherlanguage value"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 200 {object} object
// @Failure 403
// @router /create_lang_lable [post]
func (c *LangLableController) InsertLanguageLables() {
	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))
	var langLables dto.LanguageLableInsert
	if err := c.ParseForm(&langLables); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &langLables)

	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&langLables); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		return
	}

	result, lableCodes, errr := models.InsertUpdateLanugaeLables(langLables, int(userId))

	if result == 0 && errr == nil {
		lable_codesMessage := fmt.Sprintf(`[ %s ] Lable Code Already Exists , If English Language Value is not available For  [--- %s ---] lableCode  it will Insert OR Update At Same lableCode `, langLables.LableCodes, langLables.LableCodes)
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "language", lable_codesMessage))
		return
	}
	if result == 1 && errr == nil && lableCodes != "" {
		langlable := fmt.Sprintf(`[ %s ] Successfully Created Language Lable`, lableCodes)
		helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, result, helpers.TranslateMessage(c.Ctx, "success", langlable), "")
		return
	}
}

// InsertLanguageLables
// @Title Insert language lable
// @Description new langouge code
// @Param	lable_code      formData      string	      true		"lable code"
// @Param	section         formData      string	      true		"section like success or failed or errors"
// @Param	ENGlang_value   formData 	  string	      true		"here you pass original message value in english"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 200 {object} object
// @Failure 403
// @router /lang_lable_Insert [post]
func (c *LangLableController) InsertLanguageLablesUsingApi() {
	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))
	var langLables dto.LanguageLable
	if err := c.ParseForm(&langLables); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &langLables)

	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&langLables); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		return
	}
	err := models.IsLanguageLableExist(langLables.LableCodes, langLables.Sections)
	if err == nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, "Language Lable already used please insert new language lable and try again")
		return
	}
	lableCodes, err := models.InsertUpdateLanugaeLablesApi(langLables, int(userId))
	if err == nil && lableCodes != "" {
		langlable := fmt.Sprintf(`[ %s ] Successfully Created Language Lable`, lableCodes)
		helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, lableCodes, helpers.TranslateMessage(c.Ctx, "success", langlable), "")
		return
	}
	helpers.ApiFailedResponse(c.Ctx.ResponseWriter, err.Error())
}

// UpdateLanguageLables
// @Title After Login admin Can update language lable
// @Description  function after login it will work
// @Param	lable_code      formData      string	      true		"lable code"
// @Param	ENGlang_value   formData 	  string	      true		"here you pass original message value in english"
// @Param	lang_ini_code   formData 	  string	      true		"to use  for hindi [hi-IN],for gujarati [gu-IN]"
// @Param	otherlang_value formData 	  string	      true		"here you can pass ENGlanguage message  converted otherlanguage value"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 200 {object} object
// @Failure 403
// @router /update_lang_lable [post]
func (c *LangLableController) UpdateLanguageLables() {
	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))
	var langLables dto.LanguageLableUpdate
	if err := c.ParseForm(&langLables); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &langLables)
	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&langLables); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		return
	}

	result, _, errr := models.UpdateLanguageLables(langLables, int(userId))

	if result == 0 && errr == nil {
		lable_codesMessage := fmt.Sprintf(`[ %s ] Lable Code Not Exists`, langLables.LableCodes)
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "language", lable_codesMessage))
		return
	}
	if result == 1 && errr == nil {
		langlable := fmt.Sprintf(`[ %s ] Successfully Updated Language Lable`, langLables.LableCodes)
		helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, result, helpers.TranslateMessage(c.Ctx, "success", langlable), "")
		return
	}
}

// ------------------------------------ Update using Api-------------------------------------

// Update Language Lables Api
// @Title  update language lable
// @Description update lables
// @Param	lable_code      formData      string	      true		"lable code"
// @Param	section         formData      string	      true		"section like success or failed or errors"
// @Param	ENGlang_value   formData 	  string	      true		"here you pass original message value in english"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 200 {object} object
// @Failure 403
// @router /lang_lable_UpdateAPI [put]
func (c *LangLableController) UpdateLanguageLablesAPI() {
	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))
	var langLables dto.LanguageLable
	if err := c.ParseForm(&langLables); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &langLables)
	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&langLables); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		return
	}
	err := models.IsLanguageLableExist(langLables.LableCodes, langLables.Sections)
	if err != nil {
		lable_codesMessage := fmt.Sprintf(`[ %s ] Lable Code Not Exists`, langLables.LableCodes)
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "language", lable_codesMessage))
		return
	}
	result, err := models.UpdateLanguageLablesAPI(langLables, int(userId))
	if err == nil {
		langlable := fmt.Sprintf(`[ %s ] Successfully Updated Language Lable`, langLables.LableCodes)
		helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, result, helpers.TranslateMessage(c.Ctx, "success", langlable), "")
		return
	}
}

// -------------------------------------- end -------------------------------
func (c *LangLableController) FetchAllAndWriteInINIFiles() bool {
	langugeLables, _ := models.FetchAllLabels()
	languageLangLables, _ := models.FetchAllDefaultlables()

	langResult, _ := helpers.ConvertToMapSlice(langugeLables)
	helpers.CreateINIFiles(langResult)

	res, _ := helpers.ConvertToMapSlice(languageLangLables)
	helpers.CreateINIFiles(res)
	return true
}

// ImportLanguageLables
// @Title After Login admin Can import English language lable
// @Description   after login it will work
// @Param	lables_ini_file      formData      file	      true		"select language lables ini file for import"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 200 {object} object
// @Failure 403
// @router /import_language_lables [post]
func (c *LangLableController) ImportLanguageLables() {
	file, fileHeader, err := c.GetFile("lables_ini_file")
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, "FILE PARSING ERROR")
		return
	}

	languagCode := helpers.ExtractLanguageCode(fileHeader.Filename)
	uploadDir := conf.Env.BaseUploadPath + "FILES/INI/IMPORT"
	filePath, err := helpers.UploadFile(file, fileHeader, uploadDir)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "create"))
		return
	}
	dataMap, err := helpers.ParseINIFile(filePath)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, err)
		return
	}

	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))

	dataMapResult := helpers.ConvertToDataMap(dataMap)
	result, errs := models.ImportINIFiles(languagCode, userId, dataMapResult)
	if result != "" {
		if result != "en-us" {
			helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, "", "successfully imported ini file data in multilanguage lable table", "")
			return
		}
		helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, "", "successfully imported ini file data in EnglishLanguage lable table", "")
		return
	}
	message := fmt.Sprintf("INI FILE NOT IMPORT DUE TO %s", errs)
	helpers.ApiFailedResponse(c.Ctx.ResponseWriter, message)
}

// ExportFileLables
// @Title After Login User Can Export File for language lables
// @Description In this function after login user can be export language lables i
// @Param file_type  formData string true "Here only select file within [XLSX,CSV,PDF]"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Param lang query string false "use en-US, hi-IN, gu-IN, mr-IN"
// @Success 200 {object} object
// @Failure 403
// @router /export_language_lables [post]
func (c *LangLableController) ExportFileLables() {
	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))

	var fileTypes dto.LanguageLablesFileType
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
		res_data, _ := models.ExportLanguageLables()
		res_s, _ := helpers.TransformToKeyValuePairs(res_data)
		header := helpers.ExtractKeys(res_s)
		uploadedPath := conf.Env.BaseUploadPath + "LANGUAGE/FILES/XLSX"
		if create_file_type == "PDF" || create_file_type == "CSV" {
			uploadedPath = conf.Env.BaseUploadPath + "LANGUAGE/FILES/PDF"
			if create_file_type == "CSV" {
				uploadedPath = conf.Env.BaseUploadPath + "LANGUAGE/FILES/PDF"
			}
		}

		res_result, err := helpers.CreateFile(res_s, header, uploadedPath, "lang", create_file_type)

		if err != nil {
			helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "datanotfound"))
			logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
			return
		}

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
