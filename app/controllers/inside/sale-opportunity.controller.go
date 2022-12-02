package controller_inside

import (
	"crm-service-go/app/controllers"
	serviceInside "crm-service-go/app/services/inside"
	"crm-service-go/app/validation"
	"crm-service-go/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
)

type InsideSaleOpportunityController struct {
	ser *serviceInside.InsideSaleOpportunityService
}

func NewInsideSaleOpportunity(ser *serviceInside.InsideSaleOpportunityService) *InsideSaleOpportunityController {
	return &InsideSaleOpportunityController{
		ser: ser,
	}
}

func (c *InsideSaleOpportunityController) List(ctx *gin.Context) {
	var param validation.PaginationSaleOpportunity
	if err := ctx.ShouldBindQuery(&param); err != nil {
		pkg.SendErrorResponse(ctx, pkg.ResponseError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    validation.BeautyMessage(err),
		})
		return
	}

	items, err := c.ser.Pagination(ctx, param)
	if err != nil {
		pkg.SendErrorResponse(ctx, pkg.ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, items)
	return
}

func (c *InsideSaleOpportunityController) Create(ctx *gin.Context) {
	var body validation.CreateSaleOpportunityInside
	if err := ctx.ShouldBindJSON(&body); err != nil {
		pkg.SendErrorResponse(ctx, pkg.ResponseError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    validation.BeautyMessage(err),
		})
		return
	}
	item, err := c.ser.Create(ctx, body)
	if err != nil {
		pkg.SendErrorResponse(ctx, pkg.ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *InsideSaleOpportunityController) Show(ctx *gin.Context) {
	id := controllers.GetObjectIDFromPath(ctx)

	item, err := c.ser.FindById(ctx, *id, ctx.QueryArray("includes"))
	if err != nil {
		pkg.SendErrorResponse(ctx, pkg.ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *InsideSaleOpportunityController) Update(ctx *gin.Context) {
	ctx.JSON(http.StatusNoContent, nil)
}
