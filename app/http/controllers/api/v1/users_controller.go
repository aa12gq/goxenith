package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
	"goxenith/app/models/ent"
	entUser "goxenith/app/models/ent/user"
	"goxenith/app/requests"
	"goxenith/dao"
	"goxenith/pkg/auth"
	"goxenith/pkg/logger"
	"goxenith/pkg/model"
	"goxenith/pkg/response"
	"goxenith/pkg/xcopy"
	pb "goxenith/proto/app/v1"
	"strconv"
)

type UsersController struct {
	BaseAPIController
}

// CurrentUser 当前登录用户信息
func (ctrl *UsersController) CurrentUser(c *gin.Context) {
	userModel := auth.CurrentUser(c)
	response.Data(c, userModel)
}

func (c *UsersController) GetUserInfo(ctx *gin.Context) {
	idStr, _ := ctx.Params.Get("id")
	id, _ := strconv.Atoi(idStr)

	info, err := dao.DB.User.Query().Where(entUser.IDEQ(uint64(id)), entUser.DeleteEQ(model.DeletedNo)).First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			response.Abort404(ctx, fmt.Sprintf("未找到ID为 %v 的用户信息", id))
			return
		}
	}
	response.JSON(ctx, convertUserInfo(info))
}

func (c *UsersController) UpdateUserInfo(ctx *gin.Context) {
	request := pb.UpdateUserInfoRequest{}.UserInfo
	if ok := requests.Validate(ctx, &request, requests.ArticleSave); !ok {
		return
	}
	exist, err := dao.DB.User.Query().Where(entUser.UserNameEQ(request.UserName), entUser.DeleteEQ(model.DeletedNo)).Exist(ctx)
	if err != nil {
		logger.LogWarnIf("更新出错", err)
		response.Abort500(ctx, "更新出错")
		return
	}
	if exist {
		response.Abort400(ctx, fmt.Sprintf("用户名%v已被占用", request.UserName))
		return
	}

	input := &ent.User{}
	err = xcopy.Copy(request, input)
	if err != nil {
		logger.LogWarnIf("更新出错", err)
		response.Abort500(ctx, "更新出错")
		return
	}

	user, err := dao.DB.User.UpdateOneID(request.Id).SetUser(input).Save(ctx)
	if err != nil {
		logger.LogWarnIf("更新出错", err)
		response.Abort500(ctx, "更新出错")
		return
	}

	response.JSON(ctx, &pb.UpdateUserInfoReply{UserInfo: convertUserInfo(user)})
}

func convertUserInfo(user *ent.User) *pb.UserInfo {
	return &pb.UserInfo{
		Id:              user.ID,
		UserName:        user.UserName,
		RealName:        user.RealName,
		Phone:           user.Phone,
		City:            user.City,
		Age:             int32(user.Age),
		Birthday:        timestamppb.New(user.Birthday),
		PersonalProfile: user.PersonalProfile,
		Email:           user.Email,
		Avatar:          user.Avatar,
		Gender:          pb.Gender(pb.Gender_value[user.Gender.String()]),
	}
}