package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	. "mc_reverse_proxy/src/common"
)

type PacketFragment struct {
	Data []byte
	Type PacketType
}

func SplitDataframe(buffer []byte, strict bool) ([]PacketFragment, error) {
	if len(buffer) == 0 {
		return []PacketFragment{}, errors.New("empty data")
	}
	if strict {
		err, isOldProtocol, PackType := StrictValidateMCPacket(buffer)
		if err != nil {
			return []PacketFragment{{Data: buffer, Type: PackType}}, err
		}
		if isOldProtocol {
			return []PacketFragment{{Data: buffer, Type: PackType}}, nil
		}
	} else {
		err := ValidateMCPacket(buffer)
		if err != nil {
			return []PacketFragment{{Data: buffer, Type: UNKNOWN}}, err
		}
	}

	reader := bytes.NewReader(buffer)
	res := []PacketFragment{}

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
		res = append(res, PacketFragment{Data: Concat(length_byte[:n], buf), Type: MC_OTHER})
	}

	return res, nil
}
