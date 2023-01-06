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
	"os"
	"testing"
	"time"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xyerror"
	"github.com/xybor-x/xylog"
	"github.com/xybor-x/xylog/encoding"
	"github.com/xybor-x/xylog/test"
)

func TestConfigSet(t *testing.T) {
	xylog.SetTimeLayout(time.RFC3339Nano)
	xylog.SetSkipCall(3)
	xylog.SetFindCaller(false)
}

func TestLevel(t *testing.T) {
	xylog.AddLevel(130, "TEST")
	xycond.ExpectEqual(xylog.CheckLevel(130), 130).Test(t)
	xycond.ExpectPanic(xyerror.AssertionError, func() {
		xylog.CheckLevel(150)
	}).Test(t)
}

func TestSimpleConfig(t *testing.T) {
	var writer = &test.MockWriter{}
	var logger, err = xylog.SimpleConfig{
		Encoding: encoding.NewJSONEncoding(),
		Writer:   writer,
	}.Apply()

	xycond.ExpectNil(err).Test(t)
	logger.Error("foo")
	xycond.ExpectIn(`"level":"ERROR","messsage":"foo"`, writer.Captured).Test(t)
}

func TestSimpleConfigBothFilenameAndWriter(t *testing.T) {
	var writer = &test.MockWriter{}
	var _, err = xylog.SimpleConfig{
		Filename: "foo.log",
		Writer:   writer,
	}.Apply()

	xycond.ExpectError(err, xyerror.ParameterError).Test(t)
}

func TestSimpleConfigInvalidFile(t *testing.T) {
	var _, err = xylog.SimpleConfig{
		Filename: "/not/a/path",
	}.Apply()

	xycond.ExpectError(err, os.ErrNotExist).Test(t)
}
