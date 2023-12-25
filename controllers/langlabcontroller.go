package controllers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/dwarkesh2810/golang-demo/dto"
	"github.com/dwarkesh2810/golang-demo/models"
	"github.com/dwarkesh2810/golang-demo/pkg/helpers"
	"github.com/dwarkesh2810/golang-demo/pkg/validations"
)

type LangLableController struct {
	beego.Controller
}

func (c *LangLableController) InsertLanguageLables() {
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

	exists_code := models.ExistsLanguageLable(langLables.LableCode)
	if exists_code == 1 {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, fmt.Sprintf(helpers.TranslateMessage(c.Ctx, "error", "langlbl"), langLables.LableCode))
		return
	}
	result, _ := models.InsertLanguageLabels(langLables)
	if result != "" {
		helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, result, helpers.TranslateMessage(c.Ctx, "success", "langlbl"), "")
		return
	}

	helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "failedlanglbl"))
}

func (c *LangLableController) FetchAllAndWriteInINIFiles() bool {
	langugeLables, _ := models.FetchAllLabels()
	languageLangLables, _ := models.FetchAllDefaultlables()
	langLangLables, done := languageLangLables.([]orm.Params)
	if !done {
		log.Fatal("Failed to convert 'results' to []orm.Params")
		return false
	}
	langResult, _ := helpers.ConvertToMapSlice(langLangLables)
	helpers.CreateINIFiles(langResult)
	ormParams, ok := langugeLables.([]orm.Params)
	if !ok {
		log.Fatal("Failed to convert 'results' to []orm.Params")
		return false
	}
	res, _ := helpers.ConvertToMapSlice(ormParams)
	helpers.CreateINIFiles(res)
	return true
}
