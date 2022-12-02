package entities

import (
	"crm-service-go/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionLead = "Lead"

type Lead struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CustomerId string             `bson:"customerId" json:"customerId"`
	FullName   string             `bson:"fullName" json:"fullName"`
	Phone      string             `bson:"phone" json:"phone"`
	Email      string             `bson:"email" json:"email"`
	NationalId string             `bson:"nationalId" json:"nationalId"`
	PassportId string             `bson:"passportId" json:"passportId"`
	TaxId      string             `bson:"taxId" json:"taxId"`
	Address    string             `bson:"address" json:"address"`
	Province   string             `bson:"province" json:"province"`
	District   string             `bson:"district" json:"district"`
	Source     string             `bson:"source" json:"source"`
	EmployeeBy string             `bson:"employeeBy" json:"employeeBy"`
	StoreCode  string             `bson:"storeCode" json:"storeCode"`
	Type       SaleOppType        `bson:"type" json:"type"`
	Birthday   string             `bson:"birthday" json:"birthday"`
	Gender     string             `bson:"gender" json:"gender"`
	Metadata   interface{}        `bson:"metadata" json:"metadata"`
	BaseEntity `bson:"inline"`
}

func (e *Lead) DefaultValue(lead *Lead) {
	if lead.CreatedBy != "" {
		lead.CreatedBy = config.GetConfig().DefaultDataConfig.CreatedBy
	}
}
