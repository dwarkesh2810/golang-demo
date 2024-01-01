package controllers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/dwarkesh2810/golang-demo/conf"
	"github.com/dwarkesh2810/golang-demo/dto"
	"github.com/dwarkesh2810/golang-demo/models"
	"github.com/dwarkesh2810/golang-demo/pkg/helpers"
	"github.com/dwarkesh2810/golang-demo/pkg/logger"
	"github.com/dwarkesh2810/golang-demo/pkg/validations"
)

type CarController struct {
	beego.Controller
}

// AddNewCar ...
// @Title new car
// @Desciption insert car
// @swagger:parameters upload
// @Param car_name formData string true "Car name"
// @Param modified_by formData string true "modified by"
// @Param model formData string true "Car Model"
// @Param type formData string true "accepted type 'sedan','SUV','hatchback'"
// @Param file formData file true "File to be uploaded, Acepted Extension [jpeg, jpg, png, svg]"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} string
// @Failure 403
// @router /create [post]
func (c *CarController) AddNewCar() {
	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))
	var cars dto.GetNewCarRequest
	if err := c.ParseForm(&cars); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &cars)

	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&cars); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		logger.InsertAuditLogs(c.Ctx, "Error : Validation error", userId)
		return
	}

	file, fileHeader, err := c.GetFile("file")
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "filenotfound"))
		logger.InsertAuditLogs(c.Ctx, "Error :- "+err.Error(), userId)
		return
	}
	ok := validations.ValidImageType(fileHeader.Filename)
	if !ok {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "validation", "ValidImage"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+logger.LogMessage(c.Ctx, "validation.ValidImage"), userId)
		return
	}
	var carType string = string(cars.Type)
	cars.Type, err = helpers.NewCarType(carType)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "cartype"))
		logger.InsertAuditLogs(c.Ctx, "Error :- "+err.Error(), userId)
		return
	}

	uploadDir := conf.Env.BaseUploadPath + "Cars/images/"
	filepaths, err := helpers.UploadFile(file, fileHeader, uploadDir)

	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "upload"))
		logger.InsertAuditLogs(c.Ctx, "Error :- "+err.Error(), userId)
		return
	}
	cars.CreatedBy = int(userId)
	cars.CarImage = filepaths
	data, err := models.InsertNewCar(cars)

	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		logger.InsertAuditLogs(c.Ctx, "Error :- "+err.Error(), userId)
		return
	}
	logger.InsertAuditLogs(c.Ctx, fmt.Sprintf("Register new car by user : %v, new car id : %v", userId, data.Id), userId)
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, data, helpers.TranslateMessage(c.Ctx, "success", "create"), "")
}

// UpdateCar ...
// @Title update car
// @Desciption update car
// @Param car_id formData string true "Car name"
// @Param car_name formData string false "Car name"
// @Param modified_by formData string false "modified by"
// @Param model formData string false "Car Model"
// @Param type formData string false "accepted type 'sedan','SUV','hatchback'"
// @Param file formData file false "File to be uploaded, Acepted Extension [jpeg, jpg, png, svg]"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 200 {object} string
// @Failure 403
// @router /update [PUT]
func (c *CarController) UpdateCar() {
	var cars dto.UpdateCarRequest
	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))
	if err := c.ParseForm(&cars); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &cars)
	data, err := models.GetSingleCar(cars.Id)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
		return
	}
	if cars.CarName == "" {
		cars.CarName = data.CarName
	}
	if cars.ModifiedBy == "" {
		cars.ModifiedBy = data.ModifiedBy
	}
	if cars.Model == "" {
		cars.Model = data.Model
	}
	if cars.Type == "" {
		cars.Type = data.Type
	}
	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&cars); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		logger.InsertAuditLogs(c.Ctx, "Error : Validation error", userId)
		return
	}
	var carType string = string(cars.Type)
	cars.Type, err = helpers.NewCarType(carType)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "cartype"))
		logger.InsertAuditLogs(c.Ctx, "Error :- "+err.Error(), userId)
		return
	}
	file, fileheader, err := c.GetFile("file")	
	if err != nil {
		cars.UpdatedBy = int(userId)
		cars.CarImage = data.CarImage
		res, err := models.UpdateCar(cars)
		if err != nil {
			helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
			logger.InsertAuditLogs(c.Ctx, "Error :"+err.Error(), userId)
			return
		}
		helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, res, helpers.TranslateMessage(c.Ctx, "success", "update"), "")
		logger.InsertAuditLogs(c.Ctx, fmt.Sprintf("Update car by user : %v , car id : %v", userId, res.Id), userId)
		return
	}
	ok := validations.ValidImageType(fileheader.Filename)
	if !ok {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "validation", "ValidImage"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+logger.LogMessage(c.Ctx, "validation.ValidImage"), userId)
		return
	}
	uploadDir := conf.Env.BaseUploadPath + "Cars/images/"
	filepaths, err := helpers.UploadFile(file, fileheader, uploadDir)

	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "upload"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
		return
	}
	cars.UpdatedBy = int(userId)
	cars.CarImage = filepaths
	res, err := models.UpdateCar(cars)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
		return
	}
	err = os.Remove(data.CarImage)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "filenotfound"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
		return
	}
	logger.InsertAuditLogs(c.Ctx, fmt.Sprintf("Update car by user : %v , car id : %v", userId, res.Id), userId)
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, res, helpers.TranslateMessage(c.Ctx, "success", "update"), "")
}

// DeleteCar ...
// @Title remove car
// @Desciption delete car
// @Param body body dto.GetcarRequest true "delete car"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 200 {object} string
// @Failure 403
// @router /delete [delete]
func (c *CarController) DeleteCar() {
	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))
	var car dto.GetcarRequest
	if err := c.ParseForm(&car); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &car)
	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&car); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		logger.InsertAuditLogs(c.Ctx, "Error : Validation error", userId)
		return
	}

	res, err := models.GetSingleCar(car.Id)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
		return
	}

	if res.CarName == "" {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "invalidid"))
		logger.InsertAuditLogs(c.Ctx, "Error :"+logger.LogMessage(c.Ctx, "error.invalidid"), userId)
		return
	}

	data, err := models.DeleteCar(car.Id)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
		return
	}

	err = os.Remove(res.CarImage)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "filenotfound"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
		return
	}
	logger.InsertAuditLogs(c.Ctx, fmt.Sprintf("Delete car by user : %v , car id : %v", userId, data.Id), userId)
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, data, helpers.TranslateMessage(c.Ctx, "success", "delete"), "")
}

// GetAllCars ...
// @Title get cars
// @Desciption Get all car
// @Param body body dto.PaginationReq false "Insert New User"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 200 {object} string
// @Failure 403
// @router /cars [post]
func (c *CarController) GetAllCars() {
	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))
	var search dto.PaginationReq
	if err := c.ParseForm(&search); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &search)
	result, pagination_data, err := models.FetchCars(search.OpenPage, search.PageSize)

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
		logger.InsertAuditLogs(c.Ctx, "Error : "+message, userId)
		return
	}

	section_message := "read"
	section := "success"
	message := helpers.TranslateMessage(c.Ctx, section, section_message)
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, result, message, pagination_data)
	logger.InsertAuditLogs(c.Ctx, fmt.Sprintf("Get all car record by user : %v", userId), userId)
}

// GetSingleCar ...
// @Title get car
// @Desciption Get all car
// @Param body body dto.GetcarRequest true "get perticuler car"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 200 {object} string
// @Failure 403
// @router / [post]
func (c *CarController) GetSingleCar() {
	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))
	var bodyData dto.GetcarRequest
	if err := c.ParseForm(&bodyData); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &bodyData)
	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&bodyData); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		logger.InsertAuditLogs(c.Ctx, "Error : Validation error", userId)
		return
	}

	Data, err := models.GetSingleCar(bodyData.Id)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
		return
	}
	if Data.Id == 0 {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "invalidid"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+logger.LogMessage(c.Ctx, "error.invalidid"), userId)
		return
	}
	logger.InsertAuditLogs(c.Ctx, logger.LogMessage(c.Ctx, "success.read"), userId)
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, Data, helpers.TranslateMessage(c.Ctx, "success", "read"), "")
}

// Filter Car
// @Title Filter car
// @Description Fetch Filter car
// @Param search formData string true "Search car"
// @Param open_page formData int false "if you want to open specific page than give page number"
// @Param page_size formData int false "how much data you want to show at a time default it will give 10 records"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 200 {object} object
// @Failure 403
// @router /search_car [post]
func (c *CarController) FilterCars() {
	claims := helpers.GetTokenClaims(c.Ctx)
	userId := uint(claims["User_id"].(float64))
	var bodyData dto.SearchRequest
	if err := c.ParseForm(&bodyData); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &bodyData)
	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&bodyData); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		logger.InsertAuditLogs(c.Ctx, "Error : Validation error", userId)
		return
	}
	result, pagination_data, err := models.Filtercar(bodyData.Search, bodyData.OpenPage, bodyData.PageSize)

	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "searchnotfound"))
		logger.InsertAuditLogs(c.Ctx, "Error : "+err.Error(), userId)
		return
	}

	if pagination_data["pageOpen_error"] == 1 {
		current := pagination_data["current_page"]
		last := pagination_data["last_page"]
		message := fmt.Sprintf(helpers.TranslateMessage(c.Ctx, "error", "page"), current, last)
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, message)
		logger.InsertAuditLogs(c.Ctx, "Error : "+message, userId)
		return
	}
	section_message := "read"
	section := "success"
	message := helpers.TranslateMessage(c.Ctx, section, section_message)
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, result, message, pagination_data)
	logger.InsertAuditLogs(c.Ctx, message, userId)
}
