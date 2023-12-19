package main

import (
	"github.com/dwarkesh2810/golang-demo/conf"
	"github.com/dwarkesh2810/golang-demo/controllers"
	"github.com/dwarkesh2810/golang-demo/models"
	_ "github.com/dwarkesh2810/golang-demo/routers"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
)

func init() {

	conf.GetConfigMap()

	orm.RegisterDriver(conf.ConfigMaps["dbdriver"], orm.DRPostgres)
	orm.RegisterDataBase("default", conf.ConfigMaps["dbdriver"], conf.ConfigMaps["conn"])
	orm.RegisterModel(new(models.Users), new(models.HomePagesSettingTable), new(models.Car), new(models.LanguageLable), new(models.LanguageLableLang), new(models.EmailLogs))
	languageLablesFunc := controllers.LangLableController{}
	languageLablesFunc.FetchAllAndWriteInINIFiles()
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
