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
	filterResult, pagination, _, errs := helpers.FilterData(currentPage, pageSize, query, tableName, searchFields, applyPosition, countQuery, otherFieldSCount)
	if errs != nil {
		return nil, nil, errs
	}
	pagination["matchCount"] = 0

	if pagination["TotalMatches"] != nil {
		pagination["matchCount"] = pagination["TotalMatches"]
	}
	return filterResult, pagination, nil
}

/*City Models-------------------------------------------------------------------------------------------------------*/

func CityFilter(currentPage, pageSize, country_id, state_id, other_field_count int, applyPositions, countryName, stateName string, searchFields map[string]string) ([]orm.Params, map[string]interface{}, error) {
	tableName := "cities"
	query := `SELECT country_id,state_id,city_id ,city_name FROM cities`
	countQuery := `SELECT COUNT(*) as count FROM cities`
	if country_id != 0 && state_id != 0 {
		query += fmt.Sprintf(` WHERE country_id = %d AND state_id = %d `, country_id, state_id)
		countQuery += fmt.Sprintf(` WHERE country_id = %d AND state_id = %d`, country_id, state_id)
	}

	if countryName != "" || stateName != "" {
		query = fmt.Sprintf(`select city.city_name as city_name,city.city_id as city_id from countries as c
		LEFT JOIN cities as city ON city.country_id=c.country_id
		LEFT JOIN states as s ON s.state_id = city.state_id
		WHERE upper(c.country_name) = upper('%s') or upper(s.state_name)=upper('%s')
		`, countryName, stateName)

		countQuery = fmt.Sprintf(`select COUNT(*) as count from countries as c
		LEFT JOIN cities as city ON city.country_id=c.country_id
		LEFT JOIN states as s ON s.state_id = city.state_id
		WHERE upper(c.country_name) = upper('%s') or upper(s.state_name)=upper('%s')
		`, countryName, stateName)

		if countryName != "" && stateName != "" {
			query = fmt.Sprintf(`select city.city_name as city_name,city.city_id as city_id from countries as c
			LEFT JOIN cities as city ON city.country_id=c.country_id
			LEFT JOIN states as s ON s.state_id = city.state_id
			WHERE upper(c.country_name) = upper('%s') AND upper(s.state_name)=upper('%s')
			`, countryName, stateName)
			countQuery = fmt.Sprintf(`select COUNT(*) as count from countries as c
			LEFT JOIN cities as city ON city.country_id=c.country_id
			LEFT JOIN states as s ON s.state_id = city.state_id
			WHERE upper(c.country_name) = upper('%s') AND upper(s.state_name)=upper('%s')
			`, countryName, stateName)
		}
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
		pagination["matchCount"] = pagination["TotalMaches"]
	}

	return filterResult, pagination, nil
}

func FilterCountries(search string, open_page, page_size int) ([]orm.Params, map[string]interface{}, error) {
	matchCountQuery := `SELECT country_id, country_name, country_iso_code FROM countries WHERE country_name LIKE '%` + search + `%' OR country_iso_code LIKE '%` + search + `%';`
	totalRecordQuery := `SELECT COUNT(*) as totalRows FROM countries`
	mainRecordQuery := `SELECT country_id, country_name, country_iso_code FROM countries WHERE country_name LIKE '%` + search + `%' OR country_iso_code LIKE '%` + search + `%' LIMIT ? OFFSET ?`
	states, pagination, err := helpers.PaginationForSearch(open_page, page_size, totalRecordQuery, matchCountQuery, mainRecordQuery)
	if err != nil {
		return nil, nil, err
	}
	if search == "" {
		pagination["matchCount"] = 0
	}
	return states, pagination, nil
}

func GetCountry(state_id int) ([]orm.Params, error) {
	o := orm.NewOrm()
	query := `SELECT country_id, country_name, country_iso_code FROM countries WHERE country_id = ? LIMIT 1`
	var country []orm.Params
	_, err := o.Raw(query, state_id).Values(&country)
	if err != nil {
		return nil, err
	}
	return country, nil
}
