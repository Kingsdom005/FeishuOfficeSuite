package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"time"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().Immutable(),
		field.String("union_id").Unique(),
		field.String("open_id").Unique(),
		field.String("name"),
		field.String("en_name"),
		field.String("email"),
		field.String("phone"),
		field.String("avatar_url"),
		field.String("avatar_thumb"),
		field.String("avatar_middle"),
		field.String("status").Default("active"),
		field.Bool("is_activated").Default(true),
		field.Bool("is_tenant_access").Default(true),
		field.String("department_id"),
		field.Time("created_at").Default(time.Now).SchemaType(map[string]string{
			dialect.MySQL: "datetime",
		}),
		field.Time("updated_at").Default(time.Now).SchemaType(map[string]string{
			dialect.MySQL: "datetime",
		}),
	}
}

func (User) Edges() []ent.Edge {
	return nil
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email"),
		index.Fields("phone"),
		index.Fields("union_id"),
		index.Fields("department_id"),
	}
}