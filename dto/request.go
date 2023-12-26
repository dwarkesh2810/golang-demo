package dto

import (
	"github.com/dgrijalva/jwt-go"
)

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JwtClaim struct {
	Email string
	ID    int
	jwt.StandardClaims
}

type NewUserRequest struct {
	FirstName   string `json:"first_name" valid:"MaxSize(255);MinSize(3);Required"`
	LastName    string `json:"last_name" valid:"MaxSize(255);MinSize(3);Required"`
	Email       string `json:"email" valid:"MaxSize(255);Required;Email"`
	Country     int    `json:"country_id" valid:"Required"`
	Role        string `json:"role" valid:"MaxSize(255);Required"`
	PhoneNumber string `json:"phone_number" valid:"Required;InMobile"`
	Password    string `json:"password" valid:"MaxSize(25);MinSize(6);Required"`
}

type UpdateUserRequest struct {
	Id          int    `json:"user_id" valid:"Required"`
	FirstName   string `json:"first_name" valid:"MaxSize(255);MinSize(3);Required"`
	LastName    string `json:"last_name" valid:"MaxSize(255);MinSize(3);Required"`
	Email       string `json:"email" valid:"MaxSize(255);Required;Email"`
	PhoneNumber string `json:"phone_number"  valid:"InMobile;Required"`
	Country     int    `json:"country_id" valid:"Required"`
	Role        string `json:"role" valid:"MaxSize(255);Required"`
}

type ResetUserPassword struct {
	CurrentPass string `json:"current_password" valid:"Required;MaxSize(25);MinSize(6)"`
	NewPass     string `json:"new_password" valid:"Required;MaxSize(25);MinSize(6)"`
	ConfirmPass string `json:"confirm_password" valid:"Required;MaxSize(25);MinSize(6)"`
}
type VerifyEmailOTPRequest struct {
	Username string `json:"username" valid:"Required"`
	Otp      string `json:"otp" valid:"Required"`
}

type SendOtpData struct {
	Username string `json:"username"`
}

type ResetUserPasswordOtp struct {
	Username string `json:"username" valid:"Required"`
	Otp      string `json:"otp" valid:"Required"`
	NewPass  string `json:"new_password" valid:"Required;MaxSize(25);MinSize(6)"`
}

type SearchRequest struct {
	Search   string `json:"search" form:"search" valid:"Required;MinSize(3)"`
	PageSize int    `json:"page_size" form:"page_size"`
	OpenPage int    `json:"open_page" form:"open_page"`
}

type LanguageLableInsert struct {
	LableCodes            string `json:"lable_code" form:"lable_code" valid:"Required;MinSize(3)"`
	Sections              string `json:"section" form:"section" valid:"Required;MinSize(3)"`
	ENGLangValues         string `json:"ENGlang_value" form:"ENGlang_value" valid:"Required;MinSize(3)"`
	OtherINILanguageCodes string `json:"lang_ini_code" form:"lang_ini_code" valid:"Required;MinSize(5);ValidIni"`
	OtherLangValues       string `json:"otherlang_value" form:"otherlang_value" valid:"Required;MinSize(3)"`
}

type LanguageLableUpdate struct {
	LableCodes            string `json:"lable_code" form:"lable_code" valid:"Required;MinSize(3)"`
	ENGLangValues         string `json:"ENGlang_value" form:"ENGlang_value" valid:"Required;MinSize(3)"`
	OtherINILanguageCodes string `json:"lang_ini_code" form:"lang_ini_code" valid:"Required;MinSize(5);ValidIni"`
	OtherLangValues       string `json:"otherlang_value" form:"otherlang_value" valid:"Required;MinSize(3)"`
}

type HomeSeetingInsert struct {
	Section     string `json:"section" form:"section" valid:"Required;MinSize(3)"`
	DataType    string `json:"data_type" form:"data_type" valid:"Required; WithIn"`
	SettingData string `json:"setting_data" form:"setting_data"`
	LangKey     string `json:"lang_key" form:"lang_key"`
}

type HomeSeetingUpdate struct {
	Section     string `json:"section" form:"section" valid:"Required;MinSize(3)"`
	DataType    string `json:"data_type" form:"data_type" valid:"Required; WithIn"`
	SettingData string `json:"setting_data" form:"setting_data" valid:"Required"`
	SettingId   int    `json:"setting_id" form:"setting_id" valid:"Required"`
	LangKey     string `json:"lang_key" form:"lang_key"`
}

type HomeSeetingDelete struct {
	Section   string `json:"section" form:"section" valid:"Required;MinSize(3)"`
	SettingId int    `json:"setting_id" form:"setting_id" valid:"Required"`
	LangKey   string `json:"lang_key" form:"lang_key"`
}

type HomeSeetingSearch struct {
	SettingId int    `json:"setting_id" form:"setting_id"`
	LangKey   string `json:"lang_key" form:"lang_key"`
	PageSize  int    `json:"page_size" form:"page_size"`
	OpenPage  int    `json:"open_page" form:"open_page"`
}

type HomeSeetingSearchFilter struct {
	UniqueCode          string `json:"unique_code" form:"unique_code"`
	SearchParam         string `json:"search_string" form:"search_string"`
	LangKey             string `json:"lang_key" form:"lang_key"`
	PageSize            int    `json:"page_size" form:"page_size"`
	OpenPage            int    `json:"open_page" form:"open_page"`
	DataType            string `json:"data_type" form:"data_type"`
	Section             string `json:"section" form:"section"`
	SettingData         string `json:"setting_data" form:"setting_data"`
	ApplySearchPosition string `json:"apply_position" form:"apply_position"`
}

type FileType struct {
	FileType    string `json:"file_type" form:"file_type" valid:"Required; ValidType"`
	Limit       int    `json:"limit" form:"limit"`
	SratingFrom int    `json:"starting_from" form:"starting_from"`
}

type PaginationReq struct {
	PageSize int `json:"page_size" form:"page_size"`
	OpenPage int `json:"open_page" form:"open_page"`
}

type CarType string

const (
	Sedan     CarType = "sedan"
	Hatchback CarType = "hatchback"
	SUV       CarType = "SUV"
)

type GetNewCarRequest struct {
	CarName    string  `json:"car_name" form:"car_name" valid:"MaxSize(255);MinSize(3);Required"`
	CarImage   string  `json:"car_imag" form:"file"`
	ModifiedBy string  `json:"modified_by" form:"modified_by" valid:"MaxSize(255);MinSize(3);Required"`
	Model      string  `json:"model" form:"model" valid:"MaxSize(255);MinSize(3);Required"`
	Type       CarType `json:"type" form:"type"`
	CreatedBy  int
}

type UpdateCarRequest struct {
	Id         uint    `json:"car_id" form:"car_id" valid:"Required"`
	CarName    string  `json:"car_name" form:"car_name" valid:"MaxSize(255);MinSize(3)"`
	CarImage   string  `json:"car_imag" form:"file"`
	ModifiedBy string  `json:"modified_by" form:"modified_by" valid:"MaxSize(255);MinSize(3)"`
	Model      string  `json:"model" form:"model" valid:"MaxSize(255);MinSize(3)"`
	Type       CarType `json:"type" form:"type"`
	UpdatedBy  int
}

type GetcarRequest struct {
	Id uint `json:"car_id" valid:"Required"`
}

type CountryWiseState struct {
	CountryId string `json:"country_id" form:"country_id" valid:"Required"`
	PageSize  int    `json:"page_size" form:"page_size"`
	OpenPage  int    `json:"open_page" form:"open_page"`
}

type CitySearch struct {
	PageSize    int    `json:"page_size" form:"page_size"`
	OpenPage    int    `json:"open_page" form:"open_page"`
	SearchParam string `json:"search_param" form:"search_param"`
	CountryId   int    `json:"country_id" form:"country_id"`
	StateId     int    `json:"state_id" form:"state_id"`
	CountryName string `json:"country_name" form:"country_name"`
	StateName   string `json:"state_name" form:"state_name"`
}

type CountrySearch struct {
	PageSize int `json:"page_size" form:"page_size"`
	OpenPage int `json:"open_page" form:"open_page"`

	SearchParam string `json:"search_param" form:"search_param"`
}

type GetStateRequest struct {
	Id int `json:"state_id" form:"state_id" valid:"Required"`
}

type GetCountryRequest struct {
	Id int `json:"country_id" form:"country_id" valid:"Required"`
}

type LanguageLableInsertNew struct {
	LableCodes    string `json:"lable_code" form:"lable_code" valid:"Required;MinSize(2)"`
	Sections      string `json:"section" form:"section" valid:"Required;MinSize(3)"`
	ENGLangValues string `json:"ENGlang_value" form:"ENGlang_value" valid:"Required;MinSize(3)"`
}
