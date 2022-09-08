package xylog_test

import (
	"testing"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xylog"
)

func TestGetLogger(t *testing.T) {
	var names = []string{"", "foo", "foo.bar"}
	for i := range names {
		var logger1 = xylog.GetLogger(names[i])
		var logger2 = xylog.GetLogger(names[i])
		xycond.ExpectEqual(logger1, logger2).Test(t)
	}
}

func TestLoggerHandler(t *testing.T) {
	var handler = xylog.GetHandler("")
	var logger = xylog.GetLogger(t.Name())
	xycond.ExpectNotPanic(func() { logger.AddHandler(handler) }).Test(t)
}

func TestLoggerAddHandlerNil(t *testing.T) {
	var logger = xylog.GetLogger(t.Name())
	xycond.ExpectPanic(func() { logger.AddHandler(nil) }).Test(t)
}

func TestLoggerLogMethods(t *testing.T) {
	var fixedMsg = getRandomMessage()
	withLogger(t, func(logger *xylog.Logger, w *mockWriter) {
		var tests = []struct {
			methodf func(string, ...any)
			method  func(...any)
		}{
			{logger.Debugf, logger.Debug},
			{logger.Infof, logger.Info},
			{logger.Warnf, logger.Warn},
			{logger.Warningf, logger.Warning},
			{logger.Errorf, logger.Error},
			{logger.Fatalf, logger.Fatal},
			{logger.Criticalf, logger.Critical},
		}

		logger.SetLevel(xylog.DEBUG)
		for i := range tests {
			w.Reset()
			var msg = getRandomMessage()
			tests[i].method(msg)
			xycond.ExpectIn(w.Captured, msg).Test(t)

			w.Reset()
			var msgf = getRandomMessage()
			tests[i].methodf(msgf)
			xycond.ExpectIn(w.Captured, msgf).Test(t)
		}
		w.Reset()
		logger.Log(xylog.DEBUG, fixedMsg)
		xycond.ExpectIn(w.Captured, fixedMsg).Test(t)
		w.Reset()
		logger.Logf(xylog.DEBUG, fixedMsg)
		xycond.ExpectIn(w.Captured, fixedMsg).Test(t)

		logger.SetLevel(xylog.NOTLOG)
		for i := range tests {
			w.Reset()
			var msg = getRandomMessage()
			tests[i].method(msg)
			xycond.ExpectNotIn(w.Captured, msg).Test(t)

			w.Reset()
			var msgf = getRandomMessage()
			tests[i].methodf(msgf)
			xycond.ExpectNotIn(w.Captured, msgf).Test(t)
		}
		w.Reset()
		logger.Log(xylog.DEBUG, fixedMsg)
		xycond.ExpectNotIn(w.Captured, fixedMsg).Test(t)
		w.Reset()
		logger.Logf(xylog.DEBUG, fixedMsg)
		xycond.ExpectNotIn(w.Captured, fixedMsg).Test(t)
	})
}

func TestLoggerCallHandlerHierarchy(t *testing.T) {
	withLogger(t, func(logger *xylog.Logger, w *mockWriter) {
		var child = xylog.GetLogger(t.Name() + ".main")
		logger.SetLevel(xylog.INFO)

		var msg = getRandomMessage()
		child.Log(xylog.WARN, msg)
		xycond.ExpectIn(w.Captured, msg).Test(t)

		msg = getRandomMessage()
		child.Log(xylog.DEBUG, msg)
		xycond.ExpectNotIn(w.Captured, msg).Test(t)
	})
}

func TestLoggerStack(t *testing.T) {
	withLogger(t, func(logger *xylog.Logger, w *mockWriter) {
		logger.SetLevel(xylog.DEBUG)
		logger.Stack(xylog.DEBUG)
		xycond.ExpectIn(w.Captured, "xylog.(*Logger).Stack").Test(t)
	})
}

type namefilter struct {
	name string
}

func (f *namefilter) Filter(r xylog.LogRecord) bool {
	return f.name == r.Name
}

func TestLoggerFilterLog(t *testing.T) {
	withLogger(t, func(logger *xylog.Logger, w *mockWriter) {
		for _, h := range logger.Handlers() {
			h.AddFilter(&namefilter{t.Name() + ".main"})
		}
		var main = xylog.GetLogger(t.Name() + ".main")
		var other = xylog.GetLogger(t.Name() + ".other")

		var msg = getRandomMessage()
		main.Error(msg)
		xycond.ExpectIn(w.Captured, msg).Test(t)

		w.Reset()
		other.Error(msg)
		xycond.ExpectNotIn(w.Captured, msg).Test(t)
	})
}

func TestLoggerAddFields(t *testing.T) {
	withLogger(t, func(logger *xylog.Logger, w *mockWriter) {
		logger.AddField("foo", "bar")

		logger.Event("test").Error()
		xycond.ExpectIn(w.Captured, "foo=bar").Test(t)

		w.Reset()
		logger.Error("test")
		xycond.ExpectNotIn(w.Captured, "foo=bar").Test(t)
	})
}

func TestLoggerLogInvalidJSONMessage(t *testing.T) {
	withLogger(t, func(logger *xylog.Logger, w *mockWriter) {
		logger.AddField("foo", "bar")

		logger.Event("test").Field("func", func() {}).JSON().Error()
		xycond.ExpectIn(w.Captured,
			"An error occurred while formatting the message").Test(t)
	})
}

func TestLoggerHandlers(t *testing.T) {
	var handler = xylog.GetHandler("")
	var logger = xylog.GetLogger(t.Name())
	logger.AddHandler(handler)
	xycond.ExpectEqual(len(logger.Handlers()), 1).Test(t)
	xycond.ExpectEqual(logger.Handlers()[0], handler).Test(t)

	logger.RemoveHandler(handler)
	xycond.ExpectEqual(len(logger.Handlers()), 0).Test(t)
}

func TestLoggerFilters(t *testing.T) {
	var filter = &LoggerNameFilter{"foo"}
	var logger = xylog.GetLogger(t.Name())
	logger.AddFilter(filter)
	xycond.ExpectEqual(len(logger.Filters()), 1).Test(t)
	xycond.ExpectEqual(logger.Filters()[0], filter).Test(t)

	logger.RemoveFilter(filter)
	xycond.ExpectEqual(len(logger.Filters()), 0).Test(t)
}
