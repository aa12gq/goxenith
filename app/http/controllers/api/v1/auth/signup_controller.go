package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "goxenith/app/http/controllers/api/v1"
	entUser "goxenith/app/models/ent/user"
	"goxenith/pkg/database"
	"goxenith/pkg/model"
	"net/http"
)

// SignupController 注册控制器
type SignupController struct {
	v1.BaseAPIController
}

// IsPhoneExist 检测手机号是否被注册
func (sc *SignupController) IsPhoneExist(c *gin.Context) {

	type PhoneExistRequest struct {
		Phone string `json:"phone"`
	}
	request := PhoneExistRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		fmt.Println(err.Error())
		return
	}

	exist, err := database.DB.User.Query().Where(entUser.PhoneEQ(request.Phone), entUser.DeleteEQ(model.DeletedNo)).Exist(c)
	if err != nil {
		panic(err)
	}
	//  检查数据库并返回响应
	c.JSON(http.StatusOK, gin.H{
		"exist": exist,
	})
}
