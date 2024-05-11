package packet

type IPacketData interface {
	ImplPacketData()
	String() string
	Encode() ([]byte, error)
	Decode([]byte, int) error
	Length() int
	Destroy()
}
