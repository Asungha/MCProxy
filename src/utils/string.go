package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func SearchUTF8Byte(b []byte) string {
	var s string
	for _, v := range b {
		if v > 0x1F && v < 0x7F {
			s += string(v)
		} else {
			s += "."
		}
	}
	return s
}

func UTF16toUTF8(data []byte) ([]byte, error) {
	decoder := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder()

	// Convert UTF-16BE byte array to UTF-8 byte array
	utf8Bytes, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader(data), decoder))
	if err != nil {
		fmt.Println("Error converting from UTF-16BE to UTF-8:", err)
		return []byte{}, err
	}
	return utf8Bytes, nil
}
