package validation

import (
	"context"
	"crm-service-go/app/repositories"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
)

type ApiError struct {
	Field string `json:"field"`
	Msg   string `json:"msg"`
}

func BeautyMessage(err error) []ApiError {

	validationErrors := err.(validator.ValidationErrors)
	out := make([]ApiError, len(validationErrors))
	for i, fe := range validationErrors {
		out[i] = ApiError{Field: fe.Field(), Msg: msgForTag(fe.Tag())}
	}
	return out
}

func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "gte", "lte", "lt", "gt":
		return "The value invalid"
	case "datetime":
		return "The value must be a valid date"
	case "phoneExists":
		return "Phone number already exists"
	}
	return "something went wrong"
}

var LeadPhoneExists validator.Func = func(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(string)
	if len(value) > 0 {
		leadRepo := repositories.NewLeadRepository(context.Background())
		total, _ := leadRepo.BaseRepo.Count(bson.M{"phone": value})
		return total == 0
	}
	return true
}
