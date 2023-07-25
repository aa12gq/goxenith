package v1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"goxenith/app/models/ent"
	entcate "goxenith/app/models/ent/category"
	"goxenith/dao"
	"goxenith/pkg/logger"
	"goxenith/pkg/model"
	"goxenith/pkg/response"
	pb "goxenith/proto/app/v1"
	"strconv"
)

type Category struct {
	// 分类ID
	ID uint64
	// 分类名称
	Name string
	// 上级分类id
	ParentID uint64
}

type CategoryTreeNode struct {
	// 分类ID
	ID uint64
	// 分类名称
	Name string
	// 上级分类节点
	Parent *CategoryTreeNode
	// 子节点
	Children []*CategoryTreeNode
}

type CategoryTree struct {
	// 节点
	Nodes []*CategoryTreeNode
}

type CategoryController struct {
	BaseAPIController
}

func (cate *CategoryController) ValidateMaterialCategory(ctx *gin.Context, mc *pb.Category) (*pb.Category, error) {
	if mc.ParentId > 0 {
		_, err := dao.DB.Category.Get(ctx, mc.ParentId)
		if err != nil {
			if err != nil {
				if ent.IsNotFound(err) {
					return nil, errors.New("父分类不存在")
				} else {
					logger.LogWarnIf("数据访问出错", err)
					return nil, errors.New("数据访问出错")
				}
			}
		}
	}
	// 同级下是否存在相同名称分类
	dup, err := dao.DB.Category.Query().Where(
		entcate.NameEQ(mc.Name),
		entcate.ParentIDEQ(mc.ParentId),
		entcate.IDNEQ(mc.Id),
		entcate.DeleteEQ(model.DeletedNo),
	).Exist(ctx)
	if err != nil {
		logger.LogWarnIf("分类查询出错: %v", err)
		return nil, errors.New("分类查询出错")
	} else if dup {
		return nil, errors.New(fmt.Sprintf("同级分类名重复: %v", mc.Name))
	}
	return mc, nil
}

func (cate *CategoryController) CreateCategory(ctx *gin.Context) {
	request := &pb.CreateCategoryRequest{}

	if err := ctx.ShouldBind(&request.Category); err != nil {
		response.BadRequest(ctx, err, "请求解析错误，请确认请求格式是否正确。上传文件请使用 multipart 标头，参数请使用 JSON 格式。")
		return
	}

	if request.Category.Name == "" {
		response.Error(ctx, errors.New("分类名称为必填项"))
		return
	}

	category, err := cate.ValidateMaterialCategory(ctx, request.Category)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	nm, err := dao.DB.Category.Create().
		SetName(category.Name).
		SetParentID(category.ParentId).
		Save(ctx)
	if err != nil {
		logger.LogWarnIf("分类保存出错: %v", err)
		response.Error(ctx, errors.New("父类保存出错"))
		return
	}
	response.CreatedJSON(ctx, pb.CreateCategoryReply{Category: &pb.Category{
		Id:       nm.ID,
		ParentId: nm.ParentID,
		Name:     nm.Name,
	}})
}

func (cate *CategoryController) ListCategory(ctx *gin.Context) {
	request := pb.ListCategoryRequest{}
	if err := ctx.ShouldBind(&request); err != nil {
		response.BadRequest(ctx, err, "请求解析错误，请确认请求格式是否正确。上传文件请使用 multipart 标头，参数请使用 JSON 格式。")
		return
	}

	cats, err := dao.DB.Category.Query().
		Where(entcate.DeleteEQ(model.DeletedNo), entcate.ParentIDEQ(request.ParentId)).
		All(ctx)
	if err != nil {
		logger.LogWarnIf("分类查询出错: %v", err)
		response.Error(ctx, err)
		return
	}
	var mcats []*pb.Category
	for _, em := range cats {
		mcats = append(mcats, &pb.Category{
			Id:       em.ID,
			Name:     em.Name,
			ParentId: em.ParentID,
		})
	}
	response.JSON(ctx, mcats)
}

func (rp *CategoryController) GetCategory(ctx *gin.Context) {
	idStr, _ := ctx.Params.Get("id")
	idNum, _ := strconv.Atoi(idStr)
	cat, err := dao.DB.Category.Query().
		Where(entcate.DeleteEQ(model.DeletedNo), entcate.IDEQ(uint64(idNum))).
		First(ctx)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.JSON(ctx, &pb.GetCategoryReply{Category: &pb.Category{
		Id:       cat.ID,
		ParentId: cat.ParentID,
		Name:     cat.Name,
	}})
}

func (rp *CategoryController) GetMaterialCategoryTree(ctx *gin.Context) {
	cats, err := dao.DB.Category.Query().
		Where(entcate.DeleteEQ(model.DeletedNo)).
		All(ctx)
	if err != nil {
		logger.LogWarnIf("分类查询出错: %v", err)
		response.Error(ctx, errors.New("分类查询出错"))
		return
	}
	nodes := &CategoryTree{Nodes: rp.findCategoryChildren(cats)}
	response.JSON(ctx, rp.convertToPbCategoryChildren(nil, nodes.Nodes))
}

func (rp *CategoryController) findCategoryChildren(cats []*ent.Category) []*CategoryTreeNode {
	catMap := make(map[uint64][]*ent.Category)
	for _, cat := range cats {
		catMap[cat.ParentID] = append(catMap[cat.ParentID], cat)
	}

	return rp.getChildren(0, catMap)
}

func (rp *CategoryController) getChildren(parentId uint64, catMap map[uint64][]*ent.Category) []*CategoryTreeNode {
	children := make([]*CategoryTreeNode, 0)

	for _, cat := range catMap[parentId] {
		node := &CategoryTreeNode{
			ID:       cat.ID,
			Name:     cat.Name,
			Children: rp.getChildren(cat.ID, catMap),
		}
		children = append(children, node)
	}
	return children
}

func (s *CategoryController) convertToPbCategoryChildren(parent *CategoryTreeNode, nodes []*CategoryTreeNode) []*pb.CategoryTreeNode {
	var pbNodes []*pb.CategoryTreeNode
	for _, n := range nodes {
		pbNode := &pb.CategoryTreeNode{
			Id:       n.ID,
			Name:     n.Name,
			Children: s.convertToPbCategoryChildren(n, n.Children),
		}
		if parent != nil {
			pbNode.ParentId = parent.ID
		} else {
			pbNode.ParentId = 0
		}
		pbNodes = append(pbNodes, pbNode)
	}
	return pbNodes
}
