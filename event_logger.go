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

package xylog

import "sync"

var eventLoggerPool = sync.Pool{
	New: func() any {
		return &EventLogger{}
	},
}

// EventLogger is a logger wrapper supporting to compose logging message with
// key-value structure.
type EventLogger struct {
	fields []field
	lg     *Logger
}

// Field adds a key-value pair to logging message.
func (e *EventLogger) Field(key string, value any) *EventLogger {
	e.fields = append(e.fields, makeField(key, value))
	return e
}

// Debug calls Log with DEBUG level.
func (e *EventLogger) Debug() {
	defer e.free()
	if e.lg.isEnabledFor(DEBUG) {
		e.lg.log(DEBUG, e.fields...)
	}
}

// Info calls Log with INFO level.
func (e *EventLogger) Info() {
	defer e.free()
	if e.lg.isEnabledFor(INFO) {
		e.lg.log(INFO, e.fields...)
	}
}

// Warn calls Log with WARN level.
func (e *EventLogger) Warn() {
	defer e.free()
	if e.lg.isEnabledFor(WARN) {
		e.lg.log(WARN, e.fields...)
	}
}

// Warning calls Log with WARNING level.
func (e *EventLogger) Warning() {
	defer e.free()
	if e.lg.isEnabledFor(WARNING) {
		e.lg.log(WARNING, e.fields...)
	}
}

// Error calls Log with ERROR level.
func (e *EventLogger) Error() {
	defer e.free()
	if e.lg.isEnabledFor(ERROR) {
		e.lg.log(ERROR, e.fields...)
	}
}

// Fatal calls Log with FATAL level.
func (e *EventLogger) Fatal() {
	defer e.free()
	if e.lg.isEnabledFor(FATAL) {
		e.lg.log(FATAL, e.fields...)
	}
}

// Critical calls Log with CRITICAL level.
func (e *EventLogger) Critical() {
	defer e.free()
	if e.lg.isEnabledFor(CRITICAL) {
		e.lg.log(CRITICAL, e.fields...)
	}
}

// Log logs with a custom level.
func (e *EventLogger) Log(level int) {
	defer e.free()
	level = CheckLevel(level)
	if e.lg.isEnabledFor(level) {
		e.lg.log(level, e.fields...)
	}
}

// free clears fields in EventLogger and puts it to pool.
func (e *EventLogger) free() {
	e.fields = e.fields[:0]
	eventLoggerPool.Put(e)
}
