package benchmark

import (
	"testing"

	"github.com/xybor-x/xylog"
)

func BenchmarkTextFormatter(b *testing.B) {
	var record = xylog.LogRecord{}
	var formatter = addFullMacros(xylog.NewTextFormatter())
	for i := 0; i < b.N; i++ {
		formatter.Format(record)
	}
}

func BenchmarkJSONFormatter(b *testing.B) {
	var record = xylog.LogRecord{}
	var formatter = addFullMacros(xylog.NewJSONFormatter())
	for i := 0; i < b.N; i++ {
		formatter.Format(record)
	}
}
