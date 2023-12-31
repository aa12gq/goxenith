package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"goxenith/pkg/model"
	"time"
)

type Article struct {
	ent.Schema
}

func (Article) Mixin() []ent.Mixin {
	return []ent.Mixin{
		model.EntityStatMixin{},
	}
}

// Indexes of the Article.
func (Article) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("author_id"),
	}
}

// Fields of the Article.
func (Article) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").Comment("博文ID").Annotations(entsql.WithComments(true)),
		field.Uint64("author_id").Comment("博文作者ID").Annotations(entsql.WithComments(true)),
		field.String("title").Comment("博文标题").Annotations(entsql.WithComments(true)),
		field.String("summary").Comment("博文摘要").Annotations(entsql.WithComments(true)).Optional(),
		field.Text("content").Comment("博文内容").Annotations(entsql.WithComments(true)),
		field.Enum("status").Values("DRAFT", "EFFECT").Default("DRAFT").Comment("博文状态, DRAFT:草稿,EFFECT:生效").Annotations(entsql.WithComments(true)),
		field.Uint64("view_count").Default(0).Comment("浏览量").Annotations(entsql.WithComments(true)),
		field.Time("last_updated_at").Comment("最后更新时间").Annotations(entsql.WithComments(true)).Default(time.Now),
	}
}

// Edges of the Article.
func (Article) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("author", User.Type).Ref("articles").Field("author_id").Required().Unique(),
	}
}
