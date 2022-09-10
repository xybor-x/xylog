package encoding

// Buffer is a thin wrapper around a byte slice.
type Buffer struct {
	buf []byte
}

// NewBuffer returns a buffer with underlying byte slice is allocated..
func NewBuffer() *Buffer {
	return &Buffer{buf: make([]byte, 0, 64)}
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

// Copy returns a copy of Buffer with the new underlying byte slice.
func (b *Buffer) Copy() *Buffer {
	var c = &Buffer{buf: make([]byte, len(b.buf), cap(b.buf))}
	copy(c.buf, b.buf)
	return c
}
