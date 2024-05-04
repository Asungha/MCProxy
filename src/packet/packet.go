package packet

import (
	"encoding/binary"
	"fmt"
	utils "mc_reverse_proxy/src/utils"
)

type Packet[T IPacketData] struct {
	ID   int
	Data T
}

func (p *Packet[T]) Decode(b *[]byte, maxsize int) error {
	// crop := [:maxsize]
	_, n_length := binary.Varint(*b)
	if n_length == 0 {
		return fmt.Errorf("Invalid Lenght field: %d", n_length)
	}
	id, n_id := binary.Varint(*b)
	if n_id == 0 {
		return fmt.Errorf("Invalid ID field: %d", n_id)
	}

	// if p.Length != len(crop)-(n_length+n_id) {
	// 	return fmt.Errorf("Invalid packet length: %d != %d", p.Length, len(crop)-(n_length+n_id))
	// }

	// log.Printf("Length: %d, ID: 0x%x, Header size %d", length, id, n_length+n_id)
	p.ID = int(id)

	data := utils.NewHexWrapper(b, n_length+n_id, maxsize).Get()

	err := p.Data.Decode(&data)
	if err != nil {
		return err
	}
	return nil
}

func (p *Packet[T]) Encode() []byte {
	data := p.Data.Encode()
	// encoded := append(utils.IntToVarIntByte(p.Length+1), utils.IntToVarIntByte(p.ID)...)
	length := make([]byte, binary.MaxVarintLen64)
	sn := binary.PutVarint(length, int64(len(data)))

	id := make([]byte, binary.MaxVarintLen64)
	idn := binary.PutVarint(id, int64(p.ID))
	encoded := append(length[:sn], id[:idn]...)
	// log.Printf("Length %d -> %x", len(data), len(data))
	return append(encoded, data...)
}

func (p *Packet[T]) String() string {
	return fmt.Sprintf("ID: 0x%x, Data: %s", p.ID, p.Data.String())
}

func (p *Packet[T]) Length() int {
	length := make([]byte, binary.MaxVarintLen64)
	sn := binary.PutVarint(length, int64(len(p.Data.Encode())))
	return p.Data.Length() + sn
}

func NewPacket[T IPacketData](data T) *Packet[T] {
	return &Packet[T]{Data: data}
}
