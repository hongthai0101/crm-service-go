package repositories

import (
	"context"
	"crm-service-go/app/entities"
	"crm-service-go/datasources"
	"crm-service-go/pkg"
	"crm-service-go/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"time"
)

type SaleOpportunityRepository struct {
	BaseRepo *BaseRepository[entities.SaleOpportunity]
}

func NewSaleOpportunityRepository(ctx context.Context) *SaleOpportunityRepository {
	return &SaleOpportunityRepository{
		BaseRepo: &BaseRepository[entities.SaleOpportunity]{
			col: datasources.MongoDatabase.Collection(entities.CollectionSaleOpportunities),
			ctx: ctx,
		},
	}
}

func (r *SaleOpportunityRepository) GenerateCode(code string) string {
	if code != "" {
		sale, _ := r.BaseRepo.FindOne(bson.M{"code": code}, nil)
		if sale == nil {
			return code
		}
	}

	prefixCode := time.Now().Format(pkg.YYMMDD)
	suffix := strconv.Itoa(utils.Random(1000, 9999))

	code = prefixCode + suffix
	for total, _ := r.BaseRepo.Count(bson.M{"code": code}); total != 0; {
		suffix = strconv.Itoa(utils.Random(1000, 9999))
		code = prefixCode + suffix
	}

	return code
}
