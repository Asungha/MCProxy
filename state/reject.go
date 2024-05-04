package state

type RejectState struct {
	connData *ConnectionData
	err      error
}

func (r *RejectState) ImplState() {}

func (r *RejectState) Enter(connData *ConnectionData) {
	r.connData = connData
}

func (r *RejectState) Exit() {}

func (r *RejectState) Validate() error {
	return nil
}

func (r *RejectState) Update(sm *StateMachine) error {
	return r.err
}
