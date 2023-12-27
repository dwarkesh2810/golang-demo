// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact mailto:astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/dwarkesh2810/golang-demo/controllers"
	"github.com/dwarkesh2810/golang-demo/pkg/middleware"
)

func init() {
	beego.InsertFilter("*", beego.BeforeRouter, middleware.LanguageMiddlware)
	beego.InsertFilter("*", beego.BeforeRouter, middleware.RateLimiter)
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/user",
			beego.NSInclude(&controllers.UserController{}),
			beego.NSRouter("/register", &controllers.UserController{}, "Post:RegisterNewUser"),
			beego.NSRouter("/login", &controllers.UserController{}, "post:Login"),
			beego.NSRouter("/metrics", &controllers.UserController{}, "get:Get"),
			beego.NSNamespace("/secure",
				beego.NSBefore(middleware.JWTMiddleware),
				beego.NSRouter("/users", &controllers.UserController{}, "post:GetAllUsers"),
				beego.NSRouter("/update", &controllers.UserController{}, "put:UpdateUser"),
				beego.NSRouter("/reset_pass", &controllers.UserController{}, "post:ResetPassword"),
				beego.NSRouter("/delete/:id", &controllers.UserController{}, "delete:DeleteUser"),
				beego.NSRouter("/forgot_pass", &controllers.UserController{}, "post:ForgotPassword"),
				beego.NSRouter("/reset_pass_otp", &controllers.UserController{}, "post:VerifyOtpResetpassword"),
				beego.NSRouter("/search", &controllers.UserController{}, "post:SearchUser"),
				beego.NSRouter("/verify_email_otp", &controllers.UserController{}, "post:VerifyEmailOTP"),
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
		),
		beego.NSNamespace("/car",
			beego.NSInclude(&controllers.CarController{}),
			beego.NSBefore(middleware.JWTMiddleware),
			beego.NSRouter("/", &controllers.CarController{}, "post:GetSingleCar"),
			beego.NSRouter("/cars", &controllers.CarController{}, "post:GetAllCars"),
			beego.NSRouter("/create", &controllers.CarController{}, "post:AddNewCar"),
			beego.NSRouter("/update", &controllers.CarController{}, "put:UpdateCar"),
			beego.NSRouter("/delete", &controllers.CarController{}, "delete:DeleteCar"),
			beego.NSRouter("/search_car", &controllers.CarController{}, "post:FilterCars"),
		),
		beego.NSNamespace("/country",
			beego.NSInclude(&controllers.CountryController{}),
			beego.NSRouter("/list_countries", &controllers.CountryController{}, "post:FetchCountries"),
			beego.NSRouter("/filter_country", &controllers.CountryController{}, "post:FilterCountries"),
			beego.NSRouter("/filter_city", &controllers.CountryController{}, "post:FilterCity"),
			beego.NSRouter("/search_country", &controllers.CountryController{}, "post:FilterCountry"),
			beego.NSRouter("/get_country", &controllers.CountryController{}, "post:GetCountry"),
		),
		beego.NSNamespace("/state",
			beego.NSInclude(&controllers.StateController{}),
			beego.NSRouter("/list_states", &controllers.StateController{}, "post:FetchStates"),
			beego.NSRouter("/country_wise_state", &controllers.StateController{}, "post:CountryWiseState"),
			beego.NSRouter("/search_state", &controllers.StateController{}, "post:FilterStates"),
			beego.NSRouter("/get_state", &controllers.StateController{}, "post:GetState"),
		),
		beego.NSNamespace("/language_lable",
			beego.NSInclude(&controllers.LangLableController{}),
			beego.NSRouter("/update_lang_lable", &controllers.LangLableController{}, "post:UpdateLanguageLables"),
			beego.NSRouter("/create_lang_lable", &controllers.LangLableController{}, "post:InsertLanguageLables"),
			beego.NSRouter("/import_language_lables", &controllers.LangLableController{}, "post:ReadIniFile"),
			beego.NSRouter("/lang_lable_Insert", &controllers.LangLableController{}, "post:InsertLanguageLablesUsingApi"),
			beego.NSRouter("/lang_lable_UpdateAPI", &controllers.LangLableController{}, "post:UpdateLanguageLablesAPI"),
		),
	)
	beego.AddNamespace(ns)
}
