package test

import (
	"encoding/json"
	"strings"
	"testing"
)

var client = NewClient(
	"https",
	"gurujsonrpc.appspot.com",
	0,
	"guru").
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
