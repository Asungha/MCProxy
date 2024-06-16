package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
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
	// log.Printf("Readed %d", n)
	return buf[0], err
}

func ValidateDataframe(buffer []byte) (error, bool) {
	reader := bytes.NewReader(buffer)

	firstBytes := make([]byte, len(PING_HOST))
	log.Printf("%v", buffer)
	log.Printf("%v", PING_HOST)
	_, err := reader.Read(firstBytes)
	log.Printf("%v", firstBytes)
	log.Printf("%v", bytes.Equal(firstBytes, PING_HOST))
	if err == nil && bytes.Equal(firstBytes, PING_HOST) {
		// Old protocol
		log.Println("Old protocol")
		return nil, true
	} else {
		reader.Reset(buffer)
	}

	for reader.Len() > 0 {
		length, err := UvarintReader(reader)
		if err != nil {
			return fmt.Errorf("failed to read length: %v", err), false
		}
		if reader.Len() < length {
			return fmt.Errorf("length invalid"), false
		}
		buf := make([]byte, length)
		n, err := reader.Read(buf)
		if err != nil {
			return fmt.Errorf("failed to read payload: %v", err), false
		}
		if n != length {
			return fmt.Errorf("length invalid"), false
		}
		// length_n := (startPos - reader.Len())
		// log.Printf("Length n: %d", length_n)
		// startPos2 := reader.Len()
		// // Read the ID (varint)
		// id, err := VarintReader(reader)
		// if err != nil {
		// 	return fmt.Errorf("failed to read ID: %v", err)
		// }
		// if id < 0 {
		// 	return errors.New("invalid ID")
		// }
		// length_id := (startPos2 - reader.Len())
		// log.Printf("id length: %d", length_id)

		// // Calculate payload length (length of dataframe - length field - id field)
		// payloadLength := reader.Len()
		// log.Printf("Payload Length: %d", payloadLength)
		// if payloadLength < 0 || payloadLength > reader.Len() {
		// 	log.Printf("%d %d %d", payloadLength, reader.Len(), length_n)
		// 	return errors.New("invalid payload length")
		// }

		// Read the payload
		// payload := make([]byte, payloadLength)
		// n, err := reader.Read(payload)
		// if err != nil || n != payloadLength {
		// 	return fmt.Errorf("failed to read payload: %v", err)
		// }

		// Process the dataframe (length, id, payload)
		// fmt.Printf("Dataframe - Length: %d Payload: %x\n", length, buf)
	}

	return nil, false
}
