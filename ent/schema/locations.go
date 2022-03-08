package schema

import "entgo.io/ent"

// Locations holds the schema definition for the Locations entity.
type Locations struct {
	ent.Schema
}

// Fields of the Locations.
func (Locations) Fields() []ent.Field {
	return nil
}

// Edges of the Locations.
func (Locations) Edges() []ent.Edge {
	return nil
}
