// Package xylog is a logging module based on the design of python logging.
package xylog

import (
	"io/fs"
	"os"
	"strings"
	"time"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xylock"
)

func init() {
	rootLogger = newlogger("", nil)
	rootLogger.SetLevel(WARNING)
	handlerManager = make(map[string]*Handler)

	var handler = NewHandler("xybor", StderrEmitter)
	handler.SetLevel(WARNING)
	handler.SetFormatter(NewTextFormatter(
		"time=%(asctime)-30s " +
			"level=%(levelname)-8s " +
			"%(message)s",
	))

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
	CRITICAL = 50
	FATAL    = CRITICAL
	ERROR    = 40
	WARNING  = 30
	WARN     = WARNING
	INFO     = 20
	DEBUG    = 10
	NOTSET   = 0
)

// StdoutEmitter is a shortcut of NewStreamEmitter(os.Stdout)
var StdoutEmitter = NewStreamEmitter(os.Stdout)

// StderrEmitter is a shortcut of NewStreamEmitter(os.Stderr)
var StderrEmitter = NewStreamEmitter(os.Stderr)

// startTime is used as the base when calculating the relative time of events.
var startTime = time.Now().UnixMilli()

// lock is used to serialize access to shared data structures in this module.
var lock = xylock.RWLock{}

// processid is alway fixed and used to fill %(process) macro.
var processid = os.Getpid()

// rootLogger is the logger managing all loggers in program, it only should be
// used to set default handler or propagate level to all loggers.
var rootLogger *Logger

// timeLayout is the default time layout used to print asctime when logging.
var timeLayout = time.RFC3339Nano

// defaultFormatter is the formatter used to initialize handler.
var defaultFormatter = NewTextFormatter("%(message)s")

// lastHandler is used when no handler is configured to handle the log record.
var lastHandler = NewHandler("", StderrEmitter)

// handlerManager is a map to search handler by name.
var handlerManager map[string]*Handler

// fileflag is the flag to open a logging file.
var fileflag = os.O_WRONLY | os.O_APPEND | os.O_CREATE

// fileperm is the file permission when creating a logging file.
var fileperm fs.FileMode = 0666

// skipCall is the depth of Logger.log call in program, 2 by default. Increase
// this value if you want to wrap log methods of logger.
var skipCall = 2

var levelToName = map[int]string{
	CRITICAL: "CRITICAL",
	ERROR:    "ERROR",
	WARNING:  "WARNING",
	INFO:     "INFO",
	DEBUG:    "DEBUG",
	NOTSET:   "NOTSET",
}

// SetFileFlag sets the mode when open logging files. It is os.O_WRONLY |
// os.O_APPEND | os.O_CREATE by default.
func SetFileFlag(flag int) {
	lock.WLockFunc(func() { fileflag = flag })
}

// SetFilePerm sets the mode when open logging files. It is os.O_WRONLY |
// os.O_APPEND | os.O_CREATE by default.
func SetFilePerm(perm fs.FileMode) {
	lock.WLockFunc(func() { fileperm = perm })
}

// SetSkipCall sets the new skipCall value which dertermine the depth call of
// Logger.log method.
func SetSkipCall(skip int) {
	lock.WLockFunc(func() { skipCall = skip })
}

// SetTimeLayout sets the time layout to print asctime. It is time.RFC3339Nano
// by default.
func SetTimeLayout(layout string) {
	lock.WLockFunc(func() { timeLayout = layout })
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
	lock.WLockFunc(func() { levelToName[level] = levelName })
}

// GetLogger gets a logger with the specified name (channel name), creating it
// if it doesn't yet exist. This name is a dot-separated hierarchical name, such
// as "a", "a.b", "a.b.c" or similar.
//
// Leave name as empty string to get the root logger.
func GetLogger(name string) *Logger {
	if name == "" {
		return rootLogger
	}
	return lock.RWLockFunc(func() any {
		var lg = rootLogger
		for _, part := range strings.Split(name, ".") {
			if _, ok := lg.children[part]; !ok {
				lg.children[part] = newlogger(part, lg)
			}
			lg = lg.children[part]
		}
		return lg
	}).(*Logger)
}

// getLevelName returns a name associated with the given level.
func getLevelName(level int) string {
	return lock.RLockFunc(func() any {
		return levelToName[level]
	}).(string)
}

// checkLevel validates if the given level is registered or not.
func checkLevel(level int) int {
	return lock.RLockFunc(func() any {
		if _, ok := levelToName[level]; !ok {
			xycond.Panic("level %d is not registered", level)
		}
		return level
	}).(int)
}

// GetHandler returns the handler associated with the name. If no handler found,
// returns nil.
func GetHandler(name string) *Handler {
	var h, ok = handlerManager[name]
	if ok {
		return h
	}
	return nil
}

// mapHandler associates a name with a handler.
func mapHandler(name string, h *Handler) {
	if _, ok := handlerManager[name]; ok {
		xycond.Panic("do not set handler with the same name (%s)", name)
	}
	handlerManager[name] = h
}
