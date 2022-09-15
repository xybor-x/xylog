package encoding_test

import (
	"errors"
	"testing"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xylog/encoding"
)

type SStringer struct{}

func (SStringer) String() string {
	return "stringer"
}

type SGoStringer struct{}

func (SGoStringer) GoString() string {
	return "gostringer"
}

func addFullTypes(e *encoding.Encoder) {
	e.Add("string", "foo bar")
	e.Add("bool", true)
	e.Add("int", int(-1))
	e.Add("int8", int8(-2))
	e.Add("int16", int16(-3))
	e.Add("int32", int32(-4))
	e.Add("int64", int64(-5))
	e.Add("uint", uint(1))
	e.Add("uint8", uint8(2))
	e.Add("uint16", uint16(3))
	e.Add("uint32", uint32(4))
	e.Add("uint64", uint64(5))
	e.Add("float32", float32(0.1))
	e.Add("float64", float64(0.2))
	e.Add("error", errors.New("error"))
	e.Add("stringer", SStringer{})
	e.Add("gostringer", SGoStringer{})
	e.Add("default", func() {})
}

func TestTextEncoder(t *testing.T) {
	var encoder = encoding.NewEncoder(encoding.NewTextEncoding())
	addFullTypes(encoder)

	xycond.ExpectIn("string=\"foo bar\" bool=true int=-1 int8=-2 int16=-3 "+
		"int32=-4 int64=-5 uint=1 uint8=2 uint16=3 uint32=4 uint64=5 "+
		"float32=0.1 float64=0.2 error=error stringer=stringer "+
		"gostringer=gostringer",
		string(encoder.Encode())).Test(t)
}

func TestJSONEncoder(t *testing.T) {
	var encoder = encoding.NewEncoder(encoding.NewJSONEncoding())
	addFullTypes(encoder)

	xycond.ExpectIn(`"string":"foo bar","bool":true,"int":-1,"int8":-2,`+
		`"int16":-3,"int32":-4,"int64":-5,"uint":1,"uint8":2,"uint16":3,`+
		`"uint32":4,"uint64":5,"float32":0.1,"float64":0.2,"error":"error",`+
		`"stringer":"stringer","gostringer":"gostringer"`,
		string(encoder.Encode())).Test(t)
}

func TestTextEncoderClone(t *testing.T) {
	var encoder = encoding.NewEncoder(encoding.NewTextEncoding())
	encoder.Add("foo", "bar")
	var cenc = encoder.Clone()

	xycond.ExpectEqual(string(cenc.Encode()), string(encoder.Encode())).Test(t)
}

func TestTextEncoderFree(t *testing.T) {
	var encoder = encoding.NewEncoder(encoding.NewTextEncoding())
	encoder.Add("foo", "bar")
	encoder.Free()

	xycond.ExpectEmpty(encoder.Encode()).Test(t)
}

func TestJSONEncoderClone(t *testing.T) {
	var encoder = encoding.NewEncoder(encoding.NewJSONEncoding())
	encoder.Add("foo", "bar")
	var cenc = encoder.Clone()

	xycond.ExpectEqual(string(cenc.Encode()), string(encoder.Encode())).Test(t)
}

func TestJSONEncoderFree(t *testing.T) {
	var encoder = encoding.NewEncoder(encoding.NewJSONEncoding())
	encoder.Add("foo", "bar")
	encoder.Free()

	xycond.ExpectEmpty(encoder.Encode()).Test(t)
}
