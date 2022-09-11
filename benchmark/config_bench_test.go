package benchmark

import (
	"testing"

	"github.com/xybor-x/xylog"
	"github.com/xybor-x/xylog/test"
)

func BenchmarkGetSameLogger(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xylog.GetLogger("a.b.c.d")
	}
}

func BenchmarkGetRandomLogger(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xylog.GetLogger(test.GetRandomLoggerName())
	}
}

func BenchmarkGetSameHandler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xylog.GetHandler("foo")
	}
}

func BenchmarkGetRandomHandler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xylog.GetHandler(test.GetRandomLoggerName())
	}
}
