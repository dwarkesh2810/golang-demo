package test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"testing"

	"github.com/dwarkesh2810/golang-demo/controllers"
)

var car_ctrl = controllers.CarController{}

func TestGetCars(t *testing.T) {
	t.Run("Get All Cars", func(t *testing.T) {
		endPoint := "/v1/car/cars"
		tokan := LoginTokan()
		r, err := http.NewRequest("GET", endPoint, nil)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		w := RunControllerRoute(endPoint, r, &car_ctrl, tokan, "get:GetAllCars")
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
			t.Fatalf("Error in get car")
		}
		t.Log(w.Body)
	})

	t.Run("Get Single Car", func(t *testing.T) {
		endPoint := "/v1/car/"
		tokan := LoginTokan()
		var jsonStr = []byte(`{"car_id":3}`)
		r, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		w := RunControllerRoute(endPoint, r, &car_ctrl, tokan, "post:GetSingleCar")
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
			t.Fatalf("Error in get car")
		}
		t.Log(w.Body)
	})
}

func TestInsertNewcar(t *testing.T) {
	endPoint := "v1/car/create"
	tokan := LoginTokan()
	var jsonStr = []byte(`{"car_name" : "Thar","model": "LLL","modified_by": "mahindara","type":"sedan"}`)
	r, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	w := RunControllerRoute(endPoint, r, &car_ctrl, tokan, "post:AddNewCar")
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}
	t.Log(w.Body)
}

func TestUpdatecar(t *testing.T) {
	t.Run("Update Car", func(t *testing.T) {
		endPoint := "v1/car/update"
		tokan := LoginTokan()
		var jsonStr = []byte(`{"car_id":4,"car_name" : "Thar","model": "4*4","modified_by": "mahindara","type":"sedan"}`)
		r, err := http.NewRequest("PUT", endPoint, bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		w := RunControllerRoute(endPoint, r, &car_ctrl, tokan, "put:UpdateCar")
		log.Print(w.Code)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}
		t.Log(w.Body)
	})

}

func TestRemoveCar(t *testing.T) {
	endPoint := "/v1/car/delete"
	tokan := LoginTokan()
	var jsonStr = []byte(`{"car_id":4}`)
	r, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	w := RunControllerRoute(endPoint, r, &car_ctrl, tokan, "post:DeleteCar")
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
		t.Fatalf("Some Error in Delete car")
	}
	t.Log(w.Body)
}
