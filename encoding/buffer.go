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

package encoding

import (
	"strconv"
	"sync"
)

var bufferPool = sync.Pool{New: func() any {
	return &Buffer{}
}}

// Buffer is a thin wrapper around a byte slice.
type Buffer struct {
	buf []byte
}

// NewBuffer returns a buffer with underlying byte slice is allocated..
func NewBuffer() *Buffer {
	return bufferPool.Get().(*Buffer)
}

// AppendString writes a string to the Buffer.
func (b *Buffer) AppendString(s string) {
	b.buf = append(b.buf, s...)
}

// AppendByte writes a byte to the Buffer.
func (b *Buffer) AppendByte(a byte) {
	b.buf = append(b.buf, a)
}

// AppendInt writes a int to the Buffer.
func (b *Buffer) AppendInt(a int64) {
	b.buf = strconv.AppendInt(b.buf, a, 10)
}

// AppendUint writes a uint to the Buffer.
func (b *Buffer) AppendUint(a uint64) {
	b.buf = strconv.AppendUint(b.buf, a, 10)
}

// AppendBool writes a bool to the Buffer.
func (b *Buffer) AppendBool(a bool) {
	b.buf = strconv.AppendBool(b.buf, a)
}

// AppendFloat32 writes a float32 to the Buffer.
func (b *Buffer) AppendFloat32(a float32) {
	b.buf = strconv.AppendFloat(b.buf, float64(a), 'f', -1, 32)
}

// AppendFloat64 writes a float64 to the Buffer.
func (b *Buffer) AppendFloat64(a float64) {
	b.buf = strconv.AppendFloat(b.buf, a, 'f', -1, 64)
}

// Bytes returns the reference of underlying byte slice.
func (b *Buffer) Bytes() []byte {
	return b.buf
}

// Len returns the length of underlying byte slice.
func (b *Buffer) Len() int {
	return len(b.buf)
}

// Free puts the buffer into pool. DO NOT use the buffer after calling this
// method.
func (b *Buffer) Free() {
	b.buf = b.buf[:0]
	bufferPool.Put(b)
}

// Clone returns a copy of Buffer with the new underlying byte slice.
func (b *Buffer) Clone() *Buffer {
	var c = NewBuffer()
	c.buf = append(c.buf, b.buf...)
	return c
}
