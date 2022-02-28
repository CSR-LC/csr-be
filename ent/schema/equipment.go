package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/gofrs/uuid"
)

// Equipment holds the schema definition for the Equipment entity.
type Equipment struct {
	ent.Schema
}

// Fields of the Equipment.
func (Equipment) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}),
		field.String("sku").Default("unknown"),
		field.String("name").Default("unknown"),
		field.UUID("kind", uuid.UUID{}).StorageKey("kinds"),
		field.UUID("status", uuid.UUID{}).StorageKey("statuses"),
		field.Float("rate_hour"),
		field.Float("rate_day"),
		field.UUID("location", uuid.UUID{}).StorageKey("locations"),
		field.String("description").Default("unknown"),
	}
}

// Edges of the Equipment.
func (Equipment) Edges() []ent.Edge {
	return nil
}
