package validation

type StatisticsCommonRequest struct {
	AssetTypes []string `json:"assetTypes"`
	StoreCodes []string `json:"storeCodes"`
	Sources    []string `json:"sources"`
	FromDate   string   `json:"fromDate" form:"fromDate" binding:"required"`
	ToDate     string   `json:"toDate" form:"toDate" binding:"required"`
}

type StatisticsEmployeeRequest struct {
	EmployeeIds []string `json:"employeeIds" form:"employeeIds" binding:"required"`
	StatisticsCommonRequest
}
