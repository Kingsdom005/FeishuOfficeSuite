package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"time"
)

type Department struct {
	ent.Schema
}

func (Department) Fields() []ent.Field {
	return []ent.Field{
		field.String("department_id").Unique().Immutable(),
		field.String("name"),
		field.String("name_en"),
		field.String("parent_id"),
		field.Int("order").Default(0),
		field.Bool("is_root").Default(false),
		field.Int("member_count").Default(0),
		field.Time("created_at").Default(time.Now).SchemaType(map[string]string{
			dialect.MySQL: "datetime",
		}),
		field.Time("updated_at").Default(time.Now).SchemaType(map[string]string{
			dialect.MySQL: "datetime",
		}),
	}
}

func (Department) Edges() []ent.Edge {
	return nil
}

func (Department) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("parent_id"),
	}
}