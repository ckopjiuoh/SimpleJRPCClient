package logging

import (
	"net/http"
	"time"

	"github.com/motemen/go-nuts/roundtime"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

type contextKey struct {
	name string
}

var ContextKeyRequestStart = &contextKey{"RequestStart"}

var RPCLogRequest = func(req *http.Request) {
	Log.WithFields(logrus.Fields{
		"method": req.Method,
		"url":    req.URL.Host,
	}).Info("RPC Request")
}

var RPCLogResponse = func(resp *http.Response) {
	ctx := resp.Request.Context()
	if start, ok := ctx.Value(ContextKeyRequestStart).(time.Time); ok {
		Log.WithFields(logrus.Fields{
			"Status": resp.StatusCode,
			"url":    resp.Request.URL.Host,
			"time":   roundtime.Duration(time.Now().Sub(start), 1)}).Info("RPC Response")

	} else {
		Log.WithFields(logrus.Fields{
			"Status": resp.StatusCode,
			"url":    resp.Request.URL.Host,
		}).Info("RPC Response ")
	}
}
