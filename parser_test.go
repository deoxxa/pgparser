package pgparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	singleRecordString   = `({http://www.example.com/image1.png,http://www.example.com/image2.png},image/png,"a logo",123456)`
	multipleRecordString = `{"({http://www.example.com/image1.png,http://www.example.com/image2.png},image/png,\"a logo\",123456)","({http://www.example.com/banner.png},image/png,\"a banner\",123456)"}`
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
