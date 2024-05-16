package state

import (
	"context"
	"errors"
	"log"
	service "mc_reverse_proxy/src/service"
	"net"
	"runtime"
	"sync"
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
type IState interface {
	ImplState()

	Enter(*StateMachine) error
	Action() error
	Exit() IState
}

const (
	STATUS_UNKNOWN = iota
	STATUS_OK
	STATUS_EXIT
	STATUS_CANCLE
	STATUS_ERROR
)

type StateMachine struct {
	Conn            service.Connection
	currentState    IState
	states          map[string]IState
	StateChangeLock *sync.Mutex
	ctx             context.Context
	cancle          context.CancelCauseFunc
}

func (sm *StateMachine) setState(s IState) {
	sm.currentState = s
}

func (sm *StateMachine) Run() error {
	// sm.StateChangeLock.Lock()
	err := sm.currentState.Enter(sm)
	if err != nil {
		return err
	}
	err = sm.currentState.Action()
	if err != nil {
		return err
	}
	return nil
}

func (sm *StateMachine) Transition() int {
	sm.StateChangeLock.Lock()
	sm.currentState = sm.currentState.Exit()
	if sm.currentState == nil {
		// log.Printf("Exit")
		err := sm.Conn.CloseConn()
		if err != nil {
			log.Printf("[state machine] Error when transistioning state %v", err)
		}
		return STATUS_EXIT
	}
	sm.StateChangeLock.Unlock()
	err := sm.currentState.Enter(sm)
	if err != nil {
		return STATUS_ERROR
	}
	// log.Printf("Action")
	err = sm.currentState.Action()
	if err != nil {
		log.Printf("[state machine] Action failed: %v", err)
		if err.Error() == "Context Done" {
			return STATUS_EXIT
		}
		return STATUS_ERROR
	}
	// log.Printf("Action Done")
	return STATUS_OK
}

func (sm *StateMachine) Destroy() {
	sm.cancle(errors.New("force closed"))
	sm.Conn.Cancle(errors.New("force closed"))
	sm.Conn.CloseConn()
	sm.Conn.Destroy()
	runtime.GC()
}

// func (sm *StateMachine) Action() error {
// 	err := sm.currentState.Action()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

type Event struct {
	Type string
	Data map[string]string
}

func NewStateMachine(listener *net.Listener, serverList map[string]map[string]string) *StateMachine {
	ctx, cancle := context.WithCancelCause(context.Background())
	sm := &StateMachine{
		currentState:    &InitState{},
		states:          make(map[string]IState),
		StateChangeLock: &sync.Mutex{},
		ctx:             ctx,
		cancle:          cancle,
	}
	sm.Conn = service.NewConnection(sm.StateChangeLock, sm.ctx, sm.cancle, listener)
	sm.Conn.ServerList = serverList
	sm.Conn.ClientData = make(chan []byte)
	sm.Conn.ServerData = make(chan []byte)
	return sm
}
