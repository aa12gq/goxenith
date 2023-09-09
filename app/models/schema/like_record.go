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

type LikeRecord struct {
	ent.Schema
}

func (LikeRecord) Mixin() []ent.Mixin {
	return []ent.Mixin{
		model.EntityStatMixin{},
	}
}

func (LikeRecord) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "like_record",
		},
	}
}

func (LikeRecord) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "article_id").Unique(),
		index.Fields("user_id"),
		index.Fields("article_id"),
	}
}

// Fields of the LikeRecord.
func (LikeRecord) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").Comment("点赞记录表ID").Annotations(entsql.WithComments(true)),
		field.Uint64("user_id").Comment("点赞用户的ID").Annotations(entsql.WithComments(true)),
		field.Uint64("article_id").Comment("被点赞的文章ID").Annotations(entsql.WithComments(true)),
		field.Bool("is_active").Comment("是否有效的点赞，默认为true；如果用户取消点赞，则为false").Annotations(entsql.WithComments(true)),
	}
}

// Edges of the LikeRecord.
func (LikeRecord) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).Field("user_id").Required().Unique(),
		edge.To("article", Article.Type).Field("article_id").Required().Unique(),
	}
}
