package xylog

// AddHandler adds a new handler to root logger.
func AddHandler(h *Handler) {
	rootLogger.AddHandler(h)
}

// RemoveHandler removes an existed handler from root logger.
func RemoveHandler(h *Handler) {
	rootLogger.RemoveHandler(h)
}

// AddFilter adds a specified filter to root logger.
func AddFilter(f Filter) {
	rootLogger.AddFilter(f)
}

// RemoveFilter removes an existed filter from root logger.
func RemoveFilter(f Filter) {
	rootLogger.RemoveFilter(f)
}

// SetLevel sets the new logging level for root logger.
func SetLevel(level int) {
	rootLogger.SetLevel(level)
}

// Debug logs default formatting objects with DEBUG level by root logger.
func Debug(a ...any) {
	rootLogger.Debug(a...)
}

// Debugf logs a formatting message with DEBUG level by root logger.
func Debugf(msg string, a ...any) {
	rootLogger.Debugf(msg, a...)
}

// Info logs default formatting objects with INFO level by root logger.
func Info(a ...any) {
	rootLogger.Info(a...)
}

// Infof logs a formatting message with INFO level by root logger.
func Infof(msg string, a ...any) {
	rootLogger.Infof(msg, a...)
}

// Warn logs default formatting objects with WARN level by root logger.
func Warn(a ...any) {
	rootLogger.Warn(a...)
}

// Warnf logs a formatting message with WARN level by root logger.
func Warnf(msg string, a ...any) {
	rootLogger.Warnf(msg, a...)
}

// Warning logs default formatting objects with WARNING level by root logger.
func Warning(a ...any) {
	rootLogger.Warning(a...)
}

// Warningf logs a formatting message with WARNING level by root logger.
func Warningf(msg string, a ...any) {
	rootLogger.Warningf(msg, a...)
}

// Error logs default formatting objects with ERROR level by root logger.
func Error(a ...any) {
	rootLogger.Error(a...)
}

// Errorf logs a formatting message with ERROR level by root logger.
func Errorf(msg string, a ...any) {
	rootLogger.Errorf(msg, a...)
}

// Fatal logs default formatting objects with FATAL level by root logger.
func Fatal(a ...any) {
	rootLogger.Fatal(a...)
}

// Fatalf logs a formatting message with FATAL level by root logger.
func Fatalf(msg string, a ...any) {
	rootLogger.Fatalf(msg, a...)
}

// Critical logs default formatting objects with CRITICAL level by root logger.
func Critical(a ...any) {
	rootLogger.Critical(a...)
}

// Criticalf logs a formatting message with DEBUG level by root logger.
func Criticalf(msg string, a ...any) {
	rootLogger.Criticalf(msg, a...)
}

// Log logs default formatting objects with a custom level by root logger.
func Log(level int, a ...any) {
	rootLogger.Log(level, a...)
}

// Logf logs a formatting message with a custom level by root logger.
func Logf(level int, msg string, a ...any) {
	rootLogger.Logf(level, msg, a...)
}
