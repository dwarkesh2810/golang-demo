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
	conf.LoadEnv(".")
	orm.RegisterDriver(conf.Env.DbDriver, orm.DRPostgres)
	orm.RegisterDataBase("default", conf.Env.DbDriver, conf.Env.ConnString)

	orm.RegisterModel(new(models.Users), new(models.HomePagesSettingTable), new(models.Car), new(models.MultiLanguageLable), new(models.EnglishLanguageLable), new(models.EmailLogs), new(models.AuditLogs))

	// orm.RunSyncdb("default", false, true)

	languageLablesFunc := controllers.LangLableController{}
	languageLablesFunc.FetchAllAndWriteInINIFiles()

	validations.Init()
	logger.Init()
	ad.CreateTask("EmailLog", "0 */5 * * * *", ad.SendPendingEmail)
}
