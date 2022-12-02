package repositories

import (
	"context"
	"crm-service-go/app/middlewares"
	"crm-service-go/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"time"
)

type BaseRepository[T any] struct {
	col *mongo.Collection
	ctx context.Context
}

func (r *BaseRepository[T]) Find(filter interface{}, opts *options.FindOptions) ([]*T, error) {
	var results = make([]*T, 0)
	cursor, _ := r.col.Find(r.ctx, filter, opts)

	if err := cursor.All(r.ctx, &results); err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	return results, nil
}

func (r *BaseRepository[T]) Count(filter interface{}) (int64, error) {
	total, err := r.col.CountDocuments(r.ctx, filter)
	if err != nil {
		utils.Logger.Error(err)
		return 0, err
	}
	return total, nil
}

func (r *BaseRepository[T]) FindOne(filter interface{}, opts *options.FindOneOptions) (*T, error) {
	var item *T

	cursor := r.col.FindOne(r.ctx, filter, opts)
	if cursor.Err() != nil {
		utils.Logger.Error(cursor.Err())
		return nil, cursor.Err()
	}

	if err := cursor.Decode(&item); err != nil {
		utils.Logger.Error(err)
		return nil, err
	}

	return item, nil
}

func (r *BaseRepository[T]) Create(entity *T) (*T, error) {
	result, err := r.col.InsertOne(r.ctx, entity)
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}

	var item *T
	if err = r.col.FindOne(r.ctx, bson.M{"_id": result.InsertedID}).Decode(&item); err != nil {
		panic(err)
	}
	return item, nil
}

func (r *BaseRepository[T]) UpdateByID(
	Id primitive.ObjectID,
	payload bson.M,
) (*T, error) {
	_, err := r.col.UpdateByID(r.ctx, Id, bson.D{{
		"$set", payload,
	}})
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}

	var item *T
	if err = r.col.FindOne(r.ctx, bson.M{"_id": Id}).Decode(&item); err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	return item, nil
}

func (r *BaseRepository[T]) Delete(filter bson.D) (bool, error) {
	_, err := r.col.DeleteOne(r.ctx, filter)
	if err != nil {
		utils.Logger.Error(err)
		return false, err
	}
	return true, nil
}

func (r *BaseRepository[T]) FindById(Id primitive.ObjectID) (*T, error) {
	var item *T

	cursor := r.col.FindOne(r.ctx, bson.M{"_id": Id, "deletedAt": nil}, nil)
	if cursor.Err() != nil {
		utils.Logger.Error(cursor.Err())
		return nil, cursor.Err()
	}

	if err := cursor.Decode(&item); err != nil {
		utils.Logger.Error(err)
		return nil, err
	}

	return item, nil
}

func (r *BaseRepository[T]) HandleDataUpdate(s interface{}, user *middlewares.TokenPayload) (bson.M, error) {
	result := bson.M{"updatedAt": time.Now()}
	if user != nil {
		result["updatedBy"] = user.Sub
	}
	values := reflect.ValueOf(s)
	types := values.Type()
	for i := 0; i < values.NumField(); i++ {
		field := values.Field(i)
		fieldJson := types.Field(i).Tag.Get("json")
		fieldData := field.Interface()

		if !isNil(field) && fieldJson != "-" {
			result[fieldJson] = fieldData
		}
	}

	return result, nil
}

func isNil(field reflect.Value) bool {
	typeOf := reflect.TypeOf(field.Interface())
	if typeOf == reflect.TypeOf(primitive.ObjectID{}) {
		return field.Interface().(primitive.ObjectID).IsZero()
	} else {
		switch field.Kind() {
		case reflect.Ptr, reflect.Map, reflect.Chan, reflect.Slice:
			return field.IsNil()
		case reflect.String, reflect.Array:
			return field.Len() == 0
		case reflect.Struct:
			return field.IsZero()
		}
	}
	return false
}
