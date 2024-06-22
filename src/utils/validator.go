package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	packetLoggerService "mc_reverse_proxy/src/packet-logger/service"
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

func ValidateDataframe(buffer []byte) (error, bool, packetLoggerService.PacketType) {
	reader := bytes.NewReader(buffer)

	firstBytes := make([]byte, len(PING_HOST))
	_, err := reader.Read(firstBytes)
	if err == nil && bytes.Equal(firstBytes, PING_HOST) {
		return nil, true, packetLoggerService.MC_HANDSHAKE
	} else {
		reader.Reset(buffer)
	}

	err = nil

	for reader.Len() > 0 {
		length, _err := UvarintReader(reader)
		if _err != nil {
			err = fmt.Errorf("failed to read length: %v", err)
		}
		if reader.Len() < length {
			err = fmt.Errorf("length invalid")
		}
		buf := make([]byte, length)
		n, _err := reader.Read(buf)
		if _err != nil {
			err = fmt.Errorf("failed to read payload: %v", err)
		}
		if n != length {
			err = fmt.Errorf("length invalid")
		}
	}

	if err != nil {
		if IsHTTPMethod(strings.Split(string(buffer), " ")[0]) {
			return err, false, packetLoggerService.HTTP
		}
		return err, false, packetLoggerService.UNKNOWN
	}

	return nil, false, packetLoggerService.MC_OTHER
}
