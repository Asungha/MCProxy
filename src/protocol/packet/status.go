package packet

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strings"
)

type DescriptionBuilder struct {
	data []StatusDesExtra
}

func (d *DescriptionBuilder) AddWithColor(text, color string) {
	d.data = append(d.data, StatusDesExtra{Text: text, Color: &color})
}

func (d *DescriptionBuilder) Add(text string) {
	d.data = append(d.data, StatusDesExtra{Text: text})
}

func (d *DescriptionBuilder) Build() []StatusDesExtra {
	return d.data
}

type StatusDesExtra struct {
	Text  string  `json:"text"`
	Color *string `json:"color"`
}

type StatusData struct {
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
		// Sample []struct {
		// 	Name string `json:"name"`
		// 	ID   string `json:"id"`
		// } `json:"sample"`
	} `json:"players"`
	Description struct {
		Extra []StatusDesExtra `json:"extra"`
		Text  string           `json:"text"`
	} `json:"description"`
	Modinfo struct {
		Type    string   `json:"type"`
		ModList []string `json:"modList"`
	} `json:"modinfo"`
}

func (s *StatusData) JSONString() string {
	b, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(b)
}

type Status struct {
	Packet

	Json string
}

func (h *Status) ImplPacketData() {}

func (h *Status) Encode() ([]byte, error) {
	size_buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(size_buf, uint64(len(h.Json)))
	data := []byte(h.Json)
	return append(size_buf[:n], data...), nil
}

func (h *Status) Decode(data []byte, size int) error {
	// Get size of json size field
	_, n := binary.Varint(data)
	if n == 0 {
		return fmt.Errorf("Invalid size field: %d", n)
	}
	// Check if data is json
	h.Json = string((data)[n:])
	if !strings.Contains(h.Json[:20], "{\"version\"") {
		return fmt.Errorf("Invalid json data: %s", h.Json)
	}
	// h.Json = strings.ReplaceAll(h.Json, "\x00", "")
	return nil
}

func (h *Status) Length() int {
	buf := make([]byte, binary.MaxVarintLen64)
	defer func() {
		buf = nil
	}()
	n := binary.PutUvarint(buf, uint64(len(h.Json)))
	return len(h.Json) + n - 2
}

func (h *Status) String() string {
	return fmt.Sprintf("Size: %d, Json: %s", len(h.Json), h.Json)
}

func (h *Status) JSON() (StatusData, error) {
	var obj StatusData
	err := json.Unmarshal([]byte(h.Json), &obj)
	if err != nil {
		return StatusData{}, err
	}
	return obj, nil
}

func (h *Status) SetJSON(obj StatusData) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	h.Json = string(data)
	return nil
}

func (h *Status) Destroy() {
	h.Json = ""
}

type OldStatusReq struct {
	ID              byte   // 1 byte, always 0xfa
	ID_pl           byte   // 1 byte, always 0x01
	Cmd_len         []byte // 2 bytes, always 0x00 0x0b
	Cmd             string // 11 bytes, always "MC|PingHost"
	Payload_len     []byte // 2 bytes, length of the rest of the packet
	ProtocolVersion byte   // 1 byte
	Hostname_len    int    // 2 bytes
	Hostname        string // variable length
	Port            int    // 4 bytes
}

func (h *OldStatusReq) ImplPacketData() {}

func (h *OldStatusReq) Encode() ([]byte, error) {
	buf := make([]byte, 1)
	buf[0] = 0xfa
	buf = append(buf, 0x01)
	buf = append(buf, []byte{0x00, 0x0b}...)
	buf = append(buf, []byte("MC|PingHost")...)
	buf = append(buf, byte(7+len(h.Hostname)))
	buf = append(buf, h.ProtocolVersion)
	buf = append(buf, byte(h.Hostname_len>>8), byte(h.Hostname_len&0xff))
	buf = append(buf, []byte(h.Hostname)...)
	buf = append(buf, byte(h.Port>>24), byte(h.Port>>16), byte(h.Port>>8), byte(h.Port))
	return buf, nil
}

func (h *OldStatusReq) Decode(data []byte, size int) error {
	h.ID = data[0]
	h.ID_pl = data[1]
	h.Cmd_len = data[2:4]
	h.Cmd = string(data[4:15])
	h.Payload_len = data[15:17]
	h.ProtocolVersion = data[17]
	h.Hostname_len = int(binary.BigEndian.Uint16(data[18:20]))
	h.Hostname = string(data[20 : 20+h.Hostname_len])
	h.Port = int(data[20+h.Hostname_len])<<24 | int(data[21+h.Hostname_len])<<16 | int(data[22+h.Hostname_len])<<8 | int(data[23+h.Hostname_len])
	return nil
}

func (h *OldStatusReq) Length() int {
	return 24 + len(h.Hostname)
}

func (h *OldStatusReq) String() string {
	return fmt.Sprintf("ID: 0x%x, ID_pl: 0x%x, Cmd: %s, ProtocolVersion: %d, Hostname: %s, Port: %d", h.ID, h.ID_pl, h.Cmd, h.ProtocolVersion, h.Hostname, h.Port)
}

func (h *OldStatusReq) Destroy() {
	h.Cmd = ""
	h.Hostname = ""
}
