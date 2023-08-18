package v1

import (
	"github.com/gin-gonic/gin"
	entLink "goxenith/app/models/ent/link"
	"goxenith/dao"
	"goxenith/pkg/model"
	"goxenith/pkg/response"
	pb "goxenith/proto/app/v1"
)

type LinksController struct {
	BaseAPIController
}

func (ctrl *LinksController) Index(ctx *gin.Context) {
	links, err := dao.DB.Link.Query().Where(entLink.DeleteEQ(model.DeletedNo)).All(ctx)
	if err != nil {
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

	response.JSON(ctx, &pb.ListLinkReply{Links: rv})
}
