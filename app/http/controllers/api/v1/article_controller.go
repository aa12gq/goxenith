package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	entArtic "goxenith/app/models/ent/article"
	entCommunity "goxenith/app/models/ent/community"
	"goxenith/app/requests"
	"goxenith/dao"
	"goxenith/pkg/auth"
	"goxenith/pkg/logger"
	"goxenith/pkg/model"
	"goxenith/pkg/paginator"
	"goxenith/pkg/response"
	pb "goxenith/proto/app/v1"
	"strconv"
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
		SetStatus(entArtic.Status(request.Status.String())).Exec(ctx)
	if err != nil {
		response.Abort500(ctx, "保存失败")
		return
	}

	response.Success(ctx)
}

func (a *ArticleController) ListArticle(ctx *gin.Context) {
	pageParam, _ := ctx.Params.Get("page")
	page, _ := strconv.Atoi(pageParam)
	query := dao.DB.Article.Query().
		Offset(int(paginator.GetPageOffset(uint32(page), 20))).
		Limit(20).Where(entArtic.DeleteEQ(model.DeletedNo)).Where(entArtic.CommunityIDEQ(1)).WithCommunity().WithAuthor()

	total, err := dao.DB.Article.Query().Where(entArtic.DeleteEQ(model.DeletedNo)).Where(entArtic.CommunityIDEQ(1)).Count(ctx)
	if err != nil {
		response.Abort404(ctx, "未找到博文列表数据")
		return
	}

	articles, err := query.All(ctx)
	if err != nil {
		response.Abort404(ctx, "未找到博文列表数据")
		return
	}

	var rv []*pb.Article
	for _, v := range articles {
		rv = append(rv, &pb.Article{
			Id:            v.ID,
			AuthorId:      v.AuthorID,
			AuthorName:    v.Edges.Author.UserName,
			CommunityId:   v.CommunityID,
			CommunityName: v.Edges.Community.Name,
			Title:         v.Title,
			Summary:       v.Summary,
			Content:       v.Content,
			Links:         int32(v.Likes),
			Views:         int32(v.Views),
			Status:        pb.ArticleStatus(pb.ArticleStatus_value[v.Status.String()]),
		})
	}

	reply := &pb.ListArticleReply{
		Data:  rv,
		Total: uint32(total),
		Count: uint32(len(rv)),
		Page:  uint32(page),
	}
	response.JSON(ctx, reply)
}
