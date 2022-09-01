package xylog_test

import (
	"fmt"
	"os"
	"time"

	"github.com/xybor-x/xylog"
)

// NOTE: In example_test.go, xylog.StdoutEmitter is not accepted as the compared
// output. For this reason, in all examples, we must create a new one.
// In reality, you should use xylog.StdoutEmitter or xylog.StderrEmitter
// instead.

func Example() {
	// You can directly use xylog functions to log with the root logger.
	var handler = xylog.NewHandler("", xylog.NewStreamEmitter(os.Stdout))

	// Handlers in the root logger will affect to other logger, so in this
	// example, it should remove this handler from the root logger after test.
	xylog.SetLevel(xylog.DEBUG)
	defer xylog.RemoveHandler(handler)

	xylog.AddHandler(handler)
	xylog.Debug("foo")
	xylog.Infof("foo %s", "bar")

	// Output:
	// foo
	// foo bar
}

func ExampleGetLogger() {
	var handler = xylog.NewHandler("", xylog.NewStreamEmitter(os.Stdout))
	handler.SetFormatter(xylog.NewTextFormatter(
		"module=%(name)s level=%(levelname)s %(message)s"))

	var logger = xylog.GetLogger("example")
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)
	logger.AddExtra("some", "thing")
	logger.Debugf("foo %s", "bar")

	// Output:
	// module=example level=DEBUG some=thing foo bar
}

func ExampleHandler() {
	// You can use a handler throughout program without storing it in global
	// scope. All handlers can be identified by their names.
	var handlerA = xylog.NewHandler("example", xylog.StdoutEmitter)
	var handlerB = xylog.GetHandler("example")
	if handlerA == handlerB {
		fmt.Println("handlerA == handlerB")
	} else {
		fmt.Println("handlerA != handlerB")
	}

	// In case name is an empty string, it totally is a fresh handler.
	var handlerC = xylog.NewHandler("", xylog.StdoutEmitter)
	var handlerD = xylog.NewHandler("", xylog.StdoutEmitter)
	if handlerC == handlerD {
		fmt.Println("handlerC == handlerD")
	} else {
		fmt.Println("handlerC != handlerD")
	}

	// Output:
	// handlerA == handlerB
	// handlerC != handlerD
}

func ExampleNewSizeRotatingFileEmitter() {
	// Create a rotating emitter which rotates to another files if current file
	// size is over than 30 bytes. Backup maximum of two log files.
	var emitter = xylog.NewSizeRotatingFileEmitter("exampleSize.log", 30, 2)
	var handler = xylog.NewHandler("", emitter)
	handler.SetFormatter(xylog.NewTextFormatter("%(message)s"))
	var logger = xylog.GetLogger("example_file_emitter")
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
	var handler = xylog.NewHandler("", emitter)
	handler.SetFormatter(xylog.NewTextFormatter("%(message)s"))
	var logger = xylog.GetLogger("example_file_emitter")
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

func ExampleEventLogger() {
	var handler = xylog.NewHandler("", xylog.NewStreamEmitter(os.Stdout))
	handler.SetFormatter(xylog.NewTextFormatter(
		"module=%(name)s level=%(levelname)s %(message)s"))

	var logger = xylog.GetLogger("eventlogger")
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)
	logger.Event("create").Field("product", 1235).Debug()

	// Output:
	// module=eventlogger level=DEBUG event=create product=1235
}
