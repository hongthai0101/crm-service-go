package services

import (
	"crm-service-go/app/clients"
	"crm-service-go/app/repositories"
	"crm-service-go/pkg"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type Service struct {
	LeadService            *LeadService
	NoteService            *NoteService
	SaleOpportunityService *SaleOpportunityService
	StatisticsService      *StatisticsService
}

func NewService(repo *repositories.Repository, client *clients.HttpClient) *Service {
	topicService := NewTopicService()
	saleOpportunityService := NewSaleOpportunityService(
		repo.SaleRepo,
		repo.TagRepo,
		repo.LogRepo,
		repo.NoteRepo,
		NewLeadService(repo.LeadRepo),
		topicService,
		client.EmployeeClient,
		client.DigitalClient,
		client.ContractClient,
	)
	return &Service{
		LeadService:            NewLeadService(repo.LeadRepo),
		NoteService:            NewNoteService(repo.NoteRepo),
		SaleOpportunityService: saleOpportunityService,
		StatisticsService: NewStatisticsService(
			repo.SaleRepo,
			client.EmployeeClient,
			client.MasterDataClient,
		),
	}
}

func mapPolicyToFilter(ctx *gin.Context, filter bson.D) bson.D {
	if userStoreCode, ok := ctx.Get(pkg.ProjectKeyUserStoreCodes); ok {
		storeCodes := userStoreCode.([]string)
		if storeCodes[0] != "*" {
			filter = append(filter, bson.E{
				Key: "$and",
				Value: bson.A{
					bson.D{{
						"storeCode", bson.D{{
							"$in", storeCodes,
						}},
					}},
				},
			})
		}
	}
	return filter
}
