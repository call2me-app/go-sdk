package call2me

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func captureCall(t *testing.T, fn func(c *Client) (M, error)) map[string]interface{} {
	t.Helper()
	var body map[string]interface{}
	var method, path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		json.NewDecoder(r.Body).Decode(&body)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"call_id":"call_1"}`))
	}))
	defer srv.Close()
	c := New("sk_test")
	c.BaseURL = srv.URL
	out, err := fn(c)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if method != "POST" || path != "/v1/calls" {
		t.Fatalf("wrong request: %s %s", method, path)
	}
	if out["call_id"] != "call_1" {
		t.Fatalf("bad response: %v", out)
	}
	return body
}

func TestCreateCallSendsDynamicVariablesAndBeginMessage(t *testing.T) {
	body := captureCall(t, func(c *Client) (M, error) {
		return c.Calls.CreateCall(CallCreateParams{
			AgentID:          "agent_1",
			ToNumber:         "+15550001111",
			DynamicVariables: map[string]string{"name": "Ada", "order_id": "A-42"},
			BeginMessage:     "Hi {{name}}, calling about order {{order_id}}.",
		})
	})
	want := map[string]interface{}{"name": "Ada", "order_id": "A-42"}
	if !reflect.DeepEqual(body["dynamic_variables"], want) {
		t.Fatalf("dynamic_variables: %v", body["dynamic_variables"])
	}
	if body["begin_message"] != "Hi {{name}}, calling about order {{order_id}}." {
		t.Fatalf("begin_message: %v", body["begin_message"])
	}
}

func TestCreateCallOmitsUnsetFields(t *testing.T) {
	body := captureCall(t, func(c *Client) (M, error) {
		return c.Calls.CreateCall(CallCreateParams{AgentID: "agent_1", ToNumber: "+15550001111"})
	})
	want := map[string]interface{}{"agent_id": "agent_1", "to_number": "+15550001111"}
	if !reflect.DeepEqual(body, want) {
		t.Fatalf("unexpected body: %v", body)
	}
}

func TestCreateMapDropsNilAndEmpty(t *testing.T) {
	body := captureCall(t, func(c *Client) (M, error) {
		return c.Calls.Create(M{
			"agent_id":          "agent_1",
			"to_number":         "+15550001111",
			"dynamic_variables": map[string]string{"name": "Ada"},
			"begin_message":     "",
			"topic":             nil,
		})
	})
	if _, ok := body["begin_message"]; ok {
		t.Fatalf("empty begin_message sent: %v", body)
	}
	if _, ok := body["topic"]; ok {
		t.Fatalf("nil topic sent: %v", body)
	}
	if !reflect.DeepEqual(body["dynamic_variables"], map[string]interface{}{"name": "Ada"}) {
		t.Fatalf("dynamic_variables: %v", body["dynamic_variables"])
	}
}
