package services

import (
	"crm-service-go/app/clients"
	"crm-service-go/app/entities"
	"crm-service-go/app/middlewares"
	"crm-service-go/app/repositories"
	"crm-service-go/app/validation"
	"crm-service-go/config"
	"crm-service-go/pkg"
	"crm-service-go/pkg/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"
)

type SaleOpportunityPagination pkg.Pagination[entities.SaleOpportunity]
type SaleOpportunityLogsPagination pkg.Pagination[entities.Log]

type IncludeDataSaleOpportunity struct {
	Lead     map[primitive.ObjectID]*entities.Lead
	Employee *map[string]string
	Tags     map[string]*entities.Tag
}

type SaleOpportunityService struct {
	Repo           *repositories.SaleOpportunityRepository
	tagRepo        *repositories.TagRepository
	logRepo        *repositories.LogRepository
	noteRepo       *repositories.NoteRepository
	LeadService    *LeadService
	topicService   *TopicService
	employeeClient *clients.EmployeeClient
	digitalClient  *clients.DigitalClient
	contractClient *clients.ContractClient
}

var currentTime = time.Now()

func NewSaleOpportunityService(
	repo *repositories.SaleOpportunityRepository,
	tagRepo *repositories.TagRepository,
	logRepo *repositories.LogRepository,
	noteRepo *repositories.NoteRepository,
	leadService *LeadService,
	topicService *TopicService,
	employeeClient *clients.EmployeeClient,
	digitalClient *clients.DigitalClient,
	contractClient *clients.ContractClient,
) *SaleOpportunityService {
	return &SaleOpportunityService{
		Repo:           repo,
		tagRepo:        tagRepo,
		logRepo:        logRepo,
		noteRepo:       noteRepo,
		LeadService:    leadService,
		topicService:   topicService,
		employeeClient: employeeClient,
		digitalClient:  digitalClient,
		contractClient: contractClient,
	}
}

func (s *SaleOpportunityService) Pagination(
	ctx *gin.Context,
	params validation.PaginationSaleOpportunity,
	isLogged bool,
) (*SaleOpportunityPagination, error) {
	user := middlewares.LoggedUser(ctx)
	if isLogged && user == nil {
		return nil, errors.New("user not found")
	}

	s.employeeClient.Client.SetToken(ctx.GetHeader("Authorization"))

	wg := sync.WaitGroup{}
	sales := make(chan []*entities.SaleOpportunity, 1)
	filterCh := make(chan bson.D, 1)
	var (
		total int64
		items []*entities.SaleOpportunity
	)
	limit := pkg.GetLimit(params.Limit)
	skip := pkg.GetSkip(params.Skip)

	wg.Add(3)
	go func(sales chan<- []*entities.SaleOpportunity, filterCh chan<- bson.D) {
		code, lead, statuses, sources, fromDate, toDate, storeCodes, filterDateType, employeeBys, groups, onlyMe :=
			params.Code, params.Lead, params.Statuses, params.Sources, params.FromDate, params.ToDate,
			params.StoreCodes, params.FilterDateType, params.EmployeeBys, params.Groups, params.OnlyMe

		filter := bson.D{
			{
				"deletedAt", nil,
			},
		}
		if code != "" {
			codeQuery := primitive.Regex{Pattern: code, Options: "i"}
			filter = append(filter,
				bson.E{
					Key: "$or",
					Value: bson.A{
						bson.M{
							"code": codeQuery,
						},
						bson.M{
							"contractCode": codeQuery,
						},
					},
				},
			)
		}

		if lead != "" {
			leads, _ := s.findLead(lead)
			leadIds := make([]primitive.ObjectID, 0)
			if len(leads) != 0 {
				for _, item := range leads {
					leadIds = append(leadIds, item.ID)
				}
			}
			filter = append(filter, bson.E{
				Key:   "leadId",
				Value: bson.M{"$in": leadIds},
			})
		}

		if len(statuses) != 0 {
			filter = append(filter, bson.E{
				Key:   "status",
				Value: bson.M{"$in": strings.Split(statuses, ",")},
			})
		}

		if len(sources) != 0 {
			filter = append(filter, bson.E{
				Key:   "source",
				Value: bson.M{"$in": strings.Split(sources, ",")},
			})
		}

		if len(storeCodes) != 0 {
			filter = append(filter, bson.E{
				Key:   "storeCode",
				Value: bson.M{"$in": strings.Split(storeCodes, ",")},
			})
		}

		if len(employeeBys) != 0 && !onlyMe {
			filter = append(filter, bson.E{
				Key:   "employeeBy",
				Value: bson.M{"$in": strings.Split(employeeBys, ",")},
			})
		}

		if len(groups) != 0 {
			filter = append(filter, bson.E{
				Key:   "group",
				Value: bson.M{"$in": strings.Split(groups, ",")},
			})
		}

		if onlyMe {
			filter = append(filter, bson.E{
				Key:   "employeeBy",
				Value: user.Sub,
			})
		}

		if fromDate != "" && toDate != "" && filterDateType != "" {
			filter = append(filter, bson.E{
				Key:   string(filterDateType),
				Value: bson.M{"$gte": fromDate, "$lte": toDate},
			})
		}

		filter = mapPolicyToFilter(ctx, filter)
		filterCh <- filter
		findOptions := options.Find()
		findOptions.SetLimit(limit).SetSkip(skip)

		saleOpp, err := s.Repo.BaseRepo.Find(filter, findOptions)
		if err == nil && len(saleOpp) != 0 {
			sales <- saleOpp
		} else {
			sales <- []*entities.SaleOpportunity{}
		}

		wg.Done()
	}(sales, filterCh)

	go func(filterCh <-chan bson.D) {
		filter := <-filterCh
		total, _ = s.Repo.BaseRepo.Count(filter)

		wg.Done()
	}(filterCh)

	go func(sales <-chan []*entities.SaleOpportunity) {
		items = <-sales
		dataIncludes := IncludeDataSaleOpportunity{}
		includes := params.Includes
		if len(includes) != 0 && len(items) != 0 {
			dataIncludes = s.findDataIncludes(includes, items)
			items = utils.Map[*entities.SaleOpportunity](items, func(item *entities.SaleOpportunity) *entities.SaleOpportunity {
				item = mapIncludeDataToSale(includes, dataIncludes, item)
				return item
			})
		}

		wg.Done()
	}(sales)
	wg.Wait()

	return &SaleOpportunityPagination{
		Limit: limit,
		Skip:  skip,
		Total: total,
		List:  items,
	}, nil
}

func (s *SaleOpportunityService) Create(
	ctx *gin.Context,
	payload *validation.CreateSaleOpportunity,
	isLogged bool,
) (*entities.SaleOpportunity, error) {
	user := middlewares.LoggedUser(ctx)
	if isLogged && user == nil {
		return nil, errors.New("user not found")
	}
	createdBy := config.GetConfig().DefaultDataConfig.CreatedBy
	if user != nil {
		createdBy = user.Sub
	}

	entity, err := utils.TypeConverter[entities.SaleOpportunity](&payload)
	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	leadId := payload.LeadId
	var code string
	go func() {
		code = s.Repo.GenerateCode(payload.CodePrefix)
		wg.Done()
	}()

	var lead *entities.Lead
	go func() {
		lead, _ = s.LeadService.Repo.BaseRepo.FindById(leadId)
		wg.Done()
	}()
	wg.Wait()

	group := s.getSaleGroup(lead.Phone, lead)
	entity.Code = code
	entity.Group = group
	entity.CreatedBy = createdBy
	entity.UpdatedBy = createdBy

	entities.CreatingEntity(&entity.BaseEntity)

	item, err := s.Repo.BaseRepo.Create(entity)
	if err != nil {
		return nil, err
	}
	s.afterCreated(item, lead, payload)

	return item, nil
}

func (s *SaleOpportunityService) afterCreated(
	item *entities.SaleOpportunity,
	lead *entities.Lead,
	payload *validation.CreateSaleOpportunity,
) <-chan uint8 {
	r := make(chan uint8)
	go func() {
		defer close(r)

		_, _ = s.logRepo.BaseRepo.Create(&entities.Log{
			BeforeAttributes: utils.Omit(item, []string{
				"leadId", "source_refs", "code", "createdAt", "updatedAt", "ID", "updatedBy",
				"hash", "group", "lead", "employee", "created", "updated", "tagData", "id",
			}),
			AfterAttributes:   nil,
			SaleOpportunityId: item.ID,
			CreatedBy:         item.CreatedBy,
			CreatedAt:         time.Now(),
		})

		noteEntity := &entities.Note{
			Content:           payload.Note,
			SaleOpportunityId: item.ID,
			BaseEntity: entities.BaseEntity{
				CreatedBy: item.CreatedBy,
				UpdatedBy: item.UpdatedBy,
			},
		}
		entities.CreatingEntity(&noteEntity.BaseEntity)
		_, _ = s.noteRepo.BaseRepo.Create(noteEntity)

		customerId := s.findCustomerId(lead.ID)
		if customerId != "" {
			dataNotify := map[string]interface{}{
				"code":  item.Code,
				"order": utils.Pick(item, []string{"id", "code"}),
			}
			s.topicService.Send(
				config.GetConfig().TopicConfig.CustomerOrderNotification,
				map[string]interface{}{
					"data":      dataNotify,
					"receivers": []string{customerId},
				},
				map[string]string{
					"subscriptionType": string(TopicSubscriptionTypeOrderCreated),
				},
			)
		}
		r <- 1
	}()
	return r
}

func (s *SaleOpportunityService) Update(
	ctx *gin.Context,
	id primitive.ObjectID,
	payload validation.UpdateSaleOpportunity,
) (*entities.SaleOpportunity, error) {
	contractCode, status, storeCode, source := payload.ContractCode, payload.Status, payload.StoreCode, payload.Source

	if status == entities.SaleOppStatusSuccess && contractCode == "" {
		pkg.SendErrorResponse(ctx, pkg.ResponseError{
			StatusCode: http.StatusForbidden,
			Message:    "Can not update status to success without contract code",
		})
		return nil, nil
	}

	beforeItem, _ := s.Repo.BaseRepo.FindById(id)
	userRoles := ctx.MustGet(pkg.ProjectKeyUserRole).([]string)
	if !utils.Contains(userRoles, clients.RoleRegionalManager) && !utils.Contains(userRoles, clients.RoleSalesAdministrator) {
		if source == "DIRECT_OFFLINE" || storeCode == "" || beforeItem.Status == entities.SaleOppStatusSuccess {
			pkg.SendErrorResponse(ctx, pkg.ResponseError{
				StatusCode: http.StatusForbidden,
				Message:    "Can not update status to success without contract code",
			})
			return nil, nil
		}
	}

	if contractCode != "" {
		contract, err := s.contractClient.GetContract(contractCode)
		if err != nil {
			pkg.SendErrorResponse(ctx, pkg.ResponseError{
				StatusCode: http.StatusNotFound,
				Message:    "Contract code is invalid",
			})
			return nil, nil
		}
		payload.DisbursedAt = contract.DisbursedAt
		payload.DisbursedAmount = contract.LoanAmount
		payload.DisbursedAmount = contract.LoanAmount
	}

	user := middlewares.LoggedUser(ctx)
	update, err := s.Repo.BaseRepo.HandleDataUpdate(payload, user)
	if err != nil {
		return nil, err
	}

	item, _ := s.Repo.BaseRepo.UpdateByID(id, update)
	s.afterUpdated(beforeItem, item)
	return item, nil
}

func (s *SaleOpportunityService) afterUpdated(
	before *entities.SaleOpportunity,
	after *entities.SaleOpportunity,
) <-chan uint8 {
	r := make(chan uint8)
	go func() {
		defer close(r)

		beforeData := utils.Omit(before, []string{
			"leadId",
			"code",
			"createdAt",
			"updatedAt",
			"id",
			"createdBy",
			"metadata",
			"deletedAt",
			"assets",
			"source_refs",
			"updatedBy",
		})

		afterData := utils.Omit(after, []string{
			"leadId",
			"code",
			"createdAt",
			"updatedAt",
			"id",
			"createdBy",
			"metadata",
			"deletedAt",
			"assets",
			"source_refs",
			"updatedBy",
		})

		id, createdBy := after.ID, after.UpdatedBy

		var keyChange []string

		for key, val := range afterData {
			if !reflect.DeepEqual(val, beforeData[key]) {
				keyChange = append(keyChange, key)
			}
		}

		if len(keyChange) > 0 {
			beforeAttributes := utils.Pick(beforeData, keyChange)
			afterAttributes := utils.Pick(afterData, keyChange)
			if utils.Contains(keyChange, "tags") {
				beforeTags, afterTags := s.replaceTagCodes(beforeAttributes), s.replaceTagCodes(afterAttributes)
				beforeAttributes["tags"], afterAttributes["tags"] = <-beforeTags, <-afterTags
			}

			_, _ = s.logRepo.BaseRepo.Create(&entities.Log{
				ID:                primitive.ObjectID{},
				BeforeAttributes:  beforeAttributes,
				AfterAttributes:   afterAttributes,
				SaleOpportunityId: id,
				CreatedBy:         createdBy,
				CreatedAt:         time.Now(),
			})

			notification := s.findNotificationData(before, after)
			if notification != nil {
				customerId, content := notification.CustomerId, notification.Content
				s.topicService.Send(
					config.GetConfig().TopicConfig.CustomerOrderNotification,
					map[string]interface{}{
						"data": map[string]interface{}{
							"content": content,
							"order":   utils.Pick(after, []string{"id", "code"}),
							"code":    after.Code,
						},
						"receivers": []string{customerId},
					},
					map[string]string{
						"subscriptionType": string(TopicSubscriptionTypeOrderUpdated),
					},
				)
			}
		}
		r <- 1
	}()
	return r
}

type FindNotificationData struct {
	CustomerId string
	Content    string
}

func (s *SaleOpportunityService) findNotificationData(before *entities.SaleOpportunity, after *entities.SaleOpportunity) *FindNotificationData {
	beforeStatus, beforeEmployee := before.Status, before.EmployeeBy
	afterStatus, afterEmployee, leadId := after.Status, after.EmployeeBy, after.LeadId

	customerId := s.findCustomerId(leadId)

	if customerId == "" || (beforeStatus == afterStatus && beforeEmployee == afterEmployee) {
		return nil
	}

	if beforeEmployee != afterEmployee {
		employee, _ := s.employeeClient.FindById(afterEmployee)
		employeeName := "Admin"
		if employee != nil {
			employeeName = employee.DisplayName
		}
		return &FindNotificationData{
			CustomerId: customerId,
			Content:    fmt.Sprintf("Đơn hàng %s đã được chuyển cho nhân viên %s", after.Code, employeeName),
		}
	}

	statusesCheckOne := []string{string(entities.SaleOppStatusDealt), string(entities.SaleOppStatusConsulting), string(entities.SaleOppStatusUnContactable)}
	if utils.Contains(statusesCheckOne, string(beforeStatus)) && utils.Contains(statusesCheckOne, string(afterStatus)) {
		return nil
	}

	statusesCheckTwo := []string{string(entities.SaleOppStatusCancel), string(entities.SaleOppStatusDenied)}
	if utils.Contains(statusesCheckTwo, string(beforeStatus)) && utils.Contains(statusesCheckTwo, string(afterStatus)) {
		return nil
	}
	orderText := getOrderStatusDisplayText(afterStatus)
	return &FindNotificationData{
		CustomerId: customerId,
		Content:    fmt.Sprintf("Đơn hàng %s đã được chuyển sang trạng thái %s", after.Code, orderText),
	}
}

func getOrderStatusDisplayText(status entities.SaleOppStatus) string {
	switch status {
	case entities.SaleOppStatusNew:
		return "Moi"
	case entities.SaleOppStatusPending:
		return "Đã xác nhận"
	case entities.SaleOppStatusConsulting:
	case entities.SaleOppStatusDealt:
	case entities.SaleOppStatusUnContactable:
		return "Đang xử lý"
	case entities.SaleOppStatusSuccess:
		return "Thành công"
	case entities.SaleOppStatusCancel:
	case entities.SaleOppStatusDenied:
		return "Đã huỷ"
	default:
		return ""
	}
	return ""
}

func (s *SaleOpportunityService) findLead(keyword string) ([]*entities.Lead, error) {
	filter := s.LeadService.HandleLeadsKeyword(keyword)
	return s.LeadService.Repo.BaseRepo.Find(filter, options.Find().SetProjection(bson.D{
		{
			"_id", 1,
		},
	}))
}

func (s *SaleOpportunityService) findDataIncludes(
	includes []string,
	items []*entities.SaleOpportunity,
) IncludeDataSaleOpportunity {
	wg := sync.WaitGroup{}
	result := IncludeDataSaleOpportunity{}

	var (
		leadIds           []primitive.ObjectID
		employeeIds       []string
		saleIds           []primitive.ObjectID
		isIncludeLead     bool
		isIncludeEmployee bool
		isIncludeTag      bool
		totalWg           int
	)
	for _, sale := range items {
		leadIds = append(leadIds, sale.LeadId)
		saleIds = append(saleIds, sale.ID)
		employeeIds = append(employeeIds, sale.CreatedBy, sale.UpdatedBy, sale.EmployeeBy)
	}

	if utils.Contains(includes, "lead") {
		isIncludeLead = true
		totalWg++
	}
	if utils.Contains(includes, "created") || utils.Contains(includes, "updated") || utils.Contains(includes, "employee") {
		isIncludeEmployee = true
		totalWg++
	}
	if utils.Contains(includes, "tags") {
		isIncludeTag = true
		totalWg++
	}

	wg.Add(totalWg)
	if isIncludeLead {
		go func(wg *sync.WaitGroup, result *IncludeDataSaleOpportunity) {
			leads, _ := s.LeadService.Repo.BaseRepo.Find(bson.D{
				{
					"_id", bson.M{"$in": leadIds},
				},
			}, options.Find().SetProjection(bson.D{
				{
					"_id", 1,
				},
				{
					"fullName", 1,
				},
				{
					"phone", 1,
				},
				{
					"passportId", 1,
				},
			}))

			resultLead := make(map[primitive.ObjectID]*entities.Lead)
			for _, item := range leads {
				resultLead[item.ID] = item
			}
			result.Lead = resultLead
			wg.Done()
		}(&wg, &result)
	}

	if isIncludeEmployee {
		go func(wg *sync.WaitGroup, result *IncludeDataSaleOpportunity) {
			employees, _ := s.employeeClient.GetEmployees(employeeIds)
			result.Employee = employees
			wg.Done()
		}(&wg, &result)
	}

	if isIncludeTag {
		go func(wg *sync.WaitGroup, result *IncludeDataSaleOpportunity) {
			tags, _ := s.tagRepo.BaseRepo.Find(bson.D{}, options.Find().SetProjection(bson.D{
				{
					"_id", 1,
				},
				{
					"code", 1,
				},
				{
					"name", 1,
				},
			}))
			resultTag := make(map[string]*entities.Tag, len(tags))
			for _, item := range tags {
				resultTag[item.Code] = item
			}
			result.Tags = resultTag
			wg.Done()
		}(&wg, &result)
	}
	wg.Wait()
	return result
}

func (s *SaleOpportunityService) getSaleGroup(phone string, lead *entities.Lead) entities.SaleOppGroup {
	group := entities.GroupNew

	if lead == nil {
		lead, _ = s.LeadService.Repo.BaseRepo.FindOne(bson.D{{"phone", phone}}, nil)
	}
	if lead != nil {
		filter := bson.M{
			"leadId": lead.ID,
			"disbursedAt": bson.M{
				"$gte": currentTime.AddDate(0, 0, -90),
			},
			"contractCode": bson.M{
				"$ne": nil,
			},
			"disbursedAmount": bson.M{
				"$ne": 0,
			},
		}
		sale, err := s.Repo.BaseRepo.FindOne(filter, nil)
		if sale != nil && err == nil {
			group = entities.GroupOld
		}
	}

	return group
}
func (s *SaleOpportunityService) findCustomerId(id primitive.ObjectID) string {
	lead, _ := s.LeadService.Repo.BaseRepo.FindById(id)
	customerId := lead.CustomerId
	if len(customerId) == 0 {
		customer, err := s.digitalClient.FindByPhone(lead.Phone)
		if err != nil || customer == nil {
			return ""
		}
		customerId = customer.Guid
		_, _ = s.LeadService.Repo.BaseRepo.UpdateByID(lead.ID, bson.M{
			"customerId": customerId,
		})
	}
	return customerId
}

func (s *SaleOpportunityService) FindById(
	ctx *gin.Context,
	id primitive.ObjectID,
	includes []string,
) (*entities.SaleOpportunity, error) {
	s.employeeClient.Client.SetToken(ctx.GetHeader("Authorization"))

	filter := bson.D{
		{
			"_id", id,
		},
		{
			"deletedAt", nil,
		},
	}
	filter = mapPolicyToFilter(ctx, filter)
	item, err := s.Repo.BaseRepo.FindOne(filter, nil)
	if err != nil {
		return nil, err
	}
	if len(includes) > 0 {
		dataIncludes := s.findDataIncludes(includes, []*entities.SaleOpportunity{item})
		item = mapIncludeDataToSale(includes, dataIncludes, item)
	}
	return item, nil
}

func mapIncludeDataToSale(includes []string, dataIncludes IncludeDataSaleOpportunity, item *entities.SaleOpportunity) *entities.SaleOpportunity {
	if utils.Contains(includes, "lead") && dataIncludes.Lead != nil {
		item.Lead = dataIncludes.Lead[item.LeadId]
	}

	if utils.Contains(includes, "employee") && dataIncludes.Employee != nil {
		item.Employee = (*dataIncludes.Employee)[item.EmployeeBy]
	}

	if utils.Contains(includes, "created") && dataIncludes.Employee != nil {
		item.Employee = (*dataIncludes.Employee)[item.CreatedBy]
	}

	if utils.Contains(includes, "updated") && dataIncludes.Employee != nil {
		item.Employee = (*dataIncludes.Employee)[item.UpdatedBy]
	}

	if utils.Contains(includes, "tags") && dataIncludes.Tags != nil {
		var tagData []entities.Tag
		for _, tagCode := range item.Tags {
			tagData = append(tagData, *dataIncludes.Tags[tagCode])
		}
		item.TagData = tagData
	}

	return item
}

func (s *SaleOpportunityService) replaceTagCodes(payload map[string]interface{}) <-chan []string {
	tagCh := make(chan []string)

	go func() {
		defer close(tagCh)

		codes := payload["tags"]
		tags, _ := s.tagRepo.BaseRepo.Find(bson.D{{
			"code", bson.M{"$in": codes},
		}}, nil)

		var tagNames []string
		for _, tag := range tags {
			tagNames = append(tagNames, tag.Name)
		}
		tagCh <- tagNames
	}()
	return tagCh
}

func (s *SaleOpportunityService) PaginationLogs(
	ctx *gin.Context,
	id primitive.ObjectID,
	params validation.PaginationSaleOpportunityLogs,
) (*SaleOpportunityLogsPagination, error) {
	sale, err := s.Repo.BaseRepo.FindById(id)

	if err != nil || sale == nil {
		pkg.SendErrorResponse(ctx, pkg.ResponseError{
			StatusCode: http.StatusNotFound,
			Message:    "Sale opportunity not found",
		})
		return nil, nil
	}

	limit, skip := params.Limit, params.Skip
	findOptions := options.Find()
	findOptions.SetLimit(limit).SetSkip(skip)

	filter := bson.D{{
		"saleOpportunityId", id,
	}}

	logs, _ := s.logRepo.BaseRepo.Find(filter, findOptions)
	total, _ := s.logRepo.BaseRepo.Count(filter)

	return &SaleOpportunityLogsPagination{
		Total: total,
		List:  logs,
		Limit: limit,
		Skip:  skip,
	}, nil
}

func (s *SaleOpportunityService) Delete(ctx *gin.Context, id primitive.ObjectID) (*entities.SaleOpportunity, error) {
	userRoles := ctx.MustGet(pkg.ProjectKeyUserRole).([]string)
	if !utils.Contains(userRoles, clients.RoleSalesAdministrator) {
		pkg.SendErrorResponse(ctx, pkg.ResponseError{
			StatusCode: http.StatusForbidden,
			Message:    "You don't have permission to delete sale opportunity",
		})
		return nil, nil
	}
	return s.Repo.BaseRepo.UpdateByID(id, bson.M{"deletedAt": time.Now()})
}
