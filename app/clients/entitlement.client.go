package clients

import (
	"crm-service-go/config"
	"fmt"
	"net/http"
)

type PolicyResource string

const (
	PolicyResourceSaleOpportunity PolicyResource = "SALES_OPPORTUNITY"
	PolicyResourceLead            PolicyResource = "LEAD"
)

type AuthorizationAction string

const (
	AuthorizationActionReadAny   AuthorizationAction = "read:any"
	AuthorizationActionCreateAny AuthorizationAction = "create:any"
	AuthorizationActionUpdateAny AuthorizationAction = "update:any"
	AuthorizationActionDeleteAny AuthorizationAction = "delete:any"
)

const (
	RoleSalesAdministrator string = "sales_administrator"
	RoleRegionalManager    string = "sales_administrator"
)

type EntitlementSchema struct {
	Roles    []string        `json:"roles"`
	Title    []string        `json:"title"`
	Policies []*PolicySchema `json:"policies"`
}

type PolicySchema struct {
	Resource PolicyResource  `json:"resource"`
	Subject  []SubjectSchema `json:"subject"`
}

type SubjectSchema struct {
	Name    string                `json:"name"`
	Actions []AuthorizationAction `json:"actions"`
}

type EntitlementClient struct {
	Client *Client
}

func NewEntitlementClient() *EntitlementClient {
	return &EntitlementClient{
		Client: NewClient(config.GetConfig().ServiceConfig.EntitlementUrl),
	}
}

func (c *EntitlementClient) FindPolicy(id string) (*EntitlementSchema, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/policy/%s", c.Client.baseURL, id), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(HttpCtx)
	var res EntitlementSchema
	if err = c.Client.sendRequest(req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
