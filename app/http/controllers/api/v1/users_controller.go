package v1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
	"goxenith/app/models/ent"
	entArtic "goxenith/app/models/ent/article"
	"goxenith/app/models/ent/likerecord"
	entUser "goxenith/app/models/ent/user"
	"goxenith/dao"
	"goxenith/pkg/auth"
	"goxenith/pkg/logger"
	"goxenith/pkg/model"
	"goxenith/pkg/password"
	"goxenith/pkg/response"
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

	// 文章数量
	agg, _ := dao.DB.Article.Query().
		Where(entArtic.AuthorIDEQ(uint64(id)), entArtic.DeleteEQ(model.DeletedNo), entArtic.StatusEQ(entArtic.StatusEFFECT)).
		Aggregate(ent.Sum(entArtic.FieldAuthorID)).Int(ctx)

	userInfo := convertUserInfo(info)
	userInfo.ArticleTotal = int32(agg)
	response.JSON(ctx, pb.GetUserInfoReply{UserInfo: userInfo})
}

func (c *UsersController) UpdateUserInfo(ctx *gin.Context) {
	request := pb.UpdateUserInfoRequest{}
	if err := ctx.ShouldBind(&request); err != nil {
		response.BadRequest(ctx, err, "请求解析错误，请确认请求格式是否正确。上传文件请使用 multipart 标头，参数请使用 JSON 格式。")
		return
	}

	exist, err := dao.DB.User.Query().Where(entUser.UserNameEQ(request.UserName), entUser.IDNEQ(request.Id), entUser.DeleteEQ(model.DeletedNo)).Exist(ctx)
	if err != nil {
		logger.LogWarnIf("更新出错", err)
		response.Abort500(ctx, "更新出错")
		return
	}
	if exist {
		response.Abort400(ctx, fmt.Sprintf("用户名%v已被占用", request.UserName))
		return
	}

	user, err := dao.DB.User.UpdateOneID(request.Id).
		SetUserName(request.UserName).
		SetRealName(request.RealName).
		SetPhone(request.Phone).
		SetCity(request.City).
		SetGender(entUser.Gender(request.Gender.String())).
		SetAge(uint8(request.Age)).
		SetPersonalProfile(request.PersonalProfile).
		SetEmail(request.Email).
		Save(ctx)
	if err != nil {
		logger.LogWarnIf("更新出错", err)
		response.Abort500(ctx, "更新出错")
		return
	}

	response.JSON(ctx, &pb.UpdateUserInfoReply{UserInfo: convertUserInfo(user)})
}

func (c *UsersController) UpdateUserAvatar(ctx *gin.Context) {
	request := &pb.UpdateUserAvatarRequest{}
	if err := ctx.ShouldBind(request); err != nil {
		response.BadRequest(ctx, err, "请求解析错误，请确认请求格式是否正确。上传文件请使用 multipart 标头，参数请使用 JSON 格式。")
		return
	}

	exist, err := dao.DB.User.Query().Where(entUser.ID(request.Id), entUser.DeleteEQ(model.DeletedNo)).Exist(ctx)
	if err != nil {
		logger.LogWarnIf("查询用户出错", err)
		response.Abort500(ctx, "查询用户出错")
		return
	}

	if !exist {
		response.Abort404(ctx, "用户不存在")
		return
	}

	err = dao.DB.User.UpdateOneID(request.Id).
		SetAvatar(request.Avatar).
		Exec(ctx)
	if err != nil {
		logger.LogWarnIf("更新出错", err)
		response.Abort500(ctx, "更新出错")
		return
	}

	response.Success(ctx)
}

func (c *UsersController) UpdatePassword(ctx *gin.Context) {
	request := pb.UpdateUserPasswordRequest{}
	if err := ctx.ShouldBind(&request); err != nil {
		response.BadRequest(ctx, err, "请求解析错误，请确认请求格式是否正确。上传文件请使用 multipart 标头，参数请使用 JSON 格式。")
		return
	}

	if request.Id == 0 {
		response.Abort400(ctx, "缺少用户ID")
		return
	}

	user, err := dao.DB.User.Query().Where(entUser.ID(request.Id), entUser.DeleteEQ(model.DeletedNo)).First(ctx)
	if err != nil {
		response.Abort404(ctx, fmt.Sprintf("未找到ID为 %v 的用户信息", request.Id))
		return
	}

	if user.Password == "" {
		response.Abort400(ctx, "用户密码为空")
		return
	}

	if request.Password == "" {
		response.Abort400(ctx, "旧密码不能为空")
		return
	}

	if !password.BcryptPasswordMatch(request.Password, user.Password) {
		response.Abort400(ctx, "旧密码不正确")
		return
	}

	if request.NewPassword == "" {
		response.Abort400(ctx, "新密码不能为空")
		return
	}

	err = dao.DB.User.UpdateOne(user).
		SetPassword(password.BcryptPassword(request.NewPassword)).
		Exec(ctx)
	if err != nil {
		logger.LogWarnIf("更新密码出错", err)
		response.Abort500(ctx, "更新密码出错")
		return
	}

	response.Success(ctx)
}

func (a *UsersController) ListArticlesForUser(ctx *gin.Context) {
	idStr, ok := ctx.Params.Get("id")
	if !ok {
		response.BadRequest(ctx, errors.New("缺少用户id"))
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(ctx, errors.New("无效的用户id"))
		return
	}

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		response.BadRequest(ctx, errors.New("无效的页码"))
		return
	}
	pageSize, err := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))
	if err != nil {
		response.BadRequest(ctx, errors.New("无效的页面大小"))
		return
	}
	offset := (page - 1) * pageSize

	query := dao.DB.Article.Query().
		Offset(offset).
		Limit(pageSize).
		Where(entArtic.DeleteEQ(model.DeletedNo), entArtic.AuthorIDEQ(uint64(id))).
		Order(ent.Desc(entArtic.FieldCreatedAt))

	articles, err := query.All(ctx)
	if err != nil {
		response.Abort404(ctx, "未找到博文列表数据")
		return
	}

	total, err := query.Count(ctx)
	if err != nil {
		response.Abort404(ctx, "未找到博文列表数据总数")
		return
	}

	rv := make([]*pb.ListArticlesForUserReply_Article, len(articles))
	for i, v := range articles {
		likeCount, _ := dao.DB.LikeRecord.Query().Where(likerecord.ArticleIDEQ(v.ID)).Count(ctx)
		rv[i] = &pb.ListArticlesForUserReply_Article{
			Id:          v.ID,
			Title:       v.Title,
			Summary:     v.Summary,
			Likes:       int32(likeCount),
			Views:       int32(v.ViewCount),
			CreatedDate: timestamppb.New(v.CreatedAt),
			UpdatedDate: timestamppb.New(v.UpdatedAt),
		}
	}

	reply := &pb.ListArticlesForUserReply{
		Data:  rv,
		Total: uint32(total),
		Count: uint32(len(rv)),
		Page:  uint32(page),
	}
	response.JSON(ctx, reply)
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
		CreatedDate:     timestamppb.New(user.CreatedAt),
	}
}
