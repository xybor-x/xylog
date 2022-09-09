package xylog_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xylog"
)

func TestGetHandler(t *testing.T) {
	var handlerA = xylog.GetHandler(t.Name())
	var handlerB = xylog.GetHandler(t.Name())
	xycond.ExpectEqual(handlerA, handlerB).Test(t)
	xycond.ExpectEqual(handlerA.Name(), t.Name()).Test(t)
	xycond.ExpectEqual(handlerB.Name(), t.Name()).Test(t)
}

func TestGetHandlerDiff(t *testing.T) {
	var handlerA = xylog.GetHandler(t.Name() + "1")
	var handlerB = xylog.GetHandler(t.Name() + "2")
	xycond.ExpectNotEqual(handlerA, handlerB).Test(t)
}

func TestGetHandlerWithEmptyName(t *testing.T) {
	var handlerA = xylog.GetHandler("")
	var handlerB = xylog.GetHandler("")
	xycond.ExpectNotEqual(handlerA, handlerB).Test(t)
}

func TestHandlerSetFormatter(t *testing.T) {
	var handler = xylog.GetHandler("")
	var formatter = xylog.NewTextFormatter("")
	handler.SetFormatter(formatter)
	xycond.ExpectEqual(fmt.Sprint(handler.Formatter()), fmt.Sprint(formatter))
}

func TestHandlerSetLevel(t *testing.T) {
	var handler = xylog.GetHandler("")
	handler.SetLevel(xylog.WARN)
	xycond.ExpectEqual(handler.Level(), xylog.WARN).Test(t)
}

func TestHandlerFileters(t *testing.T) {
	var filter = &LoggerNameFilter{"foo"}
	var handler = xylog.GetHandler("")
	handler.AddFilter(filter)
	xycond.ExpectEqual(len(handler.Filters()), 1).Test(t)
	xycond.ExpectEqual(handler.Filters()[0], filter).Test(t)

	handler.RemoveFilter(filter)
	xycond.ExpectEqual(len(handler.Filters()), 0).Test(t)
}

func TestHandlerEmitters(t *testing.T) {
	var emitter = xylog.NewStreamEmitter(os.Stdout)
	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)
	xycond.ExpectEqual(len(handler.Emitters()), 1).Test(t)
	xycond.ExpectEqual(handler.Emitters()[0], emitter).Test(t)

	handler.RemoveEmitter(emitter)
	xycond.ExpectEqual(len(handler.Emitters()), 0).Test(t)
}
