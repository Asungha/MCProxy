package state

// import (
// 	"errors"
// 	pac "mc_reverse_proxy/src/packet"
// )

// type LoginRequestState struct {
// 	data pac.Packet[*pac.Raw]
// }

// func (l *LoginRequestState) ImplState() {}

// func (l *LoginRequestState) Enter(data pac.Packet[pac.IPacketData]) error {
// 	res, err := pac.CastPacket[*pac.Raw](data)
// 	if err != nil {
// 		return err
// 	}
// 	l.data = res
// 	return nil
// }

// func (l *LoginRequestState) Exit() {}

// func (l *LoginRequestState) Validate() error {
// 	return nil
// }

// func (l *LoginRequestState) Update(sm *StateMachine, data pac.Packet[pac.IPacketData]) error {
// 	sm.setState(&LoginResponseState{}, data)
// 	return nil
// }

// type LoginResponseState struct {
// 	data pac.Packet[*pac.Raw]
// }

// func (l *LoginResponseState) ImplState() {}

// func (l *LoginResponseState) Enter(data pac.Packet[pac.IPacketData]) error {
// 	res, err := pac.CastPacket[*pac.Raw](data)
// 	if err != nil {
// 		return err
// 	}
// 	l.data = res
// 	return nil
// }

// func (l *LoginResponseState) Exit() {}

// func (l *LoginResponseState) Validate() error {
// 	return nil
// }

// func (l *LoginResponseState) Update(sm *StateMachine, data pac.Packet[pac.IPacketData]) error {
// 	return errors.New("Not implemented")
// }
