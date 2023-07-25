package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

// CreateCategory 创建分类表单
func CreateCategory(data interface{}, c *gin.Context) map[string][]string {

	rules := govalidator.MapData{
		"name": []string{"required"},
	}
	messages := govalidator.MapData{
		"name": []string{
			"required:分类名称为必填项",
		},
	}

	errs := validate(data, rules, messages)
	return errs
}
