package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"goxenith/app/models"
	"goxenith/app/models/ent"
	"goxenith/pkg/config"
	"goxenith/pkg/database"
)

// SetupDB 初始化数据库和 ORM
func SetupDB() {

	switch config.Get("database.connection") {
	case "mysql":
		// 构建 DSN 信息
		dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&multiStatements=true&loc=Local",
			config.Get("database.mysql.username"),
			config.Get("database.mysql.password"),
			config.Get("database.mysql.host"),
			config.Get("database.mysql.port"),
			config.Get("database.mysql.database"),
			config.Get("database.mysql.charset"),
		)
		// 连接数据库
		drv, _, err := database.NewDAO("mysql", dsn)
		if err != nil {
			panic(errors.New("New Dao error "))
		}
		database.DB = ent.NewClient(ent.Driver(drv.DbDriver))
		if err := models.Migrate(context.Background(), drv, &models.MigrateOptions{
			Debug:            true,
			DropColumn:       false,
			DropIndex:        false,
			CreateForeignKey: false,
		}); err != nil {
			panic(err)
		}
	default:
		panic(errors.New("database connection not supported"))
	}

}
