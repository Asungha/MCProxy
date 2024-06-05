package enum

type Action uint

func (a *Action) String() string {
	switch *a {
	case 1:
		return "connect"
	case 2:
		return "handshake"
	case 3:
		return "status"
	case 4:
		return "login"
	case 5:
		return "disconnect"
	default:
		return "unknown"
	}
}

const (
	UNKNOWN    Action = 0
	CONNECT    Action = 1
	HANDSHAKE  Action = 2
	STATUS     Action = 3
	LOGIN      Action = 4
	DISCONNECT Action = 5
)
