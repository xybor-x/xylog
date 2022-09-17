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
	"os"
	"testing"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xylog"
	"github.com/xybor-x/xylog/encoding"
	"github.com/xybor-x/xylog/test"
)

func TestGetHandler(t *testing.T) {
	var handlerA = xylog.GetHandler(t.Name())
	var handlerB = xylog.GetHandler(t.Name())
	xycond.ExpectEqual(handlerA, handlerB).Test(t)
	xycond.ExpectEqual(handlerA.Name(), t.Name()).Test(t)
	xycond.ExpectEqual(handlerB.Name(), t.Name()).Test(t)
}

func TestGetHandlerDiff(t *testing.T) {
	var handlerA = xylog.GetHandler(t.Name() + "1")
	var handlerB = xylog.GetHandler(t.Name() + "2")
	xycond.ExpectNotEqual(handlerA, handlerB).Test(t)
}

func TestGetHandlerWithEmptyName(t *testing.T) {
	var handlerA = xylog.GetHandler("")
	var handlerB = xylog.GetHandler("")
	xycond.ExpectNotEqual(handlerA, handlerB).Test(t)
}

func TestHandlerSetLevel(t *testing.T) {
	var handler = xylog.GetHandler("")
	handler.SetLevel(xylog.WARN)
	xycond.ExpectEqual(handler.Level(), xylog.WARN).Test(t)
}

func TestHandlerFileters(t *testing.T) {
	var filter = &test.LoggerNameFilter{Name: "foo"}
	var handler = xylog.GetHandler("")
	handler.AddFilter(filter)
	xycond.ExpectEqual(len(handler.Filters()), 1).Test(t)
	xycond.ExpectEqual(handler.Filters()[0], filter).Test(t)

	handler.RemoveFilter(filter)
	xycond.ExpectEqual(len(handler.Filters()), 0).Test(t)
}

func TestHandlerEmitters(t *testing.T) {
	var emitter = xylog.NewStreamEmitter(os.Stdout)
	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)
	xycond.ExpectEqual(len(handler.Emitters()), 1).Test(t)
	xycond.ExpectEqual(handler.Emitters()[0], emitter).Test(t)

	handler.RemoveEmitter(emitter)
	xycond.ExpectEqual(len(handler.Emitters()), 0).Test(t)
}

func TestHandlerSetEncodingAfterAddField(t *testing.T) {
	test.WithHandler(t, func(h *xylog.Handler, w *test.MockWriter) {
		h.AddField("foo", "bar")
		h.AddField("name", "value")
		h.SetEncoding(encoding.NewJSONEncoding())
		h.Handle(xylog.LogRecord{})

		xycond.ExpectIn(`{"foo":"bar","name":"value"}`, w.Captured).Test(t)
	})
}

func TestHandlerInvalidMacro(t *testing.T) {
	test.WithHandler(t, func(h *xylog.Handler, w *test.MockWriter) {
		h.AddMacro("foo", "unknown")
		h.Handle(xylog.LogRecord{})

		xycond.ExpectIn("An error occurred while formatting the message",
			w.Captured).Test(t)
	})
}

func TestHandlerFullMacro(t *testing.T) {

	test.WithHandler(t, func(h *xylog.Handler, w *test.MockWriter) {
		test.AddFullMacros(h)

		h.Handle(test.FullRecord)

		xycond.ExpectEqual("asctime=ASCTIME created=1 filename=FILENAME "+
			"funcname=FUNCNAME levelname=LEVELNAME levelno=2 lineno=3 "+
			"module=MODULE msecs=4 pathname=PATHNAME process=5 "+
			"relativeCreated=6", w.Captured).Test(t)
	})
}
