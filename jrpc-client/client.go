package jrpc_client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"errors"

	"github.com/motemen/go-loghttp"
	"github.com/motemen/go-nuts/roundtime"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

type contextKey struct {
	name string
}

var ContextKeyRequestStart = &contextKey{"RequestStart"}

var RPCLogRequest = func(req *http.Request) {
	 go logrus.WithFields(logrus.Fields{
			 "method": req.Method,
			 "url":    req.URL.Host})

}

var RPCLogResponse = func(resp *http.Response) {
	ctx := resp.Request.Context()
	if start, ok := ctx.Value(ContextKeyRequestStart).(time.Time); ok {
		Log.WithFields(logrus.Fields{
			"Status": resp.StatusCode,
			"url":    resp.Request.URL.Host,
			"time": roundtime.Duration(time.Now().Sub(start), 1)}).Info("RPC Response ")

	} else {
		Log.WithFields(logrus.Fields{
			"Status": resp.StatusCode,
			"url":    resp.Request.URL.Host,
		}).Info("RPC Response ")
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
	resp, err := client.Post(c.Addr, "application/json", b)
	if err != nil {
		Log.Error(err.Error())
	}
	defer resp.Body.Close()
	r := new(RPCResponse)
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		Log.Error("error decoding response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			Log.Error("syntax error at byte offset %d", e.Offset)
		}
		Log.Error("RPC response: %q", r)
		return nil, errors.New("Not a RPC Response!")
	}
	Log.WithFields(logrus.Fields{
		"Result": r.Result,
		"Error":    r.Error,
	}).Info("RPC Response info")
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
	Error   *RPCError   `json:"error"`
	ID      int         `json:"id"`
}

type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
