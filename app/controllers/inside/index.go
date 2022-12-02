package controller_inside

import (
	"crm-service-go/app/services"
	serviceInside "crm-service-go/app/services/inside"
)

type InsideController struct {
	InsideLeadController            *InsideLeadController
	InsideSaleOpportunityController *InsideSaleOpportunityController
}

func NewInsideController(
	service *services.Service,
	insideService *serviceInside.InsideService,
) *InsideController {
	return &InsideController{
		InsideLeadController:            NewInsideLead(service.LeadService),
		InsideSaleOpportunityController: NewInsideSaleOpportunity(insideService.InsideSaleOpportunityService),
	}
}
