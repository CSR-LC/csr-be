package repositories

import (
	"context"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/ent"
)

type ActiveAreaRepository interface {
	AllActiveAreas(ctx context.Context, limit, offset int) ([]*ent.ActiveArea, error)
	TotalActiveAreas(ctx context.Context) (int, error)
}

type activeAreaRepository struct {
	client *ent.Client
}

func NewActiveAreaRepository(client *ent.Client) ActiveAreaRepository {
	return &activeAreaRepository{client: client}
}

func (r *activeAreaRepository) AllActiveAreas(ctx context.Context, limit, offset int) ([]*ent.ActiveArea, error) {
	return r.client.ActiveArea.Query().Limit(limit).Offset(offset).All(ctx)
}

func (r *activeAreaRepository) TotalActiveAreas(ctx context.Context) (int, error) {
	return r.client.ActiveArea.Query().Count(ctx)
}
