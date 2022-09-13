package encoding

import (
	"fmt"
	"strings"
)

// TextEncoder creates a buffer with format of key=value.
type TextEncoder struct {
	buf *Buffer
}

// NewTextEncoder creates a TextEncoder.
func NewTextEncoder() *TextEncoder {
	return &TextEncoder{buf: NewBuffer()}
}

// AddString adds a field of string to encoder.
func (e *TextEncoder) AddString(k, v string) {
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

// AddInt adds a field of int to encoder.
func (e *TextEncoder) AddInt(k string, v int64) {
	e.addSeperator()
	e.buf.AppendString(k)
	e.buf.AppendByte('=')
	e.buf.AppendInt(v)
}

// AddUint adds a field of uint to encoder.
func (e *TextEncoder) AddUint(k string, v uint64) {
	e.addSeperator()
	e.buf.AppendString(k)
	e.buf.AppendByte('=')
	e.buf.AppendUint(v)
}

// AddBool adds a field of bool to encoder.
func (e *TextEncoder) AddBool(k string, v bool) {
	e.addSeperator()
	e.buf.AppendString(k)
	e.buf.AppendByte('=')
	e.buf.AppendBool(v)
}

// AddFloat32 adds a field of float32 to encoder.
func (e *TextEncoder) AddFloat32(k string, v float32) {
	e.addSeperator()
	e.buf.AppendString(k)
	e.buf.AppendByte('=')
	e.buf.AppendFloat32(v)
}

// AddFloat64 adds a field of float64 to encoder.
func (e *TextEncoder) AddFloat64(k string, v float64) {
	e.addSeperator()
	e.buf.AppendString(k)
	e.buf.AppendByte('=')
	e.buf.AppendFloat64(v)
}

// Add adds an arbitrary field to encoder.
func (e *TextEncoder) Add(k string, v any) {
	switch t := v.(type) {
	case string:
		e.AddString(k, t)
	case bool:
		e.AddBool(k, t)
	case int:
		e.AddInt(k, int64(t))
	case int8:
		e.AddInt(k, int64(t))
	case int16:
		e.AddInt(k, int64(t))
	case int32:
		e.AddInt(k, int64(t))
	case int64:
		e.AddInt(k, t)
	case uint:
		e.AddUint(k, uint64(t))
	case uint16:
		e.AddUint(k, uint64(t))
	case uint32:
		e.AddUint(k, uint64(t))
	case uint64:
		e.AddUint(k, t)
	case float32:
		e.AddFloat32(k, t)
	case float64:
		e.AddFloat64(k, t)
	case fmt.Stringer:
		e.AddString(k, t.String())
	case fmt.GoStringer:
		e.AddString(k, t.GoString())
	default:
		e.AddString(k, fmt.Sprint(t))
	}
}

// Encode finishes the encoding process and returns the final byte slice.
func (e *TextEncoder) Encode() []byte {
	return e.buf.Bytes()
}

// Clone creates a new TextEncoder with the copy of underlying buffer.
func (e *TextEncoder) Clone() *TextEncoder {
	return &TextEncoder{buf: e.buf.Clone()}
}

// Free clears the buffer.
func (e *TextEncoder) Free() {
	e.buf.Free()
	e.buf = nil
}

// addSeperator adds a space if the buffer is not empty.
func (e *TextEncoder) addSeperator() {
	if e.buf.Len() > 0 {
		e.buf.AppendByte(' ')
	}
}
