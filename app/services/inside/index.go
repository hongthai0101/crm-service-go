package service_inside

import (
	"crm-service-go/app/services"
)

type InsideService struct {
	InsideSaleOpportunityService *InsideSaleOpportunityService
}

func NewInsideService(service *services.Service) *InsideService {
	return &InsideService{
		InsideSaleOpportunityService: NewInsideSaleOpportunityService(service.SaleOpportunityService),
	}
}
