package schema

import "entgo.io/ent"

// Statuses holds the schema definition for the Statuses entity.
type Statuses struct {
	ent.Schema
}

// Fields of the Statuses.
func (Statuses) Fields() []ent.Field {
	return nil
}

// Edges of the Statuses.
func (Statuses) Edges() []ent.Edge {
	return nil
}
