// Package xylog is a logging module based on the design of python logging.
package xylog

import (
	"os"
	"time"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xylock"
	"github.com/xybor-x/xylog/encoding"
)

func init() {
	rootLogger = newLogger("", nil)
	rootLogger.SetLevel(WARNING)

	var handler = GetHandler("xybor")
	handler.AddEmitter(NewStreamEmitter(os.Stderr))
	handler.SetEncoding(encoding.NewTextEncoding())
	handler.AddMacro("time", "asctime")
	handler.AddMacro("level", "levelname")
	handler.AddMacro("module", "name")

	var logger = GetLogger("xybor")
	logger.SetLevel(WARNING)
	logger.AddHandler(handler)
}

// Default levels, these can be replaced with any positive set of values having
// corresponding names. There is a pseudo-level, NOTSET, which is only really
// there as a lower limit for user-defined levels. Handlers and loggers are
// initialized with NOTSET so that they will log all messages, even at
// user-defined levels.
const (
	NOTLOG   = 1000
	CRITICAL = 50
	FATAL    = CRITICAL
	ERROR    = 40
	WARNING  = 30
	WARN     = WARNING
	INFO     = 20
	DEBUG    = 10
	NOTSET   = 0
)

// startTime is used as the base when calculating the relative time of events.
var startTime = time.Now().UnixMilli()

// globalLock is used to serialize access to shared data structures in this
// module.
var globalLock = &xylock.RWLock{}

// processid is always fixed and used to fill %(process) macro.
var processid = os.Getpid()

// rootLogger is the logger managing all loggers in program, it only should be
// used to set default handler or propagate level to all loggers.
var rootLogger *Logger

// timeLayout is the default time layout used to print asctime when logging.
var timeLayout = time.RFC3339Nano

// handlerManager is a map to search handler by name.
var handlerManager = make(map[string]*Handler)

// skipCall is the depth of Logger.log call in program, 3 by default. Increase
// this value if you want to wrap log methods of logger.
var skipCall = 3

// bufferSize is the expected size of buffer when creating a bufio.Writer from
// io.Writer.
var bufferSize = 4096

// findCaller allows finding caller information including filename, lineno,
// funcname, module.
var findCaller = false

var levelToName = map[int]string{
	CRITICAL: "CRITICAL",
	ERROR:    "ERROR",
	WARNING:  "WARNING",
	INFO:     "INFO",
	DEBUG:    "DEBUG",
	NOTSET:   "NOTSET",
}

// SetSkipCall sets the new skipCall value which dertermine the depth call of
// Logger.log method.
func SetSkipCall(skip int) {
	globalLock.WLockFunc(func() { skipCall = skip })
}

// SetTimeLayout sets the time layout to print asctime. It is time.RFC3339Nano
// by default.
func SetTimeLayout(layout string) {
	globalLock.WLockFunc(func() { timeLayout = layout })
}

// SetBufferSize sets the new expected size of buffer when creating the
// bufio.Writer.
func SetBufferSize(s int) {
	globalLock.WLockFunc(func() { bufferSize = s })
}

// SetFindCaller with true to find caller information including filename, line
// number, function name, and module.
func SetFindCaller(b bool) {
	globalLock.WLockFunc(func() { findCaller = b })
}

// AddLevel associates a log level with name. It can overwrite other log levels.
// Default log levels:
//   NOTSET       0
//   DEBUG        10
//   INFO         20
//   WARN/WARNING 30
//   ERROR/FATAL  40
//   CRITICAL     50
func AddLevel(level int, levelName string) {
	globalLock.WLockFunc(func() { levelToName[level] = levelName })
}

// CheckLevel validates if the given level is associated or not.
func CheckLevel(level int) int {
	globalLock.RLock()
	defer globalLock.RUnlock()
	xycond.AssertIn(level, levelToName)
	return level
}

// GetLevelName returns a name associated with the given level.
func GetLevelName(level int) string {
	globalLock.RLock()
	defer globalLock.RUnlock()
	return levelToName[level]
}

func makeField(key string, value any) field {
	return field{key: key, value: value}
}

type field struct {
	key   string
	value any
}

type macroField struct {
	key   string
	macro string
}
