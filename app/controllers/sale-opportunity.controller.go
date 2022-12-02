package controllers

import (
	"crm-service-go/app/services"
	"crm-service-go/app/validation"
	"crm-service-go/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SaleOpportunityController struct {
	ser *services.SaleOpportunityService
}

func NewSaleOpportunity(ser *services.SaleOpportunityService) *SaleOpportunityController {
	return &SaleOpportunityController{
		ser: ser,
	}
}

// List @Schemes
//
//	@Tags		SaleOpportunity
//	@Accept		json
//	@Produce	json
//	@Param		paginationSaleOpportunity	query		validation.PaginationSaleOpportunity	false	"paginationSaleOpportunity"
//	@Success	200							{object}	services.SaleOpportunityPagination
//	@Router		/sale-opportunities [get]
func (c *SaleOpportunityController) List(ctx *gin.Context) {
	var param validation.PaginationSaleOpportunity
	if err := ctx.ShouldBindQuery(&param); err != nil {
		pkg.SendErrorResponse(ctx, pkg.ResponseError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    validation.BeautyMessage(err),
		})
		return
	}

	items, err := c.ser.Pagination(ctx, param, true)
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

// Create @Schemes
//
//	@Tags		SaleOpportunity
//	@Accept		json
//	@Produce	json
//	@Param		SaleOpportunity	body		validation.CreateSaleOpportunity	true	"Create SaleOpportunity"
//	@Success	200				{object}	entities.SaleOpportunity
//	@Router		/sale-opportunities [post]
func (c *SaleOpportunityController) Create(ctx *gin.Context) {
	var body validation.CreateSaleOpportunity
	if err := ctx.ShouldBindJSON(&body); err != nil {
		pkg.SendErrorResponse(ctx, pkg.ResponseError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    validation.BeautyMessage(err),
		})
		return
	}
	item, err := c.ser.Create(ctx, &body, true)
	if err != nil {
		pkg.SendErrorResponse(ctx, pkg.ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, item)
}

// Show @Schemes
//
//	@Tags		SaleOpportunity
//	@Accept		json
//	@Produce	json
//	@Param		id			path		string	true	"Find SaleOpportunity By ID"
//	@Param		includes	query		array	false	"includes"
//	@Success	200			{object}	entities.SaleOpportunity
//	@Router		/sale-opportunities/{id} [get]
func (c *SaleOpportunityController) Show(ctx *gin.Context) {
	id := GetObjectIDFromPath(ctx)

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

// Update @Schemes
//
//	@Tags		SaleOpportunity
//	@Accept		json
//	@Produce	json
//	@Param		id				path		string								true	"Update By ID"
//	@Param		SaleOpportunity	body		validation.UpdateSaleOpportunity	true	"Update SaleOpportunity"
//	@Success	200				{object}	entities.SaleOpportunity
//	@Router		/sale-opportunities/{id} [put]
func (c *SaleOpportunityController) Update(ctx *gin.Context) {
	id := GetObjectIDFromPath(ctx)

	var body validation.UpdateSaleOpportunity
	if err := ctx.ShouldBindJSON(&body); err != nil {
		pkg.SendErrorResponse(ctx, pkg.ResponseError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    validation.BeautyMessage(err),
		})
		return
	}

	item, err := c.ser.Update(ctx, *id, body)
	if err != nil {
		pkg.SendErrorResponse(ctx, pkg.ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, item)
}

// Delete @Schemes
//
//	@Tags		SaleOpportunity
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"Delete By ID"
//	@Success	204
//	@Router		/sale-opportunities/{id} [delete]
func (c *SaleOpportunityController) Delete(ctx *gin.Context) {
	id := GetObjectIDFromPath(ctx)

	_, err := c.ser.Delete(ctx, *id)
	if err != nil {
		pkg.SendErrorResponse(ctx, pkg.ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

// Logs @Schemes
//
//	@Tags		SaleOpportunity
//	@Accept		json
//	@Produce	json
//	@Param		id								path		string										true	"Get Logs By SaleOpportunity ID"
//	@Param		paginationSaleOpportunityLog	query		validation.PaginationSaleOpportunityLogs	false	"paginationSaleOpportunity"
//	@Success	200								{object}	services.SaleOpportunityLogsPagination
//	@Router		/sale-opportunities/{id}/logs [get]
func (c *SaleOpportunityController) Logs(ctx *gin.Context) {
	id := GetObjectIDFromPath(ctx)

	var param validation.PaginationSaleOpportunityLogs
	if err := ctx.ShouldBindQuery(&param); err != nil {
		pkg.SendErrorResponse(ctx, pkg.ResponseError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    validation.BeautyMessage(err),
		})
		return
	}

	res, _ := c.ser.PaginationLogs(ctx, *id, param)
	ctx.JSON(http.StatusOK, res)
}
