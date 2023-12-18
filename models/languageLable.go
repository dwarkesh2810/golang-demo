package models

import (
	"log"

	"github.com/beego/beego/v2/client/orm"
	"github.com/dwarkesh2810/golang-demo/dto"
)

func InsertLanguageLabels(l dto.LanguageLableInsert) (string, error) {
	db := orm.NewOrm()
	res := LanguageLableLang{
		LableCode:     l.LableCode,
		LanguageCode:  l.Language,
		LanguageValue: l.LangValue,
		Section:       l.Section,
	}
	_, err := db.Insert(&res)
	if err != nil {
		return "", err
	}

	lastInsertedID := res.LangId

	languageLabels := []LanguageLable{
		{
			LableCode:     l.LableCode,
			LangId:        lastInsertedID,
			LanguageCode:  "hi-IN",
			Section:       l.Section,
			LanguageValue: "",
		},
		{
			LableCode:     l.LableCode,
			LangId:        lastInsertedID,
			LanguageCode:  "en-GB",
			Section:       l.Section,
			LanguageValue: "",
		},
	}

	for _, langLabel := range languageLabels {
		_, err := db.Insert(&langLabel)
		if err != nil {
			return "", err
		}
	}
	return res.LableCode, nil
}

func ExistsLanguageLable(lable_code string) int {
	db := orm.NewOrm()
	var lables LanguageLableLang
	err := db.Raw(`SELECT lang_id FROM language_lable_lang WHERE lable_code = ?`, lable_code).QueryRow(&lables)
	if err != nil {
		return 0
	}
	return 1
}

func FetchAllLabels() (interface{}, error) {
	db := orm.NewOrm()
	var labelsList []orm.Params

	_, err := db.Raw(`SELECT lable_code, language_value, language_code,section FROM language_lable`).Values(&labelsList)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	return labelsList, nil
}

func FetchAllDefaultlables() (interface{}, error) {
	db := orm.NewOrm()
	var labelsList []orm.Params

	_, err := db.Raw(`SELECT lable_code, language_value, language_code,section FROM language_lable_lang`).Values(&labelsList)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	return labelsList, nil
}
