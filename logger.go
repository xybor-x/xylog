package xylog

import (
	"fmt"
	"runtime"
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

	fullname         string
	children         map[string]*Logger
	parent           *Logger
	level            int
	handlers         map[*Handler]any
	lock             xylock.RWLock
	cache            map[int]bool
	extra            map[string]any
	persistentFields []field
}

// newlogger creates a new logger with a name and parent. The fullname of logger
// will be concatenated by the parent's fullname. This logger will not be
// automatically added to logger hierarchy. The returned logger has no child,
// no handler, and NOTSET level.
func newlogger(name string, parent *Logger) *Logger {
	var c = parent
	if c != nil && c != rootLogger {
		name = c.fullname + "." + name
	}

	return &Logger{
		f:        newfilterer(),
		fullname: name,
		children: make(map[string]*Logger),
		parent:   parent,
		level:    NOTSET,
		handlers: make(map[*Handler]any),
		lock:     xylock.RWLock{},
		cache:    make(map[int]bool),
		extra:    make(map[string]any),
	}
}

// SetLevel sets the new logging level. It also clears logging level cache of
// all loggers in program.
func (lg *Logger) SetLevel(level int) {
	lg.lock.WLockFunc(func() { lg.level = checkLevel(level) })
	rootLogger.clearCache()
}

// AddHandler adds a new handler.
func (lg *Logger) AddHandler(h *Handler) {
	xycond.AssertNotNil(h)
	lg.lock.WLockFunc(func() {
		if _, ok := lg.handlers[h]; !ok {
			lg.handlers[h] = nil
		}
	})
}

// RemoveHandler removes an existed handler.
func (lg *Logger) RemoveHandler(h *Handler) {
	lg.lock.WLockFunc(func() {
		delete(lg.handlers, h)
	})
}

// AddFilter adds a specified filter.
func (lg *Logger) AddFilter(f Filter) {
	lg.f.AddFilter(f)
}

// RemoveFilter removes an existed filter.
func (lg *Logger) RemoveFilter(f Filter) {
	lg.f.RemoveFilter(f)
}

// AddExtra adds a custom macro to logging format.
func (lg *Logger) AddExtra(key string, value any) {
	lg.extra[key] = value
}

// AddField adds a fixed key-value pair to all logging messages when using the
// EventLogger.
func (lg *Logger) AddField(key string, value any) {
	lg.persistentFields = append(lg.persistentFields, field{key, value})
}

// filter checks all filters in filterer, if there is any failed filter, it will
// returns false.
func (lg *Logger) filter(r LogRecord) bool {
	return lg.f.filter(r)
}

// Debug logs default formatting objects with DEBUG level.
func (lg *Logger) Debug(a ...any) {
	if lg.isEnabledFor(DEBUG) {
		lg.log(DEBUG, fmt.Sprint(a...))
	}
}

// Debugf logs a formatting message with DEBUG level.
func (lg *Logger) Debugf(s string, a ...any) {
	if lg.isEnabledFor(DEBUG) {
		lg.log(DEBUG, fmt.Sprintf(s, a...))
	}
}

// Info logs default formatting objects with INFO level.
func (lg *Logger) Info(a ...any) {
	if lg.isEnabledFor(INFO) {
		lg.log(INFO, fmt.Sprint(a...))
	}
}

// Infof logs a formatting message with INFO level.
func (lg *Logger) Infof(s string, a ...any) {
	if lg.isEnabledFor(INFO) {
		lg.log(INFO, fmt.Sprintf(s, a...))
	}
}

// Warn logs default formatting objects with WARN level.
func (lg *Logger) Warn(a ...any) {
	if lg.isEnabledFor(WARN) {
		lg.log(WARN, fmt.Sprint(a...))
	}
}

// Warnf logs a formatting message with WARN level.
func (lg *Logger) Warnf(s string, a ...any) {
	if lg.isEnabledFor(WARN) {
		lg.log(WARN, fmt.Sprintf(s, a...))
	}
}

// Warning logs default formatting objects with WARNING level.
func (lg *Logger) Warning(a ...any) {
	if lg.isEnabledFor(WARNING) {
		lg.log(WARNING, fmt.Sprint(a...))
	}
}

// Warningf logs a formatting message with WARNING level.
func (lg *Logger) Warningf(s string, a ...any) {
	if lg.isEnabledFor(WARNING) {
		lg.log(WARNING, fmt.Sprintf(s, a...))
	}
}

// Error logs default formatting objects with ERROR level.
func (lg *Logger) Error(a ...any) {
	if lg.isEnabledFor(ERROR) {
		lg.log(ERROR, fmt.Sprint(a...))
	}
}

// Errorf logs a formatting message with ERROR level.
func (lg *Logger) Errorf(s string, a ...any) {
	if lg.isEnabledFor(ERROR) {
		lg.log(ERROR, fmt.Sprintf(s, a...))
	}
}

// Fatal logs default formatting objects with FATAL level.
func (lg *Logger) Fatal(a ...any) {
	if lg.isEnabledFor(FATAL) {
		lg.log(FATAL, fmt.Sprint(a...))
	}
}

// Fatalf logs a formatting message with FATAL level.
func (lg *Logger) Fatalf(s string, a ...any) {
	if lg.isEnabledFor(FATAL) {
		lg.log(FATAL, fmt.Sprintf(s, a...))
	}
}

// Critical logs default formatting objects with CRITICAL level.
func (lg *Logger) Critical(a ...any) {
	if lg.isEnabledFor(CRITICAL) {
		lg.log(CRITICAL, fmt.Sprint(a...))
	}
}

// Criticalf logs a formatting message with CRITICAL level.
func (lg *Logger) Criticalf(s string, a ...any) {
	if lg.isEnabledFor(CRITICAL) {
		lg.log(CRITICAL, fmt.Sprintf(s, a...))
	}
}

// Log logs default formatting objects with a custom level.
func (lg *Logger) Log(level int, a ...any) {
	level = checkLevel(level)
	if lg.isEnabledFor(level) {
		lg.log(level, fmt.Sprint(a...))
	}
}

// Logf logs a formatting message with a custom level.
func (lg *Logger) Logf(level int, s string, a ...any) {
	level = checkLevel(level)
	if lg.isEnabledFor(level) {
		lg.log(level, fmt.Sprintf(s, a...))
	}
}

// Event creates an EventLogger which logs key-value pairs.
func (lg *Logger) Event(e string) *EventLogger {
	var elogger = &EventLogger{
		lg:     lg,
		fields: make([]field, 0, 5),
		isJSON: false,
	}
	elogger.fields = append(elogger.fields, lg.persistentFields...)
	return elogger.Field("event", e)
}

// Stack logs the stack trace.
func (lg *Logger) Stack(level int) {
	var s = string(debug.Stack())
	var lines = strings.Split(s, "\n")

	for i := range lines {
		lg.log(level, lines[i])
	}
}

// log is a low-level logging method which creates a LogRecord and then calls
// all the handlers of this logger to handle the record.
func (lg *Logger) log(level int, msg any) {
	var pc, filename, lineno, ok = runtime.Caller(skipCall)
	if !ok {
		filename = "unknown"
		lineno = -1
	}

	var record = makeRecord(lg.fullname, level, filename, lineno, msg, pc,
		lg.extra)

	lg.handle(record)
}

// handle calls the handlers for the specified record.
func (lg *Logger) handle(record LogRecord) {
	if lg.filter(record) {
		lg.callHandlers(record)
	}
}

// callHandlers passes a record to all relevant handlers.
//
// Loop through all handlers for this logger and its parents in the logger
// hierarchy. If no handler was found, output a one-off error message to
// os.Stderr.
func (lg *Logger) callHandlers(record LogRecord) {
	var c = lg
	var found = 0
	for c != nil {
		for h := range c.handlers {
			h.handle(record)
			found++
		}
		c = c.parent
	}

	if found == 0 {
		lastHandler.handle(record)
	}
}

// isEnabledFor checks if a logging level should be logged in this logger.
func (lg *Logger) isEnabledFor(level int) bool {
	var isEnabled, isCached bool
	var _ = lg.lock.RLockFunc(func() any {
		isEnabled, isCached = lg.cache[level]
		return nil
	})

	if !isCached {
		isEnabled = level >= lg.getEffectiveLevel()
		lg.lock.WLockFunc(func() { lg.cache[level] = isEnabled })
	}
	return isEnabled
}

// getEffectiveLevel gets the effective level for this logger.
//
// Loop through this logger and its parents in the logger hierarchy,
// looking for a non-zero logging level. Return the first one found.
func (lg *Logger) getEffectiveLevel() int {
	var level = lg.lock.RLockFunc(func() any { return lg.level }).(int)
	if level == NOTSET && lg.parent != nil {
		return lg.parent.getEffectiveLevel()
	}
	return level
}

// clearCache clears logging level cache of this logger and all its children.
func (lg *Logger) clearCache() {
	lg.lock.WLockFunc(func() {
		for k := range lg.cache {
			delete(lg.cache, k)
		}
	})
	for i := range lg.children {
		lg.children[i].clearCache()
	}
}

// prefixMessage adds a prefix to origin message if the prefix is not empty.
func prefixMessage(prefix, msg string) string {
	if prefix != "" {
		msg = fmt.Sprintf("%s %s", prefix, msg)
	}
	return msg
}
