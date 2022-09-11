package xylog_test

import (
	"fmt"
	"os"

	"github.com/xybor-x/xylog"
	"github.com/xybor-x/xylog/test"
)

func ExampleLogger() {
	var emitter = xylog.NewStreamEmitter(os.Stdout)
	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)

	var logger = xylog.GetLogger("example.Logger")
	defer logger.Flush()
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)
	logger.Debugf("foo %s", "bar")

	// Output:
	// messsage="foo bar"
}

func ExampleHandler() {
	// You can use a Handler throughout program because all Handlers can be
	// identified by their names.
	var handlerA = xylog.GetHandler("example.Handler")
	var handlerB = xylog.GetHandler("example.Handler")
	if handlerA == handlerB {
		fmt.Println("handlerA == handlerB")
	} else {
		fmt.Println("handlerA != handlerB")
	}

	// In case the name is an empty string, it creates different Handlers every
	// call.
	var handlerC = xylog.GetHandler("")
	var handlerD = xylog.GetHandler("")
	if handlerC == handlerD {
		fmt.Println("handlerC == handlerD")
	} else {
		fmt.Println("handlerC != handlerD")
	}

	// Output:
	// handlerA == handlerB
	// handlerC != handlerD
}

func ExampleEventLogger() {
	var emitter = xylog.NewStreamEmitter(os.Stdout)
	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)

	var logger = xylog.GetLogger("example.EventLogger")
	defer logger.Flush()
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)

	logger.Event("create").Field("product", 1235).Debug()

	// Output:
	// event="create" product="1235"
}

func ExampleJSONFormatter() {
	var emitter = xylog.NewStreamEmitter(os.Stdout)
	var formatter = xylog.NewJSONFormatter().
		AddMacro("module", "name").
		AddMacro("level", "levelname")
	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)
	handler.SetFormatter(formatter)

	var logger = xylog.GetLogger("example.JSONFormatter")
	defer logger.Flush()
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)
	logger.Event("create").Field("product", 1235).Debug()

	// Output:
	// {"event":"create","level":"DEBUG","module":"example.JSONFormatter","product":1235}
}

func ExampleTextFormatter() {
	var emitter = xylog.NewStreamEmitter(os.Stdout)
	var formatter = xylog.NewTextFormatter().
		AddMacro("module", "name").
		AddMacro("level", "levelname")
	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)
	handler.SetFormatter(formatter)

	var logger = xylog.GetLogger("example.TextFormatter")
	defer logger.Flush()
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)
	logger.Event("create").Field("employee", "david").Debug()

	// Output:
	// module="example.TextFormatter" level="DEBUG" event="create" employee="david"
}

func ExampleFilter() {
	var emitter = xylog.NewStreamEmitter(os.Stdout)
	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)
	handler.AddFilter(&test.LoggerNameFilter{Name: "example.filter.chat"})

	var logger = xylog.GetLogger("example.filter")
	defer logger.Flush()
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)

	xylog.GetLogger("example.filter.auth").Debug("auth foo")
	xylog.GetLogger("example.filter.chat").Debug("chat foo")

	// Output:
	// messsage="chat foo"
}
