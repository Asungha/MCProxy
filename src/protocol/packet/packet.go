package packet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	. "mc_reverse_proxy/src/common"
	"mc_reverse_proxy/src/utils"
	"strings"
)

// var PINGHOST = [...]byte{0x00, 0x0b, 0x00, 0x4d, 0x00, 0x43, 0x00, 0x7c, 0x00, 0x50, 0x00, 0x69, 0x00, 0x6e, 0x00, 0x67, 0x00, 0x48, 0x00, 0x6f, 0x00, 0x73, 0x00, 0x74}

type RemainingData []byte

type PacketHeader struct {
	Length        uint64
	ID            byte
	IsOldProtocol bool
}
type Packet struct {
	PacketHeader

	Payload *bytes.Reader
}

func (rp *Packet) Check() error {
	if rp.Length == 0 {
		return errors.New("Packet length not found")
	}
	if int(rp.Length) > 1+rp.Payload.Len() {
		buf := make([]byte, rp.Payload.Len())
		_, err := rp.Payload.Read(buf)
		rp.Payload.Reset(buf)
		if err != nil {
			return err
		}
		return errors.New(fmt.Sprintf("Packet length incorrect: length %d of data %x", int(rp.Length), buf))
	}
	if rp.Payload.Len() > 256 {
		buf := make([]byte, rp.Payload.Len())
		_, err := rp.Payload.Read(buf)
		rp.Payload.Reset(buf)
		if err != nil {
			return err
		}
		return errors.New(fmt.Sprintf("Packet length oversize: length %d of data %x", int(rp.Length), buf))
	}
	return nil
}

func Serialize(packet Packet) ([]byte, error) {
	if packet.Payload == nil {
		return []byte{}, nil
	}
	if packet.IsOldProtocol {
		data := make([]byte, packet.Payload.Size())
		_, err := packet.Payload.Read(data)
		if err != nil {
			return []byte{}, err
		}
		return data, nil
	}
	data := packet.Payload
	id := make([]byte, 1)
	idn := binary.PutVarint(id, int64(packet.ID))
	buf := make([]byte, data.Len())
	data.Read(buf)
	encoded := append(id[:idn], buf...)
	length := make([]byte, binary.MaxVarintLen64)
	sn := binary.PutUvarint(length, uint64(packet.Length))
	return append(length[:sn], encoded...), nil
}

func Deserialize(data []byte, isHandshake bool) (Packet, PacketType, RemainingData, error) {
	trueLength := len(data)
	if trueLength == 0 {
		return Packet{}, UNKNOWN, []byte{}, errors.New("empty data")
	}
	if isHandshake {
		packetType := UNKNOWN
		if utils.IsHTTPMethod(strings.Split(string(data), " ")[0]) {
			packetType = HTTP
		}
		err, isOldProtocol := utils.ValidateHandshake(data)
		if err != nil {
			log.Printf("[Packet deserializer] Validation Error: %v, Inferred packet type: %s", err, packetType)
			return Packet{}, packetType, []byte{}, err
		}
		if !isOldProtocol {
			length, n_length := binary.Uvarint(data)
			id := data[n_length]
			raw := Packet{PacketHeader: PacketHeader{ID: id, Length: length}, Payload: bytes.NewReader(data[n_length+1:])}
			if err := raw.Check(); err != nil {
				return Packet{}, MC_HANDSHAKE, []byte{}, err
			}
			if int(length)+n_length == trueLength {
				return raw, MC_HANDSHAKE, []byte{}, nil
			}
			remaining := data[length+1:]
			return raw, MC_HANDSHAKE, remaining, nil
		} else {
			raw := Packet{PacketHeader: PacketHeader{ID: 0x00, IsOldProtocol: true, Length: 1}, Payload: bytes.NewReader(data)}
			if err := raw.Check(); err != nil {
				return Packet{}, packetType, []byte{}, err
			}
			return raw, packetType, []byte{}, nil
		}
	} else {
		err := utils.ValidateMCPacket(data)
		if err != nil {
			return Packet{}, UNKNOWN, []byte{}, errors.New("validation failed. use strict = true for more detailed")
		}
		length, n_length := binary.Uvarint(data)
		id := data[n_length]
		raw := Packet{PacketHeader: PacketHeader{ID: id, Length: length}, Payload: bytes.NewReader(data[n_length+1:])}
		if err := raw.Check(); err != nil {
			return Packet{}, MC_OTHER, []byte{}, err
		}
		if int(length)+n_length == trueLength {
			return raw, MC_OTHER, []byte{}, nil
		}
		remaining := data[length+1:]
		return raw, MC_OTHER, remaining, nil
	}
}
