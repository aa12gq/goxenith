package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"goxenith/app/requests/validators"
	"goxenith/pkg/auth"
	pb "goxenith/proto/app/v1"
)

func UserUpdateProfile(data interface{}, c *gin.Context) map[string][]string {

	// 查询用户名重复时，过滤掉当前用户 ID
	uid := auth.CurrentUID(c)
	rules := govalidator.MapData{
		"name":             []string{"required", "alpha_num", "between:3,20" + uid},
		"city":             []string{"min_cn:2", "max_cn:20"},
		"personal_profile": []string{"max_cn:240"},
	}

	messages := govalidator.MapData{
		"name": []string{
			"required:用户名为必填项",
			"alpha_num:用户名格式错误，只允许数字和英文",
			"between:用户名长度需在 3~20 之间",
		},
		"city": []string{
			"min_cn:城市需至少 2 个字",
			"max_cn:城市不能超过 20 个字",
		},
		"personal_profile": []string{
			"max_cn:个人简介不能超过 240 个字",
		},
	}
	return validate(data, rules, messages)
}

func UserUpdatePassword(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"password":           []string{"required", "min:6"},
		"newPassword":        []string{"required", "min:6"},
		"newPasswordConfirm": []string{"required", "min:6"},
	}
	messages := govalidator.MapData{
		"password": []string{
			"required:密码为必填项",
			"min:密码长度需大于 6",
		},
		"newPassword": []string{
			"required:密码为必填项",
			"min:密码长度需大于 6",
		},
		"newPasswordConfirm": []string{
			"required:确认密码框为必填项",
			"min:确认密码长度需大于 6",
		},
	}

	// 确保 comfirm 密码正确
	errs := validate(data, rules, messages)
	_data := data.(*pb.UpdateUserPasswordRequest)
	errs = validators.ValidatePasswordConfirm(_data.NewPassword, _data.NewPasswordConfirm, errs)

	return errs
}
