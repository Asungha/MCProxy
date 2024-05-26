package statemachine

import (
	"context"
	"errors"
	"log"
	state "mc_reverse_proxy/src/state/state"
)

type State string
type TransistionPair struct {
	Source      string
	Destination string
}
type Condition struct {
	TransistionPair
	Fx func() bool
}

type IStateMachine interface {
	RegisterState(state.IState) error
	State(State) state.IState
	TransistionCondition(pair TransistionPair, condition func() bool)
	Construct() error
	Destruct() error
	SetRoot(State)
	Run() error
}

type AStateMachine struct {
	currectState state.IState
	States       map[string]state.IState
	Conditions   []Condition
	Ctx          context.Context
	Cancle       context.CancelCauseFunc
	DeferFunc    state.DeferFunction
	IStateMachine
}

func (sm *AStateMachine) RegisterState(stateName State, s state.IState) error {
	if sm.States == nil {
		sm.States = make(map[string]state.IState)
		sm.States[""] = nil
	}
	sm.States[string(stateName)] = s
	return nil
}

func (sm *AStateMachine) TransistionCondition(pair TransistionPair, condition state.ConditionFunction) {
	sm.Conditions = append(sm.Conditions, Condition{Fx: condition, TransistionPair: pair})
}

func (sm *AStateMachine) TransistionFunction(source State, fx state.TransistionFunction) {
	(sm.States[string(source)]).AddTransistionFunction(fx)
}

func (sm *AStateMachine) Construct() error {
	for _, condition := range sm.Conditions {
		c := condition
		var fx state.TransistionFunction = func(i state.IState) (state.IState, error) {
			if c.Fx() {
				return sm.States[c.Destination], nil
			}
			return nil, nil
		}
		(sm.States[c.Source]).AddTransistionFunction(fx)
	}
	return nil
}

func (sm *AStateMachine) Destruct() error {
	sm.Cancle(errors.New("Desturct Called"))
	for _, state := range sm.States {
		if state == nil {
			continue
		}
		(state).Destruct()
	}
	if sm.DeferFunc != nil {
		sm.DeferFunc()
	}
	sm.Conditions = nil
	sm.States = nil
	return nil
}

func (sm *AStateMachine) SetRoot(state State) {
	sm.currectState = sm.States[string(state)]
}

func (sm *AStateMachine) State(s State) state.IState {
	return sm.States[string(s)]
}

func (sm *AStateMachine) Run() error {
	// defer log.Println("[statemachine worker] Thread exit")
	// log.Println("[statemachine worker] start")
	defer sm.Destruct()
	defer sm.DeferFunc()
	for {
		err := (sm.currectState).Enter()
		if err != nil {
			log.Printf("run err %v", err.Error())
			return err
		}

		nextState, err := (sm.currectState).Transition()
		if err != nil {
			return err
		}
		if nextState == nil {
			return nil
		}
		sm.currectState = nextState
	}
}
