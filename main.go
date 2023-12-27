package main

import (
	ad "github.com/dwarkesh2810/golang-demo/pkg/admin"
	"github.com/dwarkesh2810/golang-demo/pkg/initialization"
	_ "github.com/dwarkesh2810/golang-demo/routers"
	"github.com/prometheus/client_golang/prometheus"

	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
)

var (
	httpRequestsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "beego_http_requests_total",
		Help: "Total number of HTTP requests processed by the Beego application.",
	})
)

func init() {
	initialization.Init()
	prometheus.MustRegister(httpRequestsTotal)
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	ad.Healthceck()
	beego.Run()
}
