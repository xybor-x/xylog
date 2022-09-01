package xylog_test

import (
	"os"
	"testing"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xyerror"
	"github.com/xybor-x/xylog"
)

type ErrorWriter struct{}

func (ew *ErrorWriter) Write(p []byte) (n int, err error) {
	return 0, xyerror.Error.New("unknown")
}

func (ew *ErrorWriter) Close() error {
	return nil
}

func TestNewStreamEmitterWithNil(t *testing.T) {
	xycond.ExpectNotPanic(func() {
		xylog.NewStreamEmitter(nil)
	}).Test(t)
}

func TestStreamEmitterEmit(t *testing.T) {
	var emitter = xylog.NewStreamEmitter(os.Stderr)
	xycond.ExpectNotPanic(func() {
		emitter.Emit(xylog.LogRecord{})
	}).Test(t)
}

func TestStreamEmitterEmitError(t *testing.T) {
	var emitter = xylog.NewStreamEmitter(&ErrorWriter{})
	xycond.ExpectPanic(func() {
		emitter.Emit(xylog.LogRecord{})
	}).Test(t)
}

func TestFileEmitter(t *testing.T) {
	var emitter = xylog.NewFileEmitter("a.log")
	xycond.ExpectNotPanic(func() {
		emitter.Emit(xylog.LogRecord{})
	}).Test(t)
}

func TestSizeRotatingFileEmitter(t *testing.T) {
	var emitter = xylog.NewSizeRotatingFileEmitter("a.log", 100, 1)
	xycond.ExpectNotPanic(func() {
		emitter.Emit(xylog.LogRecord{})
	}).Test(t)
}

func TestTimeRotatingFileEmitter(t *testing.T) {
	var emitter = xylog.NewTimeRotatingFileEmitter("a.log", 100, 1)
	xycond.ExpectNotPanic(func() {
		emitter.Emit(xylog.LogRecord{})
	}).Test(t)
}
