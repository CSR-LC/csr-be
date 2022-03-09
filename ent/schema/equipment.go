package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Equipment holds the schema definition for the Equipment entity.
type Equipment struct {
	ent.Schema
}

// Fields of the Equipment.
func (Equipment) Fields() []ent.Field {
	return []ent.Field{
		field.String("sku").Default("unknown"),
		field.String("name").Default("unknown"),
		field.Int64("rate_hour"),
		field.Int64("rate_day"),
		field.String("description").Default("unknown"),
	}

}

// Edges of the Equipment.
func (Equipment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("kinds", Kinds.Type),
		edge.To("statuses", Statuses.Type),
		edge.To("locations", Locations.Type),
	}

}
