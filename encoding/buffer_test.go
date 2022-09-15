package encoding_test

import (
	"testing"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xylog/encoding"
)

func TestBufferAppendString(t *testing.T) {
	var buf = encoding.NewBuffer()
	buf.AppendString("foo")

	xycond.ExpectEqual(string(buf.Bytes()), "foo").Test(t)
}

func TestBufferAppendInt(t *testing.T) {
	var buf = encoding.NewBuffer()
	buf.AppendInt(-10)

	xycond.ExpectEqual(string(buf.Bytes()), "-10").Test(t)
}

func TestBufferAppendUint(t *testing.T) {
	var buf = encoding.NewBuffer()
	buf.AppendUint(10)

	xycond.ExpectEqual(string(buf.Bytes()), "10").Test(t)
}

func TestBufferAppendBool(t *testing.T) {
	var buf = encoding.NewBuffer()
	buf.AppendBool(true)

	xycond.ExpectEqual(string(buf.Bytes()), "true").Test(t)
}

func TestBufferAppendByte(t *testing.T) {
	var buf = encoding.NewBuffer()
	buf.AppendByte('=')

	xycond.ExpectEqual(string(buf.Bytes()), "=").Test(t)
}

func TestBufferAppendFloat32(t *testing.T) {
	var buf = encoding.NewBuffer()
	buf.AppendFloat32(4.3)

	xycond.ExpectEqual(string(buf.Bytes()), "4.3").Test(t)
}

func TestBufferAppendFloat64(t *testing.T) {
	var buf = encoding.NewBuffer()
	buf.AppendFloat64(4.3)

	xycond.ExpectEqual(string(buf.Bytes()), "4.3").Test(t)
}

func TestBufferLen(t *testing.T) {
	var buf = encoding.NewBuffer()
	buf.AppendString("foo-bar")

	xycond.ExpectEqual(buf.Len(), 7).Test(t)
}

func TestBufferFree(t *testing.T) {
	var buf = encoding.NewBuffer()
	buf.AppendString("foo-bar")
	buf.Free()

	xycond.ExpectEmpty(buf.Bytes()).Test(t)
}

func TestBufferClone(t *testing.T) {
	var buf = encoding.NewBuffer()
	buf.AppendString("foo-bar")
	var cbuf = buf.Clone()

	xycond.ExpectEqual(string(cbuf.Bytes()), "foo-bar").Test(t)
}
