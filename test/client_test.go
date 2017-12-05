package test

import (
	"SimpleJRPCClient/jrpc-client"
	"encoding/json"
	"strings"
	"testing"
)

var client = jrpc_client.NewClient(
	"http",
	"container-1.test.automobile.ru",
	9990,
	"rpc/v1.0").
	WithPost().
	WithRPCMethod("guru.test")

func TestClient(t *testing.T) {

	r, _ := client.
		WithRPCParams([]string{"guru"}).
		Call()
	var s string
	json.Unmarshal(*r.Result, &s)
	if !strings.Contains(s, "Hello guru") {
		t.Error("Not equal!")
	}
}
