package packet

import "fmt"

type Raw struct {
	Data []byte
	L    int
}

func (h *Raw) ImplPacketData() {}

func (h *Raw) String() string {
	return fmt.Sprintf("%x", h.Data)
}

func (h *Raw) Encode() ([]byte, error) {
	return h.Data, nil
}

func (h *Raw) Length() int {
	return h.L
}

func (h *Raw) Decode(data []byte, size int) error {
	h.Data = data
	h.L = size
	return nil
}

func (h *Raw) Destroy() {
	h.Data = nil
}
