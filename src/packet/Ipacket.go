package packet

type IPacketData interface {
	ImplPacketData()
	String() string
	Encode() []byte
	Decode(*[]byte) error
	Length() int
}
