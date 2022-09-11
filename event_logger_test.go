package xylog_test

import (
	"fmt"
	"testing"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xylog"
	"github.com/xybor-x/xylog/test"
)

func TestEventLogger(t *testing.T) {
	test.WithLogger(t, func(logger *xylog.Logger, w *test.MockWriter) {
		logger.SetLevel(xylog.DEBUG)
		var msg = test.GetRandomMessage()
		var elogger = logger.Event(msg)
		var tests = []func(){
			elogger.Debug, elogger.Info, elogger.Warn, elogger.Warning,
			elogger.Error, elogger.Fatal, elogger.Critical,
		}

		for i := range tests {
			w.Reset()
			tests[i]()
			xycond.ExpectIn(fmt.Sprintf("event=\"%s\"", msg), w.Captured).
				Test(t)
		}
	})
}
