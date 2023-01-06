// Copyright (c) 2022 xybor-x
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package xylog_test

import (
	"testing"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xylog"
	"github.com/xybor-x/xylog/test"
)

func TestGetLogger(t *testing.T) {
	var names = []string{"", "foo", "foo.bar"}
	for i := range names {
		var logger1 = xylog.GetLogger(names[i])
		var logger2 = xylog.GetLogger(names[i])
		xycond.ExpectEqual(logger1, logger2).Test(t)
	}
}

func TestLoggerLogMethods(t *testing.T) {
	var fixedMsg = test.GetRandomMessage()
	test.WithLogger(t, func(logger *xylog.Logger, w *test.MockWriter) {
		var tests = []struct {
			methodf func(string, ...any)
			method  func(string)
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
			var msg = test.GetRandomMessage()
			tests[i].method(msg)
			xycond.ExpectIn(msg, w.Captured).Test(t)

			w.Reset()
			var msgf = test.GetRandomMessage()
			tests[i].methodf(msgf)
			xycond.ExpectIn(msgf, w.Captured).Test(t)
		}

		w.Reset()
		logger.Log(xylog.DEBUG, fixedMsg)
		xycond.ExpectIn(fixedMsg, w.Captured).Test(t)
		w.Reset()
		logger.Logf(xylog.DEBUG, fixedMsg)
		xycond.ExpectIn(fixedMsg, w.Captured).Test(t)

		logger.SetLevel(xylog.NOTLOG)
		for i := range tests {
			w.Reset()
			var msg = test.GetRandomMessage()
			tests[i].method(msg)
			xycond.ExpectNotIn(msg, w.Captured).Test(t)

			w.Reset()
			var msgf = test.GetRandomMessage()
			tests[i].methodf(msgf)
			xycond.ExpectNotIn(msgf, w.Captured).Test(t)
		}

		w.Reset()
		logger.Log(xylog.DEBUG, fixedMsg)
		xycond.ExpectNotIn(fixedMsg, w.Captured).Test(t)
		w.Reset()
		logger.Logf(xylog.DEBUG, fixedMsg)
		xycond.ExpectNotIn(fixedMsg, w.Captured).Test(t)
	})
}

func TestLoggerCallHandlerHierarchy(t *testing.T) {
	test.WithLogger(t, func(logger *xylog.Logger, w *test.MockWriter) {
		var child = xylog.GetLogger(t.Name() + ".main")
		logger.SetLevel(xylog.INFO)

		var msg = test.GetRandomMessage()
		child.Log(xylog.WARN, msg)
		xycond.ExpectIn(msg, w.Captured).Test(t)

		w.Reset()
		msg = test.GetRandomMessage()
		child.Log(xylog.DEBUG, msg)
		xycond.ExpectNotIn(msg, w.Captured).Test(t)
	})
}

func TestLoggerStack(t *testing.T) {
	test.WithLogger(t, func(logger *xylog.Logger, w *test.MockWriter) {
		logger.SetLevel(xylog.DEBUG)
		logger.Stack(xylog.DEBUG)
		xycond.ExpectIn("xylog.(*Logger).Stack", w.Captured).Test(t)
	})
}

func TestLoggerFilterLog(t *testing.T) {
	test.WithLogger(t, func(logger *xylog.Logger, w *test.MockWriter) {
		for _, h := range logger.Handlers() {
			h.AddFilter(&test.LoggerNameFilter{Name: t.Name() + ".main"})
		}
		var main = xylog.GetLogger(t.Name() + ".main")
		var other = xylog.GetLogger(t.Name() + ".other")

		var msg = test.GetRandomMessage()
		main.Error(msg)
		xycond.ExpectIn(msg, w.Captured).Test(t)

		w.Reset()
		other.Error(msg)
		xycond.ExpectNotIn(msg, w.Captured).Test(t)
	})
}

func TestLoggerAddField(t *testing.T) {
	test.WithLogger(t, func(logger *xylog.Logger, w *test.MockWriter) {
		logger.AddField("custom", "this is a custom field")
		logger.Event("test").Error()
		xycond.ExpectIn(`event=test custom="this is a custom field"`,
			w.Captured).Test(t)
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
	var filter = &test.LoggerNameFilter{Name: "foo"}
	var logger = xylog.GetLogger(t.Name())
	logger.AddFilter(filter)

	xycond.ExpectEqual(len(logger.Filters()), 1).Test(t)
	xycond.ExpectEqual(logger.Filters()[0], filter).Test(t)

	logger.RemoveFilter(filter)
	xycond.ExpectEqual(len(logger.Filters()), 0).Test(t)
}

func TestLoggerFindCaller(t *testing.T) {
	xylog.SetFindCaller(true)
	defer xylog.SetFindCaller(false)
	test.WithLogger(t, func(logger *xylog.Logger, w *test.MockWriter) {
		var handler = logger.Handlers()[0]
		handler.AddMacro("lineno", "lineno")
		handler.AddMacro("module", "module")
		handler.AddMacro("funcname", "funcname")

		logger.Error("foo")

		xycond.ExpectIn("lineno=182", w.Captured).Test(t)
		xycond.ExpectIn(
			"module=github.com/xybor-x/xylog_test", w.Captured).Test(t)
		xycond.ExpectIn(
			"funcname=TestLoggerFindCaller.func1", w.Captured).Test(t)
	})
}

func TestLoggerFlushParent(t *testing.T) {
	var writer = &test.MockWriter{}
	var emitter = xylog.NewDefaultEmitter(writer)
	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)

	var logger = xylog.GetLogger(t.Name())
	logger.AddHandler(handler)

	var childLogger = xylog.GetLogger(t.Name() + ".child")
	childLogger.Error("foo")
	xylog.Flush()

	xycond.ExpectIn("foo", writer.Captured).Test(t)
}
