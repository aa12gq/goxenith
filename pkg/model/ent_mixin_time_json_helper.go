package model

import (
	"context"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"time"
)

// EntMixinTimeJsonHelper 重构create_at, update_at, delete_at, json序列化时的字段名.
type EntMixinTimeJsonHelper struct {
	mixin.Schema
}

func (EntMixinTimeJsonHelper) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").Unique(),
		field.Time("created_at").Comment("创建时间").Annotations(entsql.WithComments(true)).Annotations(entsql.WithComments(true)).Immutable().Default(time.Now).StructTag(`json:"created_atP"`),
		field.Time("updated_at").Comment("更新时间").Annotations(entsql.WithComments(true)).Default(time.Now).UpdateDefault(time.Now).StructTag(`json:"updated_atP"`),
		field.Time("deleted_at").Comment("逻辑删除时间").Annotations(entsql.WithComments(true)).Optional().StructTag(`json:"deleted_atP"`),
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
func (EntMixinTimeJsonHelper) Hooks() []ent.Hook {
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
