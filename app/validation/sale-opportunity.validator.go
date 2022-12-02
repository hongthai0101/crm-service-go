package validation

import (
	"crm-service-go/app/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FilterDateType string

type CreateSaleOpportunity struct {
	Status     entities.SaleOppStatus `json:"status" binding:"required"`
	Type       entities.SaleOppType   `json:"type" binding:"required"`
	Source     string                 `json:"source" binding:"required"`
	Group      entities.SaleOppGroup  `json:"group"`
	Assets     entities.Asset         `json:"assets" binding:"required"`
	EmployeeBy string                 `json:"employeeBy"`
	StoreCode  string                 `json:"storeCode"`
	Tags       []string               `json:"tags"`
	LeadId     primitive.ObjectID     `bson:"leadId" json:"leadId" binding:"required,len=24"`
	CodePrefix string                 `json:"codePrefix"`
	Note       string                 `json:"note"`
}

type UpdateSaleOpportunity struct {
	Status          entities.SaleOppStatus `json:"status"`
	Type            entities.SaleOppType   `json:"type"`
	Source          string                 `json:"source"`
	Group           entities.SaleOppGroup  `json:"group"`
	Assets          entities.Asset         `json:"assets"`
	EmployeeBy      string                 `json:"employeeBy"`
	StoreCode       string                 `json:"storeCode"`
	ContractCode    string                 `json:"contractCode"`
	DisbursedAt     string                 `json:"-"`
	DisbursedAmount uint16                 `json:"-"`
	Tags            []string               `json:"tags"`
	LeadId          primitive.ObjectID     `bson:"leadId" json:"leadId" binding:"omitempty,len=24"`
}

type PaginationSaleOpportunity struct {
	Code           string         `json:"code" form:"code"`
	Lead           string         `json:"lead" form:"lead"`
	Statuses       string         `json:"statuses" form:"statuses"`
	Sources        string         `json:"sources" form:"sources"`
	FromDate       string         `json:"fromDate" form:"fromDate"`
	ToDate         string         `json:"toDate" form:"toDate"`
	StoreCodes     string         `json:"storeCodes" form:"storeCodes"`
	FilterDateType FilterDateType `json:"filterDateType" form:"filterDateType"`
	EmployeeBys    string         `json:"employeeBys" form:"employeeBys"`
	Groups         string         `json:"groups" form:"groups"`
	OnlyMe         bool           `json:"onlyMe" form:"onlyMe" binding:"boolean"`
	Includes       []string       `json:"includes" form:"includes"`
	Limit          int64          `json:"limit" form:"limit" binding:"gte=10"`
	Skip           int64          `json:"skip" form:"skip" binding:"gte=0"`
}

type PaginationSaleOpportunityLogs struct {
	Limit int64 `json:"limit" form:"limit" binding:"gte=10"`
	Skip  int64 `json:"skip" form:"skip" binding:"gte=0"`
}

// SaleOpportunityInside Inside
type SaleOpportunityInside struct {
	SourceRefs entities.SourceRefs    `json:"source_refs"`
	Status     entities.SaleOppStatus `json:"status"`
	Type       entities.SaleOppType   `json:"type"`
	Source     string                 `json:"source"`
	CodePrefix string                 `json:"codePrefix"`
	Assets     entities.Asset         `json:"assets" binding:"required"`
	Metadata   interface{}            `json:"metadata"`
}

// LeadInside Inside
type LeadInside struct {
	Provider     string      `json:"provider" form:"provider"`
	Source       string      `json:"source" form:"source"`
	Metadata     interface{} `json:"metadata" form:"metadata"`
	Phone        string      `json:"phone" form:"phone" binding:"required"`
	FullName     string      `json:"fullName" form:"fullName"`
	Email        string      `json:"email" form:"email"`
	CustomerId   string      `json:"customerId" form:"customerId"`
	PersonalId   string      `json:"personalId" form:"personalId"`
	District     string      `json:"district" form:"district"`
	Province     string      `json:"province" form:"province"`
	Gender       string      `json:"gender" form:"gender"`
	IsIdentified bool        `json:"isIdentified" form:"isIdentified"`
	Address      string      `json:"address" form:"address"`
}

type CreateSaleOpportunityInside struct {
	SaleOpportunity SaleOpportunityInside `json:"saleOpportunity" binding:"required"`
	Customer        LeadInside            `json:"customer" binding:"required"`
}
