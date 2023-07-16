package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"goxenith/bootstrap"
	btsConfig "goxenith/config"
	"goxenith/pkg/config"
)

func init() {
	// 加载 config 目录下的配置信息
	btsConfig.Initialize()
}

func main() {

	var env string
	flag.StringVar(&env, "env", "", "加载 .env 文件，如 --env=testing 加载的是 .env.testing 文件")
	flag.Parse()
	config.InitConfig(env)
	bootstrap.SetupLogger()
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	bootstrap.SetupDB()
	bootstrap.SetupRoute(router)

	err := router.Run(":" + config.Get("app.port"))
	if err != nil {
		fmt.Println(err.Error())
	}
}
