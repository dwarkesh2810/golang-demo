package initialization

import (
	"github.com/dwarkesh2810/golang-demo/conf"
	"github.com/dwarkesh2810/golang-demo/controllers"
	"github.com/dwarkesh2810/golang-demo/models"
	ad "github.com/dwarkesh2810/golang-demo/pkg/admin"
	"github.com/dwarkesh2810/golang-demo/pkg/logger"
	"github.com/dwarkesh2810/golang-demo/pkg/validations"
	_ "github.com/dwarkesh2810/golang-demo/routers"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/lib/pq"
)

func Init() {
	conf.GetConfigMap()

	orm.RegisterDriver(conf.ConfigMaps["dbdriver"], orm.DRPostgres)
	orm.RegisterDataBase("default", conf.ConfigMaps["dbdriver"], conf.ConfigMaps["conn"])

	orm.RegisterModel(new(models.Users), new(models.HomePagesSettingTable), new(models.Car), new(models.LanguageLable), new(models.LanguageLableLang), new(models.EmailLogs))

	languageLablesFunc := controllers.LangLableController{}
	languageLablesFunc.FetchAllAndWriteInINIFiles()

	validations.Init()
	logger.Init()
	ad.CreateTask("EmailLog", "0 */5 * * * *", ad.SendPendingEmail)
}