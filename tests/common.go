package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/dwarkesh2810/golang-demo/controllers"
	"github.com/dwarkesh2810/golang-demo/middleware"
	"github.com/dwarkesh2810/golang-demo/models"
	"github.com/dwarkesh2810/golang-demo/validations"
)

func init() {
	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase("default", "postgres", "user=root password=1234 dbname=golang_demo host=192.168.1.176 sslmode=disable")
	orm.RegisterModel(new(models.Users), new(models.HomePagesSettingTable), new(models.Car), new(models.LanguageLable), new(models.LanguageLableLang), new(models.EmailLogs))
	// orm.RunSyncdb("default", false, true)
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
	validations.Init()
}

func RunControllerRoute(endPoint string, r *http.Request, ctrl beego.ControllerInterface, tokan string, methodFuction string) *httptest.ResponseRecorder {
	r.Header.Set("Authorization", tokan)
	w := httptest.NewRecorder()
	router := beego.NewControllerRegister()
	router.InsertFilter(endPoint, beego.BeforeRouter, middleware.JWTMiddleware, beego.WithCaseSensitive(false))
	router.Add(endPoint, ctrl, beego.WithRouterMethods(ctrl, methodFuction))
	router.ServeHTTP(w, r)
	return w
}

func TruncateTable(tableName string) {
	o := orm.NewOrm()
	_, err := o.Raw("TRUNCATE TABLE " + tableName).Exec()

	if err != nil {
		fmt.Println("Failed to truncate table:", err)
		return
	}
	orm.NewOrm().Raw(`SELECT setval('"` + tableName + `_user_id_seq"', 1, false)`).Exec()
}

func LoginTokan() string {
	var user_ctrl = controllers.UserController{}
	endPoint := "/v1/user/login/"
	var jsonStr = []byte(`{
				"username" : "rideshnath.siliconithub@gmail.com",
				"password": "123456"
			}`)
	r, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router := beego.NewControllerRegister()
	router.Add(endPoint, &user_ctrl, beego.WithRouterMethods(&user_ctrl, "post:Login"))
	router.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		log.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}
	var resultMap map[string]interface{}
	err = json.Unmarshal([]byte(w.Body.Bytes()), &resultMap)
	if err != nil {
		log.Print(err.Error())
		return ""
	}
	Data := resultMap["Result"]
	dataMap, ok := Data.(map[string]interface{})
	if !ok {
		log.Print("Result is not a map[string]interface{}", resultMap["Success"])
		return ""
	}
	Tokan := dataMap["Tokan"].(string)
	return Tokan
}
