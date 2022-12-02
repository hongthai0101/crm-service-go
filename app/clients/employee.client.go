package clients

import (
	"bytes"
	"crm-service-go/config"
	"crm-service-go/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type IEmployee struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

type EmployeeClient struct {
	Client *Client
}

func NewEmployeeClient() *EmployeeClient {
	return &EmployeeClient{
		Client: NewClient(config.GetConfig().ServiceConfig.EmployeeUrl),
	}
}

func (c *EmployeeClient) findByIds(ids []string) ([]*IEmployee, error) {
	postBody, _ := json.Marshal(map[string][]string{
		"ids": ids,
	})

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/employees/list", c.Client.baseURL), bytes.NewBuffer(postBody))
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}

	req = req.WithContext(HttpCtx)
	var employees []*IEmployee
	if err := c.Client.sendRequest(req, &employees); err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	return employees, nil
}

func (c *EmployeeClient) GetEmployees(ids []string) (*map[string]string, error) {
	employees, _ := c.findByIds(ids)
	tempResult := make(map[string]string, len(employees))
	for i := 0; i < len(employees); i++ {
		tempResult[employees[i].ID] = employees[i].DisplayName
	}

	result := make(map[string]string, len(ids))
	for i := 0; i < len(ids); i++ {
		if val, ok := tempResult[ids[i]]; ok {
			result[ids[i]] = val
		} else {
			result[ids[i]] = config.GetConfig().DefaultDataConfig.CreatedName
		}
	}

	return &result, nil
}

func (c *EmployeeClient) FindById(id string) (*IEmployee, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/employees/%s", c.Client.baseURL, id), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(HttpCtx)
	var employee IEmployee
	if err = c.Client.sendRequest(req, &employee); err != nil {
		return nil, err
	}
	return &employee, nil
}
