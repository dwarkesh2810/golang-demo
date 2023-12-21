package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/dwarkesh2810/golang-demo/controllers"
)

var user_ctrl = controllers.UserController{}

// Please Run first TestLoginUser Func to generate and test login func
func TestLoginUser(t *testing.T) {
	t.Run("Login User", func(t *testing.T) {
		endPoint := "/v1/user/login/"
		var jsonStr = []byte(`{
			"username" : "rideshnath.siliconithub@gmail.com",
			"password": "123456"
		}`)
		r, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router := beego.NewControllerRegister()
		router.Add(endPoint, &user_ctrl, beego.WithRouterMethods(&user_ctrl, "post:Login"))
		router.ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}
		t.Log(w.Body)
	})
}

func TestNewUser(t *testing.T) {
	t.Run("new user", func(t *testing.T) {
		endPoint := "/v1/user/register"
		var jsonStr = []byte(`{"first_name":"Dwarkesh","last_name":"patel","email":"dwarkeshpatel.siliconithub@gmail.com","country_id":88,"role":"Developer","phone_number":"9123456789","password":"123456"}`)
		tokan := LoginTokan()
		r, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		w := RunControllerRoute(endPoint, r, &user_ctrl, tokan, "post:RegisterNewUser")
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}
		var resultMap map[string]interface{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &resultMap)
		if err != nil {
			t.Fatalf(err.Error())
			return
		}
		if resultMap["status"] == 0 {
			t.Fatalf("Error in register new user")
		}
		t.Log(w.Body)
	})
}

func TestVerifyUserEmail(t *testing.T) {
	t.Run("Verify email otp", func(t *testing.T) {
		endPoint := "/v1/user/secure/VerifyEmailOTP"
		// enter otp after getting otp from above test forget password
		var jsonStr = []byte(`{
			"username" : "rideshnath.siliconithub@gmail.com",
			"otp":"0703"
		}`)
		tokan := LoginTokan()
		r, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		w := RunControllerRoute(endPoint, r, &user_ctrl, tokan, "post:VerifyEmailOTP")
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}
		var resultMap map[string]interface{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &resultMap)
		if err != nil {
			t.Fatalf(err.Error())
			return
		}
		if resultMap["ResStatus"] == 0 {
			t.Fatalf("Error in register new user")
		}
		t.Log(w.Body)
	})
}
func TestGetusers(t *testing.T) {
	t.Run("Get All User", func(t *testing.T) {
		endPoint := "/v1/user/secure/users"
		tokan := LoginTokan()
		var jsonStr = []byte(`{
			"page_size" : 5,
			"open_page": 1
		}`)
		r, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		w := RunControllerRoute(endPoint, r, &user_ctrl, tokan, "post:GetAllUsers")
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}
		var resultMap map[string]interface{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &resultMap)
		if err != nil {
			t.Fatalf(err.Error())
			return
		}
		if resultMap["ResStatus"] == 0 {
			t.Fatalf("Data Not found")
		}
		t.Log(w.Body)
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("Update user", func(t *testing.T) {
		endPoint := "/v1/user/secure/update"
		var jsonStr = []byte(`{"user_id":2,"first_name":"Dwarkesh","last_name":"patel","email":"dwarkeshpatel@gmail.com","country_id":25,"role":"Developer","phone_number":"9123456789"}`)
		tokan := LoginTokan()
		r, err := http.NewRequest("PUT", endPoint, bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		w := RunControllerRoute(endPoint, r, &user_ctrl, tokan, "put:UpdateUser")
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}
		var resultMap map[string]interface{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &resultMap)
		if err != nil {
			t.Fatalf(err.Error())
			return
		}
		if resultMap["ResStatus"] == 0 {
			t.Fatalf("Some Error in update in user")
		}
		t.Log(w.Body)
	})
}

// Please run Indiual test t.run becouse of it is send otp on register email You can use phone number or email in the place of username
func TestForgotPassword(t *testing.T) {
	t.Run("Send otp", func(t *testing.T) {
		endPoint := "/v1/user/secure/forgot_pass"
		var jsonStr = []byte(`{
			"username" : "rideshnath.siliconithub@gmail.com"
		}`)
		tokan := LoginTokan()
		r, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		r.Header.Set("Content-Type", "application/json")
		w := RunControllerRoute(endPoint, r, &user_ctrl, tokan, "post:ForgotPassword")
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}
		t.Log(w.Body)
	})
	t.Run("Verify email otp", func(t *testing.T) {
		endPoint := "/v1/user/secure/reset_pass_otp"
		// enter otp after getting otp from above test forget password
		var jsonStr = []byte(`{
			"username" : "rideshnath.siliconithub@gmail.com",
			"otp":"0703",
			"new_password":"123456"
		}`)
		tokan := LoginTokan()
		r, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		w := RunControllerRoute(endPoint, r, &user_ctrl, tokan, "post:VerifyOtpResetpassword")
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}
		t.Log(w.Body)
	})
}

func TestSerchUser(t *testing.T) {
	t.Run("Search users", func(t *testing.T) {
		endPoint := "/v1/user/secure/search"
		var jsonStr = []byte(`{
			"search" : "d"
		}`)
		tokan := LoginTokan()
		r, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		w := RunControllerRoute(endPoint, r, &user_ctrl, tokan, "post:SearchUser")
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}
		var resultMap map[string]interface{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &resultMap)
		if err != nil {
			t.Fatalf(err.Error())
			return
		}
		if resultMap["ResStatus"] == 0 {
			t.Fatalf("Some Error in update in user")
		}
		t.Log(w.Body)
	})
}
