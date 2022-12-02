package controllers

import (
	"crm-service-go/app/repositories"
	"crm-service-go/app/services"
	"crm-service-go/pkg"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type Controller struct {
	LeadController            *LeadController
	TagController             *TagController
	NoteController            *NoteController
	SaleOpportunityController *SaleOpportunityController
	StatisticsController      *StatisticsController
}

func NewController(
	repository *repositories.Repository,
	service *services.Service,
) *Controller {
	return &Controller{
		LeadController:            NewLead(service.LeadService),
		TagController:             NewTag(repository.TagRepo),
		NoteController:            NewNote(service.NoteService),
		SaleOpportunityController: NewSaleOpportunity(service.SaleOpportunityService),
		StatisticsController:      NewStatistics(service.StatisticsService),
	}
}

func GetObjectIDFromPath(ctx *gin.Context) *primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		pkg.SendErrorResponse(ctx, pkg.ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid ObjectID",
		})
		return nil
	}
	return &id
}
