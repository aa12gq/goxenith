package model

import (
	"context"
	"io"

	schema "entgo.io/ent/dialect/sql/schema"
)

// EntSchema 适配entgo生成的Schema对象
type EntSchema interface {
	Create(ctx context.Context, opts ...schema.MigrateOption) error
	WriteTo(ctx context.Context, w io.Writer, opts ...schema.MigrateOption) error
}

func EntMigrateSchemas(ctx context.Context, schemas []EntSchema, opts ...schema.MigrateOption) error {
	for _, schema := range schemas {
		if err := schema.Create(ctx, opts...); err != nil {
			return err
		}
	}
	return nil
}
