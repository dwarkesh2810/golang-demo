package validations

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/dwarkesh2810/golang-demo/helpers"
)

func Init() {
	validation.AddCustomFunc("InMobile", indianMobile)
	validation.AddCustomFunc("WithIn", RequiredTag)
}

func ValidErr(err []*validation.Error) []string {
	message := make([]string, 0, len(err))
	for i := range err {
		message = append(message, err[i].Message)
	}
	return message
}

func GetTag(err []*validation.Error) []string {
	tags := make([]string, 0, len(err))
	for i := range err {
		tags = append(tags, err[i].Name)
	}
	return tags
}

func ValidationErrorResponse(c beego.Controller, err []*validation.Error) []string {
	errs := make([]string, 0, len(err))
	Tags := GetTag(err)
	for i := range Tags {
		var errResponse string
		switch Tags[i] {
		case "Min", "Max", "Range", "MinSize", "MaxSize", "Length", "Match", "NotMatch":
			errResponse = fmt.Sprintf("%s : "+helpers.TranslateMessage(c.Ctx, "validation", Tags[i]), err[i].Field, err[i].LimitValue)

		case "Required", "Alpha", "Numeric", "AlphaNumeric", "Email", "IP", "AlphaDash":
			errResponse = fmt.Sprintf("%s : "+helpers.TranslateMessage(c.Ctx, "validation", Tags[i]), err[i].Field)

		default:
			fields := err[i].Key
			keys := strings.Split(fields, ".")
			errResponse = fmt.Sprintf("%s : "+helpers.TranslateMessage(c.Ctx, "validation", keys[1]), keys[0])
		}
		errs = append(errs, errResponse)
	}

	log.Print(errs)
	return errs
}

func indianMobile(v *validation.Validation, obj interface{}, key string) {
	value, ok := obj.(string)
	if !ok {
		return
	}
	pattern := `^[6-9][0-9]{9}$`
	regex := regexp.MustCompile(pattern)
	if !regex.MatchString(value) {
		v.SetError(key, "Please enter a valid Indian mobile number")
	}
}

func RequiredTag(v *validation.Validation, obj interface{}, key string) {
	tags := []string{"logo", "text", "html", "banner"}

	value, ok := obj.(string)
	if !ok {
		return
	}

	ok = helpers.CheckIfExists(value, tags)
	if !ok {
		v.SetError(key, "please enter with in [logo, text, banner, html]")
	}
}


func ValidImageType(file string) bool {
	extensions := []string{"jpeg", "jpg", "png", "svg"}
	ext := helpers.GetFileExtension(file)
	return helpers.CheckIfExists(ext, extensions)
}

func ValidFileType(file string) bool {
	extensions := []string{"csv", "xlsx"}
	ext := helpers.GetFileExtension(file)
	return helpers.CheckIfExists(ext, extensions)
}
