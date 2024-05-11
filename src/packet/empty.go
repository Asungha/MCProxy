package packet

type Empty struct{}

func (h *Empty) ImplPacketData() {}

func (h *Empty) String() string {
	return string("Empty packet")
}

func (h *Empty) Encode() []byte {
	return []byte{}
}

func (h *Empty) Length() int {
	return 0
}

func (h *Empty) Decode(data []byte) error {
	return nil
}

func (h *Empty) Destroy() {}
