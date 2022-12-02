package clients

import (
	"crm-service-go/config"
	"crm-service-go/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type DigitalClient struct {
	Client *Client
}

func NewDigitalClient() *DigitalClient {
	return &DigitalClient{
		Client: NewClient(config.GetConfig().ServiceConfig.DigitalUrl),
	}
}

type DigitalCustomerSchema struct {
	Phone string `json:"phone"`
	ID    string `json:"id"`
	Guid  string `json:"guid"`
}

func (c *DigitalClient) FindByPhone(phone string) (*DigitalCustomerSchema, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/inside/users", c.Client.baseURL), nil)
	if err != nil {
		return nil, err
	}

	filter := map[string]interface{}{
		"where": map[string]interface{}{
			"phone": phone,
		},
		"limit": 1,
		"order": [1]string{"createdAt DESC"},
	}

	q := req.URL.Query()
	jsonStr, _ := json.Marshal(filter)
	q.Add("filter", string(jsonStr))
	req.URL.RawQuery = q.Encode()

	req = req.WithContext(HttpCtx)
	var res []DigitalCustomerSchema
	if err = c.Client.sendRequest(req, &res); err != nil {
		return nil, err
	}
	return utils.Get(res, 0), nil
}
