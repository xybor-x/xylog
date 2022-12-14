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
	"strings"
)

// NewTextEncoding creates a textEncoding.
func NewTextEncoding() Encoding {
	return &textEncoding{buf: NewBuffer()}
}

// textEncoding creates a buffer with format of key=value.
type textEncoding struct {
	buf *Buffer
}

// addString adds a field of string to encoder.
func (e *textEncoding) addString(k, v string) {
	e.addSeperator()
	e.buf.AppendString(k)
	e.buf.AppendByte('=')
	var ok = strings.ContainsRune(v, ' ')
	if ok {
		e.buf.AppendByte('"')
	}
	e.buf.AppendString(v)
	if ok {
		e.buf.AppendByte('"')
	}
}

// addInt adds a field of int to encoder.
func (e *textEncoding) addInt(k string, v int64) {
	e.addSeperator()
	e.buf.AppendString(k)
	e.buf.AppendByte('=')
	e.buf.AppendInt(v)
}

// addUint adds a field of uint to encoder.
func (e *textEncoding) addUint(k string, v uint64) {
	e.addSeperator()
	e.buf.AppendString(k)
	e.buf.AppendByte('=')
	e.buf.AppendUint(v)
}

// addBool adds a field of bool to encoder.
func (e *textEncoding) addBool(k string, v bool) {
	e.addSeperator()
	e.buf.AppendString(k)
	e.buf.AppendByte('=')
	e.buf.AppendBool(v)
}

// addFloat32 adds a field of float32 to encoder.
func (e *textEncoding) addFloat32(k string, v float32) {
	e.addSeperator()
	e.buf.AppendString(k)
	e.buf.AppendByte('=')
	e.buf.AppendFloat32(v)
}

// addFloat64 adds a field of float64 to encoder.
func (e *textEncoding) addFloat64(k string, v float64) {
	e.addSeperator()
	e.buf.AppendString(k)
	e.buf.AppendByte('=')
	e.buf.AppendFloat64(v)
}

// encode finishes the encoding process and returns the final byte slice.
func (e *textEncoding) encode() []byte {
	return e.buf.Bytes()
}

// clone creates a new TextEncoder with the copy of underlying buffer.
func (e *textEncoding) clone() Encoding {
	return &textEncoding{buf: e.buf.Clone()}
}

// free clears the buffer.
func (e *textEncoding) free() {
	e.buf.Free()
}

// addSeperator adds a space if the buffer is not empty.
func (e *textEncoding) addSeperator() {
	if e.buf.Len() > 0 {
		e.buf.AppendByte(' ')
	}
}
