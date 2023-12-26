package main

import (
	ad "github.com/dwarkesh2810/golang-demo/pkg/admin"
	"github.com/dwarkesh2810/golang-demo/pkg/initialization"
	_ "github.com/dwarkesh2810/golang-demo/routers"

	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
)

func init() {
	initialization.Init()
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	ad.Healthceck()
	beego.Run()
}
