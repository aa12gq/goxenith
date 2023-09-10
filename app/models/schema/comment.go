package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"goxenith/pkg/model"
)

type Comment struct {
	ent.Schema
}

func (Comment) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "comment",
		},
	}
}

func (Comment) Mixin() []ent.Mixin {
	return []ent.Mixin{
		model.EntityStatMixin{},
	}
}

// Fields of the Comment.
func (Comment) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").Comment("评论ID").Annotations(entsql.WithComments(true)),
		field.Uint64("article_id").Comment("评论所属文章").Annotations(entsql.WithComments(true)),
		field.Uint64("user_id").Comment("评论用户"),
		field.Text("content").Comment("评论内容").Annotations(entsql.WithComments(true)),
		field.Uint64("parent_id").Comment("父评论ID").Optional(),
		field.Uint64("comment_id").Comment("评论他人评论ID").Optional(),
	}
}

// Edges of the Comment.
func (Comment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).Field("user_id").Required().Unique().
			Comment("评论用户"),
		edge.To("article", Article.Type).Field("article_id").Required().Unique(),
		edge.From("replies", Comment.Type).Ref("parent").Comment("子评论"),
		edge.To("parent", Comment.Type).Unique().Field("parent_id").
			Comment("父评论").StorageKey(edge.Symbol("comment_parent")),
	}
}
