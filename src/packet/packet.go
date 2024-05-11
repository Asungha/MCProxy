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
	size, n_length := binary.Uvarint(*b)
	// log.Printf("Size: %d, n_length: %d", size, n_length)
	// if n_length == 0 {
	// 	return fmt.Errorf("Invalid Lenght field: %d", n_length)
	// }
	id, n_id := binary.Varint((*b)[n_length : n_length+1])
	// log.Printf("ID: %x, n_id: %d", id, n_id)
	// if n_id == 0 {
	// 	return fmt.Errorf("Invalid ID field: %d", n_id)
	// }

	// if p.Length != len(crop)-(n_length+n_id) {
	// 	return fmt.Errorf("Invalid packet length: %d != %d", p.Length, len(crop)-(n_length+n_id))
	// }

	// log.Printf("Length: %d, ID: 0x%x, Header size %d", length, id, n_length+n_id)
	p.ID = int(id)

	data := utils.NewHexWrapper(b, n_length+n_id, maxsize).Get()

	err := p.Data.Decode(data, int(size))
	if err != nil {
		return err
	}
	// log.Printf(p.String())/
	return nil
}

func (p *Packet[T]) Encode() ([]byte, error) {
	data, err := p.Data.Encode()
	if err != nil {
		return nil, err
	}
	// encoded := append(utils.IntToVarIntByte(p.Length+1), utils.IntToVarIntByte(p.ID)...)
	id := make([]byte, 1)
	idn := binary.PutVarint(id, int64(p.ID))

	encoded := append(id[:idn], data...)

	// log.Printf("Data length: %d", len(encoded))
	length := make([]byte, binary.MaxVarintLen64)
	sn := binary.PutUvarint(length, uint64(p.Length()+idn))
	// log.Printf("Length: %x", length[:sn])

	// log.Printf("Length %d -> %x", len(data), len(data))
	return append(length[:sn], encoded...), nil
}

func (p *Packet[T]) String() string {
	return fmt.Sprintf("ID: 0x%x, Data: %s", p.ID, p.Data.String())
}

func (p *Packet[T]) Length() int {
	length := make([]byte, binary.MaxVarintLen64)
	sn := binary.PutVarint(length, int64(p.Data.Length()))
	return p.Data.Length() + sn
}

func (p *Packet[T]) Destroy() {
	p.Data.Destroy()
}

func NewPacket[T IPacketData](data T) Packet[T] {
	return Packet[T]{Data: data}
}

func CastPacket[T IPacketData](p Packet[*Raw], output *Packet[T]) error {
	output.ID = p.ID
	data, err := p.Data.Encode()
	if err != nil {
		return err
	}
	return output.Data.Decode(data, p.Length())
}

func AsRawPacket[T IPacketData](p Packet[T], output *Packet[*Raw]) error {
	output.ID = p.ID
	data, err := p.Encode()
	if err != nil {
		return err
	}
	return output.Data.Decode(data, p.Length())
}
