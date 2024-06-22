package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	packetLoggerService "mc_reverse_proxy/src/packet-logger/service"
)

type PacketFragment struct {
	Data []byte
	Type packetLoggerService.PacketType
}

func SplitDataframe(buffer []byte) ([]PacketFragment, error) {
	if len(buffer) == 0 {
		return []PacketFragment{}, errors.New("empty data")
	}
	err, isOldProtocol, PackType := ValidateDataframe(buffer)
	if err != nil {
		return []PacketFragment{{Data: buffer, Type: PackType}}, err
	}
	if isOldProtocol {
		return []PacketFragment{{Data: buffer, Type: PackType}}, nil
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
		res = append(res, PacketFragment{Data: Concat(length_byte[:n], buf), Type: packetLoggerService.MC_OTHER})
	}

	return res, nil
}
