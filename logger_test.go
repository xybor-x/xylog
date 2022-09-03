package xylog_test

import (
	"os"
	"testing"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xylog"
)

func checkLogOutput(t *testing.T, f func(), msg string, lv, loggerlv int) {
	capturedOutput = ""
	f()
	if lv < loggerlv {
		xycond.ExpectEmpty(capturedOutput).Test(t)
	} else {
		xycond.ExpectEqual(capturedOutput, msg).Test(t)
	}
}

type CapturedEmitter struct{}

func (h *CapturedEmitter) Emit(s string) {
	capturedOutput = s
}

func (h *CapturedEmitter) SetFormatter(xylog.Formatter) {}

type NameFilter struct {
	name string
}

func (f *NameFilter) Filter(r xylog.LogRecord) bool {
	return f.name == r.Name
}

// capturedOutput is the output which CapturedHandler printed.
var capturedOutput string

// validCustomLevels will be added to xylog's level system.
var validCustomLevels = []int{-1, 25, 100}

// invalidCustomLevels will not be added to xylog's level system.
var invalidCustomLevels = []int{-10, 35, 75}

func init() {
	for i := range validCustomLevels {
		xylog.AddLevel(validCustomLevels[i], "")
	}
}

func TestLoggerValidCustomLevel(t *testing.T) {
	var logger = xylog.GetLogger(t.Name())

	xycond.ExpectNotPanic(func() {
		for i := range validCustomLevels {
			logger.SetLevel(validCustomLevels[i])
		}
	}).Test(t)
}

func TestLoggerInvalidCustomLevel(t *testing.T) {
	var logger = xylog.GetLogger(t.Name())

	for i := range invalidCustomLevels {
		xycond.ExpectPanic(func() {
			logger.SetLevel(invalidCustomLevels[i])
		}).Test(t)
	}
}

func TestLoggerHandler(t *testing.T) {
	var expectedHandler = xylog.NewHandler("", &CapturedEmitter{})
	var logger = xylog.GetLogger(t.Name())
	xycond.ExpectNotPanic(func() {
		logger.AddHandler(expectedHandler)
		logger.RemoveHandler(expectedHandler)
	}).Test(t)
}

func TestLoggerAddHandlerNil(t *testing.T) {
	var logger = xylog.GetLogger(t.Name())
	xycond.ExpectPanic(func() {
		logger.AddHandler(nil)
	}).Test(t)
}

func TestLoggerRemoveNotExistedHandler(t *testing.T) {
	var logger = xylog.GetLogger(t.Name())
	xycond.ExpectNotPanic(func() {
		logger.RemoveHandler(xylog.NewHandler("", xylog.StderrEmitter))
	}).Test(t)
}

func TestLoggerLogfMethods(t *testing.T) {
	var loggerLevel = xylog.WARN
	var logger = xylog.GetLogger(t.Name())
	logger.AddHandler(xylog.NewHandler("", &CapturedEmitter{}))
	logger.SetLevel(loggerLevel)

	var loggerMethods = map[int]func(string, ...any){
		xylog.DEBUG:    logger.Debugf,
		xylog.INFO:     logger.Infof,
		xylog.WARN:     logger.Warnf,
		xylog.ERROR:    logger.Errorf,
		xylog.CRITICAL: logger.Criticalf,
	}

	for level, method := range loggerMethods {
		checkLogOutput(t, func() { method("foo") }, "foo", level, loggerLevel)
	}
}

func TestLoggerLogMethods(t *testing.T) {
	var loggerLevel = xylog.WARN
	var logger = xylog.GetLogger(t.Name())
	logger.AddHandler(xylog.NewHandler("", &CapturedEmitter{}))
	logger.SetLevel(loggerLevel)

	var loggerMethods = map[int]func(...any){
		xylog.DEBUG:    logger.Debug,
		xylog.INFO:     logger.Info,
		xylog.WARN:     logger.Warn,
		xylog.ERROR:    logger.Error,
		xylog.CRITICAL: logger.Critical,
	}

	for level, method := range loggerMethods {
		checkLogOutput(t, func() { method("foo") }, "foo", level, loggerLevel)
	}
}

func TestLoggerCallHandlerHierarchy(t *testing.T) {
	var expectedMessage = "foo"
	var handler = xylog.NewHandler("", &CapturedEmitter{})
	var logger = xylog.GetLogger(t.Name())
	logger.SetLevel(xylog.DEBUG)
	logger.AddHandler(handler)

	logger = xylog.GetLogger(t.Name() + ".main")
	capturedOutput = ""
	logger.Info(expectedMessage)
	xycond.ExpectEqual(capturedOutput, expectedMessage).Test(t)
}

func TestLoggerLogNoHandler(t *testing.T) {
	var logger = xylog.GetLogger(t.Name())
	logger.SetLevel(xylog.DEBUG)

	xycond.ExpectNotPanic(func() {
		logger.Infof("foo")
	}).Test(t)
}

func TestLoggerLogNotSetLevel(t *testing.T) {
	var logger = xylog.GetLogger(t.Name())

	xycond.ExpectNotPanic(func() {
		logger.Fatal("foo")
	}).Test(t)
}

func TestLoggerLogInvalidCustomLevel(t *testing.T) {
	var logger = xylog.GetLogger(t.Name())
	logger.AddHandler(xylog.NewHandler("", &CapturedEmitter{}))
	logger.SetLevel(xylog.DEBUG)

	for i := range invalidCustomLevels {
		xycond.ExpectPanic(func() {
			logger.Logf(invalidCustomLevels[i], "msg")
		}).Test(t)
		xycond.ExpectPanic(func() {
			logger.Log(invalidCustomLevels[i], "msg")
		}).Test(t)
	}
}

func TestLoggerLogValidCustomLevel(t *testing.T) {
	var loggerLevel = xylog.DEBUG
	var logger = xylog.GetLogger(t.Name())
	logger.AddHandler(xylog.NewHandler("", &CapturedEmitter{}))
	logger.SetLevel(loggerLevel)

	for i := range validCustomLevels {
		checkLogOutput(t, func() { logger.Logf(validCustomLevels[i], "foo") },
			"foo", validCustomLevels[i], loggerLevel)
	}
}

func TestLoggerStack(t *testing.T) {
	var handler = xylog.NewHandler("", xylog.NewStreamEmitter(os.Stdout))
	handler.SetFormatter(xylog.NewTextFormatter("%(levelname)s %(message)s"))
	var logger = xylog.GetLogger("example.Stack")
	logger.SetLevel(xylog.DEBUG)
	logger.AddHandler(handler)

	xycond.ExpectNotPanic(func() {
		logger.Stack(xylog.DEBUG)
	}).Test(t)
}

func TestLoggerFilter(t *testing.T) {
	var expectedFilter = &NameFilter{}
	var logger = xylog.GetLogger(t.Name())
	logger.AddFilter(expectedFilter)
	xycond.ExpectNotPanic(func() {
		logger.RemoveFilter(expectedFilter)
	}).Test(t)
}

func TestLoggerFilterLog(t *testing.T) {
	var expectedMessage = "foo"
	var logger = xylog.GetLogger(t.Name())
	logger.AddHandler(xylog.NewHandler("", &CapturedEmitter{}))
	logger.SetLevel(xylog.DEBUG)

	capturedOutput = ""
	logger.AddFilter(&NameFilter{t.Name()})
	logger.Debugf(expectedMessage)
	xycond.ExpectEqual(capturedOutput, expectedMessage).Test(t)

	capturedOutput = ""
	logger.AddFilter(&NameFilter{"bar name"})
	logger.Warning(expectedMessage)
	xycond.ExpectEmpty(capturedOutput).Test(t)
}

func TestLoggerAddFields(t *testing.T) {
	var logger = xylog.GetLogger(t.Name())
	logger.AddHandler(xylog.NewHandler("", &CapturedEmitter{}))
	logger.SetLevel(xylog.DEBUG)
	logger.AddField("buzz", "bar")

	capturedOutput = ""
	logger.Event("foo").Debug()
	xycond.ExpectEqual(capturedOutput, "buzz=bar event=foo").Test(t)
}
