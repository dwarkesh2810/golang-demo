package models

import (
	"strings"
	"time"

	"github.com/beego/beego/v2/client/httplib"

	"github.com/beego/beego/v2/client/orm"
	"github.com/dwarkesh2810/golang-demo/dto"
	"github.com/dwarkesh2810/golang-demo/pkg/helpers"
)

func InsertUpdateLanugaeLables(l dto.LanguageLableInsert, userId int) (int, string, error) {
	db := orm.NewOrm()
	langIniCode := helpers.ConvertIntoIniFormateCode(l.OtherINILanguageCodes)
	existsMultiLang := existsInMultilanguageLableTable(l.LableCodes, langIniCode)

	if existsMultiLang > 0 {
		existsENG := ExistsEngDefaultValues(l.LableCodes)
		if existsENG > 0 {
			langDefualt := EnglishLanguageLable{LableCode: l.LableCodes,
				LanguageValue: l.ENGLangValues,
				UpdatedDate:   time.Now(),
				CreatedBy:     userId,
			}
			_, err := db.Update(&langDefualt, "language_value")
			if err != nil {
				return 0, "", err
			}
			return 0, "", nil
		}
		res := EnglishLanguageLable{
			LableCode:     l.LableCodes,
			LanguageCode:  "en-US",
			LanguageValue: l.ENGLangValues,
			Section:       l.Sections,
			CreatedDate:   time.Now(),
			CreatedBy:     userId,
		}
		_, err := db.Insert(&res)
		if err != nil {
			return 0, "", err
		}
		return 0, "", nil
	}

	resMulti := MultiLanguageLable{
		LableCode:     l.LableCodes,
		LanguageCode:  helpers.ConvertIntoIniFormateCode(l.OtherINILanguageCodes),
		LanguageValue: l.OtherLangValues,
		Section:       l.Sections,
		CreatedDate:   time.Now(),
		CreatedBy:     userId,
	}
	_, err := db.Insert(&resMulti)
	if err != nil {
		return 0, "", err
	}

	existSENG := ExistsEngDefaultValues(l.LableCodes)

	if existSENG > 0 {
		langDefualt := EnglishLanguageLable{LableCode: l.LableCodes,
			LanguageValue: l.ENGLangValues,
			UpdatedDate:   time.Now(),
			CreatedBy:     userId,
		}
		_, err := db.Update(&langDefualt, "language_value")
		if err != nil {
			return 0, "", err
		}
	}

	res := EnglishLanguageLable{
		LableCode:     l.LableCodes,
		LanguageCode:  "en-US",
		LanguageValue: l.ENGLangValues,
		Section:       l.Sections,
		CreatedDate:   time.Now(),
		CreatedBy:     userId,
	}
	_, errs := db.Insert(&res)
	if errs != nil {
		return 0, "", err
	}
	return 1, res.LableCode, nil

}

func UpdateLanguageLables(l dto.LanguageLableUpdate, userID int) (int, string, error) {
	db := orm.NewOrm()
	lableIniCode := helpers.ConvertIntoIniFormateCode(l.OtherINILanguageCodes)
	existsMultiLang := existsInMultilanguageLableTable(l.LableCodes, lableIniCode)
	orm.Debug = true
	if existsMultiLang > 0 {
		existsENG := ExistsEngDefaultValues(l.LableCodes)
		if existsENG > 0 {
			langDefualt := EnglishLanguageLable{
				LableCode:     l.LableCodes,
				LanguageValue: l.ENGLangValues,
				UpdatedDate:   time.Now(),
				UpdatedBy:     userID,
			}
			updateData := map[string]interface{}{
				"LanguageValue": langDefualt.LanguageValue,
				"UpdatedDate":   langDefualt.UpdatedDate,
				"UpdatedBy":     langDefualt.UpdatedBy,
			}

			_, err := db.QueryTable(new(EnglishLanguageLable)).Filter("lable_code", l.LableCodes).Update(updateData)
			if err != nil {
				return 0, "", err
			}
		}

		multilanguageUpdate := MultiLanguageLable{
			LanguageValue: l.OtherLangValues,
			UpdatedDate:   time.Now(),
			UpdatedBy:     userID,
		}

		updateData := map[string]interface{}{
			"LanguageValue": multilanguageUpdate.LanguageValue,
			"UpdatedDate":   multilanguageUpdate.UpdatedDate,
			"UpdatedBy":     multilanguageUpdate.UpdatedBy,
		}

		_, errs := db.QueryTable(new(MultiLanguageLable)).Filter("language_code", lableIniCode).Filter("lable_code", l.LableCodes).Update(updateData)

		if errs != nil {
			return 0, "", errs
		}

		return 1, "", nil
	}
	return 0, "", nil
}

func ExportLanguageLables() ([]orm.Params, error) {
	db := orm.NewOrm()
	var multiLanguageLables []orm.Params

	query := `SELECT mll.lable_id as lable_id, mll.lable_code as lable_code, mll.language_value as language_value, mll.language_code as language_code, mll.section as section, mll.created_by as created_user_id, mll.updated_by as updated_user_id,concat(us.first_name,' ',us.last_name) as updated_user_name, concat(u.first_name,' ',u.last_name) as created_user_name, mll.created_date as created_date, mll.updated_date as updated_date
FROM multi_language_lable as mll
    LEFT JOIN users as u ON u.user_id = mll.created_by
    LEFT JOIN users as us ON us.user_id = mll.updated_by
ORDER BY mll.lable_id`
	_, err := db.Raw(query).Values(&multiLanguageLables)

	if err != nil {
		return nil, err
	}

	return multiLanguageLables, nil
}

// THIS IS FOR INSERT VALUE IN MULTIPLE LANGUAGE USING API

func InsertUpdateLanugaeLablesApi(l dto.LanguageLable, userID int) (string, error) {
	db := orm.NewOrm()
	res := EnglishLanguageLable{
		LableCode:     l.LableCodes,
		LanguageCode:  "en-US",
		LanguageValue: l.ENGLangValues,
		Section:       l.Sections,
		CreatedBy:     userID,
		CreatedDate:   time.Now(),
	}
	_, err := db.Insert(&res)
	if err != nil {
		return "", err
	}
	urlstr := UrlString(l.ENGLangValues)
	err = InsertGujrati(urlstr, l.LableCodes, l.Sections, userID)
	if err != nil {
		return "", err
	}
	err = InsertHindi(urlstr, l.LableCodes, l.Sections, userID)
	if err != nil {
		return "", err
	}
	err = InsertMarathi(urlstr, l.LableCodes, l.Sections, userID)
	if err != nil {
		return "", err
	}
	return res.LableCode, nil
}

func UpdateLanguageLablesAPI(l dto.LanguageLable, userID int) (string, error) {
	db := orm.NewOrm()
	// orm.Debug = true
	langDefualt := EnglishLanguageLable{
		LableCode:     l.LableCodes,
		LanguageValue: l.ENGLangValues,
		UpdatedBy:     userID,
		UpdatedDate:   time.Now(),
	}
	updateData := map[string]interface{}{
		"LanguageValue": langDefualt.LanguageValue,
		"UpdatedDate":   langDefualt.UpdatedDate,
		"UpdatedBy":     langDefualt.UpdatedBy,
	}

	_, err := db.QueryTable(new(EnglishLanguageLable)).Filter("lable_code", l.LableCodes).Update(updateData)
	if err != nil {
		return "", err
	}
	urlstr := UrlString(l.ENGLangValues)

	err = UpdateGujrati(urlstr, l.LableCodes, userID)
	if err != nil {
		return "", err
	}
	err = UpdateHindi(urlstr, l.LableCodes, userID)
	if err != nil {
		return "", err
	}
	err = UpdateMarathi(urlstr, l.LableCodes, userID)
	if err != nil {
		return "", err
	}
	return langDefualt.LableCode, nil
}

func InsertGujrati(urlstr, lable, section string, userID int) error {
	db := orm.NewOrm()
	req := httplib.Get("https://translate.googleapis.com/translate_a/single?client=gtx&sl=en&tl=gu&dt=t&q=" + urlstr)
	str, err := req.String()
	if err != nil {
		return err
	}
	strres := GetTranslatedata(str)
	resMulti := MultiLanguageLable{
		LableCode:     lable,
		LanguageCode:  "gu-IN",
		LanguageValue: strres,
		Section:       section,
		CreatedBy:     userID,
		CreatedDate:   time.Now(),
	}
	_, err = db.Insert(&resMulti)
	return err
}

func InsertHindi(urlstr, lable, section string, userID int) error {
	db := orm.NewOrm()
	req := httplib.Get("https://translate.googleapis.com/translate_a/single?client=gtx&sl=en&tl=hi&dt=t&q=" + urlstr)
	str, err := req.String()
	if err != nil {
		return err
	}
	strres := GetTranslatedata(str)
	resMulti := MultiLanguageLable{
		LableCode:     lable,
		LanguageCode:  "hi-IN",
		LanguageValue: strres,
		Section:       section,
		CreatedBy:     userID,
		CreatedDate:   time.Now(),
	}
	_, err = db.Insert(&resMulti)
	return err
}

func InsertMarathi(urlstr, lable, section string, userID int) error {
	db := orm.NewOrm()
	req := httplib.Get("https://translate.googleapis.com/translate_a/single?client=gtx&sl=en&tl=mr&dt=t&q=" + urlstr)
	str, err := req.String()
	if err != nil {
		return err
	}
	strres := GetTranslatedata(str)
	resMulti := MultiLanguageLable{
		LableCode:     lable,
		LanguageCode:  "mr-IN",
		LanguageValue: strres,
		Section:       section,
		CreatedBy:     userID,
		CreatedDate:   time.Now(),
	}
	_, err = db.Insert(&resMulti)
	return err
}
func UpdateGujrati(urlstr, lable string, userID int) error {
	db := orm.NewOrm()
	req := httplib.Get("https://translate.googleapis.com/translate_a/single?client=gtx&sl=en&tl=gu&dt=t&q=" + urlstr)
	str, err := req.String()
	if err != nil {
		return err
	}
	strres := GetTranslatedata(str)
	multilanguageUpdate := MultiLanguageLable{
		LanguageValue: strres,
		UpdatedDate:   time.Now(),
		UpdatedBy:     userID,
	}

	updateData := map[string]interface{}{
		"LanguageValue": multilanguageUpdate.LanguageValue,
		"UpdatedDate":   multilanguageUpdate.UpdatedDate,
		"UpdatedBy":     multilanguageUpdate.UpdatedBy,
	}

	_, err = db.QueryTable(new(MultiLanguageLable)).Filter("language_code", "gu-IN").Filter("lable_code", lable).Update(updateData)
	return err
}
func UpdateHindi(urlstr, lable string, userID int) error {
	db := orm.NewOrm()
	req := httplib.Get("https://translate.googleapis.com/translate_a/single?client=gtx&sl=en&tl=hi&dt=t&q=" + urlstr)
	str, err := req.String()
	if err != nil {
		return err
	}
	strres := GetTranslatedata(str)
	multilanguageUpdate := MultiLanguageLable{
		LanguageValue: strres,
		UpdatedDate:   time.Now(),
		UpdatedBy:     userID,
	}

	updateData := map[string]interface{}{
		"LanguageValue": multilanguageUpdate.LanguageValue,
		"UpdatedDate":   multilanguageUpdate.UpdatedDate,
		"UpdatedBy":     multilanguageUpdate.UpdatedBy,
	}

	_, err = db.QueryTable(new(MultiLanguageLable)).Filter("language_code", "hi-IN").Filter("lable_code", lable).Update(updateData)
	return err
}
func UpdateMarathi(urlstr, lable string, userID int) error {
	db := orm.NewOrm()
	req := httplib.Get("https://translate.googleapis.com/translate_a/single?client=gtx&sl=en&tl=mr&dt=t&q=" + urlstr)
	str, err := req.String()
	if err != nil {
		return err
	}
	strres := GetTranslatedata(str)
	multilanguageUpdate := MultiLanguageLable{
		LanguageValue: strres,
		UpdatedDate:   time.Now(),
		UpdatedBy:     userID,
	}

	updateData := map[string]interface{}{
		"LanguageValue": multilanguageUpdate.LanguageValue,
		"UpdatedDate":   multilanguageUpdate.UpdatedDate,
		"UpdatedBy":     multilanguageUpdate.UpdatedBy,
	}

	_, err = db.QueryTable(new(MultiLanguageLable)).Filter("language_code", "mr-IN").Filter("lable_code", lable).Update(updateData)
	return err
}

func GetTranslatedata(data string) string {
	return strings.Split(data, `"`)[1]
}
func UrlString(s string) string {
	return strings.ReplaceAll(s, " ", "%20")
}

// end

func existsInMultilanguageLableTable(lable_code, iniCode string) int {
	db := orm.NewOrm()
	var lables MultiLanguageLable
	err := db.Raw(`SELECT lable_id FROM multi_language_lable  WHERE lable_code = ? AND language_code = ?`, lable_code, iniCode).QueryRow(&lables)
	if err != nil {
		return 0
	}
	return 1
}

func existsInMultilanguageLableTablesss(lable_code, iniCode, section string) int {
	db := orm.NewOrm()
	var lables MultiLanguageLable
	err := db.Raw(`SELECT lable_id FROM multi_language_lable  WHERE lable_code = ? AND language_code = ? AND section = ?`, lable_code, iniCode, section).QueryRow(&lables)
	if err != nil {
		return 0
	}
	return 1
}

func ExistsEngDefaultValues(lable_code string) int {
	db := orm.NewOrm()
	var lables EnglishLanguageLable
	err := db.Raw(`SELECT lang_id FROM english_language_lable WHERE lable_code = ?`, lable_code).QueryRow(&lables)
	if err != nil {
		return 0
	}
	return 1
}

func ExistsEngDefaultValuesss(lable_code, section string) int {
	db := orm.NewOrm()
	var lables EnglishLanguageLable
	err := db.Raw(`SELECT lang_id FROM english_language_lable WHERE lable_code = ? AND section = ?`, lable_code, section).QueryRow(&lables)
	if err != nil {
		return 0
	}
	return 1
}

func IsLanguageLableExist(lable_code, section string) error {
	db := orm.NewOrm()
	var lables EnglishLanguageLable
	err := db.Raw(`SELECT lang_id FROM english_language_lable WHERE lable_code = ? AND section = ?`, lable_code, section).QueryRow(&lables)
	return err
}

func FetchAllLabels() ([]orm.Params, error) {
	db := orm.NewOrm()
	var labelsList []orm.Params

	_, err := db.Raw(`SELECT lable_code, language_value, language_code,section FROM multi_language_lable`).Values(&labelsList)
	if err != nil {
		return nil, err
	}
	return labelsList, nil
}

func FetchAllDefaultlables() ([]orm.Params, error) {
	db := orm.NewOrm()
	var labelsList []orm.Params

	_, err := db.Raw(`SELECT lable_code, language_value, language_code,section FROM english_language_lable`).Values(&labelsList)
	if err != nil {
		return nil, err
	}
	return labelsList, nil
}

func ImportINIFiles(languageCodes string, user_id uint, dataMap map[string]map[string]map[string]string) (string, error) {
	if strings.ToUpper(languageCodes) == "EN-US" {

		for section, keys := range dataMap {
			for labelCode, languageValue := range keys {
				_, err := insertUpdateEnglishLanugaeLable(section, labelCode, languageValue["language_value"], int(user_id))
				if err != nil {
					return "", err
				}
			}
		}
		return "en-us", nil
	}

	for section, keys := range dataMap {
		for labelCode, languageValue := range keys {
			_, err := insertUpdateMultiLanugaeLabless(section, labelCode, languageValue["language_value"], languageCodes, int(user_id))
			if err != nil {
				return "", err
			}
		}
	}
	return "multi", nil
}

func insertUpdateMultiLanugaeLabless(section, labelCodes, languageValue, languageCodes string, user_id int) (int, error) {
	db := orm.NewOrm()
	langIniCode := helpers.ConvertIntoIniFormateCode(languageCodes)
	existsMultiLang := existsInMultilanguageLableTablesss(labelCodes, langIniCode, section)

	if existsMultiLang > 0 {
		multilanguageUpdate := MultiLanguageLable{
			LanguageValue: languageValue,
			UpdatedBy:     user_id,
			UpdatedDate:   time.Now(),
		}

		updateData := map[string]interface{}{
			"LanguageValue": multilanguageUpdate.LanguageValue,
			"UpdatedDate":   multilanguageUpdate.UpdatedDate,
		}

		_, errs := db.QueryTable(new(MultiLanguageLable)).Filter("language_code", langIniCode).Filter("lable_code", labelCodes).Filter("section", section).Update(updateData)

		if errs != nil {
			return 0, errs
		}
		return 1, nil
	}

	resMulti := MultiLanguageLable{
		LableCode:     labelCodes,
		LanguageCode:  helpers.ConvertIntoIniFormateCode(languageCodes),
		LanguageValue: languageValue,
		Section:       section,
		CreatedBy:     user_id,
		CreatedDate:   time.Now(),
	}
	_, err := db.Insert(&resMulti)
	if err != nil {
		return 0, err
	}

	return 1, nil
}

func insertUpdateEnglishLanugaeLable(section, labelCodes, languageValue string, user_id int) (int, error) {
	db := orm.NewOrm()
	existsEnglish := ExistsEngDefaultValuesss(labelCodes, section)

	if existsEnglish > 0 {
		englishLangLables := EnglishLanguageLable{
			LanguageValue: languageValue,
			UpdatedBy:     user_id,
			UpdatedDate:   time.Now(),
		}
		updateData := map[string]interface{}{
			"LanguageValue": englishLangLables.LanguageValue,
			"UpdatedDate":   englishLangLables.UpdatedDate,
		}

		_, errs := db.QueryTable(new(EnglishLanguageLable)).Filter("lable_code", labelCodes).Filter("section", section).Update(updateData)

		if errs != nil {
			return 0, errs
		}
		return 1, nil
	}

	englishRes := EnglishLanguageLable{
		LableCode:     labelCodes,
		LanguageCode:  "en-US",
		LanguageValue: languageValue,
		Section:       section,
		CreatedBy:     user_id,
		CreatedDate:   time.Now(),
	}
	_, err := db.Insert(&englishRes)
	if err != nil {
		return 0, err
	}

	return 1, nil
}
