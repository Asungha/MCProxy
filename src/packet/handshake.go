package packet

import (
	"fmt"
	"log"

	// hex "mc_reverse_proxy/src/utils"
	"encoding/binary"
	utils "mc_reverse_proxy/src/utils"
)

const (
	PING  = 0
	LOGIN = 1
)

type PlayerData struct {
	Name             string
	UUID             []byte
	SizeHeaderLength int
	isOldProtocol    bool
	isEmpty          bool
}

func (p PlayerData) ImplPacketData() {}

func (p *PlayerData) Encode() ([]byte, error) {
	if p.isEmpty {
		return []byte{}, nil
	}
	// log.Printf("PlayerData: %s, %s", p.Name, p.UUID)
	name := []byte(p.Name)
	name_length := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(name_length, uint64(len(name)))
	p.SizeHeaderLength = n
	raw := utils.Concat(name_length[:n], name, p.UUID)
	// log.Printf("Raw: %x", raw)
	return raw, nil
}

func (p *PlayerData) Decode(data []byte, size int) error {
	pname_l, pname_n := binary.Uvarint(data)
	if pname_n == 0 {
		p.isEmpty = true
	}
	pname := string(data[pname_n : pname_n+int(pname_l)])
	// log.Printf("Player Name: %s", pname)
	p.Name = pname

	p.UUID = data[pname_n+int(pname_l):]
	// log.Printf("Player UUID: %x", p.UUID)
	return nil
}

func (p *PlayerData) String() string {
	return fmt.Sprintf("PlayerData: %s, %s", p.Name, p.UUID)
}

func (p *PlayerData) Length() int {
	if p.isEmpty {
		return 0
	}
	return len(p.Name) + len(p.UUID) - p.SizeHeaderLength + 1
}

func (p *PlayerData) Destroy() {}

type Handshake struct {
	ProtocolVersion int
	HostnameLength  int
	Hostname        string
	Port            int
	NextState       byte
	Type            int
	// PlayerData      *Packet[*PlayerData]
	PlayerData []byte
	tail       []byte
}

func (h Handshake) ImplPacketData() {}

func (h *Handshake) encode() ([]byte, error) {
	hostname := []byte(h.Hostname)
	protocolVersion := make([]byte, binary.MaxVarintLen64)
	n_pv := binary.PutUvarint(protocolVersion, uint64(h.ProtocolVersion))
	hostname_length := make([]byte, binary.MaxVarintLen64)
	n_hl := binary.PutUvarint(hostname_length, uint64(len(hostname)))
	// next_state := make([]byte, 1)
	// binary.BigEndian.PutUint16(next_state, uint64(h.NextState))
	port := make([]byte, 2)
	binary.BigEndian.PutUint16(port, uint16(h.Port))
	raw := utils.Concat(protocolVersion[:n_pv], hostname_length[:n_hl], hostname, port, []byte{h.NextState})
	return raw, nil
}

func (h *Handshake) Encode() ([]byte, error) {
	// log.Printf("ProtocolVersion: %d, Hostname: %s, Port: %d, NextState: %d", h.ProtocolVersion, h.Hostname, h.Port, h.NextState)
	// protocol_version := hex.IntToVarIntByte(h.ProtocolVersion)
	hostname := []byte(h.Hostname)
	protocolVersion := make([]byte, binary.MaxVarintLen64)
	n_pv := binary.PutUvarint(protocolVersion, uint64(h.ProtocolVersion))
	hostname_length := make([]byte, binary.MaxVarintLen64)
	n_hl := binary.PutUvarint(hostname_length, uint64(len(hostname)))
	// next_state := make([]byte, 1)
	// binary.BigEndian.PutUint16(next_state, uint64(h.NextState))
	port := make([]byte, 2)
	binary.BigEndian.PutUint16(port, uint16(h.Port))
	// if err != nil {
	// 	log.Printf("Port Error: %v", err)
	// 	return nil, err
	// }

	var tail []byte = []byte{}
	if h.tail != nil {
		tail = h.tail
	}
	// if h.PlayerData != nil && h.NextState == 0x02 {
	// 	playerData, err := h.PlayerData.Encode()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	tail = playerData
	// } else {
	// tail = []byte{0x01, 0x00}
	// }

	raw := utils.Concat(protocolVersion[:n_pv], hostname_length[:n_hl], hostname, port, []byte{h.NextState}, tail)
	// log.Printf("Raw: %x", raw)

	return raw, nil
}

func (h *Handshake) Decode(data []byte, size int) error {
	n := len(data)
	// _, pac_n := hex.VarIntByteToInt(data)
	// // if byte_count == 0 || byte_count > 3 {
	// // 	return fmt.Errorf("Invalid Lenght field: %d", byte_count)
	// // }
	// if n != pac_n {
	// 	return fmt.Errorf("Invalid packet length: %d != %d", n, pac_n)
	// }
	// hex_version := data[1:v_length]

	// log.Printf("Data: %x", data)
	// package_length, l_length := binary.Varint(data)
	// log.Printf("Package Length: %d", package_length)

	protocolVersion, v_length := binary.Uvarint(data)
	// log.Printf("ProtocolVersion: %d, Length: %d", protocolVersion, v_length)
	if n-v_length <= 0 {
		return fmt.Errorf("invalid ProtocolVersion: %d", protocolVersion)
	}
	h_length, s_length := binary.Uvarint(data[v_length:])
	// log.Printf("Hostname Length: %d, Length: %d", h_length, s_length)
	if n-v_length-s_length <= 0 {
		return fmt.Errorf("invalid Hostname Length: %d", s_length)
	}
	h.ProtocolVersion = int(protocolVersion)
	if s_length+v_length > n-4 {
		return fmt.Errorf("invalid Hostname Length: %d", s_length)
	}
	log.Printf("ProtocolVersion: %d", h.ProtocolVersion)
	h.Hostname = string(data[s_length+v_length : s_length+v_length+int(h_length)])
	// log.Printf("Hostname: %s", h.Hostname)
	// log.Printf("Port hex: %x", data[s_length+v_length+int(h_length):s_length+v_length+int(h_length)+2])
	port := binary.BigEndian.Uint16(data[s_length+v_length+int(h_length) : s_length+v_length+int(h_length)+2])
	h.Port = int(port)
	// log.Printf("Port: %d", h.Port)
	h.NextState = data[v_length+int(h_length)+3 : v_length+int(h_length)+3+1][0]
	remainingData := data[v_length+int(h_length)+3+1:]
	if l := len(remainingData); l != 0 {
		if l == 2 {
			h.tail = remainingData
			return nil
		}
		h.PlayerData = remainingData
	}
	// if h.NextState == 2 { // login
	// 	d := data[v_length+s_length+int(h_length)+3:]
	// 	if len(d) < 5 {
	// 		h.PlayerData = nil
	// 		return nil
	// 	}
	// 	playerData := NewPacket(&PlayerData{})
	// 	log.Printf("PlayerData: %x", data[v_length+s_length+int(h_length)+3:])
	// 	err := playerData.Decode(&d, len(d))
	// 	if err != nil {
	// 		return err
	// 	}
	// 	h.PlayerData = &playerData
	// 	// log.Printf("PlayerData: %s", h.playerData.Data.String())
	// }
	// log.Printf("NextState: %d", h.NextState)

	// if i := hex.ByteToInt(hex_next_state.Get()); i > 2 {
	// 	return fmt.Errorf("Invalid NextState: %d", i)
	// } else {
	// 	h.NextState = i
	// }
	// h.NextState = hex.ByteToInt(hex_next_state.Get())
	return nil
}

func (h *Handshake) String() string {
	// log.Printf("Version: %d", h.ProtocolVersion)
	return fmt.Sprintf("ProtocolVersion: %d, Hostname: %s, Port: %d, NextState: %d, Type: %d", h.ProtocolVersion, h.Hostname, h.Port, h.NextState, h.Type)
}

func (h *Handshake) Length() int {
	data, err := h.encode()
	if err != nil {
		return -1
	}
	return len(data) - 3
}

func (h *Handshake) Destroy() {
	// h.PlayerData.Destroy()
	// h.PlayerData = nil
}

func NewHandshake() *Handshake {
	return &Handshake{}
}
