package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"goxenith/bootstrap"
	"goxenith/pkg/config"
	"goxenith/pkg/console"
	"goxenith/pkg/logger"
)

// CmdServe represents the available web sub-command.
var CmdServe = &cobra.Command{
	Use:   "serve",
	Short: "Start web server",
	Run:   runWeb,
	Args:  cobra.NoArgs,
}

func runWeb(cmd *cobra.Command, args []string) {

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	bootstrap.SetupRoute(router)
	err := router.Run(":" + config.Get("app.port"))
	if err != nil {
		logger.ErrorString("CMD", "serve", err.Error())
		console.Exit("Unable to start server, error:" + err.Error())
	}
}