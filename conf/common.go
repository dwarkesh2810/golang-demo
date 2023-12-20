package conf

import (
	"github.com/beego/beego/v2/core/config"
)

var Config config.Configer
var ConfigMaps map[string]string

func loadConfig(ext, path string) error {
	conf, err := config.NewConfig(ext, path)
	if err != nil {
		return err
	}

	Config = conf
	return nil
}

func GetConfigMap() error {
	loadConfig("ini", "conf/app.conf")
	a, err := Config.GetSection("golang-demo")

	if err != nil {
		return err
	}
	ConfigMaps = a
	return nil
}
