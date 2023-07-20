package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"goxenith/pkg/model"
)

type Project struct {
	ent.Schema
}

func (Project) Mixin() []ent.Mixin {
	return []ent.Mixin{
		model.EntityStatMixin{},
	}
}

// Fields of the Project.
func (Project) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").Comment("ProjectID"),
	}
}

// Edges of the Project.
func (Project) Edges() []ent.Edge {
	return nil
}
