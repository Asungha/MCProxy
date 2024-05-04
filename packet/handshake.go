package packet

import (
	"fmt"
	"log"
	"mc_reverse_proxy/utils"
	hex "mc_reverse_proxy/utils"
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
	NextState       int
	Type            int
}

func (h Handshake) ImplPacketData() {}

func (h *Handshake) hexProtocolVersion() string {
	return fmt.Sprintf("%x", h.ProtocolVersion)
}

func (h *Handshake) hexPort() string {
	return fmt.Sprintf("%x", h.Port)
}

func (h *Handshake) hexNextState() string {
	return fmt.Sprintf("%x", h.NextState)
}

func (h *Handshake) hexHostname() string {
	return fmt.Sprintf("%x", h.Hostname)
}

func (h *Handshake) Encode() []byte {
	// log.Printf("ProtocolVersion: %d, Hostname: %s, Port: %d, NextState: %d", h.ProtocolVersion, h.Hostname, h.Port, h.NextState)
	// protocol_version := hex.IntToVarIntByte(h.ProtocolVersion)
	hostname := hex.StrToByte(h.Hostname)
	hostname_length := hex.IntToVarIntByte(len(hostname))
	next_state := hex.IntToByte(h.NextState)
	port := hex.IntToByte(h.Port)

	// raw := append(append(append(append(append([]byte{0xfd, 0x05},hostname_length..., hostname...), port...), next_state...)), []byte{0x01}...)
	raw := append(append(append(append(append(append([]byte{0xfd, 0x05}, hostname_length...), hostname...), port...), next_state...), []byte{0x01}...))
	// log.Printf("Raw: %x", raw)

	return raw
}

func (h *Handshake) Decode(data *[]byte) error {
	n := len(*data)
	// _, pac_n := hex.VarIntByteToInt(data)
	// // if byte_count == 0 || byte_count > 3 {
	// // 	return fmt.Errorf("Invalid Lenght field: %d", byte_count)
	// // }
	// if n != pac_n {
	// 	return fmt.Errorf("Invalid packet length: %d != %d", n, pac_n)
	// }
	// hex_version := data[1:v_length]

	// log.Printf("Data: %x", data)
	if n < 3 {
		return nil
	}

	protocolVersion, v_length := hex.VarIntByteToInt(data, 0)
	// log.Printf("ProtocolVersion: %d, Length: %d", protocolVersion, v_length)
	if n-v_length <= 0 {
		return fmt.Errorf("Invalid ProtocolVersion: %d", protocolVersion)
	}
	h_length, s_length := hex.VarIntByteToInt(data, v_length)
	if n-v_length-s_length <= 0 {
		return fmt.Errorf("Invalid Hostname Length: %d", s_length)
	}
	log.Printf("Host size Length: %d", s_length)
	h.ProtocolVersion = protocolVersion
	if s_length+v_length > n-4 {
		return fmt.Errorf("Invalid Hostname Length: %d", s_length)
	}
	// hex_hostname := data[s_length+v_length : h_length+s_length+v_length]

	// hex_hostname := utils.Crop(data, s_length+v_length, h_length+s_length+v_length)
	hex_hostname := utils.HexWrapper{Data: data, Start: s_length + v_length, End: h_length + s_length + v_length}
	// log.Printf("Hostname hex: %x", hex_hostname)
	// hex_u := data[n-1 : n]
	// hex_next_state := data[n-2 : n-1]
	// hex_port := data[n-4 : n-2]
	// hex_u := utils.Crop(data, n-1, n)
	hex_next_state := utils.HexWrapper{Data: data, Start: n - 2, End: n - 1}
	hex_port := utils.HexWrapper{Data: data, Start: n - 4, End: n - 2}

	h.Hostname = hex.ByteToStr(hex_hostname.Get())
	h.Port = hex.ByteToInt(hex_port.Get())

	if i := hex.ByteToInt(hex_next_state.Get()); i > 2 {
		return fmt.Errorf("Invalid NextState: %d", i)
	} else {
		h.NextState = i
	}
	return nil
}

func (h *Handshake) String() string {
	// log.Printf("Version: %d", h.ProtocolVersion)
	return fmt.Sprintf("ProtocolVersion: %d, Hostname: %s, Port: %d, NextState: %d, Type: %d", h.ProtocolVersion, h.Hostname, h.Port, h.NextState, h.Type)
}

func (h *Handshake) Length() int {
	return len(h.Encode())
}

func NewHandshake() *Handshake {
	return &Handshake{}
}
