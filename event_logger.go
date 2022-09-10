package xylog

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
	if e.lg.isEnabledFor(DEBUG) {
		e.lg.log(DEBUG, e.fields...)
	}
}

// Info calls Log with INFO level.
func (e *EventLogger) Info() {
	if e.lg.isEnabledFor(INFO) {
		e.lg.log(INFO, e.fields...)
	}
}

// Warn calls Log with WARN level.
func (e *EventLogger) Warn() {
	if e.lg.isEnabledFor(WARN) {
		e.lg.log(WARN, e.fields...)
	}
}

// Warning calls Log with WARNING level.
func (e *EventLogger) Warning() {
	if e.lg.isEnabledFor(WARNING) {
		e.lg.log(WARNING, e.fields...)
	}
}

// Error calls Log with ERROR level.
func (e *EventLogger) Error() {
	if e.lg.isEnabledFor(ERROR) {
		e.lg.log(ERROR, e.fields...)
	}
}

// Fatal calls Log with FATAL level.
func (e *EventLogger) Fatal() {
	if e.lg.isEnabledFor(FATAL) {
		e.lg.log(FATAL, e.fields...)
	}
}

// Critical calls Log with CRITICAL level.
func (e *EventLogger) Critical() {
	if e.lg.isEnabledFor(CRITICAL) {
		e.lg.log(CRITICAL, e.fields...)
	}
}

// Log logs with a custom level.
func (e *EventLogger) Log(level int) {
	level = CheckLevel(level)
	if e.lg.isEnabledFor(level) {
		e.lg.log(level, e.fields...)
	}
}
