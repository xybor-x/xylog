package encoding

// Buffer is a thin wrapper around a byte slice.
type Buffer struct {
	buf []byte
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
