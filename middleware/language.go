package middleware

import (
	"github.com/beego/beego/v2/server/web/context"
	"github.com/dwarkesh2810/golang-demo/helpers"
)

func LanguageMiddlware(c *context.Context) {

	lang := helpers.GetLanguageFromMultipleSources(c)
	helpers.SetLanguage(c, lang)
}
