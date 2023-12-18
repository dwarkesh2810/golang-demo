package dto

type LanguageLableInsert struct {
	LableCode string `json:"lable_code" form:"lable_code"`
	LangValue string `json:"lang_value" form:"lang_value"`
	Language  string `json:"language" form:"language"`
	Section   string `json:"section" form:"section"`
}
