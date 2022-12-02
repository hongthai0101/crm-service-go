package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionNote = "Note"

type Note struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Content           string             `bson:"content" json:"content,omitempty"`
	SaleOpportunityId primitive.ObjectID `bson:"saleOpportunityId" json:"saleOpportunityId,omitempty"`
	BaseEntity        `bson:"inline"`
}
