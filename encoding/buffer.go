package encoding

import "sync"

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

// AppendSeperator writes a space to the Buffer if it is not empty.
func (b *Buffer) AppendSeperator() {
	if len(b.buf) != 0 {
		b.buf = append(b.buf, ' ')
	}
}

// AppendString writes a string to the Buffer.
func (b *Buffer) AppendString(s string) {
	b.buf = append(b.buf, s...)
}

// AppendQuotedString writes a string to the Buffer.
func (b *Buffer) AppendQuotedString(s string) {
	b.buf = append(b.buf, '"')
	b.buf = append(b.buf, s...)
	b.buf = append(b.buf, '"')
}

// AppendByte writes a byte to the Buffer.
func (b *Buffer) AppendByte(a byte) {
	b.buf = append(b.buf, a)
}

// Bytes returns the reference of underlying byte slice.
func (b *Buffer) Bytes() []byte {
	return b.buf
}

// Bytes returns the string copy of underlying byte slice.
func (b *Buffer) String() string {
	return string(b.buf)
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
