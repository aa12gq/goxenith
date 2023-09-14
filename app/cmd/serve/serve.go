package serve

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	v1 "goxenith/app/http/controllers/api/v1"
	"goxenith/bootstrap"
	"goxenith/pkg/config"
	"goxenith/pkg/console"
	"goxenith/pkg/logger"
)

func RunWeb(cmd *cobra.Command, args []string) error {
	gin.SetMode(gin.DebugMode)
	router := gin.New()
	bootstrap.SetupRoute(router)
	bootstrap.SetupOss()
	go v1.SyncArticleViewsFromRedis()

	err := router.Run(":" + config.Get("app.port"))
	if err != nil {
		logger.ErrorString("CMD", "serve", err.Error())
		console.Exit("Unable to start server, error:" + err.Error())
	}
	return nil
}
