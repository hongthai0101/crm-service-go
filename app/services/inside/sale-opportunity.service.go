package service_inside

import (
	"crm-service-go/app/entities"
	"crm-service-go/app/services"
	"crm-service-go/app/validation"
	"crm-service-go/config"
	"crm-service-go/pkg"
	"crm-service-go/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

var createBy = config.GetConfig().DefaultDataConfig.CreatedBy

type InsideSaleOpportunityService struct {
	service *services.SaleOpportunityService
}

func NewInsideSaleOpportunityService(
	service *services.SaleOpportunityService,
) *InsideSaleOpportunityService {
	return &InsideSaleOpportunityService{
		service: service,
	}
}

func (s *InsideSaleOpportunityService) Pagination(
	ctx *gin.Context,
	params validation.PaginationSaleOpportunity,
) (*services.SaleOpportunityPagination, error) {
	return s.service.Pagination(ctx, params, false)
}

func (s *InsideSaleOpportunityService) Create(
	ctx *gin.Context,
	payload validation.CreateSaleOpportunityInside,
) (*entities.SaleOpportunity, error) {
	saleOpportunity, customer := payload.SaleOpportunity, payload.Customer
	lead, err := s.createOrUpdateLead(customer)
	if lead == nil || err != nil {
		pkg.SendErrorResponse(ctx, pkg.ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    "Lead not found",
		})
	}
	var payloadSaleOpportunity, _ = utils.TypeConverter[validation.CreateSaleOpportunity](&saleOpportunity)
	payloadSaleOpportunity.LeadId = lead.ID
	status, typeSale, source := saleOpportunity.Status, saleOpportunity.Type, saleOpportunity.Source
	if len(status) == 0 {
		payloadSaleOpportunity.Status = entities.SaleOppStatusNew
	}
	if len(typeSale) == 0 {
		payloadSaleOpportunity.Type = entities.SaleOppTypeBorrower
	}
	if len(source) == 0 {
		payloadSaleOpportunity.Source = "DIRECT_ONLINE"
	}

	return s.service.Create(ctx, payloadSaleOpportunity, false)
}

func (s *InsideSaleOpportunityService) Update(
	ctx *gin.Context,
	id primitive.ObjectID,
	payload validation.UpdateSaleOpportunity,
) (*entities.SaleOpportunity, error) {
	return nil, nil
}

func (s *InsideSaleOpportunityService) FindById(
	ctx *gin.Context,
	id primitive.ObjectID,
	includes []string,
) (*entities.SaleOpportunity, error) {
	return s.service.FindById(ctx, id, includes)
}

func (s *InsideSaleOpportunityService) createOrUpdateLead(
	payload validation.LeadInside,
) (*entities.Lead, error) {
	phone, customerId := payload.Phone, payload.CustomerId
	var lead = new(entities.Lead)
	if customerId != "" {
		lead, _ = s.service.LeadService.Repo.BaseRepo.FindOne(bson.D{
			{"customerId", customerId},
			{"deletedAt", nil},
		}, nil)
	}

	if lead == nil {
		lead, _ = s.service.LeadService.Repo.BaseRepo.FindOne(bson.D{
			{"phone", phone},
			{"deletedAt", nil},
		}, nil)
		if lead == nil {
			entity, err := utils.TypeConverter[entities.Lead](&payload)
			if err != nil {
				return nil, err
			}
			entity.CreatedBy = createBy
			entity.UpdatedBy = createBy
			entities.CreatingEntity(&entity.BaseEntity)
			lead, _ = s.service.LeadService.Repo.BaseRepo.Create(entity)
		}
	}
	return lead, nil
}
