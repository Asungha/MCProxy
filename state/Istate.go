package state

import (
	"context"
	"net"
)

const (
	ACTION_UNDEFINED = iota
	ACTION_TRANSPARENT
	ACTION_REJECT
	ACTION_ACCEPT
	ACTION_CANCLE
)

type Data struct {
	bytes  *[]byte
	length int
}

type ConnectionData struct {
	TargetHostname string
	ClientAddress  string
	ClientConn     *net.Conn
	ServerConn     *net.Conn
	ClientData     chan Data
	ServerData     chan Data
	ctx            context.Context
}

type IState interface {
	ImplState()

	Enter(*ConnectionData)
	Exit()
	Validate() error
	Update(*StateMachine) error
}

type StateMachine struct {
	currentState IState
	states       map[string]IState
	connData     ConnectionData
}

func (sm *StateMachine) setState(s IState) {
	sm.currentState = s
	sm.currentState.Enter(&sm.connData)
}

func (sm *StateMachine) Transition() error {
	err := sm.currentState.Validate()
	if err != nil {
		sm.setState(&RejectState{err: err})
	}
	return sm.currentState.Update(sm)
}

func NewStateMachine(initialState IState, connData ConnectionData) *StateMachine {
	sm := &StateMachine{
		currentState: initialState,
		states:       make(map[string]IState),
	}

	sm.currentState.Enter(&connData)
	return sm
}
