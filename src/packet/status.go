package packet

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type Status struct {
	Json string
}

func (h *Status) ImplPacketData() {}

func (h *Status) Encode() []byte {
	size_buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(size_buf, uint64(len(h.Json)))
	data := []byte(h.Json)
	return append(size_buf[:n], data...)
}

func (h *Status) Decode(data *[]byte) error {
	// Get size of json size field
	_, n := binary.Varint(*data)
	if n == 0 {
		return fmt.Errorf("Invalid size field: %d", n)
	}
	// Check if data is json
	h.Json = string((*data)[1:])
	if !strings.Contains(h.Json[:20], "{\"version\"") {
		return fmt.Errorf("Invalid json data: %s", h.Json)
	}
	// h.Json = strings.ReplaceAll(h.Json, "\x00", "")
	return nil
}

func (h *Status) Length() int {
	return len(h.Json)
}

func (h *Status) String() string {
	return fmt.Sprintf("Size: %d, Json: %s", len(h.Json), h.Json)
}
