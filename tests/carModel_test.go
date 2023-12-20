package test

import (
	"testing"

	"github.com/dwarkesh2810/golang-demo/dto"
	"github.com/dwarkesh2810/golang-demo/models"
)

func TestCarModels(t *testing.T) {
	t.Run("Register New car", func(t *testing.T) {
		// TruncateTable("car")
		car := dto.GetNewCarRequest{
			CarName:    "swift",
			ModifiedBy: "suzuki",
			Model:      "swift dzire",
			Type:       "sedan",
			CarImage:   "swiftImage",
		}
		data, err := models.InsertNewCar(car)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		if data.Id == 0 {
			t.Errorf("Error in Insert")
		}
		t.Log(data)
	})

	t.Run("Update  car", func(t *testing.T) {
		// TruncateTable("car")
		car := dto.UpdateCarRequest{
			Id: 4,
			CarName:    "swift",
			ModifiedBy: "suzuki",
			Model:      "swift dzire",
			Type:       "sedan",
			CarImage:   "No Image",
		}
		data, err := models.UpdateCar(car)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		if data.Id == 0 {
			t.Errorf("Error in Update")
		}
		t.Log(data)
	})
	t.Run("Get All cars", func(t *testing.T) {
		data, _, err := models.FetchCars(1, 5)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		if len(data) == 0 {
			t.Errorf("error in fetch data")
		}
		t.Log(data)
	})
	t.Run("Get car", func(t *testing.T) {
		data, err := models.GetSingleCar(3)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		if data.CarName == "" {
			t.Errorf("error in fetch data")
		}
		t.Log(data)
	})
	t.Run("Delete car", func(t *testing.T) {
		data, err := models.DeleteCar(22)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		t.Log(data)
	})
}
