package services

import (
	"crm-service-go/app/clients"
	"crm-service-go/app/entities"
	"crm-service-go/app/repositories"
	"crm-service-go/app/validation"
	"crm-service-go/pkg"
	"crm-service-go/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

type PanicResponse struct {
	Err        error
	StatusCode int
}

type StatisticsService struct {
	saleRepo         *repositories.SaleOpportunityRepository
	employeeClient   *clients.EmployeeClient
	masterDataClient *clients.MasterDataClient
}

func NewStatisticsService(
	saleRepo *repositories.SaleOpportunityRepository,
	employeeClient *clients.EmployeeClient,
	masterDataClient *clients.MasterDataClient,
) *StatisticsService {
	return &StatisticsService{
		saleRepo:         saleRepo,
		employeeClient:   employeeClient,
		masterDataClient: masterDataClient,
	}
}

type StatisticsIndexResponse struct {
	Total           int `json:"total"`
	Processing      int `json:"processing"`
	Disbursal       int `json:"disbursal"`
	DisbursalAmount int `json:"disbursalAmount"`
	Denied          int `json:"denied"`
}

func (s *StatisticsService) StatisticsIndex(
	ctx *gin.Context,
	params *validation.StatisticsCommonRequest,
) *StatisticsIndexResponse {
	filter, findOptions := handleFilter(ctx, params, nil)
	sales, err := s.saleRepo.BaseRepo.Find(filter, findOptions)

	if err != nil {
		return nil
	}
	total := len(sales)

	processing, disbursal, denied, disbursalAmount := 0, 0, 0, 0
	for _, sale := range sales {
		status := sale.Status
		processingStatus := []string{string(entities.SaleOppStatusPending), string(entities.SaleOppStatusConsulting), string(entities.SaleOppStatusDealt)}
		if utils.Contains(processingStatus, string(status)) {
			processing++
		}

		if status == entities.SaleOppStatusSuccess && sale.ContractCode != "" {
			disbursal++
			disbursalAmount += sale.DisbursedAmount
		}

		if status == entities.SaleOppStatusDenied {
			denied++
		}
	}

	return &StatisticsIndexResponse{
		Total:           total,
		Processing:      processing,
		Disbursal:       disbursal,
		DisbursalAmount: disbursalAmount,
		Denied:          denied,
	}
}

type StatisticsDetailDataResponse struct {
	Amount int                   `json:"amount"`
	Group  entities.SaleOppGroup `json:"group"`
	Order  int                   `json:"order"`
}

type StatisticsDetailResponse struct {
	Order  int                             `json:"order"`
	Status entities.SaleOppStatus          `json:"status"`
	Data   []*StatisticsDetailDataResponse `json:"data"`
}

type StatisticsCommonResponse struct {
	Code    string                     `json:"code"`
	Details []StatisticsDetailResponse `json:"details"`
}

func (s *StatisticsService) StatisticsSources(
	ctx *gin.Context,
	params *validation.StatisticsCommonRequest,
) ([]*StatisticsCommonResponse, error) {
	s.masterDataClient.Client.SetToken(ctx.Request.Header.Get("Authorization"))

	masterDataPromise := func() <-chan *StatisticsMasterData {
		r := make(chan *StatisticsMasterData)
		go func() {
			defer close(r)
			result := s.findMasterData()
			r <- result
		}()
		return r
	}()

	salesPromise := func() <-chan []*entities.SaleOpportunity {
		r := make(chan []*entities.SaleOpportunity)
		go func() {
			defer close(r)
			filter, findOptions := handleFilter(ctx, params, nil)
			result, err := s.saleRepo.BaseRepo.Find(filter, findOptions)
			if err == nil {
				r <- result
			}
		}()
		return r
	}()

	masterData, sales := <-masterDataPromise, <-salesPromise
	response := make([]*StatisticsCommonResponse, 0)
	for key, _ := range *masterData.sources {
		itemsBySource := make([]*entities.SaleOpportunity, 0)
		for _, item := range sales {
			if item.Source == key {
				itemsBySource = append(itemsBySource, item)
			}
		}

		details := statusResponseWithWorker(itemsBySource, masterData)
		response = append(response, &StatisticsCommonResponse{
			Code:    key,
			Details: details,
		})
	}

	return response, nil
}

func (s *StatisticsService) StatisticsStores(
	ctx *gin.Context,
	params *validation.StatisticsCommonRequest,
) ([]*StatisticsCommonResponse, error) {
	s.masterDataClient.Client.SetToken(ctx.Request.Header.Get("Authorization"))

	storeCodes := params.StoreCodes
	if len(storeCodes) == 0 {
		stores := s.masterDataClient.GetStores()
		for storeCode, _ := range *stores {
			storeCodes = append(storeCodes, storeCode)
		}
	}
	if len(storeCodes) == 0 {
		return nil, nil
	}

	masterDataPromise := func() <-chan *StatisticsMasterData {
		r := make(chan *StatisticsMasterData)
		go func() {
			defer close(r)
			result := s.findMasterData()
			r <- result
		}()
		return r
	}()

	salesPromise := func() <-chan []*entities.SaleOpportunity {
		r := make(chan []*entities.SaleOpportunity)
		go func() {
			defer close(r)
			params.StoreCodes = storeCodes
			filter, findOptions := handleFilter(ctx, params, nil)
			result, err := s.saleRepo.BaseRepo.Find(filter, findOptions)
			if err == nil {
				r <- result
			}
		}()
		return r
	}()

	masterData, sales := <-masterDataPromise, <-salesPromise

	response := make([]*StatisticsCommonResponse, 0)
	for _, code := range storeCodes {
		itemsByStoreCodes := make([]*entities.SaleOpportunity, 0)
		for _, item := range sales {
			if item.StoreCode == code {
				itemsByStoreCodes = append(itemsByStoreCodes, item)
			}
		}

		details := statusResponseWithoutWorker(sales, masterData)
		response = append(response, &StatisticsCommonResponse{
			Code:    code,
			Details: details,
		})
	}

	return response, nil
}

type StatisticsEmployeeResponse struct {
	Employee struct {
		Id          string `json:"id"`
		DisplayName string `json:"displayName"`
	} `json:"employee"`
	Details []StatisticsDetailResponse `json:"details"`
}

func (s *StatisticsService) StatisticsEmployees(
	ctx *gin.Context,
	params *validation.StatisticsEmployeeRequest,
) ([]*StatisticsEmployeeResponse, error) {
	s.masterDataClient.Client.SetToken(ctx.Request.Header.Get("Authorization"))
	s.employeeClient.Client.SetToken(ctx.Request.Header.Get("Authorization"))

	masterDataPromise := func() <-chan *StatisticsMasterData {
		r := make(chan *StatisticsMasterData)
		go func() {
			defer close(r)
			result := s.findMasterData()
			r <- result
		}()
		return r
	}()

	employeeIds := params.EmployeeIds
	employeesPromise := func() <-chan *map[string]string {
		r := make(chan *map[string]string)
		go func() {
			defer close(r)

			result, err := s.employeeClient.GetEmployees(employeeIds)
			if err == nil {
				r <- result
			}
		}()
		return r
	}()

	salesPromise := func() <-chan []*entities.SaleOpportunity {
		r := make(chan []*entities.SaleOpportunity)
		go func() {
			defer close(r)

			p, _ := utils.TypeConverter[validation.StatisticsCommonRequest](params)
			fl, fo := handleFilter(ctx, p, employeeIds)
			result, err := s.saleRepo.BaseRepo.Find(fl, fo)
			if err == nil {
				r <- result
			}
		}()
		return r
	}()
	masterData, employees, sales := <-masterDataPromise, <-employeesPromise, <-salesPromise
	if masterData != nil && employees != nil && sales != nil {
		response := make([]*StatisticsEmployeeResponse, 0)
		for key, displayName := range *employees {
			details := statusResponseWithoutWorker(sales, masterData)
			response = append(response, &StatisticsEmployeeResponse{
				Employee: struct {
					Id          string `json:"id"`
					DisplayName string `json:"displayName"`
				}{Id: key, DisplayName: displayName},
				Details: details,
			})
		}

		return response, nil
	}

	return nil, nil
}

func handleFilter(
	ctx *gin.Context,
	params *validation.StatisticsCommonRequest,
	employeeIds []string,
) (bson.D, *options.FindOptions) {

	start, _ := time.Parse(pkg.DDMMYYYY, params.FromDate)
	end, _ := time.Parse(pkg.DDMMYYYY, params.ToDate)
	filter := bson.D{
		{
			"createdAt", bson.M{
			"$gte": primitive.NewDateTimeFromTime(start),
			"$lt":  primitive.NewDateTimeFromTime(end),
		},
		},
		{
			"deletedAt", nil,
		},
	}

	if len(params.AssetTypes) > 0 {
		filter = append(filter, bson.E{
			Key: "assetType", Value: bson.D{{
				"$in", params.AssetTypes,
			}},
		})
	}

	if len(params.Sources) > 0 {
		filter = append(filter, bson.E{
			Key: "source", Value: bson.D{{
				"$in", params.Sources,
			}},
		})
	}

	if len(params.StoreCodes) > 0 {
		filter = append(filter, bson.E{
			Key: "storeCode", Value: bson.D{{
				"$in", params.StoreCodes,
			}},
		})
	}

	if len(employeeIds) > 0 {
		filter = append(filter, bson.E{
			Key: "employeeId", Value: bson.D{{
				"$in", employeeIds,
			}},
		})
	}

	filter = mapPolicyToFilter(ctx, filter)

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"createdAt", -1}})
	findOptions.Projection = bson.D{
		{"status", 1},
		{"contractCode", 1},
		{"disbursedAmount", 1},
		{"assets", 1},
		{"storeCode", 1},
		{"employeeBy", 1},
		{"source", 1},
		{"group", 1},
	}

	return filter, findOptions
}

type StatisticsMasterData struct {
	sources  *map[string]string
	statuses *map[string]string
	groups   *map[string]string
}

func (s *StatisticsService) findMasterData() *StatisticsMasterData {
	wg := sync.WaitGroup{}
	wg.Add(3)

	var sources, statuses, groups *map[string]string
	go func() {
		defer wg.Done()
		sources = s.masterDataClient.GetSource()
	}()
	go func() {
		defer wg.Done()
		statuses = s.masterDataClient.GetStatuses()
	}()
	go func() {
		defer wg.Done()
		groups = s.masterDataClient.GetGroups()
	}()
	wg.Wait()

	return &StatisticsMasterData{
		sources:  sources,
		statuses: statuses,
		groups:   groups,
	}
}

func statusResponseWithoutWorker(
	items []*entities.SaleOpportunity,
	masterData *StatisticsMasterData,
) []StatisticsDetailResponse {
	response := make([]StatisticsDetailResponse, 0)
	for status, _ := range *masterData.statuses {
		itemByStatus := make([]*entities.SaleOpportunity, 0)
		for _, item := range items {
			if string(item.Status) == status {
				itemByStatus = append(itemByStatus, item)
			}
		}

		data := make([]*StatisticsDetailDataResponse, 0)

		for group, _ := range *masterData.groups {
			amount := 0
			order := 0
			for _, item := range itemByStatus {
				if string(item.Group) == group {
					order++
					if status == string(entities.SaleOppStatusSuccess) {
						amount += item.DisbursedAmount
					} else if demandLoan, ok := item.Assets.DemandLoan.(int); ok {
						amount += demandLoan
					}
				}
			}

			data = append(data, &StatisticsDetailDataResponse{
				Amount: amount,
				Group:  entities.SaleOppGroup(group),
				Order:  order,
			})
		}

		response = append(response, StatisticsDetailResponse{
			Order:  len(itemByStatus),
			Status: entities.SaleOppStatus(status),
			Data:   data,
		})
	}
	return response
}

type JobResult struct {
	data   []*StatisticsDetailDataResponse
	order  int
	status entities.SaleOppStatus
}

type JobStatistics struct {
	status     entities.SaleOppStatus
	sales      []*entities.SaleOpportunity
	masterData *StatisticsMasterData
}

func statusResponseWithWorker(
	items []*entities.SaleOpportunity,
	masterData *StatisticsMasterData,
) []StatisticsDetailResponse {
	numberWorker := 4
	jobs, jobResult :=
		make(chan *JobStatistics, 200), make(chan *JobResult, 200)

	for i := 1; i <= numberWorker; i++ {
		go worker(jobs, jobResult, fmt.Sprintf("%d", i))
	}

	response := make([]StatisticsDetailResponse, 0)
	for status, _ := range *masterData.statuses {
		go func(status entities.SaleOppStatus) {
			jobs <- &JobStatistics{
				status:     status,
				sales:      items,
				masterData: masterData,
			}
		}(entities.SaleOppStatus(status))
	}
	defer close(jobs)

	for i := 0; i < len(*masterData.statuses); i++ {
		result := <-jobResult
		response = append(response, StatisticsDetailResponse{
			Order:  result.order,
			Status: result.status,
			Data:   result.data,
		})
	}
	return response
}

func worker(
	jobs <-chan *JobStatistics,
	results chan<- *JobResult,
	name string,
) {
	wg := sync.WaitGroup{}
	for job := range jobs {
		wg.Add(1)
		go func(job *JobStatistics) {
			defer wg.Done()

			data, order := handleData(job.status, job.sales, job.masterData)
			results <- &JobResult{
				data:   data,
				order:  order,
				status: job.status,
			}

			fmt.Printf("Worker %s is handle status code %v\n", name, job.status)
		}(job)
	}
	wg.Wait()
}

func handleData(
	status entities.SaleOppStatus,
	sales []*entities.SaleOpportunity,
	masterData *StatisticsMasterData,
) ([]*StatisticsDetailDataResponse, int) {
	itemsByStatus := make([]*entities.SaleOpportunity, 0)
	for _, item := range sales {
		if item.Status == status {
			itemsByStatus = append(itemsByStatus, item)
		}
	}

	data := make([]*StatisticsDetailDataResponse, 0)
	for group, _ := range *masterData.groups {
		parseGroup := entities.SaleOppGroup(group)
		amount := 0
		order := 0
		for _, item := range itemsByStatus {
			if item.Group == parseGroup {
				order++
				if status == entities.SaleOppStatusSuccess {
					amount += item.DisbursedAmount
				} else if demandLoan, ok := item.Assets.DemandLoan.(int); ok {
					amount += demandLoan
				}
			}
		}

		data = append(data, &StatisticsDetailDataResponse{
			Amount: amount,
			Group:  parseGroup,
			Order:  order,
		})
	}
	return data, len(itemsByStatus)
}
