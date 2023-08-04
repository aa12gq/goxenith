package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"goxenith/pkg/model"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		model.EntityStatMixin{},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").Comment("用户ID"),
		field.String("user_name").Comment("用户名").Annotations(entsql.Annotation{Size: 20}, entsql.WithComments(true)),
		field.String("real_name").Comment("真实姓名").Optional().Annotations(entsql.Annotation{Size: 20}, entsql.WithComments(true)),
		field.String("phone").Comment("联系电话").Unique().Annotations(entsql.Annotation{Size: 100}, entsql.WithComments(true)),
		field.String("city").Comment("城市").Optional().Annotations(entsql.Annotation{Size: 20}, entsql.WithComments(true)),
		field.Enum("gender").Values("MALE", "FEMALE").Optional().Comment("男/女").Annotations(entsql.WithComments(true)),
		field.Uint8("age").Comment("年龄").Optional().Annotations(entsql.WithComments(true)),
		field.Time("birthday").Comment("出生日期").Optional().Annotations(entsql.WithComments(true)),
		field.String("password").Comment("密码").Annotations(entsql.Annotation{Size: 60}, entsql.WithComments(true)),
		field.String("personal_profile").Comment("个人简介").Optional().Annotations(entsql.Annotation{Size: 1024}, entsql.WithComments(true)),
		field.String("email").Comment("邮箱").Optional().Annotations(entsql.Annotation{Size: 30}, entsql.WithComments(true)),
		field.String("avatar").Comment("头像").Optional().Annotations(entsql.WithComments(true)),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("articles", Article.Type),
	}
}
