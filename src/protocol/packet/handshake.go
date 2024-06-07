package packet

import (
	"bytes"
	"fmt"

	// hex "mc_reverse_proxy/src/utils"
	"encoding/binary"
	utils "mc_reverse_proxy/src/utils"
)

const (
	PING  = 0
	LOGIN = 1
)

type Handshake struct {
	ProtocolVersion int
	HostnameLength  int
	Hostname        string
	Port            int
	NextState       byte
	Type            int
	// PlayerData      *Packet[*PlayerData]
	// PlayerData []byte
	Tail []byte
	IPacket
	PacketHeader
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
	hostname := []byte(h.Hostname)
	protocolVersion := make([]byte, binary.MaxVarintLen64)
	n_pv := binary.PutUvarint(protocolVersion, uint64(h.ProtocolVersion))
	hostname_length := make([]byte, binary.MaxVarintLen64)
	n_hl := binary.PutUvarint(hostname_length, uint64(len(hostname)))
	port := make([]byte, 2)
	binary.BigEndian.PutUint16(port, uint16(h.Port))
	raw := utils.Concat(protocolVersion[:n_pv], hostname_length[:n_hl], hostname, port, []byte{h.NextState})
	packet := Packet{PacketHeader: PacketHeader{Length: uint64(len(raw)), ID: 0x00}, Payload: bytes.NewReader(raw)}
	if err := packet.Check(); err != nil {
		return []byte{}, err
	}
	return Serialize(packet), nil
}

func (h *Handshake) Decode(data []byte) error {
	packet, remainingdate, err := Deserialize(data)
	if err != nil {
		return err
	}
	h.PacketHeader = packet.PacketHeader
	// n := int(packet.Length)
	// protocolVersion, v_length := binary.Uvarint(packet.Payload)
	protocolVersion, err := utils.UvarintReader(packet.Payload)
	if err != nil {
		return err
	}
	// if n-v_length <= 0 {
	// 	return fmt.Errorf("invalid ProtocolVersion: %d", protocolVersion)
	// }
	// h_length, s_length := binary.Uvarint(packet.Payload[v_length:])
	hostnameLength, err := utils.UvarintReader(packet.Payload)
	if err != nil {
		return err
	}
	// if n-v_length-s_length <= 0 {
	// 	return fmt.Errorf("invalid Hostname Length: %d", s_length)
	// }
	h.ProtocolVersion = int(protocolVersion)
	// if s_length+v_length > n-4 {
	// 	return fmt.Errorf("invalid Hostname Length: %d", s_length)
	// }
	hostname := make([]byte, hostnameLength)
	_, err = packet.Payload.Read(hostname)
	if err != nil {
		return err
	}

	h.Hostname = string(hostname)

	port := make([]byte, 2)
	_, err = packet.Payload.Read(port)
	if err != nil {
		return err
	}

	h.Port = int(binary.BigEndian.Uint16(port))

	// port := binary.BigEndian.Uint16(packet.Payload[s_length+v_length+int(h_length) : s_length+v_length+int(h_length)+2])
	// h.Port = int(port)
	nextState := [1]byte{}
	_, err = packet.Payload.Read(nextState[:])
	if err != nil {
		return err
	}
	h.NextState = nextState[0]
	// remainingData := packet.Payload[v_length+int(h_length)+3+1:]
	if l := len(remainingdate); l != 0 {
		h.Tail = remainingdate
	}
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
	return len(data) - 1
}

func (h *Handshake) Destroy() {
	// h.PlayerData.Destroy()
	// h.PlayerData = nil
}

func NewHandshake() *Handshake {
	return &Handshake{}
}
