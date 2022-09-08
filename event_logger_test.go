package xylog_test

import (
	"fmt"
	"testing"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xylog"
)

func TestEventLogger(t *testing.T) {
	withLogger(t, func(logger *xylog.Logger, w *mockWriter) {
		logger.SetLevel(xylog.DEBUG)
		var msg = getRandomMessage()
		var elogger = logger.Event(msg)
		var tests = []func(){
			elogger.Debug, elogger.Info, elogger.Warn, elogger.Warning,
			elogger.Error, elogger.Fatal, elogger.Critical,
		}

		for i := range tests {
			w.Reset()
			tests[i]()
			xycond.ExpectIn(w.Captured, fmt.Sprintf("event=\"%s\"", msg)).
				Test(t)
		}
	})
}

func TestEventLoggerPair(t *testing.T) {
	withLogger(t, func(logger *xylog.Logger, w *mockWriter) {
		logger.SetLevel(xylog.DEBUG)
		logger.Event("something").Field("foo", "space message").Field("bar", 1).
			Field("buzzz", true).Log(xylog.DEBUG)
		xycond.ExpectEqual(w.Captured,
			"event=something foo=\"space message\" bar=1 buzzz=true").Test(t)
	})
}

func TestEventLoggerJSON(t *testing.T) {
	withLogger(t, func(logger *xylog.Logger, w *mockWriter) {
		logger.SetLevel(xylog.DEBUG)
		logger.Event("something").Field("foo", "bar").Field("bar", 1).
			Field("buzzz", true).JSON().Info()
		xycond.ExpectEqual(w.Captured,
			`{"bar":1,"buzzz":true,"event":"something","foo":"bar"}`).Test(t)
	})
}
