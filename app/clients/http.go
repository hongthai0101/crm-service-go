package clients

import (
	"context"
	utils2 "crm-service-go/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	token      string
	baseURL    string
	HTTPClient *http.Client
}

func NewClient(BaseURL string) *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
		baseURL: BaseURL,
	}
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var HttpCtx = context.Background()

func (c *Client) sendRequest(req *http.Request, v interface{}) error {

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", c.token)
	req.Close = true

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		utils2.Logger.Error(err)
		return err
	}
	defer res.Body.Close()

	if !utils2.ContainsInt([]int{http.StatusOK, http.StatusCreated, http.StatusNoContent}, res.StatusCode) {
		bodyBytes, _ := io.ReadAll(res.Body)
		utils2.Logger.Error(string(bodyBytes))
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}
		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		utils2.Debug(err)
		return err
	}
	return nil
}

func (c *Client) SetToken(token string) {
	c.token = token
}
