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

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xylock"
)

// Logger represents a single logging channel. A "logging channel" indicates an
// area of an application. Exactly how an "area" is defined is up to the
// application developer. Since an application can have any number of areas,
// logging channels are identified by a unique string. Application areas can be
// nested (e.g. an area of "input processing" might include sub-areas "read CSV
// files", "read XLS files" and "read Gnumeric files"). To cater for this
// natural nesting, channel names are organized into a namespace hierarchy where
// levels are separated by periods. So in the instance given above, channel
// names might be "input" for the upper level, and "input.csv", "input.xls" and
// "input.gnu" for the sub-levels. There is no arbitrary limit to the depth of
// nesting.
type Logger struct {
	f *filterer

	name     string
	children map[string]*Logger
	parent   *Logger
	level    int
	handlers []*Handler
	lock     *xylock.RWLock
	cache    map[int]bool
	fields   []field
}

// GetLogger gets a logger with the specified name, creating it if it doesn't
// yet exist. This name is a dot-separated hierarchical name, such as "a",
// "a.b", "a.b.c", or similar.
//
// Leave name as empty string to get the root logger.
func GetLogger(name string) *Logger {
	if name == "" {
		return rootLogger
	}

	globalLock.Lock()
	defer globalLock.Unlock()

	var lg = rootLogger
	for _, part := range strings.Split(name, ".") {
		if _, ok := lg.children[part]; !ok {
			lg.children[part] = newLogger(part, lg)
		}
		lg = lg.children[part]
	}
	return lg
}

// Name returns the full name.
func (lg *Logger) Name() string {
	return lg.lock.RLockFunc(func() any { return lg.name }).(string)
}

// Parent returns the parent logger. If there is no parent, return nil instead.
func (lg *Logger) Parent() *Logger {
	lg.lock.RLock()
	defer lg.lock.RUnlock()
	return lg.parent
}

// Children returns direct children logger.
func (lg *Logger) Children() []*Logger {
	var children []*Logger
	lg.lock.RLock()
	defer lg.lock.RUnlock()

	for _, child := range lg.children {
		children = append(children, child)
	}

	return children
}

// Level returns the current logging level.
func (lg *Logger) Level() int {
	lg.lock.RLock()
	defer lg.lock.RUnlock()
	return lg.level
}

// SetLevel sets the new logging level.
func (lg *Logger) SetLevel(level int) {
	lg.lock.WLockFunc(func() { lg.level = level })
	rootLogger.clearCache()
}

// Handlers returns all current Handlers.
func (lg *Logger) Handlers() []*Handler {
	lg.lock.RLock()
	defer lg.lock.RUnlock()
	return lg.handlers
}

// AddHandler adds a new handler.
func (lg *Logger) AddHandler(h *Handler) {
	xycond.AssertNotNil(h)
	lg.lock.WLockFunc(func() { lg.handlers = append(lg.handlers, h) })
}

// RemoveHandler removes an existed Handler.
func (lg *Logger) RemoveHandler(h *Handler) {
	lg.lock.Lock()
	defer lg.lock.Unlock()

	for i := range lg.handlers {
		if lg.handlers[i] == h {
			lg.handlers = append(lg.handlers[:i], lg.handlers[i+1:]...)
		}
	}
}

// RemoveAllHandlers removes all existed Handlers.
func (lg *Logger) RemoveAllHandlers() {
	lg.lock.Lock()
	defer lg.lock.Unlock()
	lg.handlers = lg.handlers[:0]
}

// Filters returns all current Filters.
func (lg *Logger) Filters() []Filter {
	lg.lock.RLock()
	defer lg.lock.RUnlock()
	return lg.f.Filters()
}

// AddFilter adds a specified Filter.
func (lg *Logger) AddFilter(f Filter) {
	lg.lock.WLockFunc(func() { lg.f.AddFilter(f) })
}

// RemoveFilter removes an existed Filter.
func (lg *Logger) RemoveFilter(f Filter) {
	lg.lock.Lock()
	defer lg.lock.Unlock()

	lg.f.RemoveFilter(f)
}

// AddField adds a fixed field to all logging message of this logger.
func (lg *Logger) AddField(key string, value any) {
	lg.lock.WLockFunc(func() {
		lg.fields = append(lg.fields, makeField(key, value))
	})
}

// DeleteField deletes a fixed field from all logging message of this logger.
func (lg *Logger) DeleteField(key string) {
	lg.lock.WLockFunc(func() {
		idx := -1
		for iter, pair := range lg.fields {
			if pair.key == key {
				idx = iter
			}
		}
		if idx != -1 {
			lg.fields = append(lg.fields[:idx], lg.fields[idx+1:]...)
		}
	})
}

// UpdateField updates a fixed field to all logging message of this logger.
func (lg *Logger) UpdateField(key string, value string) {
	lg.lock.WLockFunc(func() {
		idx := -1
		for index, pair := range lg.fields {
			if pair.key == key {
				idx = index
			}
		}
		if idx != -1 {
			lg.fields[idx] = makeField(key, value)
		}
	})
}

// Debug logs default formatting objects with DEBUG level.
func (lg *Logger) Debug(s string) {
	if lg.isEnabledFor(DEBUG) {
		lg.log(DEBUG, makeField("messsage", s))
	}
}

// Debugf logs a formatting message with DEBUG level.
func (lg *Logger) Debugf(s string, a ...any) {
	if lg.isEnabledFor(DEBUG) {
		lg.log(DEBUG, makeField("messsage", fmt.Sprintf(s, a...)))
	}
}

// Info logs default formatting objects with INFO level.
func (lg *Logger) Info(s string) {
	if lg.isEnabledFor(INFO) {
		lg.log(INFO, makeField("messsage", s))
	}
}

// Infof logs a formatting message with INFO level.
func (lg *Logger) Infof(s string, a ...any) {
	if lg.isEnabledFor(INFO) {
		lg.log(INFO, makeField("messsage", fmt.Sprintf(s, a...)))
	}
}

// Warn logs default formatting objects with WARN level.
func (lg *Logger) Warn(s string) {
	if lg.isEnabledFor(WARN) {
		lg.log(WARN, makeField("messsage", s))
	}
}

// Warnf logs a formatting message with WARN level.
func (lg *Logger) Warnf(s string, a ...any) {
	if lg.isEnabledFor(WARN) {
		lg.log(WARN, makeField("messsage", fmt.Sprintf(s, a...)))
	}
}

// Warning logs default formatting objects with WARNING level.
func (lg *Logger) Warning(s string) {
	if lg.isEnabledFor(WARNING) {
		lg.log(WARNING, makeField("messsage", s))
	}
}

// Warningf logs a formatting message with WARNING level.
func (lg *Logger) Warningf(s string, a ...any) {
	if lg.isEnabledFor(WARNING) {
		lg.log(WARNING, makeField("messsage", fmt.Sprintf(s, a...)))
	}
}

// Error logs default formatting objects with ERROR level.
func (lg *Logger) Error(s string) {
	if lg.isEnabledFor(ERROR) {
		lg.log(ERROR, makeField("messsage", s))
	}
}

// Errorf logs a formatting message with ERROR level.
func (lg *Logger) Errorf(s string, a ...any) {
	if lg.isEnabledFor(ERROR) {
		lg.log(ERROR, makeField("messsage", fmt.Sprintf(s, a...)))
	}
}

// Critical logs default formatting objects with CRITICAL level.
func (lg *Logger) Critical(s string) {
	if lg.isEnabledFor(CRITICAL) {
		lg.log(CRITICAL, makeField("messsage", s))
	}
}

// Criticalf logs a formatting message with CRITICAL level.
func (lg *Logger) Criticalf(s string, a ...any) {
	if lg.isEnabledFor(CRITICAL) {
		lg.log(CRITICAL, makeField("messsage", fmt.Sprintf(s, a...)))
	}
}

// Fatal logs default formatting objects with CRITICAL level, then followed by a
// call to os.Exit(1).
func (lg *Logger) Fatal(s string) {
	lg.Critical(s)
	os.Exit(1)
}

// Fatalf logs a formatting message with CRITICAL level, then followed by a call
// to os.Exit(1).
func (lg *Logger) Fatalf(s string, a ...any) {
	lg.Criticalf(s, a...)
	os.Exit(1)
}

// Panic logs default formatting objects with CRITICAL level, then followed by a
// call to panic().
func (lg *Logger) Panic(s string) {
	lg.Critical(s)
	panic(s)
}

// Panicf logs a formatting message with CRITICAL level, then followed by a call
// to panic().
func (lg *Logger) Panicf(s string, a ...any) {
	lg.Criticalf(s, a...)
	panic(fmt.Sprintf(s, a...))
}

// Log logs default formatting objects with a custom level.
func (lg *Logger) Log(level int, s string) {
	level = CheckLevel(level)
	if lg.isEnabledFor(level) {
		lg.log(level, makeField("messsage", s))
	}
}

// Logf logs a formatting message with a custom level.
func (lg *Logger) Logf(level int, s string, a ...any) {
	level = CheckLevel(level)
	if lg.isEnabledFor(level) {
		lg.log(level, makeField("messsage", fmt.Sprintf(s, a...)))
	}
}

// Stack logs the stack trace.
func (lg *Logger) Stack(level int) {
	var s = string(debug.Stack())
	var lines = strings.Split(s, "\n")

	for i := range lines {
		lg.log(level, makeField("stack", strings.TrimSpace(lines[i])))
	}
}

// Event creates an EventLogger which logs key-value pairs.
func (lg *Logger) Event(e string) *EventLogger {
	var elogger = eventLoggerPool.Get().(*EventLogger)
	elogger.lg = lg
	elogger.Field("event", e)
	return elogger
}

// log is a low-level logging method which creates a LogRecord and then calls
// all the handlers of this logger to handle the record.
func (lg *Logger) log(level int, fields ...field) {
	fields = append(fields, lg.fields...)
	var record = makeRecord(lg.name, level, fields...)

	if lg.filter(record) {
		lg.callHandlers(record)
	}
}

// filter checks all filters in filterer, if there is any failed filter, it will
// returns false.
func (lg *Logger) filter(r LogRecord) bool {
	return lg.lock.RLockFunc(func() any { return lg.f.filter(r) }).(bool)
}

// callHandlers passes a record to all relevant handlers.
//
// Loop through all handlers for this logger and its parents in the logger
// hierarchy. If no handler was found, output a one-off error message to
// os.Stderr.
func (lg *Logger) callHandlers(record LogRecord) {
	var current = lg
	for current != nil {
		var handlers = current.Handlers()
		for i := range handlers {
			handlers[i].Handle(record)
		}
		current = current.Parent()
	}
}

// isEnabledFor checks if a logging level should be logged in this logger.
func (lg *Logger) isEnabledFor(level int) bool {
	lg.lock.RLock()
	var isEnabled, isCached = lg.cache[level]
	lg.lock.RUnlock()

	if !isCached {
		isEnabled = level >= lg.getEffectiveLevel()
		lg.lock.WLockFunc(func() { lg.cache[level] = isEnabled })
	}
	return isEnabled
}

// getEffectiveLevel gets the effective level for this logger.
//
// Loop through this logger and its parents in the logger hierarchy, looking for
// a non-zero logging level. Return the first one found.
func (lg *Logger) getEffectiveLevel() int {
	var level, parent = lg.Level(), lg.Parent()
	if level != NOTSET || parent == nil {
		return level
	}
	return parent.getEffectiveLevel()
}

// clearCache clears logging level cache of this logger and all its children.
func (lg *Logger) clearCache() {
	lg.lock.WLockFunc(func() {
		for k := range lg.cache {
			delete(lg.cache, k)
		}
	})
	for _, child := range lg.Children() {
		child.clearCache()
	}
}

// newLogger creates a Logger with a name and parent. The fullname of logger
// will be concatenated by the parent's fullname. This logger will not be
// automatically added to logger hierarchy. The returned logger has no child,
// no handler, and NOTSET level.
func newLogger(name string, parent *Logger) *Logger {
	var current = parent
	if current != nil && current != rootLogger {
		name = current.Name() + "." + name
	}

	return &Logger{
		f:        &filterer{},
		name:     name,
		children: make(map[string]*Logger),
		parent:   parent,
		level:    NOTSET,
		handlers: nil,
		lock:     &xylock.RWLock{},
		cache:    make(map[int]bool),
	}
}
