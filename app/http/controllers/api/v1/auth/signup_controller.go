package auth

import (
	"github.com/gin-gonic/gin"
	v1 "goxenith/app/http/controllers/api/v1"
	entUser "goxenith/app/models/ent/user"
	"goxenith/app/requests"
	"goxenith/pkg/database"
	"goxenith/pkg/model"
	"goxenith/pkg/response"
	pb "goxenith/proto/app/v1"
)

type SignupController struct {
	v1.BaseAPIController
}

// IsPhoneExist 检测手机号是否被注册
func (sc *SignupController) IsPhoneExist(c *gin.Context) {
	request := pb.SignupPhoneExistRequest{}
	if ok := requests.Validate(c, &request, requests.SignupPhoneExist); !ok {
		return
	}

	exist, err := database.DB.User.Query().Where(entUser.PhoneEQ(request.Phone), entUser.DeleteEQ(model.DeletedNo)).Exist(c)
	if err != nil {
		panic(err)
	}

	response.JSON(c, &pb.IsExist{Exist: exist})
}

// IsEmailExist 检测邮箱是否已注册
func (sc *SignupController) IsEmailExist(c *gin.Context) {

	request := pb.SignupEmailExistRequest{}
	if ok := requests.Validate(c, &request, requests.SignupEmailExist); !ok {
		return
	}

	exist, err := database.DB.User.Query().Where(entUser.EmailEQ(request.Email), entUser.DeleteEQ(model.DeletedNo)).Exist(c)
	if err != nil {
		panic(err)
	}

	response.JSON(c, pb.IsExist{Exist: exist})
}
