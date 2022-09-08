package xylog_test

import (
	"testing"
	"time"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xylog"
)

func TestConfigSet(t *testing.T) {
	xycond.ExpectNotPanic(func() {
		xylog.SetBufferSize(4096)
		xylog.SetTimeLayout(time.RFC3339Nano)
		xylog.SetSkipCall(2)
	}).Test(t)
}

func TestLevel(t *testing.T) {
	xylog.AddLevel(130, "TEST")
	xycond.ExpectNotPanic(func() { xylog.CheckLevel(130) }).Test(t)
	xycond.ExpectPanic(func() { xylog.CheckLevel(150) }).Test(t)
}
