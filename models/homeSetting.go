package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/dwarkesh2810/golang-demo/dto"
	"github.com/dwarkesh2810/golang-demo/helpers"
)

func RegisterSetting(c dto.HomeSeetingInsert, user_id float64, file_path interface{}) (int, error) {

	db := orm.NewOrm()
	if file_path == "" {

		file_path = c.SettingData

	}

	res := HomePagesSettingTable{
		Section:     c.Section,
		DataType:    c.DataType,
		UniqueCode:  "",
		SettingData: file_path.(string),
		CreatedBy:   int(user_id),
		UpdatedBy:   0,
		CreatedDate: time.Now(),
	}

	_, err := db.Insert(&res)
	if err != nil {
		return 0, err
	}
	lastInsertID := res.PageSettingId
	UpdateUniqueCode(lastInsertID)
	return lastInsertID, nil
}

func UpdateUniqueCode(user_id int) (int64, error) {
	db := orm.NewOrm()

	// unique_codes := helpers.UniqueCode(user_id, os.Getenv("homePageModule"))
	unique_codes := helpers.UniqueCode(user_id, "homePageModule")

	home_page_setting := HomePagesSettingTable{PageSettingId: user_id}
	if db.Read(&home_page_setting) == nil {
		home_page_setting.UniqueCode = unique_codes
		if num, err := db.Update(&home_page_setting); err == nil {
			return num, nil
		}
	}
	return 1, nil
}

func UpdateSetting(c dto.HomeSeetingUpdate, file_path interface{}, user_id float64) (int64, error) {
	db := orm.NewOrm()
	page_setting_id := c.SettingId
	homePageSetting, setting_data_type, err := FetchPageSettingByID(page_setting_id)
	if err != nil {
		return 0, err
	}

	if file_path == "" {
		file_path = c.SettingData
	}
	setting_dataType := strings.ToUpper(setting_data_type)
	if setting_dataType == "LOGO" || setting_dataType == "BANNER" {
		file_name, file_directory := helpers.SplitFilePath(homePageSetting)
		helpers.RemoveFile(file_name, file_directory)

	}
	homePageData := HomePagesSettingTable{PageSettingId: page_setting_id,
		UpdatedBy:   int(user_id),
		UpdatedDate: time.Now(),
		DataType:    c.DataType,
		Section:     c.Section,
		SettingData: file_path.(string),
	}
	if num, err := db.Update(&homePageData, "updated_by", "updated_date", "data_type", "section", "setting_data"); err == nil {
		return num, nil
	}
	return 1, nil

}

func FetchPageSettingByID(pageSettingID int) (string, string, error) {
	db := orm.NewOrm()
	var pageSetting HomePagesSettingTable
	err := db.Raw(`SELECT  setting_data,data_type FROM home_pages_setting_table WHERE page_setting_id = ?`, pageSettingID).QueryRow(&pageSetting)
	if err != nil {
		return "Some errro occured in fetch page setting by ID function", "some errror", err
	}
	return pageSetting.SettingData, pageSetting.DataType, nil
}

func DeleteSetting(page_setting_id int) int {
	db := orm.NewOrm()
	setting := HomePagesSettingTable{PageSettingId: page_setting_id}
	if _, err := db.Delete(&setting); err == nil {
		return 1
	}
	return 0

}

func HomePageSettingExistsDelete(u dto.HomeSeetingDelete) int {
	page_setting_id := u.SettingId

	page_setting_data, page_setting_type, err := FetchPageSettingByID(page_setting_id)
	if err != nil {
		return 0
	}

	if strings.ToUpper(page_setting_type) == "LOGO" || strings.ToUpper(page_setting_type) == "BANNER" {
		file_name, file_directory := helpers.SplitFilePath(page_setting_data)
		helpers.RemoveFile(file_name, file_directory)
	}

	DeleteSetting(page_setting_id)
	return 1

}

func FetchSettingPagination(current_page, pageSize int) ([]orm.Params, map[string]interface{}, error) {
	db := orm.NewOrm()
	if current_page <= 0 {
		current_page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (current_page - 1) * pageSize

	var homeResponse []orm.Params
	_, err := db.Raw(`
		SELECT hpst.section, hpst.data_type, hpst.setting_data, hpst.created_date, hpst.updated_date,
		concat(umt.first_name,' ',umt.last_name) as created_by  
		FROM home_pages_setting_table as hpst
		LEFT JOIN user_master_table as umt ON umt.user_id = hpst.created_by
		ORDER BY hpst.created_date DESC
		LIMIT ? OFFSET ?
	`, pageSize, offset).Values(&homeResponse)

	if err != nil {
		return nil, nil, err
	}

	pagination_data, pagination_err := helpers.Pagination(current_page, pageSize, "home_pages_setting_table")
	if pagination_err != nil {
		return nil, pagination_data, nil
	}
	return homeResponse, pagination_data, nil
}

func FetchSetting() (interface{}, error) {
	db := orm.NewOrm()
	var homeResponse []struct {
		Section     string    `json:"section"`
		DataType    string    `json:"data_type"`
		SettingData string    `json:"setting_data"`
		CreatedDate time.Time `json:"created_date"`
		UpdatedDate time.Time `json:"updated_date"`
		CreatedBy   string    `json:"created_by"`
	}
	_, err := db.Raw(`SELECT hpst.section, hpst.data_type, hpst.setting_data,hpst.created_date, hpst.updated_date ,concat(umt.first_name,' ',umt.last_name) as created_by  FROM home_pages_setting_table as hpst LEFT JOIN user_master_table as umt ON umt.user_id = hpst.created_by ORDER BY hpst.created_date DESC`).QueryRows(&homeResponse)

	if err != nil {
		return nil, err
	}
	if len(homeResponse) == 0 {
		return "Not Found Cars", nil
	}
	return homeResponse, nil
}

func UpdateSettings(c dto.HomeSeetingUpdate, file_path interface{}, user_id float64) (int64, error) {
	db := orm.NewOrm()
	page_setting_id := c.SettingId
	homePageSetting, setting_data_type, err := FetchPageSettingByID(page_setting_id)
	if err != nil {
		return 0, err
	}

	if file_path == "" {
		file_path = c.SettingData
	}
	setting_dataType := strings.ToUpper(setting_data_type)
	if setting_dataType == "LOGO" || setting_dataType == "BANNER" {
		file_name, file_directory := helpers.SplitFilePath(homePageSetting)
		helpers.RemoveFile(file_name, file_directory)

	}
	homePageData := HomePagesSettingTable{PageSettingId: page_setting_id,
		UpdatedBy:   int(user_id),
		UpdatedDate: time.Now(),
		DataType:    c.DataType,
		Section:     c.Section,
		SettingData: file_path.(string),
	}
	if num, err := db.Update(&homePageData, "updated_by", "updated_date", "data_type", "section", "setting_data"); err == nil {
		return num, nil
	}
	return 1, nil

}

func ExportData(limit, starting_FromRow int) (interface{}, error) {
	db := orm.NewOrm()
	var homeResponse []struct {
		PageSettingId int       `json:"page_setting_id"`
		Section       string    `json:"section"`
		DataType      string    `json:"data_type"`
		SettingData   string    `json:"setting_data"`
		CreatedDate   time.Time `json:"created_date"`
		UpdatedDate   time.Time `json:"updated_date"`
		CreatedBy     string    `json:"created_by"`
	}

	if limit <= 0 {
		limit = 20
	}

	var query string
	if starting_FromRow <= 0 {

		query = fmt.Sprintf(`SELECT hpst.page_setting_id, hpst.section, hpst.data_type, hpst.setting_data, hpst.created_date, hpst.updated_date, concat(umt.first_name,' ',umt.last_name) as created_by 
			FROM home_pages_setting_table as hpst 
			LEFT JOIN user_master_table as umt ON umt.user_id = hpst.created_by  
			ORDER BY hpst.page_setting_id 
			LIMIT %d`, limit)
	} else {

		query = fmt.Sprintf(`SELECT hpst.page_setting_id, hpst.section, hpst.data_type, hpst.setting_data, hpst.created_date, hpst.updated_date, concat(umt.first_name,' ',umt.last_name) as created_by 
			FROM home_pages_setting_table as hpst 
			LEFT JOIN user_master_table as umt ON umt.user_id = hpst.created_by  
			WHERE hpst.page_setting_id >= %d
			ORDER BY hpst.page_setting_id 
			LIMIT %d`, starting_FromRow, limit)
	}

	_, err := db.Raw(query).QueryRows(&homeResponse)

	if err != nil {
		return nil, err
	}
	if len(homeResponse) == 0 {
		return "Not Found Records", nil
	}
	return homeResponse, nil
}

func RegisterSettingBatchsss(c dto.HomeSeetingInsert, user_id float64, filePath string, rows []map[string]interface{}) ([]int, []int, error) {
	db := orm.NewOrm()
	var insertIDs, updateIDs []int
	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Fetch existing page_setting_ids for indexing
	existingIDs := make(map[int]bool)
	var existingRecords []HomePagesSettingTable
	_, err = db.QueryTable("home_pages_setting_table").All(&existingRecords)
	if err != nil && err != orm.ErrNoRows {
		tx.Rollback()
		return nil, nil, err
	}
	for _, record := range existingRecords {
		existingIDs[record.PageSettingId] = true
	}

	// Batch processing
	for _, row := range rows {
		section, ok := row["section"].(string)
		if !ok {
			tx.Rollback()
			return nil, nil, errors.New("missing 'section' in row")
		}
		pageSettingIDStr, ok := row["page_setting_id"].(string)
		if !ok {
			// If page_setting_id is missing, treat it as a new record with auto-incremented ID
			newRecord := HomePagesSettingTable{
				Section:     section,
				DataType:    row["data_type"].(string),
				SettingData: row["setting_data"].(string),
				CreatedBy:   int(user_id),
				UpdatedBy:   0,
				CreatedDate: time.Now(),
				UpdatedDate: time.Now(),
			}

			_, err := tx.Insert(&newRecord)
			if err != nil {
				tx.Rollback()
				return nil, nil, err
			}

			insertIDs = append(insertIDs, newRecord.PageSettingId)
			UpdateUniqueCode(newRecord.PageSettingId)

			continue
		}

		// Convert page_setting_id to int64
		pageSettingID, err := strconv.ParseInt(pageSettingIDStr, 10, 64)
		if err != nil {
			tx.Rollback()
			return nil, nil, errors.New("invalid 'page_setting_id' format")
		}

		if _, exists := existingIDs[int(pageSettingID)]; exists {
			// If the record exists, update it
			var existingRecord HomePagesSettingTable
			err := db.QueryTable("home_pages_setting_table").Filter("page_setting_id", int(pageSettingID)).One(&existingRecord)
			if err != nil {
				tx.Rollback()
				return nil, nil, err
			}

			existingRecord.Section = section
			existingRecord.DataType = row["data_type"].(string)
			existingRecord.SettingData = row["setting_data"].(string)
			existingRecord.UpdatedBy = int(user_id)
			existingRecord.UpdatedDate = time.Now()

			_, err = tx.Update(&existingRecord)
			if err != nil {
				tx.Rollback()
				return nil, nil, err
			}

			updateIDs = append(updateIDs, existingRecord.PageSettingId)
		} else {
			// If the record doesn't exist, insert a new one
			newRecords := HomePagesSettingTable{
				PageSettingId: int(pageSettingID),
				Section:       section,
				DataType:      row["data_type"].(string),
				UniqueCode:    "",
				SettingData:   row["setting_data"].(string),
				CreatedBy:     int(user_id),
				UpdatedBy:     0,
				CreatedDate:   time.Now(),
				UpdatedDate:   time.Now(),
			}

			_, err := tx.Insert(&newRecords)
			if err != nil {
				tx.Rollback()
				return nil, nil, err
			}

			insertIDs = append(insertIDs, newRecords.PageSettingId)
			UpdateUniqueCode(newRecords.PageSettingId)
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	helpers.RemoveFileByPath(filePath)

	return insertIDs, updateIDs, nil
}

func FetchSettingPaginations(current_page, pageSize int) ([]orm.Params, map[string]interface{}, error) {
	tableName := "home_pages_setting_table"
	query := `
	SELECT hpst.section, hpst.data_type, hpst.setting_data, hpst.created_date, hpst.updated_date,
	concat(umt.first_name,' ',umt.last_name) as created_by
	FROM home_pages_setting_table as hpst
	LEFT JOIN users as umt ON umt.user_id = hpst.created_by
	ORDER BY hpst.created_date DESC
	LIMIT ? OFFSET ?
`
	result_data, pagination, errs := helpers.FetchDataWithPaginations(current_page, pageSize, tableName, query)
	if errs != nil {
		return nil, nil, errs
	}
	return result_data, pagination, nil
}

func FilterWithPaginationFetchSettings(currentPage, pageSize int, applyPositions string, searchFields map[string]string) ([]orm.Params, map[string]interface{}, error) {

	tableName := "home_pages_setting_table"

	query := `
        SELECT hpst.section, hpst.data_type, hpst.setting_data, hpst.created_date, hpst.updated_date,
        concat(umt.first_name,' ',umt.last_name) as created_by  
        FROM home_pages_setting_table as hpst
        LEFT JOIN users as umt ON umt.user_id = hpst.created_by
    `
	countQuery := `
        SELECT COUNT(*) as count
        FROM home_pages_setting_table as hpst
        LEFT JOIN users as umt ON umt.user_id = hpst.created_by
    `

	applyPosition := ""
	if applyPositions != "" {
		applyPosition = applyPositions
	}
	filterResult, pagination, count, errs := helpers.FilterData(currentPage, pageSize, query, tableName, searchFields, applyPosition, countQuery)
	if errs != nil {
		return nil, nil, errs
	}

	pagination["matchCount"] = 0
	if count > 0 {
		pagination["matchCount"] = count
	}

	return filterResult, pagination, nil
}
