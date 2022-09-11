package test

import (
	"os"
	"testing"

	"github.com/xybor-x/xylog"
)

// WithLogger allows to use a Logger created with a MockWriter quickly.
func WithLogger(t *testing.T, f func(logger *xylog.Logger, w *MockWriter)) {
	xylog.SetBufferSize(1)
	defer xylog.SetBufferSize(4096)
	var writer = &MockWriter{}
	var emitter = xylog.NewStreamEmitter(writer)
	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)
	var logger = xylog.GetLogger(t.Name())
	logger.AddHandler(handler)

	f(logger, writer)
}

// WithStreamEmitter allows to use a StreamEmitter created with a MockWriter
// quickly.
func WithStreamEmitter(
	t *testing.T, f func(e *xylog.StreamEmitter, w *MockWriter),
) {
	xylog.SetBufferSize(1)
	defer xylog.SetBufferSize(4096)
	var writer = &MockWriter{}
	var emitter = xylog.NewStreamEmitter(writer)
	f(emitter, writer)
}

// WithBenchLogger allows to use a Logger whose output is devnull.
func WithBenchLogger(b *testing.B, f func(logger *xylog.Logger)) {
	var devnull, err = os.OpenFile(os.DevNull, os.O_WRONLY, 0666)
	if err != nil {
		b.Fail()
	}
	var emitter = xylog.NewStreamEmitter(devnull)
	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)

	var logger = xylog.GetLogger(b.Name())
	logger.AddHandler(handler)

	f(logger)
}
