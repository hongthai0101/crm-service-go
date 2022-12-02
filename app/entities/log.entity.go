package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const CollectionLog = "Log"

type BeforeAttributes struct {
}

type AfterAttributes struct {
}

type Log struct {
	ID                primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	BeforeAttributes  map[string]interface{} `bson:"beforeAttributes" json:"beforeAttributes"`
	AfterAttributes   map[string]interface{} `bson:"afterAttributes" json:"afterAttributes"`
	SaleOpportunityId primitive.ObjectID     `bson:"saleOpportunityId" json:"saleOpportunityId,omitempty"`
	CreatedBy         string                 `bson:"createdBy" json:"createdBy"`
	CreatedAt         time.Time              `bson:"createdAt" json:"createdAt"`
}
