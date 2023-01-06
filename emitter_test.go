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
	"github.com/xybor-x/xyerror"
	"github.com/xybor-x/xylog"
	"github.com/xybor-x/xylog/test"
)

func TestNewStreamEmitterWithNil(t *testing.T) {
	xycond.ExpectPanic(xyerror.AssertionError, func() {
		xylog.NewDefaultEmitter(nil)
	}).Test(t)
}

func TestStreamEmitterEmit(t *testing.T) {
	test.WithStreamEmitter(t, func(e *xylog.StreamEmitter, w *test.MockWriter) {
		var msg = test.GetRandomMessage()
		e.Emit([]byte(msg))
		xycond.ExpectIn(msg, w.Captured).Test(t)
	})
}

func TestStreamEmitterEmitError(t *testing.T) {
	test.WithStreamEmitter(t, func(e *xylog.StreamEmitter, w *test.MockWriter) {
		var msg = test.GetRandomMessage()
		w.Error = true
		e.Emit([]byte(msg))
		xycond.ExpectEmpty(w.Captured).Test(t)
	})
}
