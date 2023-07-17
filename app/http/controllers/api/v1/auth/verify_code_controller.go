package auth

import (
	"github.com/gin-gonic/gin"
	v1 "goxenith/app/http/controllers/api/v1"
	"goxenith/pkg/captcha"
	"goxenith/pkg/logger"
	"goxenith/pkg/response"
	pb "goxenith/proto/app/v1"
)

type VerifyCodeController struct {
	v1.BaseAPIController
}

// ShowCaptcha 显示图片验证码
func (vc *VerifyCodeController) ShowCaptcha(c *gin.Context) {
	id, b64s, err := captcha.NewCaptcha().GenerateCaptcha()
	logger.LogIf(err)
	response.JSON(c, pb.ShowCaptchaReply{CaptchaId: id, CaptchaImage: b64s})
}
