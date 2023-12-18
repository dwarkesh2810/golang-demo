package dto

type HomeSeetingInsert struct {
	Section     string `json:"section" form:"section"`
	DataType    string `json:"data_type" form:"data_type"`
	SettingData string `json:"setting_data" form:"setting_data"`
	LangKey     string `json:"lang_key" form:"lang_key"`
}

type HomeSeetingUpdate struct {
	Section     string `json:"section" form:"section"`
	DataType    string `json:"data_type" form:"data_type"`
	SettingData string `json:"setting_data" form:"setting_data"`
	SettingId   int    `json:"setting_id" form:"setting_id"`
	LangKey     string `json:"lang_key" form:"lang_key"`
}

type HomeSeetingDelete struct {
	Section   string `json:"section" form:"section"`
	SettingId int    `json:"setting_id" form:"setting_id"`
	LangKey   string `json:"lang_key" form:"lang_key"`
}

type HomeSeetingSearch struct {
	SettingId int    `json:"setting_id" form:"setting_id"`
	LangKey   string `json:"lang_key" form:"lang_key"`
	PageSize  int    `json:"page_size" form:"page_size"`
	OpenPage  int    `json:"open_page" form:"open_page"`
}

type HomeSeetingSearchFilter struct {
	SettingId   int    `json:"setting_id" form:"setting_id"`
	SearchParam string `json:"search_string" form:"search_string"`
	LangKey     string `json:"lang_key" form:"lang_key"`
	PageSize    int    `json:"page_size" form:"page_size"`
	OpenPage    int    `json:"open_page" form:"open_page"`
}
