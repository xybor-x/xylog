package xylog_test

import (
	"testing"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xyerror"
	"github.com/xybor-x/xylog"
)

func TestNewStreamEmitterWithNil(t *testing.T) {
	xycond.ExpectPanic(xyerror.AssertionError, func() {
		xylog.NewStreamEmitter(nil)
	}).Test(t)
}

func TestStreamEmitterEmit(t *testing.T) {
	withStreamEmitter(t, func(e *xylog.StreamEmitter, w *mockWriter) {
		var msg = getRandomMessage()
		e.Emit(msg)
		xycond.ExpectIn(msg, w.Captured).Test(t)
	})
}

func TestStreamEmitterEmitError(t *testing.T) {
	withStreamEmitter(t, func(e *xylog.StreamEmitter, w *mockWriter) {
		w.Error = true
		e.Emit(getRandomMessage())
		xycond.ExpectEmpty(w.Captured).Test(t)
	})
}

func TestStreamEmitterClose(t *testing.T) {
	withStreamEmitter(t, func(e *xylog.StreamEmitter, w *mockWriter) {
		e.Close()
		e.Emit(getRandomMessage())
		xycond.ExpectEmpty(w.Captured).Test(t)
	})
}

func TestStreamEmitterCloseTwice(t *testing.T) {
	withStreamEmitter(t, func(e *xylog.StreamEmitter, w *mockWriter) {
		e.Close()
		e.Close()
	})
}
