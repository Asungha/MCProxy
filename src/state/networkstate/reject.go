package state

// import "log"
import (
	state "mc_reverse_proxy/src/state/state"
)

type RejectState struct {
	state.State
}

// func (r *RejectState) ImplState() {}

// func (r *RejectState) Enter(sm *StateMachine) error {
// 	log.Printf("RejectState")
// 	r.sm = sm
// 	return nil
// }

// func (r *RejectState) Exit() IState {
// 	return nil
// }

// func (r *RejectState) Action() error {
// 	log.Printf("[Reject State] State Rejected : %s", r.Message)
// 	if r.sm != nil {
// 		return r.sm.Conn.CloseConn()
// 	}
// 	return nil
// }
