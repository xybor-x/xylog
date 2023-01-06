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

package test

import (
	"io"
	"testing"

	"github.com/xybor-x/xylog"
	"github.com/xybor-x/xylog/encoding"
)

// WithLogger allows using a Logger created with a MockWriter quickly.
func WithLogger(t *testing.T, f func(logger *xylog.Logger, w *MockWriter)) {
	var writer = &MockWriter{}
	var emitter = xylog.NewStreamEmitter(writer)
	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)

	var logger = xylog.GetLogger(t.Name())
	// Sometimes the testing runs multiple times and this logger will add more
	// than one handler.
	logger.RemoveAllHandlers()
	logger.AddHandler(handler)

	f(logger, writer)
}

// WithHandler allows using a Handler with MockWriter.
func WithHandler(t *testing.T, f func(h *xylog.Handler, w *MockWriter)) {
	var writer = &MockWriter{}
	var emitter = xylog.NewStreamEmitter(writer)
	var handler = xylog.GetHandler(t.Name())
	handler.AddEmitter(emitter)

	f(handler, writer)
}

// WithStreamEmitter allows using a StreamEmitter created with a MockWriter
// quickly.
func WithStreamEmitter(t *testing.T, f func(e *xylog.StreamEmitter, w *MockWriter)) {
	var writer = &MockWriter{}
	var emitter = xylog.NewStreamEmitter(writer)
	f(emitter, writer)
}

// WithBenchLogger allows using a Logger whose output is io.Discard.
func WithBenchLogger(b *testing.B, f func(logger *xylog.Logger)) {
	var emitter = xylog.NewBufferEmitter(io.Discard, 4096)
	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)
	handler.SetEncoding(encoding.NewJSONEncoding())

	var logger = xylog.GetLogger(b.Name())
	// Sometimes the testing runs multiple times and this logger will add more
	// than one handler.
	logger.RemoveAllHandlers()
	logger.AddHandler(handler)

	f(logger)
}
