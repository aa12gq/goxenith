package auth

import (
	"github.com/gin-gonic/gin"
	v1 "goxenith/app/http/controllers/api/v1"
	"goxenith/pkg/captcha"
	"goxenith/pkg/logger"
	"net/http"
)

type VerifyCodeController struct {
	v1.BaseAPIController
}

// ShowCaptcha 显示图片验证码
func (vc *VerifyCodeController) ShowCaptcha(c *gin.Context) {
	id, b64s, err := captcha.NewCaptcha().GenerateCaptcha()
	logger.LogIf(err)
	c.JSON(http.StatusOK, gin.H{
		"captcha_id":    id,
		"captcha_image": b64s,
	})
}
