package models

import (
	"strings"
	"time"

	"github.com/beego/beego/v2/client/httplib"

	"github.com/beego/beego/v2/client/orm"
	"github.com/dwarkesh2810/golang-demo/dto"
	"github.com/dwarkesh2810/golang-demo/pkg/helpers"
)

func InsertUpdateLanugaeLables(l dto.LanguageLableInsert) (int, string, error) {
	db := orm.NewOrm()
	langIniCode := helpers.ConvertIntoIniFormateCode(l.OtherINILanguageCodes)
	existsMultiLang := existsInMultilanguageLableTable(l.LableCodes, langIniCode)

	if existsMultiLang > 0 {
		existsENG := ExistsEngDefaultValues(l.LableCodes)
		if existsENG > 0 {
			langDefualt := EnglishLanguageLable{LableCode: l.LableCodes,
				LanguageValue: l.ENGLangValues,
				UpdatedDate:   time.Now(),
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
	}
	_, errs := db.Insert(&res)
	if errs != nil {
		return 0, "", err
	}
	return 1, res.LableCode, nil

}

func UpdateLanguageLables(l dto.LanguageLableUpdate) (int, string, error) {
	db := orm.NewOrm()
	lableIniCode := helpers.ConvertIntoIniFormateCode(l.OtherINILanguageCodes)
	existsMultiLang := existsInMultilanguageLableTable(l.LableCodes, lableIniCode)

	if existsMultiLang > 0 {
		existsENG := ExistsEngDefaultValues(l.LableCodes)
		if existsENG > 0 {
			langDefualt := EnglishLanguageLable{LableCode: l.LableCodes,
				LanguageValue: l.ENGLangValues,
				UpdatedDate:   time.Now(),
			}
			_, err := db.Update(&langDefualt, "language_value")
			if err != nil {
				return 0, "", err
			}
		}

		multilanguageUpdate := MultiLanguageLable{
			LanguageValue: l.OtherLangValues,
			UpdatedDate:   time.Now(),
		}

		updateData := map[string]interface{}{
			"LanguageValue": multilanguageUpdate.LanguageValue,
			"UpdatedDate":   multilanguageUpdate.UpdatedDate,
		}

		_, errs := db.QueryTable(new(MultiLanguageLable)).Filter("language_code", lableIniCode).Filter("lable_code", l.LableCodes).Update(updateData)

		if errs != nil {
			return 0, "", errs
		}

		return 1, "", nil
	}
	return 0, "", nil
}

// THIS IS FOR INSERT VALUE IN MULTIPLE LANGUAGE USING API

func InsertUpdateLanugaeLablesApi(l dto.LanguageLableInsertNew) (string, error) {
	db := orm.NewOrm()
	res := EnglishLanguageLable{
		LableCode:     l.LableCodes,
		LanguageCode:  "en-US",
		LanguageValue: l.ENGLangValues,
		Section:       l.Sections,
		CreatedDate:   time.Now(),
	}
	_, err := db.Insert(&res)
	if err != nil {
		return "", err
	}
	urlstr := UrlString(l.ENGLangValues)
	err = InsertGujrati(urlstr, l.LableCodes, l.Sections)
	if err != nil {
		return "", err
	}
	err = InsertHindi(urlstr, l.LableCodes, l.Sections)
	if err != nil {
		return "", err
	}
	err = InsertMarathi(urlstr, l.LableCodes, l.Sections)
	if err != nil {
		return "", err
	}
	return res.LableCode, nil
}

func InsertGujrati(urlstr, lable, section string) error {
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
		CreatedDate:   time.Now(),
	}
	_, err = db.Insert(&resMulti)
	return err
}

func InsertHindi(urlstr, lable, section string) error {
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
		CreatedDate:   time.Now(),
	}
	_, err = db.Insert(&resMulti)
	return err
}

func InsertMarathi(urlstr, lable, section string) error {
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
		CreatedDate:   time.Now(),
	}
	_, err = db.Insert(&resMulti)
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

func ExistsEngDefaultValues(lable_code string) int {
	db := orm.NewOrm()
	var lables EnglishLanguageLable
	err := db.Raw(`SELECT lang_id FROM english_language_lable WHERE lable_code = ?`, lable_code).QueryRow(&lables)
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
