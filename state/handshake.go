package state

// type HandshakeState struct {
// 	conn       *ConnectionData
// 	StateFlag  [2]byte
// 	ClientConn *net.Conn
// 	ServerConn *net.Conn
// }

// func (h *HandshakeState) ImplState() {}

// func (h *HandshakeState) Enter(conn *ConnectionData) {
// 	h.conn = conn
// 	select {
// 	case h.ClientConn = <-h.conn.ClientConn:
// 	}
// }
// func (h *HandshakeState) Exit() {}

// func (h *HandshakeState) Validate() error {
// 	return nil
// }

// func (h *HandshakeState) Update(sm *StateMachine) error {
// 	if h.StateFlag[0] == 0x01 { // Request Status

// 	} else if h.StateFlag[0] == 0x02 { // Request Login
// 		sm.setState(&LoginRequestState{})
// 	}
// 	return nil
// }
