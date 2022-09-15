package encoding

// NewJSONEncoding creates a new jsonEncoding.
func NewJSONEncoding() Encoding {
	var e = &jsonEncoding{buf: NewBuffer()}
	e.openNamespace()
	return e
}

// jsonEncoding creates a buffer with json format.
type jsonEncoding struct {
	buf *buffer
}

// addString adds a field of string to encoder.
func (e *jsonEncoding) addString(k, v string) {
	e.addSeperator()
	e.addKey(k)
	e.buf.AppendByte(':')
	e.buf.AppendByte('"')
	e.buf.AppendString(v)
	e.buf.AppendByte('"')
}

// addInt adds a field of int to encoder.
func (e *jsonEncoding) addInt(k string, v int64) {
	e.addSeperator()
	e.addKey(k)
	e.buf.AppendByte(':')
	e.buf.AppendInt(v)
}

// addUint adds a field of uint to encoder.
func (e *jsonEncoding) addUint(k string, v uint64) {
	e.addSeperator()
	e.addKey(k)
	e.buf.AppendByte(':')
	e.buf.AppendUint(v)
}

// addBool adds a field of bool to encoder.
func (e *jsonEncoding) addBool(k string, v bool) {
	e.addSeperator()
	e.addKey(k)
	e.buf.AppendByte(':')
	e.buf.AppendBool(v)
}

// addFloat32 adds a field of float32 to encoder.
func (e *jsonEncoding) addFloat32(k string, v float32) {
	e.addSeperator()
	e.addKey(k)
	e.buf.AppendByte(':')
	e.buf.AppendFloat32(v)
}

// addFloat64 adds a field of float64 to encoder.
func (e *jsonEncoding) addFloat64(k string, v float64) {
	e.addSeperator()
	e.addKey(k)
	e.buf.AppendByte(':')
	e.buf.AppendFloat64(v)
}

// encode finishes the encoding process and returns the final byte slice.
func (e *jsonEncoding) encode() []byte {
	if e.buf.Len() > 0 {
		e.closeNamespace()
	}
	return e.buf.Bytes()
}

// clone creates a new TextEncoder with the copy of underlying buffer.
func (e *jsonEncoding) clone() Encoding {
	return &jsonEncoding{buf: e.buf.Clone()}
}

// free clears the buffer.
func (e *jsonEncoding) free() {
	e.buf.Free()
}

func (e *jsonEncoding) addKey(k string) {
	e.buf.AppendByte('"')
	e.buf.AppendString(k)
	e.buf.AppendByte('"')
}

func (e *jsonEncoding) addSeperator() {
	if e.buf.Len() > 0 && e.buf.Bytes()[e.buf.Len()-1] != '{' {
		e.buf.AppendByte(',')
	}
}

func (e *jsonEncoding) openNamespace() {
	e.buf.AppendByte('{')
}

func (e *jsonEncoding) closeNamespace() {
	e.buf.AppendByte('}')
}
