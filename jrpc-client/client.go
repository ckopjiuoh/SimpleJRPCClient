package jrpc_client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"errors"

	. "ClientJSONRPC2/logging"
	"ClientJSONRPC2/model"
	"io/ioutil"

	"github.com/motemen/go-loghttp"
	"github.com/sirupsen/logrus"
)

type JRPCClient struct {
	Addr    string
	Method  string
	Headers map[string]string
	Body    model.RPCRequestBody

	httpClient *http.Client
}

func NewClient(addr string) *JRPCClient {
	Log = logrus.New()
	Log.Formatter = &logrus.JSONFormatter{}
	Log.SetLevel(logrus.WarnLevel)

	return &JRPCClient{
		Addr: addr,
		httpClient: &http.Client{
			Transport: &loghttp.Transport{
				LogRequest:  RPCLogRequest,
				LogResponse: RPCLogResponse,
			},
		},
		Method:  "Post",
		Headers: map[string]string{"Content-type": "application/json"},
	}
}

func (c *JRPCClient) WithHeaders(m map[string]string) *JRPCClient {
	c.Headers = m

	return c
}

func (c *JRPCClient) WithMethod(m string) *JRPCClient {
	c.Method = m

	return c
}

func (c *JRPCClient) WithBody(method string, params interface{}) *JRPCClient {
	var m = model.RPCRequestBody{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}
	c.Body = m
	return c
}

func (c *JRPCClient) Call() (*model.RPCResponse, error) {

	client := c.httpClient
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(&c.Body)

	req, err := http.NewRequest(c.Method, c.Addr, body)

	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	if err != nil {
		Log.Error(err.Error())
	}

	resp, err := client.Do(req)
	if err != nil {
		Log.Error(err.Error())
	}

	defer resp.Body.Close()
	r := new(model.RPCResponse)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&r)
	if err != nil {
		Log.Error(logrus.Fields{
			"Respnse decoding error": err.Error(),
		})
		if e, ok := err.(*json.SyntaxError); ok {
			Log.Error(logrus.Fields{
				"Syntax error": e.Error(),
				"Offset error": e.Offset,
			})
		}

		body, _ := ioutil.ReadAll(resp.Body)

		Log.Error("RPC response: %q", string(body))
		return nil, errors.New("Not a RPC Response!")
	}
	return r, nil
}
