package model

// entgo schema公用字段

import (
	"context"
	"entgo.io/ent/dialect/entsql"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type EntityStatMixin struct {
	mixin.Schema
}

const (
	DeletedNo  uint = 0
	DeletedYes uint = 1
)

func (EntityStatMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").Unique().Comment("默认主键").Annotations(entsql.WithComments(true)),
		field.Time("created_at").Comment("创建时间").Annotations(entsql.WithComments(true)).Immutable().Default(time.Now),
		field.Time("updated_at").Comment("更新时间").Annotations(entsql.WithComments(true)).Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").Comment("逻辑删除时间").Annotations(entsql.WithComments(true)).Optional(),
		field.Uint("_delete_").Comment("逻辑标识。0 未删除, 1 已删除").Annotations(entsql.WithComments(true)).
			Default(DeletedNo),
	}
}

// Hooks 定义通用的hook.
//
// 注意: 需要在对应的shema client初始化之前引入相关schema ent下的runtime包进行hook注册:
//
//	import(
//	 _ "rd.nuggets.com/project_a/app/app1/internal/data/schema1/ent/runtime"
//	)
func (EntityStatMixin) Hooks() []ent.Hook {
	return []ent.Hook{
		// 当_delete_被设置成删除状态值时，自动更新delete_at字段.
		func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				if m.Op().Is(ent.OpUpdate) || m.Op().Is(ent.OpUpdate) {
					if dv, yes := m.Field("_delete_"); yes {
						if dv.(uint) == DeletedYes {
							err := m.SetField("deleted_at", time.Now())
							if err != nil {
								panic(err)
							}
						}
					}
				}
				return next.Mutate(ctx, m)
			})
		},
	}
}
