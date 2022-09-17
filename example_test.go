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

package xylog_test

import (
	"fmt"
	"os"

	"github.com/xybor-x/xylog"
	"github.com/xybor-x/xylog/encoding"
	"github.com/xybor-x/xylog/test"
)

func ExampleLogger() {
	var emitter = xylog.NewStreamEmitter(os.Stdout)
	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)

	var logger = xylog.GetLogger("example.Logger")
	defer xylog.Flush()
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
	defer xylog.Flush()
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)

	logger.Event("create").Field("product", 1235).Debug()

	// Output:
	// event=create product=1235
}

func ExampleNewJSONEncoding() {
	var emitter = xylog.NewStreamEmitter(os.Stdout)
	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)
	handler.SetEncoding(encoding.NewJSONEncoding())
	handler.AddMacro("module", "name")
	handler.AddMacro("level", "levelname")

	var logger = xylog.GetLogger("example.JSONFormatter")
	defer xylog.Flush()
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)
	logger.Event("create").Field("product", 1235).Debug()

	// Output:
	// {"module":"example.JSONFormatter","level":"DEBUG","event":"create","product":1235}
}

func ExampleNewTextEncoding() {
	var emitter = xylog.NewStreamEmitter(os.Stdout)
	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)
	handler.SetEncoding(encoding.NewTextEncoding())
	handler.AddMacro("module", "name")
	handler.AddMacro("level", "levelname")

	var logger = xylog.GetLogger("example.TextFormatter")
	defer xylog.Flush()
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)
	logger.Event("create").Field("employee", "david").Debug()

	// Output:
	// module=example.TextFormatter level=DEBUG event=create employee=david
}

func ExampleFilter() {
	var emitter = xylog.NewStreamEmitter(os.Stdout)
	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)
	handler.AddFilter(&test.LoggerNameFilter{Name: "example.filter.chat"})

	var logger = xylog.GetLogger("example.filter")
	defer xylog.Flush()
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)

	xylog.GetLogger("example.filter.auth").Debug("auth foo")
	xylog.GetLogger("example.filter.chat").Debug("chat foo")

	// Output:
	// messsage="chat foo"
}
