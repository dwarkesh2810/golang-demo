package models

import (
	"time"

	"github.com/dwarkesh2810/golang-demo/dto"
)

type Users struct {
	UserId      int `orm:"auto;pk"`
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
	Password    string
	Isverified  int
	OtpCode     string
	Role        string
	CountryId   int
	CreatedDate time.Time
	UpdatedDate time.Time `orm:"null"`
}

type Car struct {
	Id          uint        `json:"car_id" orm:"pk;auto;column(id)"`
	CarName     string      `orm:"column(car_name)"`
	CarImage    string      `orm:"null;column(car_image)" form:"file" json:"file"`
	ModifiedBy  string      `orm:"column(modified_by)"`
	Model       string      `orm:"column(model)"`
	Type        dto.CarType `orm:"column(car_type);type(enum)"`
	CreatedDate time.Time   `orm:"null;column(ctreated_date)"`
	UpdateDate  time.Time   `orm:"null;column(updated_date)"`
	CreatedBy   int
	UpdatedBy   int `orm:"null"`
}

type HomePagesSettingTable struct {
	PageSettingId int `orm:"auto"`
	Section       string
	DataType      string `orm:"size(255)"`
	UniqueCode    string
	SettingData   string `orm:"type(text)"`
	CreatedDate   time.Time
	UpdatedDate   time.Time `orm:"null"`
	CreatedBy     int
	UpdatedBy     int
}

type EnglishLanguageLable struct {
	LangId        int `orm:"auto"`
	LanguageCode  string
	LanguageValue string
	LableCode     string `orm:"unique"`
	Section       string
	CreatedBy     int
	UpdatedBy     int
	CreatedDate   time.Time
	UpdatedDate   time.Time `orm:"null;column(updated_date)"`
}

type MultiLanguageLable struct {
	LableId       int `orm:"auto"`
	LableCode     string
	LanguageValue string
	LanguageCode  string
	Section       string
	CreatedBy     int
	UpdatedBy     int
	CreatedDate   time.Time
	UpdatedDate   time.Time `orm:"null;column(updated_date)"`
}

type EmailLogs struct {
	Id      uint   `orm:"pk;auto;column(LogId)"`
	To      string `orm:"column(emailTo)"`
	Name    string `orm:"column(name)"`
	Subject string
	Body    string
	Status  string
}

type AuditLogs struct {
	LogsId      uint `orm:"pk;auto;column(LogId)"`
	UserId      uint
	Action      string
	UserIp      string
	Discription string
	EndPoints   string
	CreatedDate time.Time
}
