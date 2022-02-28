package schema

import "entgo.io/ent"

// Kinds holds the schema definition for the Kinds entity.
type Kinds struct {
	ent.Schema
}

// Fields of the Kinds.
func (Kinds) Fields() []ent.Field {
	return nil
}

// Edges of the Kinds.
func (Kinds) Edges() []ent.Edge {
	return nil
}
