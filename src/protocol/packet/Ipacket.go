package packet

type IPacket interface {
	String() string
	Encode() ([]byte, error)
	Decode([]byte) error
	Length() int
	Destroy()
}
