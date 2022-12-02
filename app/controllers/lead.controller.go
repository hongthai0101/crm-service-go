package controllers

import (
	"crm-service-go/app/services"
	"crm-service-go/app/validation"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

type LeadController struct {
	ser *services.LeadService
}

func NewLead(ser *services.LeadService) *LeadController {
	return &LeadController{
		ser: ser,
	}
}

// List @Schemes
//
//	@Tags		Lead
//	@Accept		json
//	@Produce	json
//	@Param		paginationLead	query		validation.PaginationLead	false	"paginationLead"
//	@Success	200				{object}	services.LeadPagination
//	@Router		/leads [get]
func (c *LeadController) List(ctx *gin.Context) {
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

// Create @Schemes
//
//	@Tags		Lead
//	@Accept		json
//	@Produce	json
//	@Param		lead	body		validation.CreateLead	true	"Create Lead"
//	@Success	200		{object}	entities.Lead
//	@Router		/leads [post]
func (c *LeadController) Create(ctx *gin.Context) {
	var body validation.CreateLead
	if err := ctx.ShouldBindJSON(&body); err != nil {
		out := validation.BeautyMessage(err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": out})
		return
	}
	item, err := c.ser.Create(ctx, body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, item)
}

// Show @Schemes
//
//	@Tags		Lead
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"Find Lead By ID"
//	@Success	200	{object}	entities.Lead
//	@Router		/leads/{id} [get]
func (c *LeadController) Show(ctx *gin.Context) {
	id := GetObjectIDFromPath(ctx)
	item, err := c.ser.Repo.BaseRepo.FindById(*id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, item)
}

// Update @Schemes
//
//	@Tags		Lead
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string					true	"Update By ID"
//	@Param		lead	body		validation.UpdateLead	true	"Update Lead"
//	@Success	200		{object}	entities.Lead
//	@Router		/leads/{id} [put]
func (c *LeadController) Update(ctx *gin.Context) {
	id := GetObjectIDFromPath(ctx)

	// add custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("phoneExists", validation.LeadPhoneExists)
	}

	var body validation.UpdateLead
	if err := ctx.ShouldBindJSON(&body); err != nil {
		out := validation.BeautyMessage(err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": out})
		return
	}

	item, err := c.ser.Update(ctx, *id, body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, item)
}

// Delete @Schemes
//
//	@Tags		Lead
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"Delete By ID"
//	@Success	204
//	@Router		/leads/{id} [delete]
func (c *LeadController) Delete(ctx *gin.Context) {
	id := GetObjectIDFromPath(ctx)

	_, err := c.ser.Repo.BaseRepo.UpdateByID(*id, bson.M{"deletedAt": time.Now()})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
