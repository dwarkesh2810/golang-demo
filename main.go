package main

import (
	"github.com/dwarkesh2810/golang-demo/models"
	_ "github.com/dwarkesh2810/golang-demo/routers"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
)

func init() {
	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase("default", "postgres", "user=root password=1234 dbname=golang_demo sslmode=disable")
	orm.RegisterModel(new(models.Users), new(models.HomePagesSettingTable), new(models.Car), new(models.LanguageLable), new(models.LanguageLableLang))
	// orm.RunSyncdb("default", false, true)
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
