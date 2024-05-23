package packet

import (
	"encoding/binary"
	"fmt"
	utils "mc_reverse_proxy/src/utils"
)

type Login struct {
	Name             string
	UUID             []byte
	SizeHeaderLength int
	isOldProtocol    bool
	isEmpty          bool
	IPacket
	PacketHeader
}

func (p Login) ImplPacketData() {}

func (p *Login) Encode() ([]byte, error) {
	if p.isEmpty {
		return []byte{}, nil
	}
	// log.Printf("Login: %s, %s", p.Name, p.UUID)
	name := []byte(p.Name)
	name_length := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(name_length, uint64(len(name)))
	p.SizeHeaderLength = n
	raw := utils.Concat(name_length[:n], name, p.UUID)
	// log.Printf("Raw: %x", raw)
	packet := Packet{PacketHeader: PacketHeader{Length: uint64(len(raw) + 1), ID: 0x00}, Payload: raw}
	if err := packet.Check(); err != nil {
		return []byte{}, err
	}
	return Serialize(packet), nil
}

func (p *Login) Decode(data []byte) error {
	packet, _, err := Deserialize(data)
	if err != nil {
		return err
	}
	pname_l, pname_n := binary.Uvarint(packet.Payload)
	if pname_n == 0 {
		p.isEmpty = true
	}
	pname := string(packet.Payload[pname_n : pname_n+int(pname_l)])
	// log.Printf("Player Name: %s", pname)
	p.Name = pname

	p.UUID = packet.Payload[pname_n+int(pname_l):]
	// log.Printf("Player UUID: %x", p.UUID)
	return nil
}

func (p *Login) String() string {
	return fmt.Sprintf("Login: %s, %s", p.Name, p.UUID)
}

func (p *Login) Length() int {
	if p.isEmpty {
		return 0
	}
	return len(p.Name) + len(p.UUID) - p.SizeHeaderLength + 1
}

func (p *Login) Destroy() {}
