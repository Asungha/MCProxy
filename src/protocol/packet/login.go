package packet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	packetLoggerService "mc_reverse_proxy/src/packet-logger/service"
	utils "mc_reverse_proxy/src/utils"
)

type Login struct {
	PacketHeader

	Name             string
	UUID             []byte
	SizeHeaderLength int
	isOldProtocol    bool
	isEmpty          bool
}

func (p Login) ImplPacketData() {}

func (p *Login) Encode() ([]byte, error) {
	if p.isEmpty {
		return []byte{}, nil
	}
	name := []byte(p.Name)
	name_length := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(name_length, uint64(len(name)))
	p.SizeHeaderLength = n
	raw := utils.Concat(name_length[:n], name, p.UUID)
	packet := Packet{PacketHeader: PacketHeader{Length: uint64(len(raw) + 1), ID: 0x00}, Payload: bytes.NewReader(raw)}
	if err := packet.Check(); err != nil {
		return []byte{}, err
	}
	data, err := Serialize(packet)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func (p *Login) Decode(data []byte) (error, packetLoggerService.PacketType) {
	packet, packerType, _, err := Deserialize(data)
	if err != nil {
		return err, packerType
	}
	pname_l, err := utils.UvarintReader(packet.Payload)
	if err != nil {
		return err, packerType
	}
	if pname_l == 0 {
		p.isEmpty = true
	}
	pname := make([]byte, pname_l)
	packet.Payload.Read(pname)
	p.Name = string(pname)

	uuid := make([]byte, packet.Payload.Len())
	packet.Payload.Read(uuid)
	p.UUID = uuid
	return nil, packerType
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
