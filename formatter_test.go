package xylog_test

import (
	"fmt"
	"testing"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xyerror"
	"github.com/xybor-x/xylog"
)

var fullrecord = xylog.LogRecord{
	Asctime:         "ASCTIME",
	Created:         1,
	FileName:        "FILENAME",
	FuncName:        "FUNCNAME",
	LevelName:       "LEVELNAME",
	LevelNo:         2,
	LineNo:          3,
	Message:         "MESSAGE",
	Module:          "MODULE",
	Msecs:           4,
	Name:            "NAME",
	PathName:        "PATHNAME",
	Process:         5,
	RelativeCreated: 6,
}

func TestNewTextFormatter(t *testing.T) {
	var f = xylog.NewTextFormatter(
		"time=%(asctime)s %(levelno).3d %(module)s something")
	xycond.ExpectIn("time=%s %.3d %s something", fmt.Sprint(f)).Test(t)
}

func TestNewTextFormatterWithPercentageSign(t *testing.T) {
	var f = xylog.NewTextFormatter("%%abc)s")
	xycond.ExpectIn("%abc)s", fmt.Sprint(f)).Test(t)
}

func TestNewTextFormatterWithInvalidFormat(t *testing.T) {
	xycond.ExpectPanic(xyerror.AssertionError, func() {
		xylog.NewTextFormatter("%abc)s")
	}).Test(t)
}

func TestTextFormatterWithInvalidJSONMessage(t *testing.T) {
	var f = xylog.NewTextFormatter("%(message)s")
	var s, err = f.Format(xylog.LogRecord{
		Message: map[string]any{"bar": func() {}},
	})
	xycond.ExpectError(err, xyerror.ValueError).Test(t)
	xycond.ExpectEmpty(s).Test(t)
}

func TestTextFormatterWithInvalidMacro(t *testing.T) {
	var f = xylog.NewTextFormatter("%(unknown)d")
	var s, err = f.Format(xylog.LogRecord{})
	xycond.ExpectError(err, xyerror.ValueError).Test(t)
	xycond.ExpectEmpty(s).Test(t)
}

func TestJSONFormatterWithoutMessage(t *testing.T) {
	var f = xylog.NewJSONFormatter().AddField("message", "message")
	var s, err = f.Format(xylog.LogRecord{})
	xycond.ExpectNil(err).Test(t)
	xycond.ExpectEqual(s, `{"message":null}`).Test(t)
}

func TestJSONFormatterWithInvalidJSONMessage(t *testing.T) {
	var f = xylog.NewJSONFormatter().AddField("message", "message")
	var s, err = f.Format(xylog.LogRecord{
		Message: map[string]any{"bar": func() {}},
	})
	xycond.ExpectError(err, xyerror.ValueError).Test(t)
	xycond.ExpectEmpty(s).Test(t)
}

func TestJSONFormatterWithInvalidMacro(t *testing.T) {
	var f = xylog.NewJSONFormatter().AddField("time", "unknown")
	var s, err = f.Format(xylog.LogRecord{})
	xycond.ExpectError(err, xyerror.ValueError).Test(t)
	xycond.ExpectEmpty(s).Test(t)
}

func TestStructuredFormatterWithJSONMessage(t *testing.T) {
	var f = xylog.NewStructuredFormatter().AddField("message", "message")
	var s, err = f.Format(xylog.LogRecord{
		Message: map[string]any{"foo": "bar"},
	})
	xycond.ExpectNil(err).Test(t)
	xycond.ExpectEqual(s, `message={"foo":"bar"}`).Test(t)
}

func TestStructuredFormatterWithInvalidJSONMessage(t *testing.T) {
	var f = xylog.NewStructuredFormatter().AddField("message", "message")
	var s, err = f.Format(xylog.LogRecord{
		Message: map[string]any{"foo": func() {}},
	})
	xycond.ExpectError(err, xyerror.ValueError).Test(t)
	xycond.ExpectEmpty(s).Test(t)
}

func TestStructuredFormatterWithInvalidMacro(t *testing.T) {
	var f = xylog.NewStructuredFormatter().AddField("message", "unknown")
	var s, err = f.Format(xylog.LogRecord{})
	xycond.ExpectError(err, xyerror.ValueError).Test(t)
	xycond.ExpectEmpty(s).Test(t)
}

func TestTextFormatter(t *testing.T) {
	var formatter = xylog.NewTextFormatter(
		"%(asctime)s %(created)d %(filename)s %(funcname)s %(levelname)s " +
			"%(levelno)d %(lineno)d %(message)s %(module)s %(msecs)d " +
			"%(name)s %(pathname)s %(process)d %(relativeCreated)d")

	var s, err = formatter.Format(fullrecord)

	xycond.ExpectNil(err).Test(t)
	xycond.ExpectEqual(s, "ASCTIME 1 FILENAME FUNCNAME LEVELNAME 2 3 MESSAGE "+
		"MODULE 4 NAME PATHNAME 5 6").Test(t)
}

func TestJSONFormatter(t *testing.T) {
	var formatter = xylog.NewJSONFormatter().
		AddField("asctime", "asctime").
		AddField("created", "created").
		AddField("filename", "filename").
		AddField("funcname", "funcname").
		AddField("levelname", "levelname").
		AddField("levelno", "levelno").
		AddField("lineno", "lineno").
		AddField("message", "message").
		AddField("module", "module").
		AddField("msecs", "msecs").
		AddField("pathname", "pathname").
		AddField("process", "process").
		AddField("relativeCreated", "relativeCreated")

	var s, err = formatter.Format(fullrecord)

	xycond.ExpectNil(err).Test(t)
	xycond.ExpectEqual(s, `{"asctime":"ASCTIME","created":1,"filename":`+
		`"FILENAME","funcname":"FUNCNAME","levelname":"LEVELNAME","levelno":2,`+
		`"lineno":3,"message":"MESSAGE","module":"MODULE","msecs":4,`+
		`"pathname":"PATHNAME","process":5,"relativeCreated":6}`).Test(t)
}

func TestStructureFormatter(t *testing.T) {
	var formatter = xylog.NewStructuredFormatter().
		AddField("asctime", "asctime").
		AddField("created", "created").
		AddField("filename", "filename").
		AddField("funcname", "funcname").
		AddField("levelname", "levelname").
		AddField("levelno", "levelno").
		AddField("lineno", "lineno").
		AddField("message", "message").
		AddField("module", "module").
		AddField("msecs", "msecs").
		AddField("pathname", "pathname").
		AddField("process", "process").
		AddField("relativeCreated", "relativeCreated")

	var s, err = formatter.Format(fullrecord)

	xycond.ExpectNil(err).Test(t)
	xycond.ExpectEqual(s, "asctime=ASCTIME created=1 filename=FILENAME "+
		"funcname=FUNCNAME levelname=LEVELNAME levelno=2 lineno=3 "+
		"message=MESSAGE module=MODULE msecs=4 pathname=PATHNAME process=5 "+
		"relativeCreated=6").Test(t)
}
