package clients

import (
	"crm-service-go/config"
	"fmt"
	"net/http"
)

type ContractClient struct {
	Client *Client
}

func NewContractClient() *ContractClient {
	return &ContractClient{
		Client: NewClient(config.GetConfig().ServiceConfig.ContractUrl),
	}
}

type ContractSchema struct {
	ContractCode string `json:"contractCode"`
	DisbursedAt  string `json:"disbursedAt"`
	LoanAmount   uint16 `json:"loanAmount"`
	CreatedBy    string `json:"createdBy"`
}

func (c *ContractClient) GetContract(contractCode string) (*ContractSchema, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/inside/contracts/%s", c.Client.baseURL, contractCode), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(HttpCtx)
	var res ContractSchema
	if err = c.Client.sendRequest(req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
