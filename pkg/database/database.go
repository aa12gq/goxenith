package database

import (
	"context"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type DAO struct {
	DbDriver dialect.Driver
}

func NewDAO(driver string, dsn string) (*DAO, func(), error) {
	drv, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database, err: %v", err)
	}
	cleanup := func() {
		if err := drv.Close(); err != nil {
			fmt.Println(err.Error())
		}
	}

	sqlDrv := dialect.DebugWithContext(drv, func(ctx context.Context, i ...interface{}) {
		//l.WithContext(ctx).Info(i...)
	})

	return &DAO{
		DbDriver: sqlDrv,
	}, cleanup, nil
}
