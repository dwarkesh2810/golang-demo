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
	_, err := o.QueryTable(new(Users)).SetCond(orm.NewCondition().Or("phone_number", username).Or("email", username)).All(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func LoginUser(username string, pass string) (Users, error) {
	o := orm.NewOrm()
	var user Users
	_, err := o.QueryTable(new(Users)).SetCond(orm.NewCondition().Or("phone_number", username).Or("email", username)).Filter("password", pass).All(&user)
	if err != nil {
		return user, err
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
		return user, err
	}
	if num == 0 {
		return user, errors.New("DATABASE_ERROR")
	}
	return user, nil
}
func UpadteOtpForEmail(id int, otp string) (string, error) {
	o := orm.NewOrm()
	var user = Users{UserId: id, OtpCode: otp, Isverified: 0}
	num, err := o.Update(&user, "otp_code", "isverified")
	if err != nil {
		return "num", err
	}
	if num == 0 {
		return "user", errors.New("DATABASE_ERROR")
	}
	return "otpsent", nil
}
func GetUserDetails(id interface{}) (Users, error) {
	o := orm.NewOrm()
	var user Users
	num, err := o.QueryTable(new(Users)).Filter("user_id", id).All(&user, "first_name", "last_name", "email", "phone_number", "password")
	if err != nil {
		return user, err
	}
	if num == 0 {
		return user, errors.New("DATABASE_ERROR")
	}
	return user, nil
}
func UpdateUser(Data dto.UpdateUserRequest) (Users, error) {
	var user = Users{
		UserId:      Data.Id,
		FirstName:   Data.FirstName,
		LastName:    Data.LastName,
		CountryId:   Data.Country,
		Email:       Data.Email,
		Role:        Data.Role,
		UpdatedDate: time.Now(),
		PhoneNumber: Data.PhoneNumber,
	}
	o := orm.NewOrm()
	num, err := o.Update(&user, "user_id", "first_name", "last_name", "country_id", "email", "role", "updated_date", "phone_number")
	if err != nil {
		return user, err
	}
	if num == 0 {
		return user, errors.New("DATABASE_ERROR")
	}
	return user, nil
}

func ResetPassword(Password string, id int) (interface{}, error) {
	pass, err := helpers.HashData(Password)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	var user = Users{UserId: id, Password: pass}
	num, err := o.Update(&user, "password")
	if err != nil {
		return num, errors.New("DATABASE_ERROR")
	}
	return "PASSWORD_RESET", nil
}
func DeleteUser(id int) (string, error) {
	o := orm.NewOrm()
	var user = Users{UserId: id}
	num, err := o.Delete(&user)
	if err != nil {
		return "", err
	}
	if num == 0 {
		return "", errors.New("DATABASE_ERROR")
	}
	return "User deleted success", nil
}
func GetEmailOTP(username string, otp string) (Users, error) {
	o := orm.NewOrm()
	var user Users
	_, err := o.QueryTable(new(Users)).SetCond(orm.NewCondition().Or("phone_number", username).Or("email", username)).Filter("otp_code", otp).All(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func UpdateIsVerified(id int) error {
	o := orm.NewOrm()
	var user = Users{UserId: id, Isverified: 1, UpdatedDate: time.Now()}
	num, err := o.Update(&user, "isverified", "updated_date")
	if err != nil {
		return err
	}
	if num == 0 {
		return errors.New("DATABASE_ERROR")
	}
	return nil
}

func SearchUser(search string) ([]Users, error) {
	o := orm.NewOrm()
	var user []Users
	num, err := o.QueryTable(new(Users)).SetCond(orm.NewCondition().Or("first_name__icontains", search).Or("email__icontains", search).Or("last_name__icontains", search).Or("role__icontains", search)).All(&user)
	if err != nil {
		return nil, err
	}
	if num == 0 {
		return nil, errors.New("DATABASE_ERROR")
	}
	return user, nil
}
func VerifyEmail(email string, name string) (string, error) {
	OTP := helpers.GenerateUniqueCodeString(4)
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
			return "", err
		}
		return "failed to send otp", nil
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
		return "", err
	}
	res, err := UpadteOtpForEmail(output.UserId, OTP)
	if err != nil {
		return "", err
	}
	return res, nil
}

func FetchUsers(current_page, pageSize int) ([]orm.Params, map[string]interface{}, error) {
	tableName := "users"
	query := `SELECT u.first_name , u.last_name, u.email, u.phone_number
	FROM users as u
	ORDER BY u.user_id
	LIMIT ? OFFSET ?
`
	result_data, pagination, errs := helpers.FetchDataWithPaginations(current_page, pageSize, tableName, query)
	if errs != nil {
		return nil, nil, errs
	}
	return result_data, pagination, nil
}
