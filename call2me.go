// Package call2me provides a Go client for the Call2Me API.
package call2me

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const defaultBaseURL = "https://api.call2me.app"

type Client struct {
	APIKey  string
	BaseURL string
	HTTP    *http.Client
	Agents       *AgentsService
	Calls        *CallsService
	KnowledgeBase *KBService
	Wallet       *WalletService
	Campaigns    *CampaignsService
	Schedules    *SchedulesService
	PhoneNumbers *PhoneNumbersService
	SipTrunks    *SipTrunksService
	ApiKeys      *ApiKeysService
	Users        *UsersService
	Widgets      *WidgetsService
	Voices       *VoicesService
	Payments     *PaymentsService
}

func New(apiKey string) *Client {
	c := &Client{APIKey: apiKey, BaseURL: defaultBaseURL, HTTP: &http.Client{}}
	c.Agents = &AgentsService{c}
	c.Calls = &CallsService{c}
	c.KnowledgeBase = &KBService{c}
	c.Wallet = &WalletService{c}
	c.Campaigns = &CampaignsService{c}
	c.Schedules = &SchedulesService{c}
	c.PhoneNumbers = &PhoneNumbersService{c}
	c.SipTrunks = &SipTrunksService{c}
	c.ApiKeys = &ApiKeysService{c}
	c.Users = &UsersService{c}
	c.Widgets = &WidgetsService{c}
	c.Voices = &VoicesService{c}
	c.Payments = &PaymentsService{c}
	return c
}

type M = map[string]interface{}

func (c *Client) do(method, path string, body interface{}) ([]byte, error) {
	var reader io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reader = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, c.BaseURL+path, reader)
	if err != nil { return nil, err }
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTP.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(data))
	}
	return data, nil
}

func (c *Client) get(path string, params ...string) ([]byte, error) {
	if len(params) > 0 {
		u, _ := url.Parse(c.BaseURL + path)
		q := u.Query()
		for i := 0; i+1 < len(params); i += 2 { q.Set(params[i], params[i+1]) }
		u.RawQuery = q.Encode()
		req, _ := http.NewRequest("GET", u.String(), nil)
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
		resp, err := c.HTTP.Do(req)
		if err != nil { return nil, err }
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)
		if resp.StatusCode >= 400 { return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(data)) }
		return data, nil
	}
	return c.do("GET", path, nil)
}

func list(data []byte, err error) ([]M, error) { if err != nil { return nil, err }; var r []M; json.Unmarshal(data, &r); return r, nil }
func one(data []byte, err error) (M, error) { if err != nil { return nil, err }; var r M; json.Unmarshal(data, &r); return r, nil }

type AgentsService struct{ c *Client }
func (s *AgentsService) List() ([]M, error)                     { return list(s.c.get("/v1/agents")) }
func (s *AgentsService) Get(id string) (M, error)               { return one(s.c.do("GET", "/v1/agents/"+id, nil)) }
func (s *AgentsService) Create(data M) (M, error)               { return one(s.c.do("POST", "/v1/agents", data)) }
func (s *AgentsService) Update(id string, data M) (M, error)    { return one(s.c.do("PATCH", "/v1/agents/"+id, data)) }
func (s *AgentsService) Delete(id string) error                  { _, err := s.c.do("DELETE", "/v1/agents/"+id, nil); return err }
func (s *AgentsService) Duplicate(id string) (M, error)          { return one(s.c.do("POST", "/v1/agents/"+id+"/duplicate", nil)) }
func (s *AgentsService) Stats(id string) (M, error)              { return one(s.c.get("/v1/agents/"+id+"/stats")) }

type CallsService struct{ c *Client }
func (s *CallsService) List() ([]M, error)      { return list(s.c.get("/v1/calls")) }
func (s *CallsService) Get(id string) (M, error) { return one(s.c.do("GET", "/v1/calls/"+id, nil)) }
func (s *CallsService) End(id string) (M, error) { return one(s.c.do("POST", "/v1/calls/"+id+"/end", nil)) }

type KBService struct{ c *Client }
func (s *KBService) List() ([]M, error)                   { return list(s.c.get("/v1/knowledge-base")) }
func (s *KBService) Get(id string) (M, error)             { return one(s.c.do("GET", "/v1/knowledge-base/"+id, nil)) }
func (s *KBService) Create(data M) (M, error)             { return one(s.c.do("POST", "/v1/knowledge-base", data)) }
func (s *KBService) Delete(id string) error                { _, err := s.c.do("DELETE", "/v1/knowledge-base/"+id, nil); return err }
func (s *KBService) Query(id, q string) (M, error)        { return one(s.c.do("POST", "/v1/knowledge-base/"+id+"/query", M{"query": q})) }

type WalletService struct{ c *Client }
func (s *WalletService) Balance() (M, error)        { return one(s.c.get("/v1/wallet/balance")) }
func (s *WalletService) Transactions() ([]M, error) { return list(s.c.get("/v1/wallet/transactions")) }
func (s *WalletService) Analytics() (M, error)      { return one(s.c.get("/v1/wallet/analytics")) }

type CampaignsService struct{ c *Client }
func (s *CampaignsService) List() ([]M, error)         { return list(s.c.get("/v1/campaigns")) }
func (s *CampaignsService) Get(id string) (M, error)    { return one(s.c.do("GET", "/v1/campaigns/"+id, nil)) }
func (s *CampaignsService) Create(data M) (M, error)    { return one(s.c.do("POST", "/v1/campaigns", data)) }
func (s *CampaignsService) Start(id string) (M, error)  { return one(s.c.do("POST", "/v1/campaigns/"+id+"/action", M{"action": "start"})) }
func (s *CampaignsService) Pause(id string) (M, error)  { return one(s.c.do("POST", "/v1/campaigns/"+id+"/action", M{"action": "pause"})) }
func (s *CampaignsService) Cancel(id string) (M, error) { return one(s.c.do("POST", "/v1/campaigns/"+id+"/action", M{"action": "cancel"})) }

type SchedulesService struct{ c *Client }
func (s *SchedulesService) List() ([]M, error)      { return list(s.c.get("/v1/schedules")) }
func (s *SchedulesService) Create(data M) (M, error) { return one(s.c.do("POST", "/v1/schedules", data)) }
func (s *SchedulesService) Delete(id string) error    { _, err := s.c.do("DELETE", "/v1/schedules/"+id, nil); return err }

type PhoneNumbersService struct{ c *Client }
func (s *PhoneNumbersService) List() ([]M, error)                      { return list(s.c.get("/v1/phone-numbers")) }
func (s *PhoneNumbersService) Create(data M) (M, error)                { return one(s.c.do("POST", "/v1/phone-numbers", data)) }
func (s *PhoneNumbersService) Delete(num string) error                  { _, err := s.c.do("DELETE", "/v1/phone-numbers/"+num, nil); return err }
func (s *PhoneNumbersService) BindAgent(num, agentID string) (M, error) { return one(s.c.do("POST", "/v1/phone-numbers/"+num+"/bind", M{"agent_id": agentID})) }

type SipTrunksService struct{ c *Client }
func (s *SipTrunksService) List() ([]M, error)      { return list(s.c.get("/v1/sip-trunks")) }
func (s *SipTrunksService) Create(data M) (M, error) { return one(s.c.do("POST", "/v1/sip-trunks", data)) }
func (s *SipTrunksService) Delete(id string) error    { _, err := s.c.do("DELETE", "/v1/sip-trunks/"+id, nil); return err }
func (s *SipTrunksService) Test(id string) (M, error) { return one(s.c.do("POST", "/v1/sip-trunks/"+id+"/test", nil)) }

type ApiKeysService struct{ c *Client }
func (s *ApiKeysService) List() ([]M, error)        { return list(s.c.get("/v1/api-keys")) }
func (s *ApiKeysService) Create(data M) (M, error)   { return one(s.c.do("POST", "/v1/api-keys", data)) }
func (s *ApiKeysService) Revoke(id string) (M, error) { return one(s.c.do("PATCH", "/v1/api-keys/"+id+"/revoke", nil)) }
func (s *ApiKeysService) Delete(id string) error      { _, err := s.c.do("DELETE", "/v1/api-keys/"+id, nil); return err }

type UsersService struct{ c *Client }
func (s *UsersService) Me() (M, error)           { return one(s.c.get("/v1/users/me")) }
func (s *UsersService) Update(data M) (M, error)  { return one(s.c.do("PATCH", "/v1/users/me", data)) }
func (s *UsersService) Stats() (M, error)         { return one(s.c.get("/v1/users/me/stats")) }
func (s *UsersService) Branding() (M, error)      { return one(s.c.get("/v1/users/me/branding")) }

type WidgetsService struct{ c *Client }
func (s *WidgetsService) List() ([]M, error)                   { return list(s.c.get("/v1/widgets")) }
func (s *WidgetsService) Create(data M) (M, error)              { return one(s.c.do("POST", "/v1/widgets", data)) }
func (s *WidgetsService) Delete(id string) error                 { _, err := s.c.do("DELETE", "/v1/widgets/"+id, nil); return err }
func (s *WidgetsService) Chat(id, msg string) (M, error)        { return one(s.c.do("POST", "/v1/widgets/"+id+"/chat", M{"message": msg})) }

type VoicesService struct{ c *Client }
func (s *VoicesService) List() ([]M, error) { return list(s.c.get("/v1/voices")) }

type PaymentsService struct{ c *Client }
func (s *PaymentsService) Checkout(amount float64, currency string) (M, error) { return one(s.c.do("POST", "/v1/payments/checkout", M{"amount": amount, "currency": currency})) }
func (s *PaymentsService) History() ([]M, error)    { return list(s.c.get("/v1/payments/history")) }
func (s *PaymentsService) SavedCards() ([]M, error)  { return list(s.c.get("/v1/payments/saved-cards")) }
