package repositories

import (
	"context"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/ent"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/ent/photo"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/models"
)

type PhotoRepository interface {
	CreatePhoto(ctx context.Context, p models.Photo) (*ent.Photo, error)
	PhotoByID(ctx context.Context, id string) (*ent.Photo, error)
	DeletePhotoByID(ctx context.Context, id string) error
}

type photoRepository struct {
	client *ent.Client
}

func NewPhotoRepository(client *ent.Client) PhotoRepository {
	return &photoRepository{
		client: client,
	}
}

func (r *photoRepository) CreatePhoto(ctx context.Context, newPhoto models.Photo) (*ent.Photo, error) {
	p, err := r.client.Photo.Create().
		SetID(newPhoto.ID).
		SetURL(*newPhoto.URL).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *photoRepository) PhotoByID(ctx context.Context, id string) (*ent.Photo, error) {
	result, err := r.client.Photo.Query().Where(photo.ID(id)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *photoRepository) DeletePhotoByID(ctx context.Context, id string) error {
	_, err := r.client.Photo.Delete().Where(photo.ID(id)).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}