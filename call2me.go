// Package call2me provides a Go client for the Call2Me API.
//
// Usage:
//
//	client := call2me.New("sk_call2me_...")
//	agents, err := client.Agents.List()
package call2me

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const defaultBaseURL = "https://api.call2me.app"

type Client struct {
	APIKey  string
	BaseURL string
	HTTP    *http.Client
	Agents  *AgentsService
	Calls   *CallsService
}

func New(apiKey string) *Client {
	c := &Client{
		APIKey:  apiKey,
		BaseURL: defaultBaseURL,
		HTTP:    &http.Client{},
	}
	c.Agents = &AgentsService{client: c}
	c.Calls = &CallsService{client: c}
	return c
}

func (c *Client) do(method, path string, body interface{}) ([]byte, error) {
	var reader io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(data))
	}
	return data, nil
}

type AgentsService struct{ client *Client }

func (s *AgentsService) List() ([]map[string]interface{}, error) {
	data, err := s.client.do("GET", "/v1/agents", nil)
	if err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	json.Unmarshal(data, &result)
	return result, nil
}

func (s *AgentsService) Get(id string) (map[string]interface{}, error) {
	data, err := s.client.do("GET", "/v1/agents/"+id, nil)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	json.Unmarshal(data, &result)
	return result, nil
}

func (s *AgentsService) Create(agent map[string]interface{}) (map[string]interface{}, error) {
	data, err := s.client.do("POST", "/v1/agents", agent)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	json.Unmarshal(data, &result)
	return result, nil
}

type CallsService struct{ client *Client }

func (s *CallsService) List() ([]map[string]interface{}, error) {
	data, err := s.client.do("GET", "/v1/calls", nil)
	if err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	json.Unmarshal(data, &result)
	return result, nil
}
