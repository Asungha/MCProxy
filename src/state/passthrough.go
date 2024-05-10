package state

import (
	"errors"
)

type PassthroughState struct {
	sm *StateMachine
}

func (p *PassthroughState) ImplState() {}

func (p *PassthroughState) Enter(sm *StateMachine) error {
	// log.Printf("Passthrough Enter")
	p.sm = sm
	return nil
}
func (p *PassthroughState) Exit() IState {
	return p
}

func (p *PassthroughState) Action() error {
	select {
	case cData := <-p.sm.Conn.ClientData:
		p.sm.Conn.WriteServer(cData)
		return nil
	case sData := <-p.sm.Conn.ServerData:
		p.sm.Conn.WriteClient(sData)
		return nil
	case <-p.sm.Conn.ctx.Done():
		// log.Printf("%v", e)
		return errors.New("Context Done")
	}
}
