package models

import (
	"context"
	entschema "entgo.io/ent/dialect/sql/schema"
	ci "goxenith/app/models/ent"
	"goxenith/dao"
	"goxenith/pkg/model"
)

type MigrateOptions struct {
	Debug            bool
	DropColumn       bool
	DropIndex        bool
	CreateForeignKey bool
}

func Migrate(ctx context.Context, daoIn *dao.DAO, opt *MigrateOptions) error {
	var schemas []model.EntSchema
	if opt.Debug {
		schemas = append(schemas,
			ci.NewClient(ci.Driver(daoIn.DbDriver)).Debug().Schema,
		)
	} else {
		schemas = append(schemas,
			ci.NewClient(ci.Driver(daoIn.DbDriver)).Schema,
		)
	}
	return model.EntMigrateSchemas(ctx, schemas,
		entschema.WithAtlas(true),
		entschema.WithDropColumn(opt.DropColumn),
		entschema.WithDropIndex(opt.DropIndex),
		entschema.WithForeignKeys(opt.CreateForeignKey),
	)
}
