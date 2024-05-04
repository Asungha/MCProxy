package utils

import (
	"encoding/hex"
	"fmt"
)

func ByteToInt(b []byte) int {
	var n int
	for _, v := range b {
		n = n << 8
		n = n | int(v)
	}
	return n
}

func IntToByte(n int) []byte {
	var b []byte
	for n > 0 {
		b = append([]byte{byte(n & 0xFF)}, b...)
		n = n >> 8
	}
	return b
}

func VarIntByteToInt(b *[]byte, start int) (value int, bytecount int) {
	var v int = 0
	var position int = start
	var currentByte byte
	var i int = 1

	for {
		currentByte = (*b)[position]
		v |= int(currentByte&0x7F) << (position)

		if (currentByte & 0x80) == 0 {
			break
		}
		i++
		position += 7
		if position > 35 {
			return 0, 0
		}
	}
	return v, i
}

func IntToVarIntByte(n int) []byte {
	var b []byte
	for n > 0 {
		v := byte(n & 0x7F)
		n = n >> 7
		if n > 0 {
			v = v | 0x80
		}
		b = append([]byte{v}, b...)
	}
	if len(b) == 0 {
		return []byte{0}
	}
	return b
}

// func VarIntLen(b []byte) int {
// 	var n int
// 	for _, v := range b {
// 		n++
// 		if v&0x80 == 0 {
// 			break
// 		}
// 	}
// 	return n
// }

func ByteToHex(b []byte) string {
	return fmt.Sprintf("%x", b)
}

func HexToByte(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}

func ByteToStr(b []byte) string {
	return string(b)
}

func StrToByte(s string) []byte {
	return []byte(s)
}

func SwapEndianness(b []byte) []byte {
	n := len(b)
	for i := 0; i < n/2; i++ {
		b[i], b[n-i-1] = b[n-i-1], b[i]
	}
	return b
}
