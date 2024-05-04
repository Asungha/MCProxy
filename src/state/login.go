package state

import (
	"errors"
)

type LoginRequestState struct {
	connData *ConnectionData
}

func (l *LoginRequestState) ImplState() {}

func (l *LoginRequestState) Enter(connData *ConnectionData) {
	l.connData = connData
}

func (l *LoginRequestState) Exit() {}

func (l *LoginRequestState) Validate() error {
	return nil
}

func (l *LoginRequestState) Update(sm *StateMachine) error {
	sm.setState(&LoginResponseState{})
	return nil
}

type LoginResponseState struct {
	connData *ConnectionData
}

func (l *LoginResponseState) ImplState() {}

func (l *LoginResponseState) Enter(connData *ConnectionData) {
	l.connData = connData
}

func (l *LoginResponseState) Exit() {}

func (l *LoginResponseState) Validate() error {
	return nil
}

func (l *LoginResponseState) Update(sm *StateMachine) error {
	return errors.New("Not implemented")
}
