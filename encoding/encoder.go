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

import "fmt"

// Encoding instances allow encoding with a specified format.
type Encoding interface {
	// addString adds a field of string to the Encoder.
	addString(k, v string)

	// addInt adds a field of int to the Encoder.
	addInt(k string, v int64)

	// addUint adds a field of uint to the Encoder.
	addUint(k string, v uint64)

	// addBool adds a field of bool to the Encoder.
	addBool(k string, v bool)

	// addFloat32 adds a field of float32 to the Encoder.
	addFloat32(k string, v float32)

	// addFloat64 adds a field of float64 to the Encoder.
	addFloat64(k string, v float64)

	// encode finishes the encoding process and returns the final byte slice.
	encode() []byte

	// clone creates a new TextEncoder with the copy of underlying buffer.
	clone() Encoding

	// free clears the buffer.
	free()
}

// Encoder is a wrapper struct of Encoding.
type Encoder struct {
	encoding Encoding
}

// NewEncoder returns a Encoder with a specified Encoding.
func NewEncoder(e Encoding) *Encoder {
	return &Encoder{encoding: e}
}

// Add is a generic method using to encode a generic value.
func (encoder *Encoder) Add(k string, v any) {
	switch t := v.(type) {
	case string:
		encoder.encoding.addString(k, t)
	case bool:
		encoder.encoding.addBool(k, t)
	case int:
		encoder.encoding.addInt(k, int64(t))
	case int8:
		encoder.encoding.addInt(k, int64(t))
	case int16:
		encoder.encoding.addInt(k, int64(t))
	case int32:
		encoder.encoding.addInt(k, int64(t))
	case int64:
		encoder.encoding.addInt(k, t)
	case uint:
		encoder.encoding.addUint(k, uint64(t))
	case uint8:
		encoder.encoding.addUint(k, uint64(t))
	case uint16:
		encoder.encoding.addUint(k, uint64(t))
	case uint32:
		encoder.encoding.addUint(k, uint64(t))
	case uint64:
		encoder.encoding.addUint(k, t)
	case float32:
		encoder.encoding.addFloat32(k, t)
	case float64:
		encoder.encoding.addFloat64(k, t)
	case error:
		encoder.encoding.addString(k, t.Error())
	case fmt.Stringer:
		encoder.encoding.addString(k, t.String())
	case fmt.GoStringer:
		encoder.encoding.addString(k, t.GoString())
	default:
		encoder.encoding.addString(k, fmt.Sprint(t))
	}
}

// Encode finishes the encoding process and returns the final byte slice.
func (encoder *Encoder) Encode() []byte {
	return encoder.encoding.encode()
}

// Clone creates a new TextEncoder with the copy of underlying buffer.
func (encoder *Encoder) Clone() *Encoder {
	return &Encoder{encoding: encoder.encoding.clone()}
}

// Free clears the buffer.
func (encoder *Encoder) Free() {
	encoder.encoding.free()
}
