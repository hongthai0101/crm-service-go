package controllers

import (
	"crm-service-go/app/services"
	"crm-service-go/app/validation"
	"crm-service-go/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type NoteController struct {
	ser *services.NoteService
}

func NewNote(ser *services.NoteService) *NoteController {
	return &NoteController{
		ser: ser,
	}
}

// List @Schemes
//
//	@Tags		Note
//	@Accept		json
//	@Produce	json
//	@Param		paginationNote	query		validation.PaginationNote	false	"paginationNote"
//	@Success	200				{object}	services.NotePagination
//	@Router		/notes [get]
func (c *NoteController) List(ctx *gin.Context) {
	var param validation.PaginationNote
	if err := ctx.ShouldBindQuery(&param); err != nil {
		utils.Debug(err)
		out := validation.BeautyMessage(err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": out})
		return
	}

	items, err := c.ser.Pagination(param)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, items)
	return
}

// Create @Schemes
//
//	@Tags		Note
//	@Accept		json
//	@Produce	json
//	@Param		Note	body		validation.CreateNote	true	"Create Note"
//	@Success	200		{object}	entities.Note
//	@Router		/notes [post]
func (c *NoteController) Create(ctx *gin.Context) {
	var body validation.CreateNote
	if err := ctx.ShouldBindJSON(&body); err != nil {
		out := validation.BeautyMessage(err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": out})
		return
	}
	item, err := c.ser.Create(body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, item)
}

// Show @Schemes
//
//	@Tags		Note
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"Find Note By Id"
//	@Success	200	{object}	entities.Note
//	@Router		/notes/{id} [get]
func (c *NoteController) Show(ctx *gin.Context) {
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

// Update @Schemes
//
//	@Tags		Note
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string					true	"Update By ID"
//	@Param		Note	body		validation.UpdateNote	true	"Update Note"
//	@Success	200		{object}	entities.Note
//	@Router		/notes/{id} [put]
func (c *NoteController) Update(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var body validation.UpdateNote
	if err = ctx.ShouldBindJSON(&body); err != nil {
		out := validation.BeautyMessage(err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": out})
		return
	}

	item, err := c.ser.Update(ctx, id, body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, item)
}

// Delete @Schemes
//
//	@Tags		Note
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"Delete By ID"
//	@Success	204
//	@Router		/notes/{id} [delete]
func (c *NoteController) Delete(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	_, err = c.ser.Repo.BaseRepo.UpdateByID(id, bson.M{"deletedAt": time.Now()})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
