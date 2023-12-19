package controllers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/dwarkesh2810/golang-demo/dto"
	"github.com/dwarkesh2810/golang-demo/helpers"
	"github.com/dwarkesh2810/golang-demo/models"
	"github.com/dwarkesh2810/golang-demo/validations"
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
// @Param file formData file true "File to be uploaded"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} string
// @Failure 403
// @router /create [post]
func (c *CarController) AddNewCar() {
	var cars dto.GetNewCarRequest
	if err := c.ParseForm(&cars); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &cars)

	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&cars); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		return
	}

	file, fileheader, err := c.GetFile("file")
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "filenotfound"))
		return
	}
	var carType string = string(cars.Type)
	cars.Type, err = helpers.NewCarType(carType)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "cartype"))
		return
	}
	uploadDir := "./uploads/Cars/images/"
	filepaths, err := helpers.UploadFile(file, fileheader, uploadDir)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "upload"))
		return
	}
	cars.CarImage = filepaths
	data, err := models.InsertNewCar(cars)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		return
	}
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
// @Param file formData file false "File to be uploaded"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} string
// @Failure 403
// @router /update [PUT]
func (c *CarController) UpdateCar() {
	var cars dto.UpdateCarRequest
	if err := c.ParseForm(&cars); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &cars)

	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&cars); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		return
	}

	data, err := models.GetSingleCar(cars.Id)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		return
	}

	if data.CarName == "" {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "datanotfound"))
		return
	}

	file, fileheader, err := c.GetFile("file")
	if err != nil {
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
		var carType string = string(cars.Type)
		cars.Type, err = helpers.NewCarType(carType)
		if err != nil {
			helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "cartype"))
			return
		}
		cars.CarImage = data.CarImage
		res, err := models.UpdateCar(cars)
		if err != nil {
			helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
			return
		}
		helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, res, helpers.TranslateMessage(c.Ctx, "success", "update"), "")
		return
	}
	var carType string = string(cars.Type)
	cars.Type, err = helpers.NewCarType(carType)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "cartype"))
		return
	}
	uploadDir := "./uploads/Cars/images/"
	filepaths, err := helpers.UploadFile(file, fileheader, uploadDir)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "upload"))
		return
	}
	cars.CarImage = filepaths
	res, err := models.UpdateCar(cars)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		return
	}
	err = os.Remove(data.CarImage)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "filenotfound"))
		return
	}
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, res, helpers.TranslateMessage(c.Ctx, "success", "update"), "")
}

// DeleteCar ...
// @Title remove car
// @Desciption delete car
// @Param body body dto.GetcarRequest true "delete car"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} string
// @Failure 403
// @router /delete [delete]
func (c *CarController) DeleteCar() {
	var car dto.GetcarRequest
	if err := c.ParseForm(&car); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &car)

	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&car); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		return
	}

	res, err := models.GetSingleCar(car.Id)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		return
	}

	if res.CarName == "" {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "datanotfound"))
		return
	}

	data, err := models.DeleteCar(car.Id)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
		return
	}
	err = os.Remove(res.CarImage)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "filenotfound"))
		return
	}
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, data, helpers.TranslateMessage(c.Ctx, "success", "delete"), "")
}

// GetAllCars ...
// @Title get cars
// @Desciption Get all car
// @Param body body dto.PaginationReq false "Insert New User"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} string
// @Failure 403
// @router /cars [post]
func (c *CarController) GetAllCars() {
	var search dto.PaginationReq
	if err := c.ParseForm(&search); err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "parsing"))
		return
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &search)

	valid := validation.Validation{}
	if isValid, _ := valid.Valid(&search); !isValid {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, validations.ValidationErrorResponse(c.Controller, valid.Errors))
		return
	}

	result, pagination_data, _ := models.FetchCars(search.OpenPage, search.PageSize)
	if pagination_data["pageOpen_error"] == 1 {
		current := pagination_data["current_page"]
		last := pagination_data["last_page"]
		message := fmt.Sprintf("PAGE NUMBER %d IS NOT EXISTS , LAST PAGE NUMBER IS %d", current, last)
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, message)
		return
	}

	if result != nil {
		section_message := ""
		section := ""
		message := helpers.TranslateMessage(c.Ctx, section, section_message)
		helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, result, message, pagination_data)
		return
	}
	helpers.ApiFailedResponse(c.Ctx.ResponseWriter, helpers.TranslateMessage(c.Ctx, "error", "db"))
}

// GetSingleCar ...
// @Title get car
// @Desciption Get all car
// @Param body body models.GetcarRequest true "get perticuler car"
// @Param   Authorization   header  string  true  "Bearer YourAccessToken"
// @Success 201 {object} string
// @Failure 403
// @router / [post]
func (c *CarController) GetSingleCar() {
	var bodyData dto.GetcarRequest
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

	Data, err := models.GetSingleCar(bodyData.Id)
	if err != nil {
		helpers.ApiFailedResponse(c.Ctx.ResponseWriter, err.Error())
	}
	helpers.ApiSuccessResponse(c.Ctx.ResponseWriter, Data, helpers.TranslateMessage(c.Ctx, "success", "read"), "")
}
