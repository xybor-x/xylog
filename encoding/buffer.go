package encoding

import (
	"strconv"
	"sync"
)

var bufferPool = sync.Pool{New: func() any {
	return &buffer{}
}}

// buffer is a thin wrapper around a byte slice.
type buffer struct {
	buf []byte
}

// NewBuffer returns a buffer with underlying byte slice is allocated..
func NewBuffer() *buffer {
	return bufferPool.Get().(*buffer)
}

// AppendString writes a string to the Buffer.
func (b *buffer) AppendString(s string) {
	b.buf = append(b.buf, s...)
}

// AppendByte writes a byte to the Buffer.
func (b *buffer) AppendByte(a byte) {
	b.buf = append(b.buf, a)
}

// AppendInt writes a int to the Buffer.
func (b *buffer) AppendInt(a int64) {
	b.buf = strconv.AppendInt(b.buf, a, 10)
}

// AppendUint writes a uint to the Buffer.
func (b *buffer) AppendUint(a uint64) {
	b.buf = strconv.AppendUint(b.buf, a, 10)
}

// AppendBool writes a bool to the Buffer.
func (b *buffer) AppendBool(a bool) {
	b.buf = strconv.AppendBool(b.buf, a)
}

// AppendFloat32 writes a float32 to the Buffer.
func (b *buffer) AppendFloat32(a float32) {
	b.buf = strconv.AppendFloat(b.buf, float64(a), 'f', -1, 32)
}

// AppendFloat64 writes a float64 to the Buffer.
func (b *buffer) AppendFloat64(a float64) {
	b.buf = strconv.AppendFloat(b.buf, a, 'f', -1, 64)
}

// Bytes returns the reference of underlying byte slice.
func (b *buffer) Bytes() []byte {
	return b.buf
}

// Len returns the lenght of underlying byte slice.
func (b *buffer) Len() int {
	return len(b.buf)
}

// Free puts the buffer into pool. DO NOT use the buffer after calling this
// method.
func (b *buffer) Free() {
	b.buf = b.buf[:0]
	bufferPool.Put(b)
}

// Clone returns a copy of Buffer with the new underlying byte slice.
func (b *buffer) Clone() *buffer {
	var c = NewBuffer()
	c.buf = append(c.buf, b.buf...)
	return c
}
