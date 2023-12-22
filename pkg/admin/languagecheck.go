package admin

import (
	"errors"

	"github.com/dwarkesh2810/golang-demo/controllers"
	_ "github.com/dwarkesh2810/golang-demo/routers"

	_ "github.com/lib/pq"
)

type LanguageCheck struct {
}

func (lc *LanguageCheck) Check() error {
	if lc.isSuccess() {
		return nil
	} else {
		return errors.New("not translate messages")
	}
}

func (lc *LanguageCheck) isSuccess() bool {
	languageLablesFunc := controllers.LangLableController{}
	return languageLablesFunc.FetchAllAndWriteInINIFiles()
}
