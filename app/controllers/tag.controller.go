package controllers

import (
	"crm-service-go/app/entities"
	"crm-service-go/app/repositories"
	"crm-service-go/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

type TagController struct {
	repo *repositories.TagRepository
}

func NewTag(repo *repositories.TagRepository) *TagController {
	return &TagController{
		repo: repo,
	}
}

// All @Schemes
//
//	@Tags		Tag
//	@Accept		json
//	@Produce	json
//	@Success	200	{array}	entities.Tag
//	@Router		/tags [get]
func (c *TagController) All(ctx *gin.Context) {
	items, err := c.repo.BaseRepo.Find(bson.M{}, nil)
	if err != nil {
		utils.Debug(err)
		items = []*entities.Tag{}
	}
	ctx.JSON(http.StatusOK, items)
}
