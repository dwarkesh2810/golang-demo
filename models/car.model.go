package models

import (
	"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/dwarkesh2810/golang-demo/dto"
	"github.com/dwarkesh2810/golang-demo/pkg/helpers"
)

func InsertNewCar(data dto.GetNewCarRequest) (Car, error) {
	o := orm.NewOrm()
	var car = Car{
		CarName:     data.CarName,
		CarImage:    data.CarImage,
		ModifiedBy:  data.ModifiedBy,
		Model:       data.Model,
		Type:        data.Type,
		CreatedDate: time.Now(),
		CreatedBy:   data.CreatedBy,
	}
	_, err := o.Insert(&car)
	if err != nil {
		return car, err
	}
	return car, nil
}

func GetSingleCar(id uint) (Car, error) {
	o := orm.NewOrm()
	var car Car
	_, err := o.QueryTable(new(Car)).Filter("id", id).All(&car)
	if err != nil {
		return car, err
	}
	return car, nil
}

func UpdateCar(data dto.UpdateCarRequest) (Car, error) {
	o := orm.NewOrm()
	var car = Car{
		Id:         data.Id,
		CarName:    data.CarName,
		ModifiedBy: data.ModifiedBy,
		Model:      data.Model,
		Type:       data.Type,
		CarImage:   data.CarImage,
		UpdateDate: time.Now(),
		UpdatedBy:  data.UpdatedBy,
	}

	num, err := o.Update(&car, "id", "car_name", "modified_by", "model", "car_type", "car_image", "updated_date", "updated_by")
	if err != nil {
		return car, err
	}
	if num == 0 {
		return car, errors.New("NOT_FOUND")
	}
	return car, nil
}

func DeleteCar(id uint) (Car, error) {
	o := orm.NewOrm()
	var car = Car{Id: id}
	num, err := o.Delete(&car)
	if err != nil {
		return car, err
	}
	if num == 0 {
		return car, errors.New("NOT_FOUND")
	}
	return car, nil
}

func FetchCars(current_page, pageSize int) ([]orm.Params, map[string]interface{}, error) {
	tableName := "car"
	query := `SELECT c.car_name , c.modified_by, c.model, c.car_type
	FROM car as c
	ORDER BY c.id
	LIMIT ? OFFSET ?`
	result_data, pagination, errs := helpers.FetchDataWithPaginations(current_page, pageSize, tableName, query)
	if errs != nil {
		return nil, nil, errs
	}
	return result_data, pagination, nil
}

func Filtercar(search string, open_page, page_size int) ([]orm.Params, map[string]interface{}, error) {
	matchCountQuery := `SELECT id, car_name FROM car WHERE car_name LIKE '%` + search + `%' OR modified_by LIKE '%` + search + `%' OR model LIKE '%` + search + `%' OR car_type LIKE '%` + search + `%';`
	totalRecordQuery := `SELECT COUNT(*) as totalRows FROM car`
	mainRecordQuery := `SELECT id ,car_name ,modified_by, model,car_type FROM car WHERE car_name LIKE '%` + search + `%' OR modified_by LIKE '%` + search + `%' OR model LIKE '%` + search + `%' OR car_type LIKE '%` + search + `%' LIMIT ? OFFSET ?`
	states, pagination, err := helpers.PaginationForSearch(open_page, page_size, totalRecordQuery, matchCountQuery, mainRecordQuery)
	if err != nil {
		return nil, nil, err
	}
	if search == "" {
		pagination["matchCount"] = 0
	}
	return states, pagination, nil
}
