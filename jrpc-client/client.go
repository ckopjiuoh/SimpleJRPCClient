package jrpc_client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"errors"

	"github.com/motemen/go-loghttp"
	"github.com/motemen/go-nuts/roundtime"
	log "github.com/sirupsen/logrus"
)

type contextKey struct {
	name string
}

var ContextKeyRequestStart = &contextKey{"RequestStart"}

var RPCLogRequest = func(req *http.Request) {
	log.Info("--> %s %s %s", req.Method, req.URL, req.Body)
}

var RPCLogResponse = func(resp *http.Response) {
	ctx := resp.Request.Context()
	if start, ok := ctx.Value(ContextKeyRequestStart).(time.Time); ok {
		log.Info("<-- %d %s (%s)", resp.StatusCode, resp.Request.URL, roundtime.Duration(time.Now().Sub(start), 2))
	} else {
		log.Info("<-- %d %s", resp.StatusCode, resp.Request.URL)
	}
}

type JRPCClient struct {
	Addr   string
	Client *http.Client
}

func NewClient(addr string) *JRPCClient {
	return &JRPCClient{
		Addr: addr,
		Client: &http.Client{
			Transport: &loghttp.Transport{
				LogRequest:  RPCLogRequest,
				LogResponse: RPCLogResponse,
			},
		},
	}
}

func (c JRPCClient) Call(method string, params interface{}) (*RPCResponse, error) {
	client := c.Client
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(&RPCRequest{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	})
	resp, err := client.Post(c.Addr, "application/json; charset=utf-8", b)
	if err != nil {
		log.Println(err.Error())
	}
	defer resp.Body.Close()
	r := new(RPCResponse)
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		log.Error("error decoding response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Error("syntax error at byte offset %d", e.Offset)
		}
		log.Error("RPC response: %q", r)
		return nil, errors.New("Not a RPC Response!")
	}
	log.Info(r)
	return r, nil
}

type RPCRequest struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

type RPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   RPCError    `json:"error"`
	ID      int         `json:"id"`
}

type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
