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
		xylog.NewStreamEmitter(nil)
	}).Test(t)
}

func TestStreamEmitterEmit(t *testing.T) {
	test.WithStreamEmitter(t, func(e *xylog.StreamEmitter, w *test.MockWriter) {
		var msg = test.GetRandomMessage()
		e.Emit(msg)
		xycond.ExpectIn(msg, w.Captured).Test(t)
	})
}

func TestStreamEmitterEmitError(t *testing.T) {
	test.WithStreamEmitter(t, func(e *xylog.StreamEmitter, w *test.MockWriter) {
		w.Error = true
		e.Emit(test.GetRandomMessage())
		xycond.ExpectEmpty(w.Captured).Test(t)
	})
}

func TestStreamEmitterClose(t *testing.T) {
	test.WithStreamEmitter(t, func(e *xylog.StreamEmitter, w *test.MockWriter) {
		e.Close()
		e.Emit(test.GetRandomMessage())
		xycond.ExpectEmpty(w.Captured).Test(t)
	})
}

func TestStreamEmitterCloseTwice(t *testing.T) {
	test.WithStreamEmitter(t, func(e *xylog.StreamEmitter, w *test.MockWriter) {
		e.Close()
		e.Close()
	})
}
