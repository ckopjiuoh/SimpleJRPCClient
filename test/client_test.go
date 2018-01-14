package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ckopjiuoh/SimpleJRPCClient/jrpc-client"
	"github.com/ckopjiuoh/SimpleJRPCClient/model"
)

func TestClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respBody := model.RPCResponse{
			JSONRPC: "2.0",
			Result:  "1",
			ID:      1,
			Error: nil,
		}

		b, _ := json.Marshal(respBody)
		w.Write(b)
	}))

	defer ts.Close()

	var client = jrpc_client.NewClient(
		"",
		ts.URL,
		0,
		"").
		WithPost()

	r, _ := client.Call()

	if r.Result.(string) != "1" {
		t.Error("Not equal!")
	}
}
