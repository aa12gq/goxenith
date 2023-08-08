package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
	"goxenith/app/models/ent"
	entArtic "goxenith/app/models/ent/article"
	"goxenith/app/requests"
	"goxenith/dao"
	"goxenith/pkg/auth"
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

	err := dao.DB.Article.Create().
		SetAuthorID(currentUser.ID).
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
	pageParam := ctx.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageParam)
	pageSizeParam := ctx.DefaultQuery("pageSize", "20")
	pageSize, _ := strconv.Atoi(pageSizeParam)
	offset := int(paginator.GetPageOffset(uint32(page), uint32(pageSize)))

	tx, err := dao.DB.BeginTx(ctx, nil)
	if err != nil {
		response.Abort404(ctx, "未找到博文列表数据")
		return
	}

	query := tx.Article.Query().
		Offset(offset).
		Limit(pageSize).
		Where(entArtic.DeleteEQ(model.DeletedNo)).
		WithAuthor()

	articles, err := query.All(ctx)
	if err != nil {
		response.Abort404(ctx, "未找到博文列表数据")
		return
	}

	total, err := query.Count(ctx)
	if err != nil {
		response.Abort404(ctx, "未找到博文列表数据")
		return
	}

	rv := make([]*pb.Article, 0, len(articles))
	for _, v := range articles {
		rv = append(rv, convertArticle(v))
	}

	reply := &pb.ListArticleReply{
		Data:  rv,
		Total: uint32(total),
		Count: uint32(len(rv)),
		Page:  uint32(page),
	}
	response.JSON(ctx, reply)
}

func (a *ArticleController) GetArticle(ctx *gin.Context) {
	idStr, _ := ctx.Params.Get("id")
	id, _ := strconv.Atoi(idStr)

	article, err := dao.DB.Article.Query().Where(entArtic.IDEQ(uint64(id)), entArtic.DeleteEQ(model.DeletedNo)).WithAuthor().First(ctx)
	if err != nil {
		response.Abort404(ctx, fmt.Sprintf("未找到ID为 %v 的博文数据", id))
		return
	}

	// 作者的文章数量
	agg, _ := dao.DB.Article.Query().
		Where(entArtic.AuthorIDEQ(article.AuthorID), entArtic.DeleteEQ(model.DeletedNo)).
		Aggregate(ent.Sum(entArtic.FieldAuthorID)).Int(ctx)
	reply := &pb.GetArticleReply{Article: convertArticle(article)}

	reply.Article.Author.ArticleTotal = int32(agg)

	response.JSON(ctx, reply)
}

func (a *ArticleController) UpdateArticle(ctx *gin.Context) {
	request := pb.UpdateArticleRequest{}.Article

	if ok := requests.Validate(ctx, request, requests.ArticleSave); !ok {
		return
	}
	currentUser := auth.CurrentUser(ctx)
	if currentUser.ID == 0 {
		return
	}

	article, err := dao.DB.Article.Query().Where(entArtic.IDEQ(request.Id), entArtic.DeleteEQ(model.DeletedNo)).First(ctx)
	if err != nil {
		response.Abort500(ctx, "保存失败")
		return
	}

	if article.AuthorID != currentUser.ID {
		response.Abort403(ctx)
		return
	}

	err = article.Update().
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

func (a *ArticleController) DeleteArticle(ctx *gin.Context) {
	idStr, _ := ctx.Params.Get("id")
	id, _ := strconv.Atoi(idStr)

	currentUser := auth.CurrentUser(ctx)
	if currentUser.ID == 0 {
		return
	}

	article, err := dao.DB.Article.Query().Where(entArtic.IDEQ(uint64(id)), entArtic.DeleteEQ(model.DeletedNo)).First(ctx)
	if err != nil {
		response.Abort500(ctx, "删除失败")
		return
	}

	if article.AuthorID != currentUser.ID {
		response.Abort403(ctx)
		return
	}
	err = article.Update().SetDelete(model.DeletedYes).Exec(ctx)
	if err != nil {
		response.Abort500(ctx, "删除失败")
		return
	}

	response.Success(ctx)
}

func convertArticle(article *ent.Article) *pb.Article {
	return &pb.Article{
		Id: article.ID,
		Author: &pb.Article_Author{
			Id:     article.AuthorID,
			Name:   article.Edges.Author.UserName,
			Avatar: article.Edges.Author.Avatar,
		},
		Title:       article.Title,
		Summary:     article.Summary,
		Content:     article.Content,
		Links:       int32(article.Likes),
		Views:       int32(article.Views),
		Status:      pb.ArticleStatus(pb.ArticleStatus_value[article.Status.String()]),
		CreatedDate: timestamppb.New(article.CreatedAt),
		UpdatedDate: timestamppb.New(article.UpdatedAt),
	}
}
