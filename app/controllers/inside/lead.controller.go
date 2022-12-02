package controller_inside

import (
	"crm-service-go/app/services"
	"crm-service-go/app/validation"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type InsideLeadController struct {
	ser *services.LeadService
}

func NewInsideLead(ser *services.LeadService) *InsideLeadController {
	return &InsideLeadController{
		ser: ser,
	}
}

func (c *InsideLeadController) List(ctx *gin.Context) {
	var param validation.PaginationLead
	if err := ctx.ShouldBindQuery(&param); err != nil {
		out := validation.BeautyMessage(err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": out})
		return
	}

	items, err := c.ser.Pagination(ctx, param)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, items)
	return
}

// Show @Schemes
//
//	@Tags		Lead
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"Find Lead By ID"
//	@Success	200	{object}	entities.Lead
//	@Router		/leads/{id} [get]
func (c *InsideLeadController) Show(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	item, err := c.ser.Repo.BaseRepo.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, item)
}
