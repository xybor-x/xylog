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
	"io"
	"os"
	"time"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xyerror"
	"github.com/xybor-x/xylock"
	"github.com/xybor-x/xylog/encoding"
)

func init() {
	rootLogger = newLogger("", nil)
	rootLogger.SetLevel(WARNING)

	var handler = GetHandler("xybor")
	handler.AddEmitter(NewDefaultEmitter(os.Stderr))
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

// timeLayout is the default time layout used to print asctime when logging.
var timeLayout = time.RFC3339Nano

// rootLogger is the parent Logger of all Loggers in program.
var rootLogger *Logger

// handlerManager is a map to search Handler by name.
var handlerManager = make(map[string]*Handler)

// emitterManager is a list containing all created Emitters in program.
var emitterManager []Emitter

// skipCall is the depth of Logger.log call in program, 3 by default. Increase
// this value if you want to wrap log methods of logger.
var skipCall = 3

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

// Flush writes unflushed buffered data to outputs.
func Flush() {
	globalLock.RLock()
	defer globalLock.RUnlock()

	for i := range emitterManager {
		emitterManager[i].Flush()
	}
}

// SimpleConfig supports to quickly create a Logger without configurating
// Emitter and Handler.
type SimpleConfig struct {
	// Name is the name of Logger. It can be used later with GetLogger function.
	// Default to an empty name (the root logger).
	Name string

	// Use the specified encoding to format the output. Default to TextEncoding.
	Encoding encoding.Encoding

	// Specify that Logger will write the output to a file. Do NOT use together
	// with Writer.
	Filename string

	// Specify the mode to open file. Default to APPEND | CREATE | WRONLY.
	Filemode int

	// Specify the permission when creating the file. Default to 0666.
	Fileperm os.FileMode

	// The logging level. Default to WARNING.
	Level int

	// The time layout when format the time string. Default to RFC3339Nano.
	TimeLayout string

	// Specify that Logger will write the output to a file. Do NOT use together
	// with Filename.
	Writer io.Writer

	macros []macroField
}

// AddMacro adds a macro value to output format.
func (cfg *SimpleConfig) AddMacro(name, value string) *SimpleConfig {
	cfg.macros = append(cfg.macros, macroField{key: name, macro: value})
	return cfg
}

// Apply creates a Logger based on the configuration.
func (cfg SimpleConfig) Apply() (*Logger, error) {
	if cfg.TimeLayout != "" {
		SetTimeLayout(cfg.TimeLayout)
	}

	if cfg.Filename != "" && cfg.Writer != nil {
		return nil, xyerror.ParameterError.New("do not set both filename and writer")
	}

	var filemode = cfg.Filemode
	if filemode == 0 {
		filemode = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	}

	var fileperm = cfg.Fileperm
	if fileperm == 0 {
		fileperm = 0666
	}

	var writer = cfg.Writer
	if writer == nil {
		if cfg.Filename == "" {
			writer = os.Stderr
		} else {
			var err error
			writer, err = os.OpenFile(cfg.Filename, filemode, fileperm)
			if err != nil {
				return nil, err
			}
		}
	}

	var emitter = NewDefaultEmitter(writer)

	var enc = cfg.Encoding
	if enc == nil {
		enc = encoding.NewTextEncoding()
	}

	var macros = cfg.macros
	if macros == nil {
		macros = append(macros, macroField{key: "time", macro: "asctime"})
		macros = append(macros, macroField{key: "level", macro: "levelname"})
	}

	var handler = GetHandler("")
	handler.AddEmitter(emitter)
	handler.SetEncoding(enc)
	for i := range macros {
		handler.AddMacro(macros[i].key, macros[i].macro)
	}

	var level = cfg.Level
	if level == 0 {
		level = WARNING
	}

	var logger = GetLogger(cfg.Name)
	logger.AddHandler(handler)
	logger.SetLevel(level)

	return logger, nil
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
