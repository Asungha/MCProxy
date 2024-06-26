package packet

import (
	"bytes"
)

var EMPTY_DATA = []byte{}
var NON_RELEVANT_DATA = []byte{0x01, 0x02, 0x03, 0x04}

var SAMPLE_STATUS_DATA_1 = []byte{0x10, 0x00, 0xfd, 0x05, 0x09, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74, 0x63, 0xdd, 0x01}
var reader_1 = bytes.NewReader(SAMPLE_STATUS_DATA_1[2:])
var SAMPLE_STATUS_PACKET_1 = Packet{
	PacketHeader: PacketHeader{
		Length: 16,
		ID:     0x00,
	},
	Payload: reader_1,
}
var SAMPLE_STATUS_PAYLOAD_PACKET_1 = Handshake{
	PacketHeader: PacketHeader{
		Length: 16,
		ID:     0x00,
	},
	hostname:        "localhost",
	ProtocolVersion: 765,
	HostnameLength:  9,
	Port:            25565,
	NextState:       0x01,
	IsFML:           false,
	Tail:            []byte{},
}

//////

var SAMPLE_STATUS_DATA_2 = []byte{0x15, 0x00, 0xfd, 0x05, 0x0e, 0x31, 0x39, 0x32, 0x2e, 0x31, 0x36, 0x38, 0x2e, 0x31, 0x30, 0x30, 0x2e, 0x36, 0x33, 0x63, 0xdd, 0x01}
var reader_2 = bytes.NewReader(SAMPLE_STATUS_DATA_2[2:])
var SAMPLE_STATUS_PACKET_2 = Packet{
	PacketHeader: PacketHeader{
		Length: 21,
		ID:     0x00,
	},
	Payload: reader_2,
}

var SAMPLE_STATUS_PAYLOAD_PACKET_2 = Handshake{
	PacketHeader: PacketHeader{
		Length: 21,
		ID:     0x00,
	},
	hostname:        "192.168.100.63",
	ProtocolVersion: 765,
	HostnameLength:  14,
	Port:            25565,
	NextState:       0x01,
	IsFML:           false,
	Tail:            []byte{},
}

//////

var SAMPLE_STATUS_DATA_LENGTH_OVER = []byte{0x16, 0x00, 0xfd, 0x05, 0x0e, 0x31, 0x39, 0x32, 0x2e, 0x31, 0x36, 0x38, 0x2e, 0x31, 0x30, 0x30, 0x2e, 0x36, 0x33, 0x63, 0xdd, 0x01}
var reader_err = bytes.NewReader(SAMPLE_STATUS_DATA_LENGTH_OVER[2:])
var SAMPLE_STATUS_PACKET_LENGTH_OVER = Packet{
	PacketHeader: PacketHeader{
		Length: 23,
		ID:     0x00,
	},
	Payload: reader_err,
}

var SAMPLE_STATUS_PAYLOAD_PACKET_LENGTH_OVER = Handshake{
	PacketHeader: PacketHeader{
		Length: 23,
		ID:     0x00,
	},
	hostname:        "192.168.100.63",
	ProtocolVersion: 765,
	HostnameLength:  14,
	Port:            25565,
	NextState:       0x01,
	IsFML:           false,
	Tail:            []byte{},
}

///////

var SAMPLE_STATUS_DATA_TAIL = []byte{0x01, 0x00}
var SAMPLE_STATUS_DATA_TAILLING = []byte{0x15, 0x00, 0xfd, 0x05, 0x0e, 0x31, 0x39, 0x32, 0x2e, 0x31, 0x36, 0x38, 0x2e, 0x31, 0x30, 0x30, 0x2e, 0x36, 0x33, 0x63, 0xdd, 0x01, 0x01, 0x00}
var reader_tailling = bytes.NewReader(SAMPLE_STATUS_DATA_2[2:])
var SAMPLE_STATUS_PACKET__TAILLING = Packet{
	PacketHeader: PacketHeader{
		Length: 21,
		ID:     0x00,
	},
	Payload: reader_tailling,
}

var SAMPLE_STATUS_PAYLOAD_PACKET__TAILLING = Handshake{
	PacketHeader: PacketHeader{
		Length: 21,
		ID:     0x00,
	},
	hostname:        "192.168.100.63",
	ProtocolVersion: 765,
	HostnameLength:  14,
	Port:            25565,
	NextState:       0x01,
	IsFML:           false,
	Tail:            []byte{},
}

var LEGACY_STATUS_REQ = []byte{0xfe, 0x01, 0xfa, 0x00, 0x0b, 0x00, 0x4d, 0x00, 0x43, 0x00, 0x7c, 0x00, 0x50, 0x00, 0x69, 0x00, 0x6e, 0x00, 0x67, 0x00, 0x48, 0x00, 0x6f, 0x00, 0x73, 0x00, 0x74, 0x00, 0x19, 0x7f, 0x00, 0x09, 0x00, 0x31, 0x00, 0x32, 0x00, 0x37, 0x00, 0x2e, 0x00, 0x30, 0x00, 0x2e, 0x00, 0x30, 0x00, 0x2e, 0x00, 0x31, 0x00, 0x00, 0x63, 0xdd}
var NORMAL_STATUS_REQ = []byte{0x10, 0x00, 0xfd, 0x05, 0x09, 0x31, 0x32, 0x37, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x63, 0xdd, 0x01}

func Get_NORMAL_STATUS_REQ_PACKET() Packet {
	return Packet{
		PacketHeader: PacketHeader{
			Length:        16,
			ID:            0x00,
			IsOldProtocol: false,
		},
		Payload: bytes.NewReader(NORMAL_STATUS_REQ[4:]),
	}
}

func Get_LEGACY_STATUS_REQ_PACKET() Packet {
	return Packet{PacketHeader: PacketHeader{
		Length:        1,
		ID:            0x00,
		IsOldProtocol: true,
	},
		Payload: bytes.NewReader(LEGACY_STATUS_REQ)}
}

func Get_LEGACY_STATUS_REQ_HANDSHAKE() *Handshake {
	h := NewHandshake()
	h.Decode(LEGACY_STATUS_REQ, true)
	return h
}
