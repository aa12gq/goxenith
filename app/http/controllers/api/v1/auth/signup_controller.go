package auth

import (
	"github.com/gin-gonic/gin"
	v1 "goxenith/app/http/controllers/api/v1"
	"goxenith/app/models/ent"
	entUser "goxenith/app/models/ent/user"
	"goxenith/app/requests"
	"goxenith/dao"
	"goxenith/pkg/auth"
	"goxenith/pkg/logger"
	"goxenith/pkg/model"
	"goxenith/pkg/password"
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

	exist, err := dao.DB.User.Query().
		Where(entUser.PhoneEQ(request.Phone),
			entUser.DeleteEQ(model.DeletedNo)).Exist(c)
	if err != nil {
		logger.LogWarnIf("查询出错", err)
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

	exist, err := dao.DB.User.Query().
		Where(entUser.EmailEQ(request.Email),
			entUser.DeleteEQ(model.DeletedNo)).Exist(c)
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

	_user, res := sc.createUserByPhone(c, &request)
	if !res {
		return
	}

	token := auth.NewJWT().IssueToken(_user.ID, _user.UserName)

	response.JSON(c, &pb.SignupUserUsingPhoneReply{
		Data: &pb.SignupUserUsingPhoneReply_Data{
			Id:   _user.ID,
			Name: _user.UserName,
		},
		Token: token,
	})

}

func (sc *SignupController) SignupUsingEmail(c *gin.Context) {
	request := pb.SignupUsingEmailRequest{}
	if ok := requests.Validate(c, &request, requests.SignupUsingEmail); !ok {
		return
	}

	if exist := sc.emailExists(c, request.Email); exist {
		return
	}

	if exist := sc.nameExists(c, request.Name); exist {
		return
	}

	_user, res := sc.createUserByEmail(c, &request)
	if !res {
		return
	}
	token := auth.NewJWT().IssueToken(_user.ID, _user.UserName)
	response.JSON(c, &pb.SignupUsingEmailReply{
		Data: &pb.SignupUsingEmailReply_Data{
			Id:   _user.ID,
			Name: _user.UserName,
		},
		Token: token,
	})
}

func (sc *SignupController) phoneExists(c *gin.Context, phone string) bool {
	exist, err := dao.DB.User.Query().Where(entUser.PhoneEQ(phone)).Exist(c)
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
	nameIsExist, err := dao.DB.User.Query().Where(entUser.UserNameEQ(name)).Exist(c)
	if err != nil {
		response.Abort500(c, "查询用户名失败，请稍后尝试~")
		return false
	}

	if nameIsExist {
		response.Abort400(c, "用户名已存在")
	}
	return nameIsExist
}

func (sc *SignupController) createUserByPhone(c *gin.Context, request *pb.SignupUserUsingPhoneRequest) (*ent.User, bool) {
	_user, err := dao.DB.User.Create().
		SetUserName(request.Name).SetPhone(request.Phone).SetPassword(password.BcryptPassword(request.Password)).Save(c)
	if err != nil || _user.ID <= 0 {
		logger.LogWarnIf("创建用户失败", err)
		response.Abort500(c, "创建用户失败，请稍后尝试~")
		return nil, false
	}
	return _user, true
}

func (sc *SignupController) createUserByEmail(c *gin.Context, request *pb.SignupUsingEmailRequest) (*ent.User, bool) {
	_user, err := dao.DB.User.Create().
		SetUserName(request.Name).SetPhone(request.Email).SetPassword(password.BcryptPassword(request.Password)).Save(c)
	if err != nil || _user.ID <= 0 {
		logger.LogWarnIf("创建用户失败", err)
		response.Abort500(c, "创建用户失败，请稍后尝试~")
		return nil, false
	}
	return _user, true
}

func (sc *SignupController) emailExists(c *gin.Context, phone string) bool {
	exist, err := dao.DB.User.Query().Where(entUser.EmailEQ(phone)).Exist(c)
	if err != nil {
		response.Abort500(c, "查询邮箱失败，请稍后尝试~")
		return false
	}

	if exist {
		response.Abort400(c, "该邮箱已存在")
	}
	return exist
}
