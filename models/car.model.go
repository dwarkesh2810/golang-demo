package models

import (
	"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/dwarkesh2810/golang-demo/dto"
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
	num, err := o.QueryTable(new(Car)).Filter("id", id).All(&car)
	if err != nil {
		return car, err
	}
	if num == 0 {
		return car, errors.New("DATABASE_ERROR")
	}
	return car, nil
}

func UpdateCar(data dto.UpdateCarRequest) (interface{}, error) {
	o := orm.NewOrm()
	var car = Car{
		Id:         data.Id,
		CarName:    data.CarName,
		ModifiedBy: data.ModifiedBy,
		Model:      data.Model,
		Type:       data.Type,
		CarImage:   data.CarImage,
		UpdateDate: time.Now(),
	}

	num, err := o.Update(&car, "id", "car_name", "modified_by", "model", "car_type", "car_image", "updated_at")
	if err != nil {
		return num, err
	}
	if num == 0 {
		return nil, errors.New("NOT_FOUND")
	}
	return car, nil
}

func DeleteCar(id uint) (interface{}, error) {
	o := orm.NewOrm()
	var car = Car{Id: id}
	num, err := o.Delete(&car)
	if err != nil {
		return nil, err
	}
	if num == 0 {
		return nil, errors.New("NOT_FOUND")
	}
	return car, nil
}