package controllers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/dwarkesh2810/golang-demo/conf"
	"github.com/dwarkesh2810/golang-demo/dto"
	"github.com/dwarkesh2810/golang-demo/models"
	"github.com/dwarkesh2810/golang-demo/pkg/helpers"
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
// @Param   Authorization   header        string        true        "Bearer YourAccessToken"
// @Success 200 {object} object
// @Failure 403
// @router /create_lang_lable [post]
func (u *LangLableController) InsertLanguageLables() {
	var langLables dto.LanguageLableInsert
	if err := u.ParseForm(&langLables); err != nil {
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "parsing"))
		return
	}

	json.Unmarshal(u.Ctx.Input.RequestBody, &langLables)

	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&langLables); !isValid {
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, validations.ValidationErrorResponse(u.Controller, valid.Errors))
		return
	}

	result, lableCodes, errr := models.InsertUpdateLanugaeLables(langLables)

	if result == 0 && errr == nil {
		lable_codesMessage := fmt.Sprintf(`[ %s ] Lable Code Already Exists , If English Language Value is not available For  [--- %s ---] lableCode  it will Insert OR Update At Same lableCode `, langLables.LableCodes, langLables.LableCodes)
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "language", lable_codesMessage))
		return
	}
	if result == 1 && errr == nil && lableCodes != "" {
		langlable := fmt.Sprintf(`[ %s ] Successfully Created Language Lable`, lableCodes)
		helpers.ApiSuccessResponse(u.Ctx.ResponseWriter, result, helpers.TranslateMessage(u.Ctx, "success", langlable), "")
		return
	}
}

// InsertLanguageLables
// @Title Insert language lable
// @Description new langouge code
// @Param	lable_code      formData      string	      true		"lable code"
// @Param	section         formData      string	      true		"section like success or failed or errors"
// @Param	ENGlang_value   formData 	  string	      true		"here you pass original message value in english"
// @Success 200 {object} object
// @Failure 403
// @router /lang_lable_Insert [post]
func (c *LangLableController) InsertLanguageLablesUsingApi() {
	var langLables dto.LanguageLableInsertNew
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
	lableCodes, err := models.InsertUpdateLanugaeLablesApi(langLables)
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
// @Param   Authorization   header        string        true        "Bearer YourAccessToken"
// @Success 200 {object} object
// @Failure 403
// @router /update_lang_lable [post]
func (u *LangLableController) UpdateLanguageLables() {
	var langLables dto.LanguageLableUpdate
	if err := u.ParseForm(&langLables); err != nil {
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "error", "parsing"))
		return
	}

	json.Unmarshal(u.Ctx.Input.RequestBody, &langLables)
	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&langLables); !isValid {
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, validations.ValidationErrorResponse(u.Controller, valid.Errors))
		return
	}

	result, _, errr := models.UpdateLanguageLables(langLables)

	if result == 0 && errr == nil {
		lable_codesMessage := fmt.Sprintf(`[ %s ] Lable Code Not Exists`, langLables.LableCodes)
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, helpers.TranslateMessage(u.Ctx, "language", lable_codesMessage))
		return
	}
	if result == 1 && errr == nil {
		langlable := fmt.Sprintf(`[ %s ] Successfully Updated Language Lable`, langLables.LableCodes)
		helpers.ApiSuccessResponse(u.Ctx.ResponseWriter, result, helpers.TranslateMessage(u.Ctx, "success", langlable), "")
		return
	}
}

func (c *LangLableController) FetchAllAndWriteInINIFiles() bool {
	langugeLables, _ := models.FetchAllLabels()
	languageLangLables, _ := models.FetchAllDefaultlables()

	langResult, _ := helpers.ConvertToMapSlice(langugeLables)
	helpers.CreateINIFiles(langResult)

	res, _ := helpers.ConvertToMapSlice(languageLangLables)
	helpers.CreateINIFiles(res)
	return true
}

// ReadIniFile
// @Title After Login admin Can import language lable
// @Description   after login it will work
// @Param	getIniFile      formData      file	      true		"lable code"
// @Param   Authorization   header        string        true        "Bearer YourAccessToken"
// @Success 200 {object} object
// @Failure 403
// @router /import_language_lables [post]
func (c *LangLableController) ReadIniFile() {

	file, fileHeader, err := c.GetFile("getIniFile")
	if err != nil {
		c.Ctx.WriteString("Error uploading file")
		return
	}

	languagCode := helpers.ExtractLanguageCode(fileHeader.Filename)

	// ok := validations.ImportValidFileType(fileHeader.Filename)
	// if !ok {
	// 	helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "validation", "ValidFile"))
	// 	return
	// }

	uploadDir := conf.ConfigMaps["basepath"] + "FILES/INI/IMPORT"
	filePath, err := helpers.UploadFile(file, fileHeader, uploadDir)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "create"))
		return
	}
	dataMap, err := helpers.ParseINIFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	dataMapResult := helpers.ConvertToDataMap(dataMap)
	result := models.ProcessMapData(languagCode, dataMapResult)
	if result == 1 {

		helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, "", "successfully imported ini file data in multilanguage lanble table", "")
		return
	}
	helpers.ApiFailedResponse(c.Ctx.ResponseWriter, "INI FILE NOT IMPORT PLEASE TRY AGAIN")

}
