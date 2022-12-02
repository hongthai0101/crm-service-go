package repositories

import (
	"context"
	"crm-service-go/app/entities"
	"crm-service-go/datasources"
)

type TagRepository struct {
	BaseRepo *BaseRepository[entities.Tag]
}

func NewTagRepository(ctx context.Context) *TagRepository {
	return &TagRepository{
		BaseRepo: &BaseRepository[entities.Tag]{
			col: datasources.MongoDatabase.Collection(entities.CollectionTag),
			ctx: ctx,
		},
	}
}
