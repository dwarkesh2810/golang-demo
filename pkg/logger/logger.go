package logger

import (
	"github.com/beego/beego/v2/core/logs"
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
