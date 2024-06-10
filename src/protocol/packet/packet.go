package packet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"mc_reverse_proxy/src/utils"
)

var PINGHOST = [...]byte{0x00, 0x0b, 0x00, 0x4d, 0x00, 0x43, 0x00, 0x7c, 0x00, 0x50, 0x00, 0x69, 0x00, 0x6e, 0x00, 0x67, 0x00, 0x48, 0x00, 0x6f, 0x00, 0x73, 0x00, 0x74}

type RemainingData []byte

type PacketHeader struct {
	Length uint64
	ID     byte
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
		if err != nil {
			return err
		}
		return errors.New(fmt.Sprintf("Packet length incorrect: length %d of data %x", int(rp.Length), buf))
	}
	if rp.Payload.Len() > 256 {
		buf := make([]byte, rp.Payload.Len())
		_, err := rp.Payload.Read(buf)
		if err != nil {
			return err
		}
		return errors.New(fmt.Sprintf("Packet length oversize: length %d of data %x", int(rp.Length), buf))
	}
	return nil
}

func Serialize(packet Packet) []byte {
	data := packet.Payload
	id := make([]byte, 1)
	idn := binary.PutVarint(id, int64(packet.ID))
	buf := make([]byte, data.Len())
	data.Read(buf)
	encoded := append(id[:idn], buf...)
	length := make([]byte, binary.MaxVarintLen64)
	sn := binary.PutUvarint(length, uint64(packet.Length+uint64(idn)))
	return append(length[:sn], encoded...)
}

func Deserialize(data []byte) (Packet, RemainingData, error) {
	trueLength := len(data)
	length, n_length := binary.Uvarint(data)
	// if n_length+1 >= trueLength || length >= uint64(trueLength) {
	// 	return Packet{}, data, errors.New("invalid data length")
	// }
	if err := utils.ValidateDataframe(data); err != nil {
		log.Printf("[Packet deserializer] Error: %v", err)
		return Packet{}, []byte{}, err
	}
	id := data[n_length]
	raw := Packet{PacketHeader: PacketHeader{ID: id, Length: length}, Payload: bytes.NewReader(data[n_length+1:])}
	if err := raw.Check(); err != nil {
		return Packet{}, []byte{}, err
	}
	if int(length)+n_length == trueLength {
		return raw, []byte{}, nil
	}
	remaining := data[length+1:]
	return raw, remaining, nil
}
