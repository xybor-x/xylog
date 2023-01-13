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

		w.Reset()
		logger.Event(msg).Debug()
		xycond.ExpectIn(fmt.Sprintf("event=\"%s\"", msg), w.Captured).
			Test(t)

		w.Reset()
		logger.Event(msg).Info()
		xycond.ExpectIn(fmt.Sprintf("event=\"%s\"", msg), w.Captured).
			Test(t)

		w.Reset()
		logger.Event(msg).Warn()
		xycond.ExpectIn(fmt.Sprintf("event=\"%s\"", msg), w.Captured).
			Test(t)

		w.Reset()
		logger.Event(msg).Warning()
		xycond.ExpectIn(fmt.Sprintf("event=\"%s\"", msg), w.Captured).
			Test(t)

		w.Reset()
		logger.Event(msg).Error()
		xycond.ExpectIn(fmt.Sprintf("event=\"%s\"", msg), w.Captured).
			Test(t)

		w.Reset()
		xycond.ExpectPanic(nil, logger.Event(msg).Panic).Test(t)
		xycond.ExpectIn(fmt.Sprintf("event=\"%s\"", msg), w.Captured).
			Test(t)

		w.Reset()
		logger.Event(msg).Critical()
		xycond.ExpectIn(fmt.Sprintf("event=\"%s\"", msg), w.Captured).
			Test(t)
	})
}

func TestEventLoggerLog(t *testing.T) {
	test.WithLogger(t, func(logger *xylog.Logger, w *test.MockWriter) {
		logger.Event("foo").Log(xylog.CRITICAL)
		xycond.ExpectIn("event=foo", w.Captured).Test(t)
	})
}
