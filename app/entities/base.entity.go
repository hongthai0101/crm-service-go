package entities

import (
	"crm-service-go/config"
	time "time"
)

type BaseEntity struct {
	CreatedBy string      `bson:"createdBy" json:"createdBy,omitempty"`
	UpdatedBy string      `bson:"updatedBy" json:"updatedBy,omitempty"`
	CreatedAt time.Time   `bson:"createdAt" json:"createdAt,omitempty"`
	UpdatedAt time.Time   `bson:"updatedAt" json:"updatedAt,omitempty"`
	DeletedAt interface{} `bson:"deletedAt" json:"deletedAt,omitempty"`
}

func CreatingEntity(entity *BaseEntity) {
	entity.UpdatedAt = time.Now()
	entity.CreatedAt = time.Now()
	if len(entity.CreatedBy) == 0 {
		entity.CreatedBy = config.GetConfig().DefaultDataConfig.CreatedBy
	}
	if len(entity.UpdatedBy) == 0 {
		entity.UpdatedBy = config.GetConfig().DefaultDataConfig.CreatedBy
	}
	entity.DeletedAt = nil
}

func UpdatingEntity(entity *BaseEntity) {
	entity.UpdatedAt = time.Now()
	if len(entity.UpdatedBy) == 0 {
		entity.UpdatedBy = config.GetConfig().DefaultDataConfig.CreatedBy
	}
}

func DeletingEntity(entity *BaseEntity) {
	entity.DeletedAt = time.Now()
}
