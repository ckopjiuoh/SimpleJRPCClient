package test

import (
	"github.com/ckopjiuoh/SimpleJRPCClient/jrpc-client"
	"strings"
	"testing"
)

var client = jrpc_client.NewClient(
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

	if !strings.Contains(r.Result.(string), "Hello guru") {
		t.Error("Not equal!")
	}
}
