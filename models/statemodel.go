package models

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/dwarkesh2810/golang-demo/helpers"
)

func FetchStateList(current_page, pageSize int) ([]orm.Params, map[string]interface{}, error) {
	tableName := "states"
	query := `
SELECT state_id ,state_name ,state_iso_code FROM states 
	LIMIT ? OFFSET ?
`
	result_data, pagination, errs := helpers.FetchDataWithPaginations(current_page, pageSize, tableName, query)
	if errs != nil {
		return nil, nil, errs
	}
	return result_data, pagination, nil
}

func CountryWiseState(open_page, page_size int, country string) ([]orm.Params, map[string]interface{}, error) {
	db := orm.NewOrm()
	if page_size <= 0 {
		page_size = 10
	}
	if open_page <= 0 {
		open_page = 1
	}
	var count []orm.Params
	query := `SELECT state_id ,state_name ,state_iso_code FROM states where country_id = ` + country
	_, err := db.Raw(query).Values(&count)
	if err != nil {
		return nil, nil, err
	}
	pagination, err := helpers.PaginationForSearch(open_page, page_size, len(count))
	if err != nil {
		return nil, nil, err
	}
	offset := (open_page - 1) * page_size
	var states []orm.Params
	query = `SELECT state_id ,state_name ,state_iso_code FROM states where country_id = ` + country + ` LIMIT ? OFFSET ?`
	_, err = db.Raw(query, page_size, offset).Values(&states)
	if err != nil {
		return nil, nil, err
	}
	pagination["matchCount"] = len(count)
	return states, pagination, nil
}

func FilterStates(search string, open_page, page_size int) ([]orm.Params, map[string]interface{}, error) {
	db := orm.NewOrm()
	if page_size <= 0 {
		page_size = 10
	}
	if open_page <= 0 {
		open_page = 1
	}
	var count []orm.Params
	query := `SELECT state_id, state_name, state_iso_code FROM states WHERE state_name LIKE '%` + search + `%' OR state_iso_code LIKE '%` + search + `%';`
	_, err := db.Raw(query).Values(&count)
	if err != nil {
		return nil, nil, err
	}

	pagination, err := helpers.PaginationForSearch(open_page, page_size, len(count))
	if err != nil {
		return nil, nil, err
	}
	offset := (open_page - 1) * page_size
	var states []orm.Params
	query = `SELECT state_id ,state_name ,state_iso_code FROM states WHERE state_name LIKE '%` + search + `%' OR state_iso_code LIKE '%` + search + `%' LIMIT ? OFFSET ?`
	_, err = db.Raw(query, page_size, offset).Values(&states)
	if err != nil {
		return nil, nil, err
	}
	pagination["matchCount"] = len(count)
	return states, pagination, nil
}
