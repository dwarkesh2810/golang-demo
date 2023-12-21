package controllers

import (
	"encoding/json"
	"log"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/dwarkesh2810/golang-demo/dto"
	"github.com/dwarkesh2810/golang-demo/helpers"
	"github.com/dwarkesh2810/golang-demo/models"
	"github.com/dwarkesh2810/golang-demo/validations"
)

type LangLableController struct {
	beego.Controller
}

func (u *LangLableController) InsertLanguageLables() {
	var langLables dto.LanguageLableInsert
	if err := u.ParseForm(&langLables); err != nil {
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, "Parsing Data Error")
		return
	}

	json.Unmarshal(u.Ctx.Input.RequestBody, &langLables)

	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&langLables); !isValid {
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, validations.ValidationErrorResponse(u.Controller, valid.Errors))
		return
	}

	exists_code := models.ExistsLanguageLable(langLables.LableCode)
	if exists_code == 1 {
		helpers.ApiFailedResponse(u.Ctx.ResponseWriter, `[`+langLables.LableCode+`] Already Exists Please Use Different Lable Code !`)
		return
	}
	result, _ := models.InsertLanguageLabels(langLables)
	if result != "" {
		helpers.ApiSuccessResponse(u.Ctx.ResponseWriter, result, " Successfully Created Language Lable", "")
		return
	}

	helpers.ApiFailedResponse(u.Ctx.ResponseWriter, "Not Create Please Try Again !")
}

func (u *LangLableController) FetchAllAndWriteInINIFiles() {
	langugeLables, _ := models.FetchAllLabels()
	languageLangLables, _ := models.FetchAllDefaultlables()
	langLangLables, done := languageLangLables.([]orm.Params)
	if !done {
		log.Fatal("Failed to convert 'results' to []orm.Params")
	}
	langResult, _ := helpers.ConvertToMapSlice(langLangLables)
	helpers.CreateINIFiles(langResult)
	ormParams, ok := langugeLables.([]orm.Params)
	if !ok {
		log.Fatal("Failed to convert 'results' to []orm.Params")
	}
	res, _ := helpers.ConvertToMapSlice(ormParams)
	helpers.CreateINIFiles(res)
}
