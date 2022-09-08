package xylog_test

import (
	"fmt"
	"os"

	"github.com/xybor-x/xylog"
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
	// foo bar
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
	logger.AddField("boss", "foo")

	logger.Event("create").Field("product", 1235).Debug()
	logger.Event("use").Field("product", "bar").JSON().Debug()

	// Output:
	// event=create boss=foo product=1235
	// {"boss":"foo","event":"use","product":"bar"}
}

func ExampleTextFormatter() {
	var emitter = xylog.NewStreamEmitter(os.Stdout)
	var formatter = xylog.NewTextFormatter(
		"module=%(name)s level=%(levelname)s %(message)s custom=%(custom)s")

	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)
	handler.SetFormatter(formatter)

	var logger = xylog.GetLogger("example.Formatter")
	defer logger.Flush()
	logger.AddExtra("custom", "something")
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)
	logger.Debug("product=1235")

	// Output:
	// module=example.Formatter level=DEBUG product=1235 custom=something
}

func ExampleJSONFormatter() {
	var emitter = xylog.NewStreamEmitter(os.Stdout)
	var formatter = xylog.NewJSONFormatter().
		AddField("module", "name").
		AddField("level", "levelname").
		AddField("", "message")
	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)
	handler.SetFormatter(formatter)

	var logger = xylog.GetLogger("example.JSONFormatter")
	defer logger.Flush()
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)
	logger.Event("create").Field("product", 1235).JSON().Debug()

	// Output:
	// {"event":"create","level":"DEBUG","module":"example.JSONFormatter","product":1235}
}

func ExampleStructuredFormatter() {
	var emitter = xylog.NewStreamEmitter(os.Stdout)
	var formatter = xylog.NewStructuredFormatter().
		AddField("module", "name").
		AddField("level", "levelname").
		AddField("", "message")
	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)
	handler.SetFormatter(formatter)

	var logger = xylog.GetLogger("example.StructuredFormatter")
	defer logger.Flush()
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)
	logger.Event("create").Field("employee", "david").Debug()

	// Output:
	// module=example.StructuredFormatter level=DEBUG event=create employee=david
}

func ExampleFilter() {
	var emitter = xylog.NewStreamEmitter(os.Stdout)
	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)
	handler.AddFilter(&LoggerNameFilter{"example.filter.chat"})

	var logger = xylog.GetLogger("example.filter")
	defer logger.Flush()
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)

	xylog.GetLogger("example.filter.auth").Debug("auth foo")
	xylog.GetLogger("example.filter.chat").Debug("chat foo")

	// Output:
	// chat foo
}
