package benchmark

import (
	"testing"

	"github.com/xybor-x/xylog"
)

func BenchmarkGetSameLogger(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xylog.GetLogger("a.b.c.d")
	}
}

func BenchmarkGetRandomLogger(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xylog.GetLogger(getRandomLoggerName())
	}
}

func BenchmarkGetSameHandler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xylog.GetHandler("foo")
	}
}

func BenchmarkGetRandomHandler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xylog.GetHandler(loggerNames[i%len(loggerNames)])
	}
}
