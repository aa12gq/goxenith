package bootstrap

import (
	"errors"
	"fmt"
	"goxenith/app/ent"
	"goxenith/pkg/config"
	"goxenith/pkg/database"
)

var DB *ent.Client

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
		DB = ent.NewClient(ent.Driver(drv.DbDriver))
	default:
		panic(errors.New("database connection not supported"))
	}

}
