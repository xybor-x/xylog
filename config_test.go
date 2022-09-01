package xylog_test

import (
	"os"
	"testing"
	"time"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xylog"
)

var registeredHandlerNames = []string{"foo", "bar"}
var notRegisteredHandlerNames = []string{"foobar", "barfoo"}

func init() {
	for i := range registeredHandlerNames {
		xylog.NewHandler(registeredHandlerNames[i], xylog.StdoutEmitter)
	}
}

func TestGetLogger(t *testing.T) {
	var names = []string{"", "foo", "foo.bar"}
	for i := range names {
		var logger1 = xylog.GetLogger(names[i])
		var logger2 = xylog.GetLogger(names[i])
		xycond.ExpectEqual(logger1, logger2).Test(t)
	}
}
func TestGetHandler(t *testing.T) {
	for i := range registeredHandlerNames {
		var handlerA = xylog.GetHandler(registeredHandlerNames[i])
		var handlerB = xylog.GetHandler(registeredHandlerNames[i])
		xycond.ExpectEqual(handlerA, handlerB).Test(t)
	}
}

func TestGetHandlerDiff(t *testing.T) {
	var handlerA = xylog.GetHandler(registeredHandlerNames[0])
	var handlerB = xylog.GetHandler(registeredHandlerNames[1])
	xycond.ExpectNotEqual(handlerA, handlerB).Test(t)
}

func TestGetHandlerNotRegisterBefore(t *testing.T) {
	for i := range notRegisteredHandlerNames {
		var handler = xylog.GetHandler(notRegisteredHandlerNames[i])
		xycond.ExpectNil(handler).Test(t)
	}
}

func TestSetTimeLayout(t *testing.T) {
	xycond.ExpectNotPanic(func() {
		xylog.SetTimeLayout("123")
		xylog.SetTimeLayout(time.RFC3339Nano)
	}).Test(t)
}

func TestSetFileFlag(t *testing.T) {
	xycond.ExpectNotPanic(func() {
		xylog.SetFileFlag(os.O_WRONLY | os.O_APPEND | os.O_CREATE)
	}).Test(t)
}

func TestSetFilePerm(t *testing.T) {
	xycond.ExpectNotPanic(func() {
		xylog.SetFilePerm(0666)
	}).Test(t)
}

func TestSetSkipCall(t *testing.T) {
	xycond.ExpectNotPanic(func() {
		xylog.SetSkipCall(2)
	}).Test(t)
}
