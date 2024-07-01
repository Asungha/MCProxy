package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	. "mc_reverse_proxy/src/common"
	"strings"
)

var PING_HOST = []byte{0xfe, 0x01, 0xfa, 0x00, 0x0b, 0x00, 0x4D, 0x00, 0x43, 0x00, 0x7C, 0x00, 0x50, 0x00, 0x69, 0x00, 0x6E, 0x00, 0x67, 0x00, 0x48, 0x00, 0x6F, 0x00, 0x73, 0x00, 0x74}

func UvarintReader(r io.Reader) (int, error) {
	varint, err := binary.ReadUvarint(&ByteReader{r})
	if err != nil {
		return 0, err
	}
	return int(varint), nil
}

func VarintReader(r io.Reader) (int, error) {
	varint, err := binary.ReadVarint(&ByteReader{r})
	if err != nil {
		return 0, err
	}
	return int(varint), nil
}

// byteReader wraps an io.Reader to implement io.ByteReader
type ByteReader struct {
	r io.Reader
}

func (b *ByteReader) ReadByte() (byte, error) {
	var buf [1]byte
	_, err := b.r.Read(buf[:])
	return buf[0], err
}

func StrictValidateMCPacket(buffer []byte) (error, bool, PacketType) {
	reader := bytes.NewReader(buffer)

	firstBytes := make([]byte, len(PING_HOST))
	_, err := reader.Read(firstBytes)
	if err == nil && bytes.Equal(firstBytes, PING_HOST) {
		pacType := UNKNOWN
		if IsHTTPMethod(strings.Split(string(buffer), " ")[0]) {
			pacType = HTTP
		}
		return nil, true, pacType
	} else {
		reader.Reset(buffer)
	}
	return ValidateMCPacket(buffer), false, MC_OTHER
}

func ValidateHandshake(buffer []byte) (error, bool) {
	reader := bytes.NewReader(buffer)

	firstBytes := make([]byte, len(PING_HOST))
	_, err := reader.Read(firstBytes)
	if err == nil && bytes.Equal(firstBytes, PING_HOST) {
		return nil, true
	} else {
		reader.Reset(buffer)
	}
	return ValidateMCPacket(buffer), false
}

func ValidateMCPacket(buffer []byte) error {
	reader := bytes.NewReader(buffer)
	for reader.Len() > 0 {
		length, err := UvarintReader(reader)
		if err != nil {
			return fmt.Errorf("failed to read length: %v", err)
		}
		if reader.Len() < length {
			return fmt.Errorf("length invalid")
		}
		buf := make([]byte, length)
		n, err := reader.Read(buf)
		if err != nil {
			return fmt.Errorf("failed to read payload: %v", err)
		}
		if n != length {
			return fmt.Errorf("length invalid")
		}
	}
	return nil
}
