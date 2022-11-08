package repository

import (
	"context"

	"onthemat/pkg/ent"
)

type ImageRepository interface {
	Create(ctx context.Context, image *ent.Image, userID int) error
}

type imageRepository struct {
	db *ent.Client
}

func NewImageRepository(db *ent.Client) ImageRepository {
	return &imageRepository{
		db: db,
	}
}

func (repo *imageRepository) Create(ctx context.Context, image *ent.Image, userID int) error {
	return repo.db.Image.Create().
		SetName(image.Name).
		SetPath(image.Path).
		SetSize(image.Size).
		SetContentType(image.ContentType).
		SetType(image.Type).
		SetUserID(userID).Exec(ctx)
}
