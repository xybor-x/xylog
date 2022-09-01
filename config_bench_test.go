package xylog_test

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/xybor-x/xylog"
)

var randomNames []string

func init() {
	rand.Seed(time.Now().UnixNano())
	var letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := 0; i < 100000; i++ {
		var npart = rand.Intn(5)
		var name = make([]string, npart)
		for j := 0; j < npart; j++ {
			var p = rand.Intn(len(letters))
			name[j] = letters[p : p+1]
		}
		randomNames = append(randomNames, strings.Join(name, "."))
	}
}

func BenchmarkGetSameLogger(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xylog.GetLogger("a.b.c.d")
	}
}

func BenchmarkGetRandomLogger(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xylog.GetLogger(randomNames[i%len(randomNames)])
	}
}

func BenchmarkGetSameHandler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xylog.GetHandler("foo")
	}
}

func BenchmarkGetRandomHandler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xylog.GetHandler(randomNames[i%len(randomNames)])
	}
}
