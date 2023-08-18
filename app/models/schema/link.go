package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"goxenith/pkg/model"
)

type Link struct {
	ent.Schema
}

func (Link) Mixin() []ent.Mixin {
	return []ent.Mixin{
		model.EntityStatMixin{},
	}
}

// Fields of the Link.
func (Link) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").Comment("友情链接ID"),
		field.String("name").Comment("友情链接名称"),
		field.String("url").Comment("友情链接URL"),
		field.String("img_path").Comment("友情链接图片链接"),
	}
}

// Edges of the Link.
func (Link) Edges() []ent.Edge {
	return nil
}
