package xylog_test

import (
	"testing"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xylog"
)

func testRootLogger(t *testing.T, f func(int)) {
	var handler = xylog.GetHandler("")
	handler.AddEmitter(&CapturedEmitter{})
	handler.SetLevel(xylog.DEBUG)
	xylog.AddHandler(handler)

	var loggerLevel = xylog.INFO
	xylog.SetLevel(loggerLevel)

	f(loggerLevel)

}

func TestRootHandler(t *testing.T) {
	var handler = xylog.GetHandler("")
	xycond.ExpectNotPanic(func() {
		xylog.AddHandler(handler)
	}).Test(t)
}

func TestRootFilter(t *testing.T) {
	var filter = &NameFilter{}
	xycond.ExpectNotPanic(func() {
		xylog.AddFilter(filter)
	}).Test(t)
}

func TestRootSetLevel(t *testing.T) {
	var levels = []int{
		xylog.NOTSET,
		xylog.DEBUG,
		xylog.INFO,
		xylog.WARN,
		xylog.WARNING,
		xylog.ERROR,
		xylog.FATAL,
		xylog.CRITICAL,
	}

	for i := range levels {
		xycond.ExpectNotPanic(func() {
			xylog.SetLevel(levels[i])
		}).Test(t)
	}
}

func TestRootLogfMethods(t *testing.T) {
	var methods = map[int]func(string, ...any){
		xylog.DEBUG:    xylog.Debugf,
		xylog.INFO:     xylog.Infof,
		xylog.WARN:     xylog.Warnf,
		xylog.ERROR:    xylog.Errorf,
		xylog.CRITICAL: xylog.Criticalf,
	}

	testRootLogger(t, func(loggerLevel int) {
		for level, method := range methods {
			checkLogOutput(t, func() { method("foo") }, "foo", level, loggerLevel)
		}
	})
}

func TestRootLogMethods(t *testing.T) {
	var methods = map[int]func(...any){
		xylog.DEBUG:    xylog.Debug,
		xylog.INFO:     xylog.Info,
		xylog.WARN:     xylog.Warn,
		xylog.ERROR:    xylog.Error,
		xylog.CRITICAL: xylog.Critical,
	}

	testRootLogger(t, func(loggerLevel int) {
		for level, method := range methods {
			checkLogOutput(t, func() { method("foo") }, "foo", level, loggerLevel)
		}
	})
}

func TestRootLoggerEvent(t *testing.T) {
	testRootLogger(t, func(loggerLevel int) {
		checkLogOutput(t, func() {
			xylog.Event("foo").Field("bar", "buzz").Info()
		}, "event=foo bar=buzz", xylog.INFO, xylog.INFO)
	})
}

func TestRootLoggerStack(t *testing.T) {
	xycond.ExpectNotPanic(func() {
		xylog.Stack(xylog.ERROR)
	}).Test(t)
}
