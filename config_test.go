package xylog_test

import (
	"testing"
	"time"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xyerror"
	"github.com/xybor-x/xylog"
)

func TestConfigSet(t *testing.T) {
	xylog.SetBufferSize(4096)
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
