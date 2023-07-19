package serve

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	appCmd "goxenith/app/cmd"
	"goxenith/app/models/ent"
	"goxenith/bootstrap"
	"goxenith/dao"
	"goxenith/pkg/config"
	"goxenith/pkg/console"
	"goxenith/pkg/logger"
	"goxenith/pkg/redis"
)

type App struct {
	Engine *gin.Engine
	Redis  *redis.RedisClient
	DB     *ent.Client
}

func newApp(engine *gin.Engine, redis *redis.RedisClient, driver *dao.DAO) *App {
	db := ent.NewClient(ent.Driver(driver.DbDriver))
	return &App{
		Engine: engine,
		Redis:  redis,
		DB:     db,
	}
}

func RunWeb(cmd *cobra.Command, args []string) error {
	var err error

	config.InitConfig(appCmd.Env)
	if err != nil {
		panic(err)
	}
	bootstrap.SetupLogger()

	gin.SetMode(gin.DebugMode)
	router := gin.New()
	bootstrap.SetupRoute(router)
	err = router.Run(":" + config.Get("app.port"))
	if err != nil {
		logger.ErrorString("CMD", "serve", err.Error())
		console.Exit("Unable to start server, error:" + err.Error())
	}
	return nil
}
