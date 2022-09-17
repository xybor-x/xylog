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

package encoding_test

import (
	"testing"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xylog/encoding"
)

func TestBufferAppendString(t *testing.T) {
	var buf = encoding.NewBuffer()
	buf.AppendString("foo")

	xycond.ExpectEqual(string(buf.Bytes()), "foo").Test(t)
}

func TestBufferAppendInt(t *testing.T) {
	var buf = encoding.NewBuffer()
	buf.AppendInt(-10)

	xycond.ExpectEqual(string(buf.Bytes()), "-10").Test(t)
}

func TestBufferAppendUint(t *testing.T) {
	var buf = encoding.NewBuffer()
	buf.AppendUint(10)

	xycond.ExpectEqual(string(buf.Bytes()), "10").Test(t)
}

func TestBufferAppendBool(t *testing.T) {
	var buf = encoding.NewBuffer()
	buf.AppendBool(true)

	xycond.ExpectEqual(string(buf.Bytes()), "true").Test(t)
}

func TestBufferAppendByte(t *testing.T) {
	var buf = encoding.NewBuffer()
	buf.AppendByte('=')

	xycond.ExpectEqual(string(buf.Bytes()), "=").Test(t)
}

func TestBufferAppendFloat32(t *testing.T) {
	var buf = encoding.NewBuffer()
	buf.AppendFloat32(4.3)

	xycond.ExpectEqual(string(buf.Bytes()), "4.3").Test(t)
}

func TestBufferAppendFloat64(t *testing.T) {
	var buf = encoding.NewBuffer()
	buf.AppendFloat64(4.3)

	xycond.ExpectEqual(string(buf.Bytes()), "4.3").Test(t)
}

func TestBufferLen(t *testing.T) {
	var buf = encoding.NewBuffer()
	buf.AppendString("foo-bar")

	xycond.ExpectEqual(buf.Len(), 7).Test(t)
}

func TestBufferFree(t *testing.T) {
	var buf = encoding.NewBuffer()
	buf.AppendString("foo-bar")
	buf.Free()

	xycond.ExpectEmpty(buf.Bytes()).Test(t)
}

func TestBufferClone(t *testing.T) {
	var buf = encoding.NewBuffer()
	buf.AppendString("foo-bar")
	var cbuf = buf.Clone()

	xycond.ExpectEqual(string(cbuf.Bytes()), "foo-bar").Test(t)
}
