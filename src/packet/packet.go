package packet

import (
	"encoding/binary"
	"errors"
	"fmt"
)

var PINGHOST = [...]byte{0x00, 0x0b, 0x00, 0x4d, 0x00, 0x43, 0x00, 0x7c, 0x00, 0x50, 0x00, 0x69, 0x00, 0x6e, 0x00, 0x67, 0x00, 0x48, 0x00, 0x6f, 0x00, 0x73, 0x00, 0x74}

type RemainingData []byte

type PacketHeader struct {
	Length uint64
	ID     byte
}
type Packet struct {
	Payload []byte
	PacketHeader
}

func (rp *Packet) Check() error {
	if rp.Length == 0 {
		return errors.New("Packet length not found")
	}
	if int(rp.Length) > 1+len(rp.Payload) {
		return errors.New(fmt.Sprintf("Packet length incorrect: length %d of data %x", int(rp.Length), string(rp.Payload)))
	}
	return nil
}

func Serialize(packet Packet) []byte {
	data := packet.Payload
	id := make([]byte, 1)
	idn := binary.PutVarint(id, int64(packet.ID))
	encoded := append(id[:idn], data...)
	length := make([]byte, binary.MaxVarintLen64)
	sn := binary.PutUvarint(length, uint64(len(packet.Payload)+idn))
	return append(length[:sn], encoded...)
}

func Deserialize(data []byte) (Packet, RemainingData, error) {
	length, n_length := binary.Uvarint(data)
	id := data[n_length]
	raw := Packet{PacketHeader: PacketHeader{ID: id, Length: length}, Payload: data[n_length+1:]}
	if err := raw.Check(); err != nil {
		return Packet{}, []byte{}, err
	}
	remaining := data[length+1:]
	return raw, remaining, nil
}
