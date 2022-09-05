package xylog_test

import (
	"testing"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xylog"
)

func TestGetHandlerWithEmptyName(t *testing.T) {
	var handlerA = xylog.GetHandler("")
	var handlerB = xylog.GetHandler("")
	xycond.ExpectNotEqual(handlerA, handlerB).Test(t)
}

func TestHandlerSetFormatter(t *testing.T) {
	var handler = xylog.GetHandler("")
	xycond.ExpectNotPanic(func() {
		handler.SetFormatter(xylog.NewTextFormatter(""))
	}).Test(t)
}

func TestHandlerFilter(t *testing.T) {
	var expectedFilter = &NameFilter{}
	var handler = xylog.GetHandler("")
	xycond.ExpectNotPanic(func() {
		handler.AddFilter(expectedFilter)
	}).Test(t)
}

func TestHandlerFilterLog(t *testing.T) {
	var expectedMessage = "foo foo"
	var tests = []struct {
		name       string
		filterName string
	}{
		{t.Name() + "1", t.Name() + "1"},
		{t.Name() + "2", "foobar"},
	}

	for i := range tests {
		var handler = xylog.GetHandler("")
		handler.AddEmitter(&CapturedEmitter{})
		handler.AddFilter(&NameFilter{tests[i].filterName})
		handler.SetLevel(xylog.DEBUG)

		var logger = xylog.GetLogger(tests[i].name)
		logger.SetLevel(xylog.DEBUG)
		logger.AddHandler(handler)

		capturedOutput = ""
		logger.Warningf(expectedMessage)
		if tests[i].filterName != tests[i].name {
			xycond.ExpectEmpty(capturedOutput).Test(t)
		} else {
			xycond.ExpectEqual(capturedOutput, expectedMessage).Test(t)
		}
	}
}

func TestHandlerLevel(t *testing.T) {
	var expectedMessage = "foo foo"
	var loggerLevel = xylog.INFO
	var tests = []struct {
		name  string
		level int
	}{
		{t.Name() + "1", xylog.DEBUG},
		{t.Name() + "2", xylog.ERROR},
	}

	for i := range tests {
		var handler = xylog.GetHandler(tests[i].name)
		handler.AddEmitter(&CapturedEmitter{})
		handler.SetLevel(tests[i].level)

		var logger = xylog.GetLogger(tests[i].name)
		logger.SetLevel(xylog.DEBUG)
		logger.AddHandler(handler)
		capturedOutput = ""
		logger.Logf(loggerLevel, expectedMessage)
		if loggerLevel < tests[i].level {
			xycond.ExpectEmpty(capturedOutput).Test(t)
		} else {
			xycond.ExpectEqual(capturedOutput, expectedMessage).Test(t)
		}
	}
}
