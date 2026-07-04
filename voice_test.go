package call2me

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVoiceSessionCreate(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"token":"t","url":"wss://x","room_name":"agent_a_1","session_limit_sec":3600}`))
	}))
	defer srv.Close()

	c := New("sk_test")
	c.BaseURL = srv.URL
	out, err := c.VoiceSessions.CreateSession("agent_abc", M{"name": "Ada"})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if gotMethod != "POST" || gotPath != "/v1/voice/sessions" {
		t.Fatalf("wrong request: %s %s", gotMethod, gotPath)
	}
	if gotBody["agent_id"] != "agent_abc" {
		t.Fatalf("agent_id not sent: %v", gotBody)
	}
	if out["room_name"] != "agent_a_1" {
		t.Fatalf("bad response: %v", out)
	}
}
