package repositories

import (
	"context"
	"crm-service-go/app/entities"
	"crm-service-go/datasources"
)

type LogRepository struct {
	BaseRepo *BaseRepository[entities.Log]
}

func NewLogRepository(ctx context.Context) *LogRepository {
	return &LogRepository{
		BaseRepo: &BaseRepository[entities.Log]{
			col: datasources.MongoDatabase.Collection(entities.CollectionLog),
			ctx: ctx,
		},
	}
}
