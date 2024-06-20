package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

func SplitDataframe(buffer []byte) ([][]byte, error) {
	if len(buffer) == 0 {
		return [][]byte{}, errors.New("empty data")
	}
	_, isOldProtocol := ValidateDataframe(buffer)
	if isOldProtocol {
		return [][]byte{buffer}, nil
	}
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
		length_byte := make([]byte, binary.MaxVarintLen64)
		n = binary.PutUvarint(length_byte, uint64(length))
		res = append(res, Concat(length_byte[:n], buf))
	}

	return res, nil
}
