// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/dwarkesh2810/golang-demo/controllers"
	"github.com/dwarkesh2810/golang-demo/middleware"
)

func init() {
	beego.InsertFilter("*", beego.BeforeRouter, middleware.LanguageMiddlware)
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/user",
			beego.NSInclude(&controllers.UserController{}),
			beego.NSRouter("/register", &controllers.UserController{}, "Post:RegisterNewUser"),
			beego.NSRouter("/login", &controllers.UserController{}, "post:Login"),
			beego.NSNamespace("/secure",
				beego.NSBefore(middleware.JWTMiddleware),
				// beego.NSRouter("/Users", &controllers.UserController{}, "post:GetAllUser"),
			),
		),
		beego.NSNamespace("/homepage",
			beego.NSBefore(middleware.JWTMiddleware),
			beego.NSInclude(&controllers.HomeSettingController{}),
			beego.NSRouter("/register_settings", &controllers.HomeSettingController{}, "post:RegisterSettings"),
			beego.NSRouter("/update_settings", &controllers.HomeSettingController{}, "post:UpdateSettings"),
			beego.NSRouter("/fetch_settings", &controllers.HomeSettingController{}, "post:FetchSettings"),
			beego.NSRouter("/filter_hpst", &controllers.HomeSettingController{}, "post:FiltersFetchSettings"),
			beego.NSRouter("/export", &controllers.HomeSettingController{}, "post:ExportFile"),
			beego.NSRouter("/delete_settings", &controllers.HomeSettingController{}, "post:DeleteSetting"),
			beego.NSRouter("/import", &controllers.HomeSettingController{}, "post:ImportFile"),
			beego.NSRouter("/create_lang_lable", &controllers.LangLableController{}, "post:InsertLanguageLables"),
		),
	)
	beego.AddNamespace(ns)
}
