package validation

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateNote struct {
	Content           string             `json:"content" form:"content" binding:"required"`
	SaleOpportunityId primitive.ObjectID `json:"saleOpportunityId" form:"saleOpportunityId" binding:"required"`
}

type UpdateNote struct {
	Content string `json:"content" form:"content"`
}

type PaginationNote struct {
	SaleOpportunityId string `json:"saleOpportunityId" form:"saleOpportunityId" binding:"required,len=24"`
	Limit             int64  `json:"limit" form:"limit" binding:"numeric,gte=10"`
	Skip              int64  `json:"skip" form:"skip" binding:"numeric,gte=0"`
}
