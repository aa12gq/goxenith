package serve

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"goxenith/bootstrap"
	"goxenith/pkg/config"
	"goxenith/pkg/console"
	"goxenith/pkg/logger"
)

func RunWeb(cmd *cobra.Command, args []string) error {
	gin.SetMode(gin.DebugMode)
	router := gin.New()
	bootstrap.SetupRoute(router)

	err := router.Run(":" + config.Get("app.port"))
	if err != nil {
		logger.ErrorString("CMD", "serve", err.Error())
		console.Exit("Unable to start server, error:" + err.Error())
	}
	return nil
}
