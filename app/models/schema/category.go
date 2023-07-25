package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"goxenith/pkg/model"
)

// Category 分类
type Category struct {
	ent.Schema
}

func (Category) Mixin() []ent.Mixin {
	return []ent.Mixin{
		model.EntityStatMixin{},
	}
}

func (Category) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "category",
		},
	}
}

// Fields of the Category.
func (Category) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Comment("分类名称").
			Annotations(entsql.Annotation{Size: 100}).NotEmpty(),
		field.Uint64("parent_id").Optional().Comment("上级分类id"),
	}
}

// Edges of the Category.
func (Category) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("children", Category.Type).Ref("parent"),
		edge.To("parent", Category.Type).Unique().Field("parent_id").Comment("上级分类id").
			StorageKey(edge.Symbol("parent")),
	}
}

func (Category) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name", "parent_id").Unique(),
	}
}
