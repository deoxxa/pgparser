package pgparser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	singleRecordString   = `({http://www.example.com/image1.png,http://www.example.com/image2.png},image/png,"a logo",123456)`
	multipleRecordString = `{"({http://www.example.com/image1.png,http://www.example.com/image2.png},image/png,\"a logo\",123456)","({http://www.example.com/banner.png},image/png,\"a banner\",123456)"}`
	unquotedString       = `this is a test`
	emptyArray           = `{}`
)

type record struct {
	URLs []string
	Type string
	Name string
	Size int
}

func TestSingleRecord(t *testing.T) {
	a := assert.New(t)

	var r record

	if !a.NoError(Unmarshal(singleRecordString, &r)) {
		return
	}

	a.Equal(record{
		URLs: []string{
			"http://www.example.com/image1.png",
			"http://www.example.com/image2.png",
		},
		Type: "image/png",
		Name: "a logo",
		Size: 123456,
	}, r)
}

func TestMultipleRecords(t *testing.T) {
	a := assert.New(t)

	var r []record

	if !a.NoError(Unmarshal(multipleRecordString, &r)) {
		return
	}

	a.Equal([]record{
		{
			URLs: []string{
				"http://www.example.com/image1.png",
				"http://www.example.com/image2.png",
			},
			Type: "image/png",
			Name: "a logo",
			Size: 123456,
		},
		{
			URLs: []string{
				"http://www.example.com/banner.png",
			},
			Type: "image/png",
			Name: "a banner",
			Size: 123456,
		},
	}, r)
}

func TestByteSlice(t *testing.T) {
	a := assert.New(t)

	var r []byte

	if !a.NoError(Unmarshal(unquotedString, &r)) {
		return
	}

	a.Equal([]byte(unquotedString), r)
}

func TestByteSliceQuoted(t *testing.T) {
	a := assert.New(t)

	var r []byte

	if !a.NoError(Unmarshal("\""+unquotedString+"\"", &r)) {
		return
	}

	a.Equal([]byte(unquotedString), r)
}

func TestString(t *testing.T) {
	a := assert.New(t)

	var r string

	if !a.NoError(Unmarshal(unquotedString, &r)) {
		return
	}

	a.Equal(unquotedString, r)
}

func TestStringQuoted(t *testing.T) {
	a := assert.New(t)

	var r string

	if !a.NoError(Unmarshal("\""+unquotedString+"\"", &r)) {
		return
	}

	a.Equal(unquotedString, r)
}

func TestBadTarget(t *testing.T) {
	a := assert.New(t)

	var r interface{}

	a.Error(Unmarshal("t", &r))
}

func TestNonPointerTarget(t *testing.T) {
	a := assert.New(t)

	var r string

	a.Error(Unmarshal("t", r))
}

func TestBadSourceForString(t *testing.T) {
	a := assert.New(t)

	var s string

	a.Error(Unmarshal(emptyArray, &s))
}

func TestBadSourceForByteSlice(t *testing.T) {
	a := assert.New(t)

	var s []byte

	a.Error(Unmarshal(emptyArray, &s))
}

func TestBadSourceForArray(t *testing.T) {
	a := assert.New(t)

	var l []string

	a.Error(Unmarshal("x", &l))
}

func TestBadSourceForTuple(t *testing.T) {
	a := assert.New(t)

	var s struct{ A, B, C string }

	a.Error(Unmarshal("x", &s))
}

func TestUnclosedArray(t *testing.T) {
	a := assert.New(t)

	var l []string

	a.Error(Unmarshal("{a,b,c", &l))
}

func TestUnclosedTuple(t *testing.T) {
	a := assert.New(t)

	var s struct{ A, B, C string }

	a.Error(Unmarshal("(a,b,c", &s))
}

func TestMismatchedArrayDelimiters(t *testing.T) {
	a := assert.New(t)

	var l []string

	a.Error(Unmarshal("{a,b,c)", &l))
}

func TestMismatchedTupleDelimiters(t *testing.T) {
	a := assert.New(t)

	var s struct{ A, B, C string }

	a.Error(Unmarshal("(a,b,c}", &s))
}

func TestUnfinishedTuple(t *testing.T) {
	a := assert.New(t)

	var s struct{ A, B, C string }

	a.Error(Unmarshal("(a,b)", &s))
}

func TestGoodInt(t *testing.T) {
	a := assert.New(t)

	var i int

	if a.NoError(Unmarshal("2147483647", &i)) {
		a.Equal("2147483647", fmt.Sprintf("%d", i))
	}
}

func TestGoodUint(t *testing.T) {
	a := assert.New(t)

	var i uint

	if a.NoError(Unmarshal("4294967295", &i)) {
		a.Equal("4294967295", fmt.Sprintf("%d", i))
	}
}

func TestGoodInt8(t *testing.T) {
	a := assert.New(t)

	var i int8

	if a.NoError(Unmarshal("127", &i)) {
		a.Equal("127", fmt.Sprintf("%d", i))
	}
}

func TestGoodUint8(t *testing.T) {
	a := assert.New(t)

	var i uint8

	if a.NoError(Unmarshal("255", &i)) {
		a.Equal("255", fmt.Sprintf("%d", i))
	}
}

func TestGoodInt16(t *testing.T) {
	a := assert.New(t)

	var i int16

	if a.NoError(Unmarshal("32767", &i)) {
		a.Equal("32767", fmt.Sprintf("%d", i))
	}
}

func TestGoodUint16(t *testing.T) {
	a := assert.New(t)

	var i uint16

	if a.NoError(Unmarshal("65535", &i)) {
		a.Equal("65535", fmt.Sprintf("%d", i))
	}
}

func TestGoodInt32(t *testing.T) {
	a := assert.New(t)

	var i int32

	if a.NoError(Unmarshal("2147483647", &i)) {
		a.Equal("2147483647", fmt.Sprintf("%d", i))
	}
}

func TestGoodUint32(t *testing.T) {
	a := assert.New(t)

	var i uint32

	if a.NoError(Unmarshal("4294967295", &i)) {
		a.Equal("4294967295", fmt.Sprintf("%d", i))
	}
}

func TestGoodInt64(t *testing.T) {
	a := assert.New(t)

	var i int64

	if a.NoError(Unmarshal("9223372036854775807", &i)) {
		a.Equal("9223372036854775807", fmt.Sprintf("%d", i))
	}
}

func TestGoodUint64(t *testing.T) {
	a := assert.New(t)

	var i uint64

	if a.NoError(Unmarshal("18446744073709551615", &i)) {
		a.Equal("18446744073709551615", fmt.Sprintf("%d", i))
	}
}

func TestTooLongInt(t *testing.T) {
	a := assert.New(t)

	var i int

	a.Error(Unmarshal("2147483648", &i))
}

func TestTooLongUint(t *testing.T) {
	a := assert.New(t)

	var i uint

	a.Error(Unmarshal("4294967296", &i))
}

func TestTooLongInt8(t *testing.T) {
	a := assert.New(t)

	var i int8

	a.Error(Unmarshal("128", &i))
}

func TestTooLongUint8(t *testing.T) {
	a := assert.New(t)

	var i uint8

	a.Error(Unmarshal("256", &i))
}

func TestTooLongInt16(t *testing.T) {
	a := assert.New(t)

	var i int16

	a.Error(Unmarshal("32768", &i))
}

func TestTooLongUint16(t *testing.T) {
	a := assert.New(t)

	var i uint16

	a.Error(Unmarshal("65536", &i))
}

func TestTooLongInt32(t *testing.T) {
	a := assert.New(t)

	var i int32

	a.Error(Unmarshal("2147483648", &i))
}

func TestTooLongUint32(t *testing.T) {
	a := assert.New(t)

	var i uint32

	a.Error(Unmarshal("4294967296", &i))
}

func TestTooLongInt64(t *testing.T) {
	a := assert.New(t)

	var i int64

	a.Error(Unmarshal("9223372036854775808", &i))
}

func TestTooLongUint64(t *testing.T) {
	a := assert.New(t)

	var i uint64

	a.Error(Unmarshal("18446744073709551616", &i))
}

func TestBadSourceForInt(t *testing.T) {
	a := assert.New(t)

	var i int

	a.Error(Unmarshal(emptyArray, &i))
}

func TestBadSourceForUint(t *testing.T) {
	a := assert.New(t)

	var i uint

	a.Error(Unmarshal(emptyArray, &i))
}

func TestBadSourceForInt8(t *testing.T) {
	a := assert.New(t)

	var i int8

	a.Error(Unmarshal(emptyArray, &i))
}

func TestBadSourceForUint8(t *testing.T) {
	a := assert.New(t)

	var i uint8

	a.Error(Unmarshal(emptyArray, &i))
}

func TestBadSourceForInt16(t *testing.T) {
	a := assert.New(t)

	var i int16

	a.Error(Unmarshal(emptyArray, &i))
}

func TestBadSourceForUint16(t *testing.T) {
	a := assert.New(t)

	var i uint16

	a.Error(Unmarshal(emptyArray, &i))
}

func TestBadSourceForInt32(t *testing.T) {
	a := assert.New(t)

	var i int32

	a.Error(Unmarshal(emptyArray, &i))
}

func TestBadSourceForUint32(t *testing.T) {
	a := assert.New(t)

	var i uint32

	a.Error(Unmarshal(emptyArray, &i))
}

func TestBadSourceForInt64(t *testing.T) {
	a := assert.New(t)

	var i int64

	a.Error(Unmarshal(emptyArray, &i))
}

func TestBadSourceForUint64(t *testing.T) {
	a := assert.New(t)

	var i uint64

	a.Error(Unmarshal(emptyArray, &i))
}
