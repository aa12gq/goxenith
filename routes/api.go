package routes

import (
	"github.com/gin-gonic/gin"
	"goxenith/app/http/controllers/api/v1/auth"
)

func RegisterAPIRoutes(r *gin.Engine) {

	v1 := r.Group("/v1")
	{
		authGroup := v1.Group("/auth")
		{
			suc := new(auth.SignupController)
			// 判断手机是否已注册
			authGroup.POST("/signup/phone/exist", suc.IsPhoneExist)
			// 判断邮箱是否已注册
			authGroup.POST("/signup/email/exist", suc.IsEmailExist)
			// 用户注册
			authGroup.POST("/signup/using-phone", suc.SignupUsingPhone)
			// 发送验证码
			vcc := new(auth.VerifyCodeController)
			// 图片验证码，需要加限流
			authGroup.POST("/verify-codes/captcha", vcc.ShowCaptcha)
			authGroup.POST("/verify-codes/phone", vcc.SendUsingPhone)
			lgc := new(auth.LoginController)
			// 使用手机号，短信验证码进行登录
			authGroup.POST("/login/using-phone", lgc.LoginByPhone)
			// 支持手机号，Email 和 用户名
			authGroup.POST("/login/using-password", lgc.LoginByPassword)
			authGroup.POST("/login/refresh-token", lgc.RefreshToken)
			
		}
	}
}
