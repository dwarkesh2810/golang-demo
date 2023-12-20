package models

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/dwarkesh2810/golang-demo/helpers"
)

func FetchCountriesList(current_page, pageSize int) ([]orm.Params, map[string]interface{}, error) {
	tableName := "mod_country_master"
	query := `
	SELECT country_id ,name as country_name FROM mod_country_master 
	LIMIT ? OFFSET ?
`
	result_data, pagination, errs := helpers.FetchDataWithPaginations(current_page, pageSize, tableName, query)
	if errs != nil {
		return nil, nil, errs
	}
	return result_data, pagination, nil
}
