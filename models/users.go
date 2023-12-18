package models

import (
	"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/dwarkesh2810/golang-demo/dto"
	"github.com/dwarkesh2810/golang-demo/helpers"
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

func InsertNewUser(Data dto.NewUserRequest) (Users, error) {
	o := orm.NewOrm()
	var user Users
	pass, err := helpers.HashData(Data.Password)
	if err != nil {
		return user, err
	}
	user = Users{
		FirstName:   Data.FirstName,
		LastName:    Data.LastName,
		CountryId:   Data.Country,
		Email:       Data.Email,
		PhoneNumber: Data.PhoneNumber,
		Password:    pass,
		Role:        Data.Role,
		CreatedDate: time.Now(),
	}
	num, err := o.Insert(&user)
	if err != nil {
		return user, errors.New("DATABASE_ERROR")
	}
	if num == 0 {
		return user, errors.New("DATABASE_ERROR")
	}
	return user, nil
}
func UpadteOtpForEmail(id int, otp string) (string, error) {
	o := orm.NewOrm()
	var user = Users{UserId: id, OtpCode: otp, Isverified: 0}
	num, err := o.Update(&user, "otp", "verified")
	if err != nil {
		return "num", errors.New("DATABASE_ERROR")
	}
	if num == 0 {
		return "user", errors.New("DATABASE_ERROR")
	}
	return "OTP_SENT", nil
}
func VerifyEmail(email string, name string) (string, error) {
	OTP := helpers.GenerateOtp()
	subject := "Verify your email"
	body := `<div style="font-family: Helvetica,Arial,sans-serif;min-width:1000px;overflow:auto;line-height:2">
	<div style="margin:50px auto;width:70%;padding:20px 0">
	<div style="border-bottom:1px solid #eee">
			<a href="" style="font-size:1.4em;color: #00466a;text-decoration:none;font-weight:600">Hello, I am Ridesh</a>
		</div>
		<p style="font-size:1.1em">Hi, ` + name + `</p>
		<p>Thank you for Register in this app . Use the following OTP to verify your email. OTP is valid for 5 minutes</p>
		<h2 style="background: #00466a;margin: 0 auto;width: max-content;padding: 0 10px;color: #fff;border-radius: 4px;">` + OTP + `</h2>
		<p style="font-size:0.9em;">Regards,<br />Er. Ridesh Nath</p>
		<hr style="border:none;border-top:1px solid #eee" />
		<div style="float:right;padding:8px 0;color:#aaa;font-size:0.8em;line-height:1;font-weight:300">
			<p>Ridesh Nath</p>
			<p>Burhanpur M.P</p>
			<p>India</p>
		</div>
	</div>
</div>`
	o := orm.NewOrm()
	sendemail := EmailLogs{}
	_, err := helpers.SendMailOTp(email, name, subject, body)
	if err != nil {
		sendemail = EmailLogs{
			To:      email,
			Name:    name,
			Subject: subject,
			Body:    body,
			Status:  "pending",
		}
		_, err := o.Insert(&sendemail)
		if err != nil {
			return "", errors.New("DATABASE_ERROR")
		}
	}
	sendemail = EmailLogs{
		To:      email,
		Name:    name,
		Subject: subject,
		Body:    body,
		Status:  "success",
	}
	_, err = o.Insert(&sendemail)
	if err != nil {
		return "", err
	}
	output, err := GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", errors.New("DATABASE_ERROR")
	}
	res, err := UpadteOtpForEmail(output.UserId, OTP)
	if err != nil {
		return "", err
	}
	return res, nil
}
