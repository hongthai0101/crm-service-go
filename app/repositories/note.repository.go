package repositories

import (
	"context"
	"crm-service-go/app/entities"
	"crm-service-go/datasources"
)

type NoteRepository struct {
	BaseRepo *BaseRepository[entities.Note]
}

func NewNoteRepository(ctx context.Context) *NoteRepository {
	return &NoteRepository{
		BaseRepo: &BaseRepository[entities.Note]{
			col: datasources.MongoDatabase.Collection(entities.CollectionNote),
			ctx: ctx,
		},
	}
}
