package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func SplitDataframe(buffer []byte) ([][]byte, error) {
	// log.Printf("%x", buffer)
	reader := bytes.NewReader(buffer)
	res := [][]byte{}
	for reader.Len() > 0 {
		length, err := UvarintReader(reader)
		if err != nil {
			return res, fmt.Errorf("failed to read length: %v", err)
		}
		if reader.Len() < length {
			return res, fmt.Errorf("length invalid. expect %d got %d", reader.Len(), length)
		}
		buf := make([]byte, length)
		n, err := reader.Read(buf)
		if err != nil {
			return res, fmt.Errorf("failed to read payload: %v", err)
		}
		if n != length {
			return res, fmt.Errorf("length invalid")
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
		length_byte := make([]byte, binary.MaxVarintLen64)
		n = binary.PutUvarint(length_byte, uint64(length))
		res = append(res, Concat(length_byte[:n], buf))
	}

	return res, nil
}
