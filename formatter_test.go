package xylog_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xylog"
)

func TestNewTextFormatter(t *testing.T) {
	var f = xylog.NewTextFormatter(
		"time=%(asctime)s %(levelno).3d %(module)s something")
	xycond.ExpectTrue(
		strings.Contains(fmt.Sprint(f), "time=%s %.3d %s something")).Test(t)
}

func TestNewTextFormatterWithPercentageSign(t *testing.T) {
	var f = xylog.NewTextFormatter(
		"%%abc)s")
	xycond.ExpectTrue(strings.Contains(fmt.Sprint(f), "%abc)s")).Test(t)
}

func TestTextFormatter(t *testing.T) {
	var formatter = xylog.NewTextFormatter(
		"%(asctime)s %(created)d %(filename)s %(funcname)s %(levelname)s " +
			"%(levelno)d %(lineno)d %(message)s %(module)s %(msecs)d " +
			"%(name)s %(pathname)s %(process)d %(relativeCreated)d")

	var s = formatter.Format(xylog.LogRecord{
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
	})

	xycond.ExpectEqual(s, "ASCTIME 1 FILENAME FUNCNAME LEVELNAME 2 3 MESSAGE "+
		"MODULE 4 NAME PATHNAME 5 6").Test(t)
}
