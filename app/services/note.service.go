package services

import (
	"crm-service-go/app/entities"
	"crm-service-go/app/middlewares"
	"crm-service-go/app/repositories"
	"crm-service-go/app/validation"
	"crm-service-go/pkg"
	"crm-service-go/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NotePagination pkg.Pagination[entities.Note]

type NoteService struct {
	Repo *repositories.NoteRepository
}

func NewNoteService(repo *repositories.NoteRepository) *NoteService {
	return &NoteService{Repo: repo}
}

func (s *NoteService) Pagination(param validation.PaginationNote) (NotePagination, error) {
	var items []*entities.Note

	saleOpportunitiesId, _ := primitive.ObjectIDFromHex(param.SaleOpportunityId)
	filter := bson.D{
		{
			"deletedAt", nil,
		},
		{
			"saleOpportunitiesId", saleOpportunitiesId,
		},
	}

	limit := pkg.GetLimit(param.Limit)
	skip := pkg.GetSkip(param.Skip)
	findOptions := options.Find()
	findOptions.SetLimit(limit).SetSkip(skip)

	items, err := s.Repo.BaseRepo.Find(filter, findOptions)
	if err != nil {
		return NotePagination{}, err
	}
	total, _ := s.Repo.BaseRepo.Count(filter)

	return NotePagination{
		Limit: limit,
		Skip:  skip,
		Total: total,
		List:  items,
	}, nil
}

func (s *NoteService) Create(payload validation.CreateNote) (*entities.Note, error) {
	entity, err := utils.TypeConverter[entities.Note](&payload)
	if err != nil {
		return nil, err
	}
	entities.CreatingEntity(&entity.BaseEntity)
	return s.Repo.BaseRepo.Create(entity)
}

func (s *NoteService) Update(ctx *gin.Context, id primitive.ObjectID, payload validation.UpdateNote) (*entities.Note, error) {
	user := middlewares.LoggedUser(ctx)
	update, err := s.Repo.BaseRepo.HandleDataUpdate(payload, user)
	if err != nil {
		return nil, err
	}
	return s.Repo.BaseRepo.UpdateByID(id, update)
}
