package xylog_test

import (
	"testing"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xylog"
)

func TestNewHandlerWithEmptyName(t *testing.T) {
	var handlerA = xylog.NewHandler("", xylog.StdoutEmitter)
	var handlerB = xylog.NewHandler("", xylog.StdoutEmitter)
	xycond.ExpectNotEqual(handlerA, handlerB).Test(t)
}

func TestHandlerSetFormatter(t *testing.T) {
	var handler = xylog.NewHandler(t.Name(), xylog.StdoutEmitter)
	xycond.ExpectNotPanic(func() {
		handler.SetFormatter(xylog.NewTextFormatter(""))
	}).Test(t)
}

func TestHandlerFilter(t *testing.T) {
	var expectedFilter = &NameFilter{}
	var handler = xylog.NewHandler(t.Name(), xylog.StdoutEmitter)
	handler.AddFilter(expectedFilter)
	xycond.ExpectNotPanic(func() {
		handler.RemoveFilter(expectedFilter)
	}).Test(t)
}

func TestHandlerFilterLog(t *testing.T) {
	var expectedMessage = "foo foo"
	var tests = []struct {
		handlerName string
		filterName  string
	}{
		{t.Name() + "1", t.Name()},
		{t.Name() + "2", "foobar"},
	}

	for i := range tests {
		var handler = xylog.NewHandler(tests[i].handlerName, &CapturedEmitter{})
		handler.AddFilter(&NameFilter{tests[i].filterName})
		handler.SetLevel(xylog.DEBUG)

		var logger = xylog.GetLogger(t.Name())
		logger.SetLevel(xylog.DEBUG)
		logger.AddHandler(handler)

		capturedOutput = ""
		logger.Warningf(expectedMessage)
		if tests[i].filterName != t.Name() {
			xycond.ExpectEmpty(capturedOutput).Test(t)
		} else {
			xycond.ExpectEqual(capturedOutput, expectedMessage).Test(t)
		}
		logger.RemoveHandler(handler)
	}
}

func TestHandlerLevel(t *testing.T) {
	var expectedMessage = "foo foo"
	var loggerLevel = xylog.INFO
	var tests = []struct {
		handlerName string
		level       int
	}{
		{t.Name() + "1", xylog.DEBUG},
		{t.Name() + "2", xylog.ERROR},
	}

	for i := range tests {
		var handler = xylog.NewHandler(tests[i].handlerName, &CapturedEmitter{})
		handler.SetLevel(tests[i].level)

		var logger = xylog.GetLogger(t.Name())
		logger.SetLevel(xylog.DEBUG)
		logger.AddHandler(handler)
		capturedOutput = ""
		logger.Logf(loggerLevel, expectedMessage)
		if loggerLevel < tests[i].level {
			xycond.ExpectEmpty(capturedOutput).Test(t)
		} else {
			xycond.ExpectEqual(capturedOutput, expectedMessage).Test(t)
		}
		logger.RemoveHandler(handler)
	}
}
