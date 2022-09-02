package xylog_test

import (
	"testing"

	"github.com/xybor-x/xylog"
)

func BenchmarkTextFormatter(b *testing.B) {
	var record = xylog.LogRecord{}
	var formatter = xylog.NewTextFormatter(
		"time=%(asctime)s " +
			"source=%(filename)s.%(funcname)s:%(lineno)d " +
			"level=%(levelname)s " +
			"module=%(module)s " +
			"%(message)s",
	)
	for i := 0; i < b.N; i++ {
		formatter.Format(record)
	}
}

func BenchmarkJSONFormatter(b *testing.B) {
	var record = xylog.LogRecord{}
	var formatter = xylog.NewJSONFormatter().
		AddField("asctime", "asctime").
		AddField("created", "created").
		AddField("filename", "filename").
		AddField("funcname", "funcname").
		AddField("levelname", "levelname")

	for i := 0; i < b.N; i++ {
		formatter.Format(record)
	}
}
