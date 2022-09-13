package encoding

import (
	"fmt"
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

// Append writes an arbitrary object to the Buffer.
func (b *Buffer) Append(a any) {
	switch t := a.(type) {
	case string:
		b.AppendString(t)
	case byte:
		b.AppendByte(t)
	case bool:
		b.AppendBool(t)
	case int:
		b.AppendInt(int64(t))
	case int8:
		b.AppendInt(int64(t))
	case int16:
		b.AppendInt(int64(t))
	case int32:
		b.AppendInt(int64(t))
	case int64:
		b.AppendInt(t)
	case uint:
		b.AppendUint(uint64(t))
	case uint16:
		b.AppendUint(uint64(t))
	case uint32:
		b.AppendUint(uint64(t))
	case uint64:
		b.AppendUint(t)
	case float32:
		b.AppendFloat32(t)
	case float64:
		b.AppendFloat64(t)
	case fmt.Stringer:
		b.AppendString(t.String())
	case fmt.GoStringer:
		b.AppendString(t.GoString())
	default:
		b.AppendString(fmt.Sprint(t))
	}
}

// Bytes returns the reference of underlying byte slice.
func (b *Buffer) Bytes() []byte {
	return b.buf
}

// Len returns the lenght of underlying byte slice.
func (b *Buffer) Len() int {
	return len(b.buf)
}

// Free clears the underlying byte slice and put the buffer into pool.
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
