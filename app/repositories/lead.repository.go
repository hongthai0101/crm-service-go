package repositories

import (
	"context"
	"crm-service-go/app/entities"
	"crm-service-go/datasources"
)

type LeadRepository struct {
	BaseRepo *BaseRepository[entities.Lead]
}

func NewLeadRepository(ctx context.Context) *LeadRepository {
	return &LeadRepository{
		BaseRepo: &BaseRepository[entities.Lead]{
			col: datasources.MongoDatabase.Collection(entities.CollectionLead),
			ctx: ctx,
		},
	}
}
