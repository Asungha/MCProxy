package packet

type Login struct {
	Header []byte
	Player string
	loaded bool
}

func (h *Login) ImplPacket() {}

func (h *Login) Serialize() []byte {
	return append(h.Header, []byte(h.Player)...)
}

func (h *Login) Deserialize(data []byte) error {
	h.loaded = true
	h.Header = data[:7]
	h.Player = string(data[7:])
	return nil
}

func (h *Login) IsLoaded() bool {
	return h.loaded
}

func (h *Login) String() string {
	return string(h.Header) + h.Player
}

func (h *Login) Length() int {
	return len(h.Header) + len(h.Player)
}
