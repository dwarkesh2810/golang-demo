package logger

import (
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/beego/i18n"
	"github.com/dwarkesh2810/golang-demo/models"
)

var log *logs.BeeLogger

func Init() {
	log = logs.NewLogger(100)
	logs.SetLogger(logs.AdapterMultiFile, `{"filename":"pkg/logger/logs/app.log","maxlines":10000000,"separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`)

}

func Error(message string, err error) {
	log.Error("%s -------------------------------------------> %v", message, err)
}

func Info(message string, info string) {
	log.Info("%s---------------------------------------------> %s", message, info)
}

func Debug(message string) {
	log.Debug(message)
}

func Alert(message string, alert string) {
	log.Alert("%s -------------------------------------------> %v", message, alert)
}

func Warning(message string, warning string) {
	log.Warning("%s ------------------------------------------> %v", message, warning)
}

func InsertAuditLogs(c *context.Context, discription string, userId uint) error {

	action := c.Request.Method
	ip := c.Input.IP()
	endPoint := c.Request.URL.Path

	log := models.AuditLogs{
		UserId:      userId,
		Action:      action,
		UserIp:      ip,
		Discription: discription,
		EndPoints:   endPoint,
		CreatedDate: time.Now(),
	}
	return models.InsertAuditLog(log)
}

func LogMessage(c *context.Context, key string) string {
	return i18n.Tr("en-US", key)
}