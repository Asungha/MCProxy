package common

type PacketType string

const (
	UNKNOWN      PacketType = "unknown"
	MC_HANDSHAKE PacketType = "mc_handshake"
	MC_LOGIN     PacketType = "mc_login"
	MC_GAMMPLAY  PacketType = "mc_gameplay"
	MC_OTHER     PacketType = "mc_other"
	HTTP         PacketType = "http"
)
