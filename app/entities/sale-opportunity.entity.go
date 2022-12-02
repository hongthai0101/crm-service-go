package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const CollectionSaleOpportunities = "SaleOpportunities"

type SaleOppGroup string

const (
	GroupOld SaleOppGroup = "OLD"
	GroupNew SaleOppGroup = "NEW"
)

type SaleOppStatus string

const (
	SaleOppStatusNew           SaleOppStatus = "NEW"
	SaleOppStatusSuccess       SaleOppStatus = "SUCCESS"
	SaleOppStatusPending       SaleOppStatus = "PENDING"
	SaleOppStatusConsulting    SaleOppStatus = "CONSULTING"
	SaleOppStatusDealt         SaleOppStatus = "DEALT"
	SaleOppStatusDenied        SaleOppStatus = "DENIED"
	SaleOppStatusCancel        SaleOppStatus = "CANCEL"
	SaleOppStatusUnContactable SaleOppStatus = "UNCONTACTABLE"
)

type SaleOppType string

const (
	SaleOppTypeBorrower   SaleOppType = "BORROWER"
	SaleOppTypePartner    SaleOppType = "PARTNER"
	SaleOppTypeInvestment SaleOppType = "INVESTMENT"
)

type AssetMedia struct {
	Url      string `bson:"url" json:"url,omitempty"`
	MimeType string `bson:"mimeType" json:"mimeType,omitempty"`
}

type Asset struct {
	Description string       `bson:"description" json:"description"`
	Media       []AssetMedia `bson:"media" json:"media"`
	AssetType   string       `bson:"assetType" json:"assetType"`
	DemandLoan  interface{}  `bson:"demandLoan" json:"demandLoan"`
	LoanTerm    interface{}  `bson:"loanTerm" json:"loanTerm"`
}

type SourceRefs struct {
	Source     string      `bson:"source" json:"source"`
	RefId      string      `bson:"ref_id" json:"refId"`
	CustomerId interface{} `bson:"customerId" json:"customerId"`
}

type SaleOpportunity struct {
	ID              primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	SourceRefs      SourceRefs             `bson:"source_refs" json:"source_refs"`
	Code            string                 `bson:"code" json:"code"`
	Status          SaleOppStatus          `bson:"status" json:"status"`
	Type            SaleOppType            `bson:"type" json:"type"`
	Source          string                 `bson:"source" json:"source"`
	Group           SaleOppGroup           `bson:"group" json:"group"`
	Assets          Asset                  `bson:"assets" json:"assets"`
	EmployeeBy      string                 `bson:"employeeBy" json:"employeeBy"`
	StoreCode       string                 `bson:"storeCode" json:"storeCode"`
	DisbursedAt     *time.Time             `bson:"disbursedAt" json:"disbursedAt"`
	ContractCode    string                 `bson:"contractCode" json:"contractCode"`
	Tags            []string               `bson:"tags" json:"tags"`
	DisbursedAmount int                    `bson:"disbursedAmount" json:"disbursedAmount"`
	LeadId          primitive.ObjectID     `bson:"leadId" json:"leadId,omitempty"`
	Hash            string                 `bson:"hash" json:"hash,omitempty"`
	Metadata        map[string]interface{} `bson:"metadata" json:"metadata,omitempty"`
	BaseEntity      `bson:"inline"`

	// Include data when get sale opportunity
	Lead     *Lead  `bson:"-" json:"lead,omitempty"`
	Created  string `bson:"-" json:"created,omitempty"`
	Updated  string `bson:"-" json:"updated,omitempty"`
	Employee string `bson:"-" json:"employee,omitempty"`
	TagData  []Tag  `bson:"-" json:"tagData,omitempty"`
}
