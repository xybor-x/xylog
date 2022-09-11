package xylog_test

import (
	"testing"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xyerror"
	"github.com/xybor-x/xylog"
	"github.com/xybor-x/xylog/test"
)

func TestJSONFormatterWithoutFields(t *testing.T) {
	var f = xylog.NewJSONFormatter()
	var s, err = f.Format(xylog.LogRecord{})
	xycond.ExpectNil(err).Test(t)
	xycond.ExpectEqual(s, `{}`).Test(t)
}

func TestJSONFormatterWithInvalidField(t *testing.T) {
	var f = xylog.NewJSONFormatter().AddField("foo", func() {})
	var s, err = f.Format(xylog.LogRecord{})
	xycond.ExpectError(err, xyerror.ValueError).Test(t)
	xycond.ExpectEmpty(s).Test(t)
}

func TestJSONFormatterWithInvalidMacro(t *testing.T) {
	var f = xylog.NewJSONFormatter().AddMacro("time", "unknown")
	var s, err = f.Format(xylog.LogRecord{})
	xycond.ExpectError(err, xyerror.ValueError).Test(t)
	xycond.ExpectEmpty(s).Test(t)
}

func TestTextFormatterWithInvalidMacro(t *testing.T) {
	var f = xylog.NewTextFormatter().AddMacro("message", "unknown")
	var s, err = f.Format(xylog.LogRecord{})
	xycond.ExpectError(err, xyerror.ValueError).Test(t)
	xycond.ExpectEmpty(s).Test(t)
}

func TestJSONFormatter(t *testing.T) {
	var formatter = test.AddFullMacros(xylog.NewJSONFormatter())
	var s, err = formatter.Format(test.FullRecord)

	xycond.ExpectNil(err).Test(t)
	xycond.ExpectEqual(s, `{"asctime":"ASCTIME","created":1,"filename":`+
		`"FILENAME","funcname":"FUNCNAME","levelname":"LEVELNAME","levelno":2,`+
		`"lineno":3,"module":"MODULE","msecs":4,"pathname":"PATHNAME",`+
		`"process":5,"relativeCreated":6}`).Test(t)
}

func TestStructureFormatter(t *testing.T) {
	var formatter = test.AddFullMacros(xylog.NewTextFormatter())
	var s, err = formatter.Format(test.FullRecord)

	xycond.ExpectNil(err).Test(t)
	xycond.ExpectEqual(s, `asctime="ASCTIME" created="1" filename="FILENAME" `+
		`funcname="FUNCNAME" levelname="LEVELNAME" levelno="2" lineno="3" `+
		`module="MODULE" msecs="4" pathname="PATHNAME" process="5" `+
		`relativeCreated="6"`).Test(t)
}
