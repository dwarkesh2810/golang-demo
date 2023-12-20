package test

import (
	"testing"

	"github.com/dwarkesh2810/golang-demo/dto"
	"github.com/dwarkesh2810/golang-demo/models"
)

func TestUserModels(t *testing.T) {
	t.Run("Register new user", func(t *testing.T) {
		// TruncateTable("users")
		var user = dto.NewUserRequest{
			FirstName:   "Ridesh",
			LastName:    "Nath",
			Email:       "rideshnath.siliconithub@gmail.com",
			PhoneNumber: "1234567890",
			Role:        "developer",
			Password:    "123456",
		}
		data, err := models.InsertNewUser(user)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		if data.UserId == 0 {
			t.Errorf("error in register")
			return
		}
		t.Log(data)
	})
	t.Run("Get All User", func(t *testing.T) {
		data, _, err := models.FetchUsers(1, 5)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		if len(data) == 0 {
			t.Errorf("error :- No data found")
		}
		t.Log(data)
	})
	t.Run("Get UserByEmail", func(t *testing.T) {
		data, err := models.GetUserByEmail("rideshnath.siliconithub@gmail.com")
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		if data.Email == "" {
			t.Errorf("error :- No data found")
		}
		t.Log(data)
	})
	t.Run("Get UserId", func(t *testing.T) {
		data, err := models.GetUserDetails(1)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		if data.Email == "" {
			t.Errorf("error :- No data found")
		}
		t.Log(data)
	})
	t.Run("Login user", func(t *testing.T) {
		username := "1234567890"
		password := "123456"
		data, err := models.LoginUser(username, password)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		if data.Email == "" {
			t.Errorf("error :- No data found")
		}
		t.Log(data)
	})
	t.Run("Update user", func(t *testing.T) {
		var user = dto.UpdateUserRequest{
			Id:          1,
			FirstName:   "Devendra",
			LastName:    "pohekar",
			Email:       "rideshnath.siliconithub@gmail.com",
			PhoneNumber: "1234567890",
			Role:        "deve",
			Country:     5,
		}
		data, err := models.UpdateUser(user)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		t.Log(data)
	})
	t.Run("serch user", func(t *testing.T) {
		search := "dev"
		data, err := models.SearchUser(search)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		if len(data) == 0 {
			t.Errorf("error :- No data found")
		}
		t.Log(data)
	})
	t.Run("Delete User", func(t *testing.T) {
		data, err := models.DeleteUser(10)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		if data != "User deleted success" {
			t.Errorf("error in delete")
		}
		t.Log(data)
	})
}
