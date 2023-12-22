package models

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/dwarkesh2810/golang-demo/pkg/helpers"
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
	matchCountQuery := `SELECT state_id ,state_name ,state_iso_code FROM states where country_id = ` + country
	totalRecordQuery := `SELECT COUNT(*) as totalRows FROM states`
	mainRecordQuery := `SELECT state_id ,state_name ,state_iso_code FROM states where country_id = ` + country + ` LIMIT ? OFFSET ?`
	states, pagination, err := helpers.PaginationForSearch(open_page, page_size, totalRecordQuery, matchCountQuery, mainRecordQuery)
	if err != nil {
		return nil, nil, err
	}
	return states, pagination, nil
}

func FilterStates(search string, open_page, page_size int) ([]orm.Params, map[string]interface{}, error) {
	matchCountQuery := `SELECT state_id, state_name, state_iso_code FROM states WHERE state_name LIKE '%` + search + `%' OR state_iso_code LIKE '%` + search + `%';`
	totalRecordQuery := `SELECT COUNT(*) as totalRows FROM states`
	mainRecordQuery := `SELECT state_id ,state_name ,state_iso_code FROM states WHERE state_name LIKE '%` + search + `%' OR state_iso_code LIKE '%` + search + `%' LIMIT ? OFFSET ?`
	states, pagination, err := helpers.PaginationForSearch(open_page, page_size, totalRecordQuery, matchCountQuery, mainRecordQuery)
	if err != nil {
		return nil, nil, err
	}
	return states, pagination, nil
}

func GetState(state_id int) ([]orm.Params, error) {
	o := orm.NewOrm()
	query := `SELECT state_id, state_name, state_iso_code FROM states WHERE state_id = ? LIMIT 1`
	var state []orm.Params
	_, err := o.Raw(query, state_id).Values(&state)
	if err != nil {
		return nil, err
	}
	return state, nil
}
