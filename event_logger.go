package xylog

import (
	"fmt"
	"strings"
)

// EventLogger is a logger wrapper supporting to compose logging message with
// key-value pair.
type EventLogger struct {
	msg string
	lg  *Logger
}

// Field adds a key-value pair to logging message.
func (e *EventLogger) Field(key string, value any) *EventLogger {
	var s = fmt.Sprint(value)
	if strings.Contains(s, " ") {
		s = fmt.Sprintf("\"%s\"", s)
	}
	e.msg = prefixMessage(e.msg, key+"="+s)
	return e
}

// Debug calls Log with DEBUG level.
func (e *EventLogger) Debug() {
	if e.lg.isEnabledFor(DEBUG) {
		e.lg.log(DEBUG, e.msg)
	}
}

// Info calls Log with INFO level.
func (e *EventLogger) Info() {
	if e.lg.isEnabledFor(INFO) {
		e.lg.log(INFO, e.msg)
	}
}

// Warn calls Log with WARN level.
func (e *EventLogger) Warn() {
	if e.lg.isEnabledFor(WARN) {
		e.lg.log(WARN, e.msg)
	}
}

// Warning calls Log with WARNING level.
func (e *EventLogger) Warning() {
	if e.lg.isEnabledFor(WARNING) {
		e.lg.log(WARNING, e.msg)
	}
}

// Error calls Log with ERROR level.
func (e *EventLogger) Error() {
	if e.lg.isEnabledFor(ERROR) {
		e.lg.log(ERROR, e.msg)
	}
}

// Fatal calls Log with FATAL level.
func (e *EventLogger) Fatal() {
	if e.lg.isEnabledFor(FATAL) {
		e.lg.log(FATAL, e.msg)
	}
}

// Critical calls Log with CRITICAL level.
func (e *EventLogger) Critical() {
	if e.lg.isEnabledFor(CRITICAL) {
		e.lg.log(CRITICAL, e.msg)
	}
}

// Log logs with a custom level.
func (e *EventLogger) Log(level int) {
	level = checkLevel(level)
	if e.lg.isEnabledFor(level) {
		e.lg.log(level, e.msg)
	}
}
