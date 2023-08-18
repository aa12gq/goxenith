package v1

import (
	"github.com/gin-gonic/gin"
	entLink "goxenith/app/models/ent/link"
	"goxenith/dao"
	"goxenith/pkg/cache"
	"goxenith/pkg/helpers"
	"goxenith/pkg/model"
	"goxenith/pkg/response"
	pb "goxenith/proto/app/v1"
	"time"
)

type LinksController struct {
	BaseAPIController
}

func (ctrl *LinksController) Index(ctx *gin.Context) {
	cacheKey := "links:all"
	expireTime := 120 * time.Minute

	var cachedLinks []*pb.Link
	// 从缓存中获取数据
	cache.GetObject(cacheKey, &cachedLinks)
	if !helpers.Empty(cachedLinks) {
		response.JSON(ctx, &pb.ListLinkReply{Links: cachedLinks})
		return
	}

	links, err := dao.DB.Link.Query().Where(entLink.DeleteEQ(model.DeletedNo)).All(ctx)
	if err != nil || helpers.Empty(links) {
		response.Abort404(ctx)
		return
	}

	rv := make([]*pb.Link, 0, len(links))
	for _, v := range links {
		rv = append(rv, &pb.Link{
			Id:      v.ID,
			Name:    v.Name,
			Url:     v.URL,
			ImgPath: v.ImgPath,
		})
	}

	// 设置缓存
	cache.Set(cacheKey, rv, expireTime)

	response.JSON(ctx, &pb.ListLinkReply{Links: rv})
}
