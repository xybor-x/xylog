package xylog

import (
	"fmt"
	"strings"
)

type field struct {
	key   string
	value any
}

// EventLogger is a logger wrapper supporting to compose logging message with
// key-value structure.
type EventLogger struct {
	fields []field
	lg     *Logger
	isJSON bool
}

// JSON parses fields as the json format.
func (e *EventLogger) JSON() *EventLogger {
	e.isJSON = true
	return e
}

// Field adds a key-value pair to logging message.
func (e *EventLogger) Field(key string, value any) *EventLogger {
	e.fields = append(e.fields, field{key, value})
	return e
}

// Debug calls Log with DEBUG level.
func (e *EventLogger) Debug() {
	if e.lg.isEnabledFor(DEBUG) {
		e.lg.log(DEBUG, e.createMessage())
	}
}

// Info calls Log with INFO level.
func (e *EventLogger) Info() {
	if e.lg.isEnabledFor(INFO) {
		e.lg.log(INFO, e.createMessage())
	}
}

// Warn calls Log with WARN level.
func (e *EventLogger) Warn() {
	if e.lg.isEnabledFor(WARN) {
		e.lg.log(WARN, e.createMessage())
	}
}

// Warning calls Log with WARNING level.
func (e *EventLogger) Warning() {
	if e.lg.isEnabledFor(WARNING) {
		e.lg.log(WARNING, e.createMessage())
	}
}

// Error calls Log with ERROR level.
func (e *EventLogger) Error() {
	if e.lg.isEnabledFor(ERROR) {
		e.lg.log(ERROR, e.createMessage())
	}
}

// Fatal calls Log with FATAL level.
func (e *EventLogger) Fatal() {
	if e.lg.isEnabledFor(FATAL) {
		e.lg.log(FATAL, e.createMessage())
	}
}

// Critical calls Log with CRITICAL level.
func (e *EventLogger) Critical() {
	if e.lg.isEnabledFor(CRITICAL) {
		e.lg.log(CRITICAL, e.createMessage())
	}
}

// Log logs with a custom level.
func (e *EventLogger) Log(level int) {
	level = checkLevel(level)
	if e.lg.isEnabledFor(level) {
		e.lg.log(level, e.createMessage())
	}
}

// createMessage returns the message based on fields.
func (e *EventLogger) createMessage() any {
	if e.isJSON {
		var data = make(map[string]any)
		for i := range e.fields {
			data[e.fields[i].key] = e.fields[i].value
		}
		return data
	}

	var msg string
	for i := range e.fields {
		var s = fmt.Sprint(e.fields[i].value)
		if strings.ContainsRune(s, ' ') {
			s = "\"" + s + "\""
		}
		msg = prefixMessage(msg, e.fields[i].key+"="+s)
	}
	return msg
}
