package services

import (
	"cloud.google.com/go/pubsub"
	"context"
	"crm-service-go/config"
	"crm-service-go/pkg/utils"
	"encoding/json"
)

type TopicSubscriptionType string

const (
	TopicSubscriptionTypeOrderCreated TopicSubscriptionType = "customer.order.created"
	TopicSubscriptionTypeOrderUpdated TopicSubscriptionType = "customer.order.updated"
)

type TopicService struct{}

func NewTopicService() *TopicService {
	return &TopicService{}
}

func (s *TopicService) Send(topicName string, body interface{}, attributes map[string]string) bool {

	utils.Logger.Info(map[string]interface{}{
		"topic":      topicName,
		"body":       body,
		"attributes": attributes,
	})

	ctx := context.Background()
	projectID := config.GetConfig().GCSConfig.ProjectId
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		utils.Logger.Error(err.Error())
		return false
	}
	topic := client.Topic(topicName)

	data, _ := json.Marshal(body)
	var msg = &pubsub.Message{
		Data:       data,
		Attributes: attributes,
	}
	if _, err = topic.Publish(ctx, msg).Get(ctx); err != nil {
		utils.Logger.Error(err.Error())
		return false
	}
	return true
}
