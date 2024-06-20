package packet

import (
	"bytes"
	"fmt"
	"strings"

	// hex "mc_reverse_proxy/src/utils"
	"encoding/binary"
	utils "mc_reverse_proxy/src/utils"
)

const (
	PING  = 0
	LOGIN = 1
)

type Handshake struct {
	PacketHeader

	ProtocolVersion int
	HostnameLength  int
	hostname        string
	Port            int
	NextState       byte
	Type            int
	IsFML           bool
	Tail            []byte
	RawData         []byte
}

func (h Handshake) ImplPacketData() {}

func (h *Handshake) GetHostname() string {
	if strings.Contains(h.hostname, "FML") {
		host := strings.Split(h.hostname, "FML")[0]
		return host[:len(host)-1]
	} else {
		return h.hostname
	}
}

func (h *Handshake) encode() []byte {
	hostname := []byte(h.hostname)
	protocolVersion := make([]byte, binary.MaxVarintLen64)
	n_pv := binary.PutUvarint(protocolVersion, uint64(h.ProtocolVersion))
	hostname_length := make([]byte, binary.MaxVarintLen64)
	n_hl := binary.PutUvarint(hostname_length, uint64(len(hostname)))
	// next_state := make([]byte, 1)
	// binary.BigEndian.PutUint16(next_state, uint64(h.NextState))
	port := make([]byte, 2)
	binary.BigEndian.PutUint16(port, uint16(h.Port))
	raw := utils.Concat(protocolVersion[:n_pv], hostname_length[:n_hl], hostname, port, []byte{h.NextState})
	return raw
}

func (h *Handshake) Encode() ([]byte, error) {
	hostname := []byte(h.hostname)
	protocolVersion := make([]byte, binary.MaxVarintLen64)
	n_pv := binary.PutUvarint(protocolVersion, uint64(h.ProtocolVersion))
	hostname_length := make([]byte, binary.MaxVarintLen64)
	n_hl := binary.PutUvarint(hostname_length, uint64(len(hostname)))
	port := make([]byte, 2)
	binary.BigEndian.PutUint16(port, uint16(h.Port))
	var raw []byte
	if h.IsOldProtocol {
		raw = h.RawData
	} else {
		raw = utils.Concat(protocolVersion[:n_pv], hostname_length[:n_hl], hostname, port, []byte{h.NextState})
	}
	packet := Packet{PacketHeader: h.PacketHeader, Payload: bytes.NewReader(raw)}
	if err := packet.Check(); err != nil {
		return []byte{}, err
	}
	data, err := Serialize(packet)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func (h *Handshake) Decode(data []byte) error {
	packet, remainingdate, err := Deserialize(data)
	if err != nil {
		return err
	}
	if packet.IsOldProtocol {
		h.PacketHeader = packet.PacketHeader
		reader := bytes.NewReader(data)
		reader.Read(make([]byte, len(utils.PING_HOST)))

		// length_b := make([]byte, 2)
		_, err := reader.Read(make([]byte, 2))
		if err != nil {
			return err
		}
		// _ := binary.BigEndian.Uint64(length_b)

		protocol_b := make([]byte, 1)
		_, err = reader.Read(protocol_b[:])
		if err != nil {
			return err
		}
		protocol := int(protocol_b[0])

		str_length_b := make([]byte, 2)
		_, err = reader.Read(str_length_b[:])
		if err != nil {
			return err
		}
		str_length := int(str_length_b[1]) | int(str_length_b[0])<<8

		hostname_b := make([]byte, str_length*2)
		_, err = reader.Read(hostname_b[:])
		if err != nil {
			return err
		}
		hostname_utf8, err := utils.UTF16toUTF8(hostname_b)
		if err != nil {
			return err
		}
		hostname := fmt.Sprintf("%s", hostname_utf8)

		port_b := make([]byte, 4)
		_, err = reader.Read(hostname_b[:])
		if err != nil {
			return err
		}
		port := int(port_b[1]) | int(port_b[0])<<8

		h.hostname = hostname
		h.Port = int(port)
		h.ProtocolVersion = int(protocol)
		h.HostnameLength = int(str_length)
		h.NextState = 0x01
		h.RawData = data
		return nil
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

	if strings.Contains(string(hostname), `\0FML\0`) {
		h.IsFML = true
	}

	h.hostname = string(hostname)

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
	return fmt.Sprintf("ProtocolVersion: %d, Hostname: %s, Port: %d, NextState: %d, Type: %d", h.ProtocolVersion, h.hostname, h.Port, h.NextState, h.Type)
}

func (h *Handshake) Length() int {
	data := h.encode()
	return len(data) - 1
}

func (h *Handshake) Destroy() {
	// h.PlayerData.Destroy()
	// h.PlayerData = nil
}

func NewHandshake() *Handshake {
	return &Handshake{}
}
