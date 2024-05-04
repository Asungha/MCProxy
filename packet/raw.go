package packet

type Raw struct {
	Data *[]byte
}

func (h *Raw) ImplPacketData() {}

func (h *Raw) String() string {
	return string(*h.Data)
}

func (h *Raw) Encode() []byte {
	return *h.Data
}

func (h *Raw) Length() int {
	return len(*h.Data)
}

func (h *Raw) Decode(data *[]byte) error {
	h.Data = data
	return nil
}
