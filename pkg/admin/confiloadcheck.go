package admin

import (
	"errors"

	"github.com/dwarkesh2810/golang-demo/conf"
	_ "github.com/dwarkesh2810/golang-demo/routers"

	_ "github.com/lib/pq"
)

type ConfigFileCheck struct {
}

func (lc *ConfigFileCheck) Check() error {
	if lc.isSuccess() {
		return nil
	} else {
		return errors.New("not translate messages")
	}
}

func (lc *ConfigFileCheck) isSuccess() bool {
	err := conf.GetConfigMap()
	return !(err != nil)
}
