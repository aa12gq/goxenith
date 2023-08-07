package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"goxenith/pkg/model"
)

type Community struct {
	ent.Schema
}

func (Community) Mixin() []ent.Mixin {
	return []ent.Mixin{
		model.EntityStatMixin{},
	}
}

// Fields of the Community.
func (Community) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").Comment("社区ID").Annotations(entsql.WithComments(true)),
		field.String("name").Comment("社区名称").Annotations(entsql.WithComments(true)),
		field.String("logo").Comment("社区logo").Annotations(entsql.WithComments(true)).Optional(),
		field.String("introduce").Comment("社区介绍").Annotations(entsql.WithComments(true)),
	}
}

// Edges of the Community.
func (Community) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("articles", Article.Type),
	}
}
