// MIT License
//
// Copyright (c) 2022 xybor-x
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package benchmark

import (
	"testing"

	"github.com/xybor-x/xylog"
	"github.com/xybor-x/xylog/encoding"
	"github.com/xybor-x/xylog/test"
)

var DevnullEmitter *xylog.StreamEmitter

func BenchmarkLoggerDisable(b *testing.B) {
	test.WithBenchLogger(b, func(logger *xylog.Logger) {
		logger.SetLevel(xylog.ERROR)
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				logger.Debug(test.GetRandomMessage())
			}
		})
	})
}

func BenchmarkLoggerWithoutHandler(b *testing.B) {
	test.WithBenchLogger(b, func(logger *xylog.Logger) {
		logger.RemoveAllHandlers()
		logger.SetLevel(xylog.DEBUG)
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				logger.Debug(test.GetRandomMessage())
			}
		})
	})
}

func BenchmarkLoggerTextEncoding(b *testing.B) {
	test.WithBenchLogger(b, func(logger *xylog.Logger) {
		test.AddFullMacros(logger.Handlers()[0])
		logger.Handlers()[0].SetEncoding(encoding.NewTextEncoding())
		logger.SetLevel(xylog.DEBUG)
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				logger.Debug(test.GetRandomMessage())
			}
		})
	})
}

func BenchmarkLoggerJSONEncoding(b *testing.B) {
	test.WithBenchLogger(b, func(logger *xylog.Logger) {
		test.AddFullMacros(logger.Handlers()[0])
		logger.Handlers()[0].SetEncoding(encoding.NewJSONEncoding())
		logger.SetLevel(xylog.DEBUG)
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				logger.Debug(test.GetRandomMessage())
			}
		})
	})
}
