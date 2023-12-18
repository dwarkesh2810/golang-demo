package models

import (
	"errors"

	"github.com/beego/beego/v2/client/orm"
)

func GetUserByEmail(username string) (Users, error) {
	o := orm.NewOrm()
	var user Users
	// orm.Debug = true
	num, err := o.QueryTable(new(Users)).SetCond(orm.NewCondition().Or("phone_number", username).Or("email", username)).All(&user)
	if err != nil {
		return user, errors.New("DATABASE_ERROR")
	}
	if num == 0 {
		return user, errors.New("LOGIN_ERROR")
	}
	return user, nil
}

func LoginUser(username string, pass string) (Users, error) {
	o := orm.NewOrm()
	var user Users
	num, err := o.QueryTable(new(Users)).SetCond(orm.NewCondition().Or("phone_number", username).Or("email", username)).Filter("password", pass).All(&user)
	if err != nil {
		return user, errors.New("DATABASE_ERROR")
	}
	if num == 0 {
		return user, errors.New("DATABASE_ERROR")
	}
	return user, nil
}
