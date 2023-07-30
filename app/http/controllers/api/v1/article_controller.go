package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goxenith/app/models/ent/article"
	entCommunity "goxenith/app/models/ent/community"
	"goxenith/app/requests"
	"goxenith/dao"
	"goxenith/pkg/auth"
	"goxenith/pkg/logger"
	"goxenith/pkg/model"
	"goxenith/pkg/response"
	pb "goxenith/proto/app/v1"
)

type ArticleController struct {
	BaseAPIController
}

func (a *ArticleController) CreateArticle(ctx *gin.Context) {
	request := &pb.CreateArticleRequest{}

	if ok := requests.Validate(ctx, request, requests.ArticleSave); !ok {
		return
	}
	currentUser := auth.CurrentUser(ctx)
	if currentUser.ID == 0 {
		return
	}
	exist, err := dao.DB.Community.Query().Where(entCommunity.IDEQ(request.CommunityId), entCommunity.DeleteEQ(model.DeletedNo)).Exist(ctx)
	if err != nil {
		logger.LogWarnIf("博文保存出错: %v", err)
		response.Abort400(ctx, "博文保存出错")
		return
	}

	if !exist {
		logger.LogWarnIf(fmt.Sprintf("未找到ID为 %v 的社区", request.CommunityId), err)
		response.Abort400(ctx, "博文保存出错")
	}

	err = dao.DB.Article.Create().
		SetAuthorID(currentUser.ID).
		SetCommunityID(request.CommunityId).
		SetTitle(request.Title).
		SetSummary(request.Summary).
		SetContent(request.Content).
		SetStatus(article.Status(request.Status.String())).Exec(ctx)
	if err != nil {
		response.Abort500(ctx, "保存失败")
		return
	}

	response.Success(ctx)
}
