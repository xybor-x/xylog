package xylog_test

import (
	"fmt"
	"os"
	"time"

	"github.com/xybor-x/xylog"
)

// NOTE: In example_test.go, xylog.StdoutEmitter is not accepted as the standard
// output. For this reason, in all examples, we must create a new one.
// In reality, you should use xylog.StdoutEmitter or xylog.StderrEmitter
// instead.

func Example() {
	// You can directly use xylog functions to log with the root logger.
	var handler = xylog.GetHandler("")
	handler.AddEmitter(xylog.NewStreamEmitter(os.Stdout))
	xylog.AddHandler(handler)
	xylog.SetLevel(xylog.DEBUG)

	xylog.Debug("foo")
	xylog.Debugf("foo %s", "bar")

	// Output:
	// foo
	// foo bar
}

func ExampleGetLogger() {
	var handler = xylog.GetHandler("")
	handler.AddEmitter(xylog.NewStreamEmitter(os.Stdout))

	var logger = xylog.GetLogger("example.GetLogger")
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)
	logger.Debugf("foo %s", "bar")

	// Output:
	// foo bar
}

func ExampleHandler() {
	// You can use a handler throughout program without storing it in global
	// scope. All handlers can be identified by their names.
	var handlerA = xylog.GetHandler("example.Handler")
	var handlerB = xylog.GetHandler("example.Handler")
	if handlerA == handlerB {
		fmt.Println("handlerA == handlerB")
	} else {
		fmt.Println("handlerA != handlerB")
	}

	// In case name is an empty string, it totally is a fresh handler.
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

func ExampleFormatter() {
	var handler = xylog.GetHandler("")
	handler.AddEmitter(xylog.NewStreamEmitter(os.Stdout))
	handler.SetFormatter(xylog.NewTextFormatter(
		"module=%(name)s level=%(levelname)s %(message)s custom=%(custom)s"))

	var logger = xylog.GetLogger("example.Formatter")
	logger.AddExtra("custom", "something")
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)
	logger.Debug("product=1235")

	// Output:
	// module=example.Formatter level=DEBUG product=1235 custom=something
}

func ExampleEventLogger() {
	var handler = xylog.GetHandler("")
	handler.AddEmitter(xylog.NewStreamEmitter(os.Stdout))
	handler.SetFormatter(xylog.NewTextFormatter(
		"module=%(name)s level=%(levelname)s %(message)s"))

	var logger = xylog.GetLogger("example.EventLogger")
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)
	logger.AddField("boss", "foo")

	logger.Event("create").Field("product", 1235).Debug()
	logger.Event("use").Field("product", "bar").JSON().Debug()

	// Output:
	// module=example.EventLogger level=DEBUG boss=foo event=create product=1235
	// module=example.EventLogger level=DEBUG {"boss":"foo","event":"use","product":"bar"}
}

func ExampleJSONFormatter() {
	var formatter = xylog.NewJSONFormatter().
		AddField("module", "name").
		AddField("level", "levelname").
		AddField("", "message")
	var handler = xylog.GetHandler("")
	handler.AddEmitter(xylog.NewStreamEmitter(os.Stdout))
	handler.SetFormatter(formatter)

	var logger = xylog.GetLogger("example.JSONFormatter")
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)
	logger.Event("create").Field("product", 1235).JSON().Debug()

	// Output:
	// {"event":"create","level":"DEBUG","module":"example.JSONFormatter","product":1235}
}

func ExampleStructureFormatter() {
	var formatter = xylog.NewStructureFormatter().
		AddField("module", "name").
		AddField("level", "levelname").
		AddField("", "message")
	var handler = xylog.GetHandler("")
	handler.AddEmitter(xylog.NewStreamEmitter(os.Stdout))
	handler.SetFormatter(formatter)

	var logger = xylog.GetLogger("example.StructureFormatter")
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)
	logger.Event("create").Field("employee", "david").Debug()

	// Output:
	// module=example.StructureFormatter level=DEBUG event=create employee=david
}

func ExampleNewSizeRotatingFileEmitter() {
	// Create a rotating emitter which rotates to another files if current file
	// size is over than 30 bytes. Backup maximum of two log files.
	var emitter = xylog.NewSizeRotatingFileEmitter("exampleSize.log", 30, 2)
	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)
	handler.SetFormatter(xylog.NewTextFormatter("%(message)s"))
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
	handler.SetFormatter(xylog.NewTextFormatter("%(message)s"))
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
