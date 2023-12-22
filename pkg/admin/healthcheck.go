package admin

import "github.com/beego/beego/v2/core/admin"

func Healthceck() {
	admin.AddHealthCheck("Database connection", &DatabaseCheck{})
	admin.AddHealthCheck("Language lables", &LanguageCheck{})
	admin.AddHealthCheck("Load config", &ConfigFileCheck{})
}
