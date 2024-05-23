package state

import (
	state "mc_reverse_proxy/src/state/state"
)

type PassthroughState struct {
	state.State
}

// func (p *PassthroughState) ImplState() {}

// func (p *PassthroughState) Enter(sm *StateMachine) error {
// 	p.sm = sm
// 	return nil
// }
// func (p *PassthroughState) Exit() IState {
// 	return p
// }

// func (p *PassthroughState) Action() error {
// 	select {
// 	case cData := <-p.sm.Conn.ClientData:
// 		p.sm.StateChangeLock.Lock()
// 		p.sm.Conn.WriteServer(cData)
// 		p.sm.StateChangeLock.Unlock()
// 		return nil
// 	case sData := <-p.sm.Conn.ServerData:
// 		p.sm.StateChangeLock.Lock()
// 		p.sm.Conn.WriteClient(sData)
// 		p.sm.StateChangeLock.Unlock()
// 		return nil
// 	case <-p.sm.Conn.Ctx.Done():
// 		// log.Printf("%v", e)
// 		return errors.New("context Done")
// 	}
// }
