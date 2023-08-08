package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

func ArticleSave(data interface{}, c *gin.Context) map[string][]string {

	rules := govalidator.MapData{
		"title":   []string{"required", "min_cn:3", "max_cn:40"},
		"summary": []string{"required"},
		"content": []string{"required", "min_cn:3", "max_cn:50000"},
	}
	messages := govalidator.MapData{
		"title": []string{
			"required:博文标题为必填项",
			"min_cn:标题长度需大于 3",
			"max_cn:标题长度需小于 40",
		},
		"summary": []string{
			"required:博文摘要为必填项",
		},
		"content": []string{
			"required:博文内容为必填项",
			"min_cn:长度需大于 10",
		},
	}
	return validate(data, rules, messages)
}
