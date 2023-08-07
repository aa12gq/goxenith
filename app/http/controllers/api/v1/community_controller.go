package v1

import (
	"github.com/gin-gonic/gin"
	entComm "goxenith/app/models/ent/community"
	"goxenith/dao"
	"goxenith/pkg/logger"
	"goxenith/pkg/model"
	"goxenith/pkg/response"
	pb "goxenith/proto/app/v1"
)

type CommunityController struct {
	BaseAPIController
}

func (c *CommunityController) ListCommunity(ctx *gin.Context) {
	communitys, err := dao.DB.Community.Query().
		Where(entComm.DeleteEQ(model.DeletedNo)).
		All(ctx)
	if err != nil {
		logger.LogWarnIf("社区列表出错: %v", err)
		response.Error(ctx, err)
		return
	}
	var cs []*pb.Community
	for _, em := range communitys {
		cs = append(cs, &pb.Community{
			Id:        em.ID,
			Name:      em.Name,
			Logo:      em.Logo,
			Introduce: em.Introduce,
		})
	}
	response.JSON(ctx, cs)
}
