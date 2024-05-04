package state

import "errors"

type OKState struct {
	connData *ConnectionData
}

func (o *OKState) ImplState() {}

func (o *OKState) Enter(connData *ConnectionData) {
	o.connData = connData
}

func (o *OKState) Exit() {}

func (o *OKState) Validate() error {
	return nil
}

func (o *OKState) Update(sm *StateMachine) error {
	return errors.New("State Done")
}
