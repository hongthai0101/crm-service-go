package services

import (
	"crm-service-go/app/entities"
	"crm-service-go/app/middlewares"
	"crm-service-go/app/repositories"
	"crm-service-go/app/validation"
	"crm-service-go/pkg"
	"crm-service-go/pkg/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LeadPagination pkg.Pagination[entities.Lead]

type LeadService struct {
	Repo *repositories.LeadRepository
}

func NewLeadService(repo *repositories.LeadRepository) *LeadService {
	return &LeadService{Repo: repo}
}

func (s *LeadService) Pagination(ctx *gin.Context, param validation.PaginationLead) (LeadPagination, error) {
	var items []*entities.Lead

	keyword := param.Keyword
	filter := bson.D{
		{
			"deletedAt", nil,
		},
	}
	if keyword != "" {
		filter = s.HandleLeadsKeyword(keyword)
	}
	filter = mapPolicyToFilter(ctx, filter)

	limit := pkg.GetLimit(param.Limit)
	skip := pkg.GetSkip(param.Skip)
	findOptions := options.Find()
	findOptions.SetLimit(limit).SetSkip(skip)

	items, err := s.Repo.BaseRepo.Find(filter, findOptions)
	if err != nil {
		return LeadPagination{}, err
	}
	total, _ := s.Repo.BaseRepo.Count(filter)

	return LeadPagination{
		Limit: limit,
		Skip:  skip,
		Total: total,
		List:  items,
	}, nil
}

func (s *LeadService) Create(ctx *gin.Context, payload validation.CreateLead) (*entities.Lead, error) {
	user := middlewares.LoggedUser(ctx)
	nationalId, passportId, phone := payload.NationalId, payload.PassportId, payload.Phone
	filter := bson.D{
		{
			"deletedAt", nil,
		},
	}
	if len(nationalId) != 0 || len(passportId) != 0 || len(phone) != 0 {
		values := bson.A{}
		if len(nationalId) != 0 {
			values = append(values, bson.M{
				"nationalId": nationalId,
			})
		}
		if len(passportId) != 0 {
			values = append(values, bson.M{
				"passportId": passportId,
			})
		}
		if len(phone) != 0 {
			values = append(values, bson.M{
				"phone": phone,
			})
		}
		if len(values) != 0 {
			filter = append(filter, bson.E{
				Key:   "$or",
				Value: values,
			})
		}
	}
	exists, err := s.Repo.BaseRepo.Count(filter)
	if err != nil {
		return nil, err
	}
	if exists != 0 {
		return nil, errors.New("lead already exists")
	}

	entity, err := utils.TypeConverter[entities.Lead](&payload)
	if err != nil {
		return nil, err
	}
	if user != nil {
		entity.CreatedBy = user.Sub
		entity.UpdatedBy = user.Sub
	}
	entities.CreatingEntity(&entity.BaseEntity)
	return s.Repo.BaseRepo.Create(entity)
}

func (s *LeadService) Update(ctx *gin.Context, id primitive.ObjectID, payload validation.UpdateLead) (*entities.Lead, error) {
	user := middlewares.LoggedUser(ctx)
	update, err := s.Repo.BaseRepo.HandleDataUpdate(payload, user)
	if err != nil {
		return nil, err
	}
	return s.Repo.BaseRepo.UpdateByID(id, update)
}

func (s *LeadService) HandleLeadsKeyword(keyword string) bson.D {
	keywordQuery := primitive.Regex{Pattern: keyword, Options: "i"}
	filter := bson.D{
		{
			"deletedAt", nil,
		},
		{
			"$or", bson.A{
				bson.M{
					"phone": keywordQuery,
				},
				bson.M{
					"passportId": keywordQuery,
				},
				bson.M{
					"nationalId": keywordQuery,
				},
				bson.M{
					"fullName": keywordQuery,
				},
			},
		},
	}
	return filter
}
