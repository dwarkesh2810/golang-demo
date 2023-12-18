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

func LoadConfig() error {
	return loadConfig("ini", "/home/silicon/Dwarkesh/golang/golang-demo/golang-demo/conf/app.conf")
}

func GetValue(key string) (string, error) {
	LoadConfig()
	val, err := Config.String(key)
	if err != nil {
		return "", err
	}
	return val, nil
}

func GetConfigMap() error {
	a, err := Config.GetSection("golang-demo")

	if err != nil {
		return err
	}

	ConfigMaps = a
	return nil
}
