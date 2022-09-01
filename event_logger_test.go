package xylog_test

import (
	"testing"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xylog"
)

func TestNewEventLogger(t *testing.T) {
	var logger = xylog.GetLogger(t.Name())
	logger.SetLevel(xylog.DEBUG)
	xycond.ExpectNotPanic(func() {
		var elogger = logger.Event("event")
		elogger.Field("foo", "bar")
		elogger.Debug()
		elogger.Info()
		elogger.Warn()
		elogger.Warning()
		elogger.Error()
		elogger.Fatal()
		elogger.Critical()
		elogger.Log(validCustomLevels[1])
	}).Test(t)
}
