package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

func UvarintReader(r io.Reader) (int, error) {
	varint, err := binary.ReadUvarint(&byteReader{r})
	if err != nil {
		return 0, err
	}
	return int(varint), nil
}

func VarintReader(r io.Reader) (int, error) {
	varint, err := binary.ReadVarint(&byteReader{r})
	if err != nil {
		return 0, err
	}
	return int(varint), nil
}

// byteReader wraps an io.Reader to implement io.ByteReader
type byteReader struct {
	r io.Reader
}

func (b *byteReader) ReadByte() (byte, error) {
	var buf [1]byte
	_, err := b.r.Read(buf[:])
	// log.Printf("Readed %d", n)
	return buf[0], err
}

func ValidateDataframe(buffer []byte) error {
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
		// fmt.Printf("Dataframe - Length: %d, ID: %d, Payload: %x\n", length, id, payload)
	}

	return nil
}
