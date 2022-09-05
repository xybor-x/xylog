package xylog_test

import (
	"fmt"
	"os"
	"time"

	"github.com/xybor-x/xylog"
)

func ExampleLogger() {
	var emitter = xylog.NewStreamEmitter(os.Stdout)
	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)

	var logger = xylog.GetLogger("example.Logger")
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
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)
	logger.AddField("boss", "foo")

	logger.Event("create").Field("product", 1235).Debug()
	logger.Event("use").Field("product", "bar").JSON().Debug()

	// Output:
	// boss=foo event=create product=1235
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
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)
	logger.Event("create").Field("employee", "david").Debug()

	// Output:
	// module=example.StructuredFormatter level=DEBUG event=create employee=david
}

// LoggerNameFilter only logs out records belongs to a specified logger.
type LoggerNameFilter struct {
	name string
}

func (f *LoggerNameFilter) Filter(r xylog.LogRecord) bool {
	return f.name == r.Name
}

func ExampleFilter() {
	var emitter = xylog.NewStreamEmitter(os.Stdout)
	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)
	handler.AddFilter(&LoggerNameFilter{"example.filter.chat"})

	var logger = xylog.GetLogger("example.filter")
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)

	xylog.GetLogger("example.filter.auth").Debug("auth foo")
	xylog.GetLogger("example.filter.chat").Debug("chat foo")

	// Output:
	// chat foo
}

func ExampleNewSizeRotatingFileEmitter() {
	// Create a rotating emitter which rotates to another files if current file
	// size is over than 30 bytes. Backup maximum of two log files.
	var emitter = xylog.NewSizeRotatingFileEmitter("exampleSize.log", 30, 2)
	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)

	var logger = xylog.GetLogger("example.SizeRotatingFileEmitter")
	logger.SetLevel(xylog.DEBUG)
	logger.AddHandler(handler)

	for i := 0; i < 20; i++ {
		// logger will write 80 bytes (including newlines).
		logger.Debug("foo")
	}

	if _, err := os.Stat("exampleSize.log"); err == nil {
		fmt.Println("Created exampleSize.log")
	}

	if _, err := os.Stat("exampleSize.log.1"); err == nil {
		fmt.Println("Created exampleSize.log.1")
	}

	if _, err := os.Stat("exampleSize.log.2"); err == nil {
		fmt.Println("Created exampleSize.log.2")
	}

	// Output:
	// Created exampleSize.log
	// Created exampleSize.log.1
	// Created exampleSize.log.2
}

func ExampleNewTimeRotatingFileEmitter() {
	// Create a rotating emitter which rotates to another files if logger spent
	// one second to log.
	var emitter = xylog.NewTimeRotatingFileEmitter(
		"exampleTime.log", time.Second, 2)
	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)
	var logger = xylog.GetLogger("example.TimeRotatingFileEmitter")
	logger.SetLevel(xylog.DEBUG)
	logger.AddHandler(handler)

	for i := 0; i < 20; i++ {
		// logger will write for 4s.
		time.Sleep(200 * time.Millisecond)
		logger.Debug("foo")
	}

	if _, err := os.Stat("exampleTime.log"); err == nil {
		fmt.Println("Created exampleTime.log")
	}

	if _, err := os.Stat("exampleTime.log.1"); err == nil {
		fmt.Println("Created exampleTime.log.1")
	}

	if _, err := os.Stat("exampleTime.log.2"); err == nil {
		fmt.Println("Created exampleTime.log.2")
	}

	// Output:
	// Created exampleTime.log
	// Created exampleTime.log.1
	// Created exampleTime.log.2
}
