package xylog_test

import (
	"testing"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xylog"
)

func TestEventLogger(t *testing.T) {
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

func TestEventLoggerPair(t *testing.T) {
	var handler = xylog.NewHandler("", &CapturedEmitter{})
	handler.SetLevel(xylog.DEBUG)
	var logger = xylog.GetLogger(t.Name())
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)

	var elogger = logger.Event("something").
		Field("foo", "bar").Field("bar", 1).Field("buzzz", true)

	capturedOutput = ""
	elogger.Debug()
	xycond.ExpectEqual(
		capturedOutput, "event=something foo=bar bar=1 buzzz=true").Test(t)
}

func TestEventLoggerJSON(t *testing.T) {
	var handler = xylog.NewHandler("", &CapturedEmitter{})
	handler.SetLevel(xylog.DEBUG)
	var logger = xylog.GetLogger(t.Name())
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)

	var elogger = logger.Event("something").
		Field("foo", "bar").Field("bar", 1).Field("buzzz", true).JSON()

	capturedOutput = ""
	elogger.Debug()
	xycond.ExpectEqual(capturedOutput,
		`{"bar":1,"buzzz":true,"event":"something","foo":"bar"}`).Test(t)
}
