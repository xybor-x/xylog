package benchmark

import (
	"testing"

	"github.com/xybor-x/xylog"
	"github.com/xybor-x/xylog/test"
)

func BenchmarkTextFormatter(b *testing.B) {
	var formatter = test.AddFullMacros(xylog.NewTextFormatter())
	for i := 0; i < b.N; i++ {
		formatter.Format(test.FullRecord)
	}
}

func BenchmarkJSONFormatter(b *testing.B) {
	var formatter = test.AddFullMacros(xylog.NewJSONFormatter())
	for i := 0; i < b.N; i++ {
		formatter.Format(test.FullRecord)
	}
}
