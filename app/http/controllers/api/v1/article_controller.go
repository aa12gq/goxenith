package v1

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
	"goxenith/app/models/ent"
	entArtic "goxenith/app/models/ent/article"
	"goxenith/app/models/ent/likerecord"
	"goxenith/app/requests"
	"goxenith/dao"
	"goxenith/pkg/auth"
	"goxenith/pkg/cache"
	"goxenith/pkg/helpers"
	"goxenith/pkg/logger"
	"goxenith/pkg/model"
	"goxenith/pkg/paginator"
	"goxenith/pkg/response"
	pb "goxenith/proto/app/v1"
	"strconv"
	"strings"
	"time"
)

const maxRetries = 3
const retryDelay = 5 * time.Second

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
		response.Abort403(ctx)
		return
	}

	summary := getSummary(request.Content, 100)
	err := dao.DB.Article.Create().
		SetAuthorID(currentUser.ID).
		SetTitle(request.Title).
		SetSummary(summary).
		SetContent(request.Content).
		SetStatus(entArtic.Status(request.Status.String())).
		SetLastUpdatedAt(time.Now()).Exec(ctx)
	if err != nil {
		logger.LogWarnIf("保存失败", err)
		response.Abort500(ctx, "保存失败")
		return
	}

	response.Success(ctx)
}

// 确保不会在中间的多字节字符上截断，并移除 Markdown 语法
func getSummary(content string, length int) string {
	strippedContent := helpers.RemoveMarkdown(content)

	if len(strippedContent) <= length {
		return strippedContent
	}
	utf8Content := []rune(strippedContent)
	if len(utf8Content) <= length {
		return strippedContent
	}
	return string(utf8Content[:length])
}

func (a *ArticleController) ListArticle(ctx *gin.Context) {
	pageParam := ctx.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageParam)
	pageSizeParam := ctx.DefaultQuery("pageSize", "10")
	pageSize, _ := strconv.Atoi(pageSizeParam)
	offset := int(paginator.GetPageOffset(uint32(page), uint32(pageSize)))
	sortType := ctx.DefaultQuery("sortType", "latest")
	title := ctx.DefaultQuery("title", "")

	tx, err := dao.DB.BeginTx(ctx, nil)
	if err != nil {
		logger.LogWarnIf("开启事务出错", err)
		response.Abort404(ctx, "未找到博文列表数据")
		return
	}

	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
		} else {
			err := tx.Commit()
			if err != nil {
				return
			}
		}
	}()

	countQuery := tx.Article.Query().Where(entArtic.DeleteEQ(model.DeletedNo))
	total, err := countQuery.Count(ctx)
	if err != nil {
		logger.LogWarnIf("查询博文列表数量出错", err)
		response.Abort404(ctx, "未找到博文列表数据")
		return
	}

	articlesQuery := tx.Article.Query().
		Offset(offset).
		Limit(pageSize).
		Where(entArtic.DeleteEQ(model.DeletedNo)).Where(func(selector *sql.Selector) {
		selector.Where(sql.Like(selector.C(entArtic.FieldTitle), fmt.Sprintf("%%%v%%", title)))
	}).
		WithAuthor()

	if sortType == "latest" {
		articlesQuery = articlesQuery.Order(ent.Desc(entArtic.FieldCreatedAt))
	} else {
		//articlesQuery = articlesQuery.Order(ent.Desc(entArtic.FieldLikes))
	}

	articles, err := articlesQuery.All(ctx)
	if err != nil {
		logger.LogWarnIf("查询博文列表出错", err)
		response.Abort404(ctx, "未找到博文列表数据")
		return
	}

	rv := make([]*pb.Article, 0, len(articles))
	for _, v := range articles {
		likeCount, _ := dao.DB.LikeRecord.Query().Where(likerecord.ArticleIDEQ(v.ID)).Count(ctx)
		rv = append(rv, convertArticle(v, likeCount))
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
	likeCount, _ := dao.DB.LikeRecord.Query().Where(likerecord.ArticleIDEQ(article.ID)).Count(ctx)

	reply := &pb.GetArticleReply{Article: convertArticle(article, likeCount)}

	reply.Article.Author.ArticleTotal = int32(agg)

	response.JSON(ctx, reply)
}

func (a *ArticleController) UpdateArticle(ctx *gin.Context) {
	request := &pb.UpdateArticleRequest{}

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

	summary := getSummary(request.Content, 100)
	err = article.Update().
		SetTitle(request.Title).
		SetSummary(summary).
		SetContent(request.Content).
		SetStatus(entArtic.Status(request.Status.String())).
		SetLastUpdatedAt(time.Now()).Exec(ctx)
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

func (a *ArticleController) LikeArticle(ctx *gin.Context) {
	request := &pb.LikeArticleRequest{}
	if err := ctx.ShouldBind(request); err != nil {
		response.BadRequest(ctx, err, "请求解析错误，请确认请求格式是否正确。上传文件请使用 multipart 标头，参数请使用 JSON 格式。")
		return
	}

	currentUser := auth.CurrentUser(ctx)
	if currentUser.ID == 0 {
		return
	}

	likeKey := fmt.Sprintf("article:%d:like:user:%d", request.Id, currentUser)
	ttl := time.Hour * 24

	if cache.Has(likeKey) {
		response.JSON(ctx, gin.H{
			"message": "您已经点过赞了",
		})
		return
	}

	cache.Set(likeKey, 1, ttl)

	articleExist, err := dao.DB.Article.Query().Where(entArtic.IDEQ(request.Id)).Exist(ctx)
	if err != nil {
		logger.LogWarnIf("查询文章出错", err)
		response.Abort500(ctx, "文章点赞失败")
		return
	}

	if !articleExist {
		response.Abort404(ctx, "文章不存在")
		return
	}

	// 查询是否有点赞记录
	record, err := dao.DB.LikeRecord.Query().
		Where(likerecord.ArticleID(uint64(request.Id)), likerecord.UserID(currentUser.ID)).
		Only(ctx)

	if err != nil && !ent.IsNotFound(err) {
		logger.LogWarnIf("查询点赞记录出错", err)
		response.Abort500(ctx, "文章点赞失败")
		return
	}

	if record != nil {
		if record.IsActive {
			// 已经点过赞，将其取消
			_, err = dao.DB.LikeRecord.UpdateOneID(record.ID).SetIsActive(false).Save(ctx)
			if err == nil {
				cache.Forget(likeKey)
				response.JSON(ctx, gin.H{
					"message": "取消点赞成功",
				})
				return
			}
		} else {
			// 之前取消过点赞，现在重新点赞
			_, err = dao.DB.LikeRecord.UpdateOneID(record.ID).SetIsActive(true).Save(ctx)
		}
	} else {
		// 之前从未点赞，创建新记录
		_, err = dao.DB.LikeRecord.Create().
			SetArticleID(uint64(request.Id)).
			SetUserID(currentUser.ID).
			SetIsActive(true).
			Save(ctx)
	}

	if err != nil {
		logger.LogWarnIf("保存点赞记录失败", err)
		response.Abort500(ctx, "文章点赞失败")
		return
	}

	response.JSON(ctx, gin.H{
		"message": "点赞成功",
	})
}

func (a *ArticleController) CheckLikeStatus(ctx *gin.Context) {
	currentUser := auth.CurrentUser(ctx)
	if currentUser.ID == 0 {
		response.Abort403(ctx)
		return
	}

	articleIDStr, _ := ctx.Params.Get("id")
	articleID, _ := strconv.Atoi(articleIDStr)
	if articleID == 0 {
		response.Abort400(ctx, "无效的文章ID")
		return
	}

	// 检查点赞记录
	likeRecord, err := dao.DB.LikeRecord.Query().
		Where(likerecord.ArticleID(uint64(articleID)), likerecord.UserID(currentUser.ID)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			response.JSON(ctx, gin.H{
				"liked": false,
			})
			return
		}
		logger.LogWarnIf("查询点赞记录出错", err)
		response.Abort500(ctx, "检查点赞状态失败")
		return
	}

	response.JSON(ctx, gin.H{
		"liked": likeRecord.IsActive,
	})
}

// SyncArticleViewsFromRedis 定时从Redis更新数据库
func SyncArticleViewsFromRedis() {
	ticker := time.NewTicker(30 * time.Minute)

	for _ = range ticker.C {
		for i := 0; i < maxRetries; i++ {
			articleIDs, views, _ := getAllArticleViewsFromRedis()
			if len(articleIDs) == 0 && len(views) == 0 {
				break
			}

			err := updateArticleViewsToDB(articleIDs, views)
			if err == nil {
				break
			}

			if i < maxRetries-1 {
				logger.LogWarnIf("同步失败，准备重试", err)
				time.Sleep(retryDelay)
			} else {
				logger.LogWarnIf("达到最大重试次数，同步失败", err)
			}
		}
	}
}

func getAllArticleViewsFromRedis() ([]int, []int, error) {
	var articleIDs []int
	var views []int

	keys := cache.GetStringSlice("article:*:view")
	for _, key := range keys {
		idStr := strings.Split(key, ":")[1]
		articleID, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, nil, err
		}

		viewCount := cache.GetInt(key)
		articleIDs = append(articleIDs, articleID)
		views = append(views, viewCount)
	}

	return articleIDs, views, nil
}

func updateArticleViewsToDB(articleIDs []int, views []int) error {
	tx, err := dao.DB.Tx(context.Background())
	if err != nil {
		logger.LogWarnIf("开始事务失败", err)
		return err
	}

	for i, id := range articleIDs {
		err := tx.Article.UpdateOneID(uint64(id)).AddViewCount(int64(views[i])).Exec(context.Background())
		if err != nil {
			logger.LogWarnIf(fmt.Sprintf("更新文章ID %d 的浏览量失败", id), err)
			_ = tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		logger.LogWarnIf("提交事务失败", err)
		return err
	}

	resetRedisCounters(articleIDs)
	return nil
}

func resetRedisCounters(articleIDs []int) {
	// TODO: 重置Redis中的文章浏览量计数器
	for _, id := range articleIDs {
		key := fmt.Sprintf("article:%d:view", id)
		cache.Forget(key)
	}
}

func (a *ArticleController) ViewArticle(ctx *gin.Context) {
	request := &pb.UpdateArticleViewsRequest{}
	if err := ctx.ShouldBind(request); err != nil {
		response.BadRequest(ctx, err, "请求解析错误，请确认请求格式是否正确。上传文件请使用 multipart 标头，参数请使用 JSON 格式。")
		return
	}

	article, err := dao.DB.Article.Query().Where(entArtic.IDEQ(uint64(request.Id)),
		entArtic.DeleteEQ(model.DeletedNo)).First(ctx)
	if err != nil {
		response.Abort404(ctx, fmt.Sprintf("未找到ID为 %v 的博文数据", request.Id))
		return
	}

	clientIP := ctx.ClientIP()
	cookieName := "user_track"
	cookieValue, err := ctx.Cookie(cookieName)
	if err != nil || cookieValue == "" {
		// 如果cookie不存在，为用户创建一个新的
		cookieValue = fmt.Sprintf("%v-%s", time.Now().UnixNano(), clientIP)
		ctx.SetCookie(cookieName, cookieValue, 3600*24*365, "/", "", false, true) // 设置一年有效期
	}

	viewKey := fmt.Sprintf("article:%d:view:ip:%s:cookie:%s", request.Id, clientIP, cookieValue)
	viewTTL := time.Minute * 30
	// 如果此cookie/IP组合在30分钟内没有访问过该文章，则增加浏览量
	if !cache.Has(viewKey) {
		err = article.Update().AddViewCount(1).Exec(ctx)
		if err != nil {
			logger.LogWarnIf("增加浏览量失败", err)
			response.Abort500(ctx, "增加浏览量失败")
			return
		}
		cache.Set(viewKey, 1, viewTTL) // 记录该IP和cookie已经访问了此文章
	}

	article, _ = dao.DB.Article.Query().Where(entArtic.ID(article.ID)).Select(entArtic.FieldViewCount).Only(ctx)

	response.JSON(ctx, pb.UpdateArticleViewsReply{Views: int32(article.ViewCount)})
}

func (a *ArticleController) ToggleCollectArticle(ctx *gin.Context) {
	currentUser := auth.CurrentUser(ctx)
	if currentUser.ID == 0 {
		response.Abort403(ctx)
		return
	}

	articleIDStr, _ := ctx.Params.Get("id")
	articleID, _ := strconv.Atoi(articleIDStr)
	if articleID == 0 {
		response.Abort400(ctx, "无效的文章ID")
		return
	}

	record, err := dao.DB.LikeRecord.Query().
		Where(likerecord.ArticleID(uint64(articleID)), likerecord.UserID(currentUser.ID)).
		Only(ctx)

	var isNowCollected bool
	if err == nil && record != nil {
		// 反转收藏状态
		_, err = dao.DB.LikeRecord.UpdateOneID(record.ID).SetIsCollect(!record.IsCollect).Save(ctx)
		isNowCollected = !record.IsCollect // 使用反转后的值
	} else {
		_, err = dao.DB.LikeRecord.Create().
			SetArticleID(uint64(articleID)).
			SetUserID(currentUser.ID).
			SetIsCollect(true).
			Save(ctx)
		isNowCollected = true
	}

	if err != nil {
		logger.LogWarnIf("保存收藏记录失败", err)
		response.Abort500(ctx, "操作失败")
		return
	}

	// 根据当前的收藏状态返回相应的消息
	if isNowCollected {
		response.JSON(ctx, gin.H{
			"message": "收藏成功",
		})
	} else {
		response.JSON(ctx, gin.H{
			"message": "取消收藏成功",
		})
	}
}

func (a *ArticleController) GetCollectedArticles(ctx *gin.Context) {
	currentUser := auth.CurrentUser(ctx)
	if currentUser.ID == 0 {
		response.Abort403(ctx)
		return
	}

	articles, err := dao.DB.LikeRecord.Query().
		Where(likerecord.UserID(currentUser.ID), likerecord.IsCollect(true)).WithArticle(func(query *ent.ArticleQuery) {
		query.WithAuthor()
	}).
		All(ctx)

	if err != nil {
		logger.LogWarnIf("获取收藏文章失败", err)
		response.Abort500(ctx, "获取收藏文章失败")
		return
	}

	responseArticles := make([]*pb.Article, 0)
	for _, article := range articles {
		likeCount, _ := dao.DB.LikeRecord.Query().Where(likerecord.ArticleIDEQ(article.ID)).Count(ctx)
		responseArticles = append(responseArticles, convertArticle(article.Edges.Article, int(likeCount)))
	}

	response.JSON(ctx, gin.H{
		"articles": responseArticles,
	})
}

func convertArticle(article *ent.Article, likes int) *pb.Article {
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
		Likes:       int32(likes),
		Views:       int32(article.ViewCount),
		Status:      pb.ArticleStatus(pb.ArticleStatus_value[article.Status.String()]),
		CreatedDate: timestamppb.New(article.CreatedAt),
		UpdatedDate: timestamppb.New(article.LastUpdatedAt),
	}
}
