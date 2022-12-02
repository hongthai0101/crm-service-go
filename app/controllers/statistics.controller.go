package controllers

import (
	"crm-service-go/app/services"
	"crm-service-go/app/validation"
	"crm-service-go/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
)

type StatisticsController struct {
	ser *services.StatisticsService
}

func NewStatistics(service *services.StatisticsService) *StatisticsController {
	return &StatisticsController{
		ser: service,
	}
}

// Index @Schemes
//
//	@Tags		Statistics
//	@Accept		json
//	@Produce	json
//	@Param		StatisticsIndex	query		validation.StatisticsCommonRequest	false	"StatisticsIndex"
//	@Success	200				{object}	services.StatisticsIndexResponse
//	@Router		/statistics/index [get]
func (c *StatisticsController) Index(ctx *gin.Context) {
	var params validation.StatisticsCommonRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		pkg.SendErrorResponse(ctx, pkg.ResponseError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    validation.BeautyMessage(err),
		})
		return
	}

	res := c.ser.StatisticsIndex(ctx, &params)
	ctx.JSON(http.StatusOK, res)
	return
}

// Source @Schemes
//
//	@Tags		Statistics
//	@Accept		json
//	@Produce	json
//	@Param		StatisticsSource	query	validation.StatisticsCommonRequest	false	"StatisticsSource"
//	@Success	200					{array}	services.StatisticsCommonResponse
//	@Router		/statistics/sources [get]
func (c *StatisticsController) Source(ctx *gin.Context) {
	var params validation.StatisticsCommonRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		pkg.SendErrorResponse(ctx, pkg.ResponseError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    validation.BeautyMessage(err),
		})
		return
	}

	res, _ := c.ser.StatisticsSources(ctx, &params)
	ctx.JSON(http.StatusOK, res)
	return
}

// Store @Schemes
//
//	@Tags		Statistics
//	@Accept		json
//	@Produce	json
//	@Param		StatisticsStore	query	validation.StatisticsCommonRequest	false	"StatisticsStore"
//	@Success	200					{array}	services.StatisticsCommonResponse
//	@Router		/statistics/stores [get]
func (c *StatisticsController) Store(ctx *gin.Context) {
	var params validation.StatisticsCommonRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		pkg.SendErrorResponse(ctx, pkg.ResponseError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    validation.BeautyMessage(err),
		})
		return
	}

	res, _ := c.ser.StatisticsStores(ctx, &params)
	ctx.JSON(http.StatusOK, res)
	return
}

// Employee @Schemes
//
//	@Tags		Statistics
//	@Accept		json
//	@Produce	json
//	@Param		StatisticsSource	query	validation.StatisticsEmployeeRequest	false	"StatisticsSource"
//	@Success	200					{array}	services.StatisticsEmployeeResponse
//	@Router		/statistics/employees [get]
func (c *StatisticsController) Employee(ctx *gin.Context) {
	var params validation.StatisticsEmployeeRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		pkg.SendErrorResponse(ctx, pkg.ResponseError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    validation.BeautyMessage(err),
		})
		return
	}

	res, _ := c.ser.StatisticsEmployees(ctx, &params)
	ctx.JSON(http.StatusOK, res)
	return
}
