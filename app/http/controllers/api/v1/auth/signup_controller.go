package auth

import (
	"github.com/gin-gonic/gin"
	v1 "goxenith/app/http/controllers/api/v1"
	entUser "goxenith/app/models/ent/user"
	"goxenith/app/requests"
	"goxenith/pkg/database"
	"goxenith/pkg/logger"
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

// SignupUsingPhone 使用手机和验证码进行注册
func (sc *SignupController) SignupUsingPhone(c *gin.Context) {
	request := pb.SignupUserUsingPhoneRequest{}
	if ok := requests.Validate(c, &request, requests.SignupUsingPhone); !ok {
		return
	}

	if exist := sc.phoneExists(c, request.Phone); exist {
		return
	}

	if exist := sc.nameExists(c, request.Name); exist {
		return
	}

	if !sc.createUser(c, &request) {
		return
	}

	response.Success(c)
}

func (sc *SignupController) phoneExists(c *gin.Context, phone string) bool {
	exist, err := database.DB.User.Query().Where(entUser.PhoneEQ(phone)).Exist(c)
	if err != nil {
		response.Abort500(c, "查询手机号失败，请稍后尝试~")
		return false
	}

	if exist {
		response.Abort400(c, "该手机号已存在")
	}
	return exist
}

func (sc *SignupController) nameExists(c *gin.Context, name string) bool {
	nameIsExist, err := database.DB.User.Query().Where(entUser.UserNameEQ(name)).Exist(c)
	if err != nil {
		response.Abort500(c, "查询用户名失败，请稍后尝试~")
		return false
	}

	if nameIsExist {
		response.Abort400(c, "用户名已存在")
	}
	return nameIsExist
}

func (sc *SignupController) createUser(c *gin.Context, request *pb.SignupUserUsingPhoneRequest) bool {
	_user, err := database.DB.User.Create().
		SetUserName(request.Name).SetPhone(request.Phone).SetPassword(request.Password).Save(c)
	if err != nil || _user.ID <= 0 {
		logger.LogWarnIf("创建用户失败", err)
		response.Abort500(c, "创建用户失败，请稍后尝试~")
		return false
	}
	return true
}
