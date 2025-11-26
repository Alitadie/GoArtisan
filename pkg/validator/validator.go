package validator

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

// 全局变量存放翻译器
var trans ut.Translator

// Init 初始化验证器翻译 (在 main 或 bootstrap 中调用)
func Init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 1. 注册 Tag Name 函数
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		// 2. 注册翻译器
		zhT := zh.New()
		uni := ut.New(zhT, zhT)
		trans, _ = uni.GetTranslator("zh")

		// 3. 注册中文翻译
		_ = zh_translations.RegisterDefaultTranslations(v, trans)
	}
}

// TranslateError 将校验错误转换为 Map
func Translate(err error) map[string]string {
	result := make(map[string]string)

	errors, ok := err.(validator.ValidationErrors)
	if !ok {
		result["error"] = err.Error()
		return result
	}

	for _, e := range errors {
		// e.Field() 已经是我们处理过 json tag 后的字段名了
		result[e.Field()] = e.Translate(trans)
	}
	return result
}
