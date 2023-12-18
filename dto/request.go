package dto

import "github.com/dgrijalva/jwt-go"

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
	FirstName   string `json:"first_name" valid:"MaxSize(255);MinSize(3);Required;Match(/^[a-zA-Z][a-zA-Z0-9._-]{0,31}$/)"`
	LastName    string `json:"last_name" valid:"MaxSize(255);MinSize(3);Required;Match(/^[a-zA-Z][a-zA-Z0-9._-]{0,31}$/)"`
	Email       string `json:"email" valid:"MaxSize(255);Required;Email"`
	Country     int    `json:"country_id" valid:"Required"`
	Role        string `json:"role" valid:"MaxSize(255);Required"`
	PhoneNumber string `json:"phone_number" valid:"Required;IsMobile"`
	Password    string `json:"password" valid:"MaxSize(25);MinSize(6);Required"`
}

type UpdateUserRequest struct {
	Id          int    `json:"user_id" valid:"Required"`
	FirstName   string `json:"first_name" valid:"MaxSize(255);MinSize(3);Required;Match(/^[a-zA-Z][a-zA-Z0-9._-]{0,31}$/)"`
	LastName    string `json:"last_name" valid:"MaxSize(255);MinSize(3);Required;Match(/^[a-zA-Z][a-zA-Z0-9._-]{0,31}$/)"`
	Email       string `json:"email" valid:"MaxSize(255);Required;Email"`
	PhoneNumber string `json:"phone_number"  valid:"IsMobile;Required"`
	Country     int    `json:"country_id" valid:"Required"`
	Role        string `json:"role" valid:"MaxSize(255);Required"`
}

type ResetUserPassword struct {
	CurrentPass string `json:"current_password"`
	NewPass     string `json:"new_password"`
	ConfirmPass string `json:"confirm_password"`
}
type VerifyEmailOTPRequest struct {
	Username string `json:"username"`
	Otp      string `json:"otp"`
}

type SendOtpData struct {
	Username string `json:"username"`
}

type ResetUserPasswordOtp struct {
	Username string `json:"username"`
	Otp      string `json:"otp"`
	NewPass  string `json:"new_password"`
}
