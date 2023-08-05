package bootstrap

import (
	"database/sql"
	entSql "entgo.io/ent/dialect/sql"
	"fmt"
	"goxenith/app/models/ent"
	"goxenith/dao"
	"goxenith/pkg/config"
	"time"
)

func SetupDB() {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&multiStatements=true&loc=Local",
		config.Get("database.mysql.username"),
		config.Get("database.mysql.password"),
		config.Get("database.mysql.host"),
		config.Get("database.mysql.port"),
		config.Get("database.mysql.database"),
		config.Get("database.mysql.charset"),
	)
	db, err := sql.Open(config.Get("database.connection"), dsn)
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(config.GetInt("database.mysql.max_idle_connections"))
	db.SetMaxOpenConns(config.GetInt("database.mysql.max_open_connections"))
	db.SetConnMaxLifetime(time.Duration(config.GetInt("database.mysql.max_life_seconds")) * time.Second)
	drv := entSql.OpenDB("mysql", db)
	dao.DB = ent.NewClient(ent.Driver(drv))
}
