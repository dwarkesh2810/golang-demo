package models

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	"github.com/dwarkesh2810/golang-demo/pkg/helpers"
)

/*Country Models----------------------------------------------------------------------------------------*/
func FetchCountriesList(current_page, pageSize int) ([]orm.Params, map[string]interface{}, error) {
	tableName := "countries"
	query := `
	SELECT country_id ,country_name  FROM countries 
	LIMIT ? OFFSET ?
`
	result_data, pagination, errs := helpers.FetchDataWithPaginations(current_page, pageSize, tableName, query)
	if errs != nil {
		return nil, nil, errs
	}
	return result_data, pagination, nil
}
func CountryFilter(currentPage, pageSize int, applyPositions string, searchFields map[string]string) ([]orm.Params, map[string]interface{}, error) {
	tableName := "countries"

	query := `
        SELECT country_id,country_name
        FROM countries
    `

	countQuery := `
        SELECT COUNT(*) as count
        FROM countries
    `

	applyPosition := ""
	if applyPositions != "" {
		applyPosition = applyPositions
	}
	var otherFieldSCount int = 0
	filterResult, pagination, count, errs := helpers.FilterData(currentPage, pageSize, query, tableName, searchFields, applyPosition, countQuery, otherFieldSCount)
	if errs != nil {
		return nil, nil, errs
	}
	pagination["matchCount"] = 0
	if count > 0 {
		pagination["matchCount"] = count
	}
	return filterResult, pagination, nil
}

/*City Models-------------------------------------------------------------------------------------------------------*/
func CityFilter(currentPage, pageSize, country_id, state_id, other_field_count int, applyPositions string, searchFields map[string]string) ([]orm.Params, map[string]interface{}, error) {
	tableName := "cities"
	query := `SELECT country_id,state_id,city_id ,city_name FROM cities`
	countQuery := `SELECT COUNT(*) as count FROM cities`

	if country_id != 0 && state_id != 0 {
		query += fmt.Sprintf(` WHERE country_id = %d AND state_id = %d `, country_id, state_id)
		countQuery += fmt.Sprintf(` WHERE country_id = %d AND state_id = %d`, country_id, state_id)
	}

	applyPosition := ""
	if applyPositions != "" {
		applyPosition = applyPositions
	}

	filterResult, pagination, count, errs := helpers.FilterData(currentPage, pageSize, query, tableName, searchFields, applyPosition, countQuery, other_field_count)

	if errs != nil {
		return nil, nil, errs
	}

	pagination["matchCount"] = 0
	if count > 0 {
		pagination["matchCount"] = count
	}

	return filterResult, pagination, nil
}

/*state_models------------------------------------------------------*/
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
