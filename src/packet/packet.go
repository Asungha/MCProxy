package packet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	utils "mc_reverse_proxy/src/utils"
)

var PINGHOST = [...]byte{0x00, 0x0b, 0x00, 0x4d, 0x00, 0x43, 0x00, 0x7c, 0x00, 0x50, 0x00, 0x69, 0x00, 0x6e, 0x00, 0x67, 0x00, 0x48, 0x00, 0x6f, 0x00, 0x73, 0x00, 0x74}

type Packet[T IPacketData] struct {
	ID   int
	Data T
}

func (p *Packet[T]) Decode(b *[]byte, maxsize int) error {
	// crop := [:maxsize]
	size, n_length := binary.Uvarint(*b)
	if int(size) > len(*b)-n_length {
		log.Printf("[Packet Decoder] Data malformed. Try another method to serialize.")
		if len(*b) >= 40 {
			log.Printf("%x, %x %x %x", (*b)[:2], (*b)[3:27], PINGHOST, (*b)[27:])
			if bytes.Equal((*b)[:2], []byte{0xfe, 0x01}) && bytes.Equal((*b)[3:27], PINGHOST[:]) {
				log.Printf("[Packet Decoder] Old handshake detected")
				log.Printf("[Packet Decoder] Data length: %x", (*b)[27:29])
				size := binary.BigEndian.Uint16(utils.Concat((*b)[27:29]))
				if len((*b)[27+2:]) == int(size) {
					pv := binary.BigEndian.Uint16(utils.Concat([]byte{0x00}, (*b)[29:30]))
					log.Printf("[Packet Decoder] Data length: %d", size)
					protocolVersion := make([]byte, 64)
					pv_n := binary.PutUvarint(protocolVersion, uint64(pv))
					log.Printf("[Packet Decoder] Protocol version %x (%d)", (*b)[29:30], pv)
					host_l := binary.BigEndian.Uint16((*b)[30:32])
					log.Printf("[Packet Decoder] Hostname length %x (%d)", (*b)[30:32], host_l)
					hostNameLength := make([]byte, 256)
					h_n := binary.PutUvarint(hostNameLength, uint64(host_l))
					hostname_b, err := utils.UTF16toUTF8((*b)[32 : len(*b)-4])
					if err != nil {
						log.Printf("[Packet Decoder] Error: %v", err)
						return err
					}
					port := (*b)[len(*b)-2:]
					data := utils.Concat(protocolVersion[:pv_n], hostNameLength[:h_n], hostname_b, port, []byte{0x01})
					log.Printf("[Packet Decoder] Reformatted data: %x", data)
					return p.Data.Decode(data, len(data))
				}
			} else {
				msg := "Data malformed and can't be serialized"
				log.Printf("[Packet Decoder] %s", msg)
				return errors.New(msg)
			}
		}
		msg := fmt.Sprintf("Data length not correct (%d !> 40)", len(*b))
		log.Printf("[Packet Decoder] %s", msg)
		return errors.New(msg)
	}
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
