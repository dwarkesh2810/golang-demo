package models

import (
	"github.com/beego/beego/v2/client/orm"
)

func InsertAuditLog(Data AuditLogs) error {
	o := orm.NewOrm()
	_, err := o.Insert(&Data)
	if err != nil {
		return err
	}
	return nil
}

// func FetchLogs(current_page, pageSize int) ([]orm.Params, map[string]interface{}, error) {
// 	tableName := "audit_logs"
// 	query := `SELECT u.first_name , u.last_name, u.email, u.phone_number
// 	FROM users as u
// 	ORDER BY u.user_id
// 	LIMIT ? OFFSET ?
// `
// 	result_data, pagination, errs := helpers.FetchDataWithPaginations(current_page, pageSize, tableName, query)
// 	if errs != nil {
// 		return nil, nil, errs
// 	}
// 	return result_data, pagination, nil
// }
