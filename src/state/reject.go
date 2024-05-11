package state

import "log"

type RejectState struct {
	sm      *StateMachine
	Message string
}

func (r *RejectState) ImplState() {}

func (r *RejectState) Enter(sm *StateMachine) error {
	log.Printf("RejectState")
	r.sm = sm
	return nil
}

func (r *RejectState) Exit() IState {
	return nil
}

func (r *RejectState) Action() error {
	log.Printf("[Reject State] State Rejected : %s", r.Message)
	return r.sm.Conn.CloseConn()
}
