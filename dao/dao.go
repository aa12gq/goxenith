package dao

import (
	"context"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"goxenith/app/models/ent"
	"goxenith/pkg/config"
)

var DB *ent.Client

type DAO struct {
	DbDriver dialect.Driver
}

func NewDAO() (*DAO, func(), error) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&multiStatements=true&loc=Local",
		config.Get("database.mysql.username"),
		config.Get("database.mysql.password"),
		config.Get("database.mysql.host"),
		config.Get("database.mysql.port"),
		config.Get("database.mysql.database"),
		config.Get("database.mysql.charset"),
	)
	drv, err := sql.Open(config.Get("database.connection"), dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database, err: %v", err)
	}
	cleanup := func() {
		if err := drv.Close(); err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
	}

	sqlDrv := dialect.DebugWithContext(drv, func(ctx context.Context, i ...interface{}) {
		fmt.Println(i)
	})

	return &DAO{
		DbDriver: sqlDrv,
	}, cleanup, nil
}
