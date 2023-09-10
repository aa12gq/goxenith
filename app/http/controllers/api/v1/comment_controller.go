package v1

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
	"goxenith/app/models/ent"
	"goxenith/app/models/ent/comment"
	"goxenith/dao"
	"goxenith/pkg/auth"
	"goxenith/pkg/logger"
	"goxenith/pkg/model"
	"goxenith/pkg/response"
	pb "goxenith/proto/app/v1"
	"strconv"
)

const (
	PageSize     = 20
	MaxNestLevel = 2 // 限制为2层嵌套
)

type CommentController struct {
	BaseAPIController
}

func (a *CommentController) AddComment(ctx *gin.Context) {
	request := &pb.AddCommentRequest{}
	if err := ctx.ShouldBind(request); err != nil {
		response.BadRequest(ctx, err, "请求解析错误")
		return
	}

	currentUser := auth.CurrentUser(ctx)
	if currentUser.ID == 0 {
		response.Abort403(ctx)
		return
	}

	builder := dao.DB.Comment.Create().
		SetContent(request.Content).
		SetArticleID(request.ArticleId).
		SetUserID(currentUser.ID)

	if request.ParentId != 0 {
		builder = builder.SetParentID(request.ParentId)
	}

	cmt, err := builder.Save(ctx)

	if err != nil {
		logger.LogWarnIf("添加评论失败", err)
		response.Abort500(ctx, "添加评论失败")
		return
	}

	response.JSON(ctx, gin.H{
		"message": "评论成功",
		"comment": cmt,
	})
}

func fetchNestedComments(ctx *gin.Context, cmt *ent.Comment) *pb.Comment {
	rpcComment := convertComment(cmt)
	childComments, _ := dao.DB.Comment.Query().Where(comment.ParentID(cmt.ID)).WithUser().All(ctx)
	for _, child := range childComments {
		rpcComment.ChildComments = append(rpcComment.ChildComments, fetchNestedComments(ctx, child))
	}
	return rpcComment
}

// GetComments 仅加载一级评论
func (a *CommentController) GetComments(ctx *gin.Context) {
	articleIDStr, _ := ctx.Params.Get("articleId")
	articleID, _ := strconv.Atoi(articleIDStr)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	offset := (page - 1) * PageSize

	// 仅加载顶级评论
	topLevelComments, err := dao.DB.Comment.Query().
		Where(comment.ArticleID(uint64(articleID)), comment.DeleteEQ(model.DeletedNo), comment.ParentIDEQ(0)).
		WithUser().
		WithArticle().
		Offset(offset).
		Limit(PageSize).
		Order(ent.Desc(comment.FieldCreatedAt)).
		All(ctx)

	if err != nil {
		logger.LogWarnIf("获取评论失败", err)
		response.Abort500(ctx, "获取评论失败")
		return
	}

	responseComments := make([]*pb.Comment, len(topLevelComments))
	for i, cmt := range topLevelComments {
		responseComments[i] = convertComment(cmt)
	}

	resp := &pb.TopLevelCommentsResponse{
		TopLevelComments: responseComments,
	}
	response.JSON(ctx, resp)
}

// GetChildComments 获取指定父评论下的子评论
func (a *CommentController) GetChildComments(ctx *gin.Context) {
	parentIDStr, _ := ctx.Params.Get("parentId")
	parentID, _ := strconv.Atoi(parentIDStr)

	childComments, err := dao.DB.Comment.Query().
		Where(comment.ParentIDEQ(uint64(parentID)), comment.DeleteEQ(model.DeletedNo)).
		WithUser().
		Order(ent.Desc(comment.FieldCreatedAt)).
		All(ctx)

	if err != nil {
		logger.LogWarnIf("获取子评论失败", err)
		response.Abort500(ctx, "获取子评论失败")
		return
	}

	responseComments := make([]*pb.Comment, len(childComments))
	for i, cmt := range childComments {
		responseComments[i] = convertComment(cmt)
	}

	resp := &pb.ChildCommentsResponse{
		ChildComments: responseComments,
	}
	response.JSON(ctx, resp)
}

// GetFullCommentTree 获取整个评论树
func (a *CommentController) GetFullCommentTree(ctx *gin.Context) {
	articleIDStr, _ := ctx.Params.Get("articleId")
	articleID, _ := strconv.Atoi(articleIDStr)

	// 获取所有与该文章ID相关的评论
	allComments, err := dao.DB.Comment.Query().
		Where(comment.ArticleID(uint64(articleID)), comment.DeleteEQ(model.DeletedNo)).
		WithUser().
		Order(ent.Desc(comment.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		logger.LogWarnIf("获取评论失败", err)
		response.Abort500(ctx, "获取评论失败")
		return
	}

	// 创建评论树
	var commentTree []*pb.Comment
	for _, cmt := range allComments {
		if cmt.ParentID == 0 {
			commentTree = append(commentTree, fetchNestedComments(ctx, cmt))
		}
	}

	resp := &pb.FullCommentTreeResponse{
		Comments: commentTree,
	}
	response.JSON(ctx, resp)
}

func convertComment(comment *ent.Comment) *pb.Comment {
	return &pb.Comment{
		Id: comment.ID,
		Author: &pb.Author{
			Id:     comment.Edges.User.ID,
			Name:   comment.Edges.User.UserName,
			Avatar: comment.Edges.User.Avatar,
		},
		Content:     comment.Content,
		ParentId:    comment.ParentID,
		CreatedDate: timestamppb.New(comment.CreatedAt),
	}
}

func (a *CommentController) DeleteComment(ctx *gin.Context) {
	commentIDStr, _ := ctx.Params.Get("id")
	commentID, _ := strconv.Atoi(commentIDStr)

	currentUser := auth.CurrentUser(ctx)
	if currentUser.ID == 0 {
		response.Abort403(ctx)
		return
	}

	cmm, err := dao.DB.Comment.Query().Where(comment.ID(uint64(commentID)),
		comment.DeleteEQ(model.DeletedNo)).WithUser().Only(ctx)
	if err != nil || cmm.Edges.User.ID != currentUser.ID {
		response.Abort403(ctx, "您没有权限删除此评论")
		return
	}

	err = dao.DB.Comment.DeleteOneID(uint64(commentID)).Exec(ctx)
	if err != nil {
		logger.LogWarnIf("删除评论失败", err)
		response.Abort500(ctx, "删除评论失败")
		return
	}

	response.JSON(ctx, gin.H{
		"message": "评论删除成功",
	})
}
