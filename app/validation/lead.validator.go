package validation

import (
	"crm-service-go/app/entities"
)

type CreateLead struct {
	FullName   string               `json:"fullName" form:"fullName" binding:"required"`
	Phone      string               `json:"phone" form:"phone" binding:"required"`
	Email      string               `json:"email" form:"email" binding:"required"`
	NationalId string               `json:"nationalId" form:"nationalId" binding:"required"`
	PassportId string               `json:"passportId" form:"passportId" binding:"required"`
	TaxId      string               `json:"taxId" form:"taxId" binding:"required"`
	Address    string               `json:"address" form:"address" binding:"required"`
	Province   string               `json:"province" form:"province" binding:"required"`
	District   string               `json:"district" form:"district" binding:"required"`
	Source     string               `json:"source" form:"source" binding:"required"`
	EmployeeBy string               `json:"employeeBy" form:"employeeBy" binding:"required"`
	StoreCode  string               `json:"storeCode" form:"storeCode" binding:"required"`
	Type       entities.SaleOppType `json:"type" form:"type" binding:"required"`
	Birthday   string               `json:"birthday" form:"birthday" binding:"required"`
	Gender     string               `json:"gender" form:"gender" binding:"required"`
}

type UpdateLead struct {
	FullName   string               `json:"fullName" form:"fullName" binding:"required"`
	Phone      string               `json:"phone" form:"phone" binding:"phoneExists"`
	Email      string               `json:"email" form:"email"`
	NationalId string               `json:"nationalId" form:"nationalId"`
	PassportId string               `json:"passportId" form:"passportId"`
	TaxId      string               `json:"taxId" form:"taxId"`
	Address    string               `json:"address" form:"address"`
	Province   string               `json:"province" form:"province"`
	District   string               `json:"district" form:"district"`
	Source     string               `json:"source" form:"source"`
	EmployeeBy string               `json:"employeeBy" form:"employeeBy"`
	StoreCode  string               `json:"storeCode" form:"storeCode"`
	Type       entities.SaleOppType `json:"type" form:"type"`
	Birthday   string               `json:"birthday" form:"birthday"`
	Gender     string               `json:"gender" form:"gender"`
}

type PaginationLead struct {
	Keyword string `json:"keyword" form:"keyword" binding:"required"`
	Limit   int64  `json:"limit" form:"limit" binding:"gte=10"`
	Skip    int64  `json:"skip" form:"skip" binding:"gte=0"`
}
