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
	Events       *EventsService
	VoiceSessions *VoiceSessionsService
	EndUsers      *EndUsersService
	Webhooks      *WebhooksService
	// 5 Eki 2026'da eklendi: 278 uçtan yalnız 83'ü kapsanıyordu ve
	// tercümanın 11 ucunun HİÇBİRİ yoktu.
	Interpreters  *InterpretersService
	Numbers       *NumbersService
	Sms           *SmsService
	Extension     *ExtensionService
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
	c.Events = &EventsService{c}
	c.VoiceSessions = &VoiceSessionsService{c}
	c.EndUsers = &EndUsersService{c}
	c.Webhooks = &WebhooksService{c}
	c.Interpreters = &InterpretersService{c}
	c.Numbers = &NumbersService{c}
	c.Sms = &SmsService{c}
	c.Extension = &ExtensionService{c}
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

type VoiceSessionsService struct{ c *Client }
func (s *VoiceSessionsService) CreateSession(agentID string, context M) (M, error) {
	return one(s.c.do("POST", "/v1/voice/sessions", M{"agent_id": agentID, "context": context}))
}
// Get fetches a voice session's detail + transcript by room name.
func (s *VoiceSessionsService) Get(roomName string) (M, error) {
	return one(s.c.do("GET", "/v1/voice/sessions/"+roomName, nil))
}

type EndUsersService struct{ c *Client }
// CreateToken mints an ephemeral end-user token (eut_) billed to the calling tenant.
func (s *EndUsersService) CreateToken(externalID string, opts M) (M, error) {
	return one(s.c.do("POST", "/v1/end-users/"+externalID+"/tokens", opts))
}

type WebhooksService struct{ c *Client }
// Set replaces this tenant's webhook URL + secret (secret auto-generated if omitted).
func (s *WebhooksService) Set(opts M) (M, error) { return one(s.c.do("PUT", "/v1/webhooks", opts)) }
func (s *WebhooksService) Get() (M, error)       { return one(s.c.do("GET", "/v1/webhooks", nil)) }

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
// Create places a REAL outbound call — to_number is dialled and charged.
//
// `topic` is why the call is being placed; it reaches the agent as
// call_purpose and is woven into both the opening line and the system
// prompt, so the agent states its reason instead of a generic greeting.
//
// Missing from this SDK until 5 Oct 2026 despite the README advertising
// outbound calling.
func (s *CallsService) Create(data M) (M, error) { return one(s.c.do("POST", "/v1/calls", clean(data))) }
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
func (s *ApiKeysService) Revoke(id string) (M, error) { return one(s.c.do("DELETE", "/v1/api-keys/"+id, nil)) }
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
func (s *PaymentsService) SavedCards() ([]M, error)  { return list(s.c.get("/v1/payments/methods")) }

// EventsService — report operational / business events to Call2Me.
// POST /v1/events is public (no auth required) but sending the API key
// lifts the per-minute ceiling from 10 (anon) to 100 (authenticated).
// GET is admin-only.
type EventsService struct{ c *Client }

// EventReport is the payload for Events.Report.
// Fields mirror the server's IncomingEvent model — only Type, Source, and
// Message are required.
type EventReport struct {
	Type        string                 `json:"type"`
	Source      string                 `json:"source"`
	Message     string                 `json:"message"`
	Severity    string                 `json:"severity,omitempty"`
	Tenant      string                 `json:"tenant,omitempty"`
	SessionID   string                 `json:"session_id,omitempty"`
	Fingerprint string                 `json:"fingerprint,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
}

func (s *EventsService) Report(r EventReport) (M, error) {
	if r.Source == "" { r.Source = "api" }
	return one(s.c.do("POST", "/v1/events", r))
}

func (s *EventsService) Query(severity, typ, fingerprint string, hours, limit int) ([]M, error) {
	params := []string{}
	if severity != "" { params = append(params, "severity", severity) }
	if typ != "" { params = append(params, "type", typ) }
	if fingerprint != "" { params = append(params, "fingerprint", fingerprint) }
	if hours > 0 { params = append(params, "hours", fmt.Sprintf("%d", hours)) }
	if limit > 0 { params = append(params, "limit", fmt.Sprintf("%d", limit)) }
	data, err := s.c.get("/v1/events", params...)
	if err != nil { return nil, err }
	var r struct{ Items []M `json:"items"` }
	json.Unmarshal(data, &r)
	return r.Items, nil
}

// clean drops keys whose value is nil.
//
// Optional fields are passed as nil, but sending an explicit null is not
// the same as omitting the field — some endpoints reject it with a 422.
func clean(m M) M {
	out := M{}
	for k, v := range m {
		if v == nil {
			continue
		}
		if s, ok := v.(string); ok && s == "" {
			continue
		}
		out[k] = v
	}
	return out
}

// InterpretersService — live interpreter: two people each speak their own
// language. Covers both delivery paths: a phone call that dials both sides,
// and a browser session with a shareable link. None of this existed in this
// SDK before 5 Oct 2026 even though the interpreter is a headline product.
type InterpretersService struct{ c *Client }

func (s *InterpretersService) List(limit, offset int) ([]M, error) {
	return list(s.c.get("/v1/interpreters", "limit", fmt.Sprintf("%d", limit), "offset", fmt.Sprintf("%d", offset)))
}
func (s *InterpretersService) Get(id string) (M, error) {
	return one(s.c.do("GET", "/v1/interpreters/"+id, nil))
}

// Create — name and both languages are required. Without phoneNumber only
// the browser path works; the phone path needs a line to dial out from.
func (s *InterpretersService) Create(data M) (M, error) {
	return one(s.c.do("POST", "/v1/interpreters", clean(data)))
}
func (s *InterpretersService) Update(id string, data M) (M, error) {
	return one(s.c.do("PATCH", "/v1/interpreters/"+id, clean(data)))
}
func (s *InterpretersService) Delete(id string) error {
	_, err := s.c.do("DELETE", "/v1/interpreters/"+id, nil)
	return err
}

// Calls — sessions run through any interpreter on the account.
func (s *InterpretersService) Calls(limit, offset int) ([]M, error) {
	return list(s.c.get("/v1/interpreters/calls", "limit", fmt.Sprintf("%d", limit), "offset", fmt.Sprintf("%d", offset)))
}

// Call places a REAL interpreted phone call — both sides are dialled and
// the call is charged. targetNumber is required.
func (s *InterpretersService) Call(id, targetNumber, initiatorNumber string) (M, error) {
	return one(s.c.do("POST", "/v1/interpreters/"+id+"/call", clean(M{
		"target_number": targetNumber, "initiator_number": initiatorNumber,
	})))
}

// EnableWeb opens or closes the browser session link. Places no phone call.
func (s *InterpretersService) EnableWeb(id string, enabled bool, requiresPasscode interface{}) (M, error) {
	return one(s.c.do("POST", "/v1/interpreters/"+id+"/web", clean(M{
		"enabled": enabled, "requires_passcode": requiresPasscode,
	})))
}
func (s *InterpretersService) EndWeb(id string) (M, error) {
	return one(s.c.do("POST", "/v1/interpreters/"+id+"/web/end", nil))
}

// WebStatus — who is connected to the browser session right now.
func (s *InterpretersService) WebStatus(id string) (M, error) {
	return one(s.c.do("GET", "/v1/interpreters/"+id+"/web/live", nil))
}

// CreatePasscode mints a single-use join passcode for the browser session.
func (s *InterpretersService) CreatePasscode(id string) (M, error) {
	return one(s.c.do("POST", "/v1/interpreters/"+id+"/web/passcodes", nil))
}

// NumbersService — searching and buying phone numbers. Distinct from
// PhoneNumbers, which manages numbers you already own.
type NumbersService struct{ c *Client }

func (s *NumbersService) AllowedCountries() ([]M, error) {
	return list(s.c.get("/v1/numbers/allowed-countries"))
}

// Search lists purchasable numbers. Does NOT buy anything.
func (s *NumbersService) Search(country string, params ...string) ([]M, error) {
	return list(s.c.get("/v1/numbers/search", append([]string{"country", country}, params...)...))
}

// Purchase BUYS the number: the balance is charged and monthly rent starts.
func (s *NumbersService) Purchase(data M) (M, error) {
	return one(s.c.do("POST", "/v1/numbers/purchase", clean(data)))
}

// Checkout returns a payment link; it does not charge by itself.
func (s *NumbersService) Checkout(data M) (M, error) {
	return one(s.c.do("POST", "/v1/numbers/checkout", clean(data)))
}

// Release PERMANENTLY releases the number; it leaves the account.
func (s *NumbersService) Release(number string) error {
	_, err := s.c.do("DELETE", "/v1/numbers/"+number, nil)
	return err
}

type SmsService struct{ c *Client }

func (s *SmsService) Send(to, text, from string) (M, error) {
	return one(s.c.do("POST", "/v1/sms", clean(M{"to": to, "text": text, "from": from})))
}
func (s *SmsService) List(params ...string) ([]M, error) {
	return list(s.c.get("/v1/sms", params...))
}

// ExtensionService — Chrome extension, live translation of a browser tab.
// Device-scoped: the extension links a browser to the account, then runs
// sessions against that link.
type ExtensionService struct{ c *Client }

func (s *ExtensionService) Config(deviceID string) (M, error) {
	return one(s.c.get("/v1/ext/config", "device_id", deviceID))
}

// Usage — minutes used and remaining for the account.
func (s *ExtensionService) Usage() (M, error) { return one(s.c.get("/v1/ext/usage")) }

func (s *ExtensionService) LinkRequest(deviceID, email string) (M, error) {
	return one(s.c.do("POST", "/v1/ext/link/request", M{"device_id": deviceID, "email": email}))
}
func (s *ExtensionService) LinkConfirm(token string) (M, error) {
	return one(s.c.do("POST", "/v1/ext/link/confirm", M{"token": token}))
}
func (s *ExtensionService) LinkStatus(deviceID string) (M, error) {
	return one(s.c.do("POST", "/v1/ext/link/status", M{"device_id": deviceID}))
}
func (s *ExtensionService) LinkAttach(deviceID string) (M, error) {
	return one(s.c.do("POST", "/v1/ext/link/attach", M{"device_id": deviceID}))
}
func (s *ExtensionService) LinkDetach(deviceID string) (M, error) {
	return one(s.c.do("POST", "/v1/ext/link/detach", M{"device_id": deviceID}))
}
func (s *ExtensionService) SessionStart(data M) (M, error) {
	return one(s.c.do("POST", "/v1/ext/session/start", clean(data)))
}
func (s *ExtensionService) SessionHeartbeat(deviceID, sessionID string, elapsedSeconds int) (M, error) {
	return one(s.c.do("POST", "/v1/ext/session/heartbeat", M{
		"device_id": deviceID, "session_id": sessionID, "elapsed_seconds": elapsedSeconds,
	}))
}
func (s *ExtensionService) SessionEnd(deviceID, sessionID, reason string) (M, error) {
	return one(s.c.do("POST", "/v1/ext/session/end", clean(M{
		"device_id": deviceID, "session_id": sessionID, "reason": reason,
	})))
}
